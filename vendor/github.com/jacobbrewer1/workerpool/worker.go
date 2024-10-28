package workerpool

import (
	"errors"
	"runtime"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// ErrQueueFull is returned when the worker pool queue is full.
	ErrQueueFull = errors.New("worker pool queue is full")

	// ErrWorkerPoolStopped is returned when the worker pool is stopped.
	ErrWorkerPoolStopped = errors.New("worker pool stopped")
)

const (
	minWorkers     = 1
	minQueueLength = 0
)

type WorkerPool struct {
	// wg is the wait group for the worker pool.
	wg *sync.WaitGroup

	// done is the channel to signal the worker pool to stop.
	done chan struct{}

	// tasks is the channel to send tasks to the worker pool.
	tasks chan Runnable

	// delayedStart is the flag to indicate if the worker pool should start immediately.
	//
	// If delayedStart is set to true, the worker pool will start only when the first task is scheduled.
	// This is useful when you want to start the worker pool only when you have tasks to schedule.
	delayedStart bool

	// started is the flag to indicate if the worker pool has started.
	started bool

	// totalWorkers is the total number of workers to deploy to the pool.
	totalWorkers int

	// maxQueueLength is the maximum number of tasks that can be scheduled.
	maxQueueLength int
}

// NewWorkerPool creates a new worker pool with the given options.
//
// By default, the worker pool will have the same number of workers as the number of CPUs.
// The maximum queue length is set to the total number of workers multiplied by 1000.
//
// You can customize the worker pool by providing the options.
func NewWorkerPool(opts ...WorkerOption) *WorkerPool {
	pool := &WorkerPool{
		totalWorkers:   runtime.NumCPU(),
		maxQueueLength: runtime.NumCPU() * 1000,
		wg:             new(sync.WaitGroup),
		done:           make(chan struct{}),
		delayedStart:   false,
		started:        false,
	}

	for _, opt := range opts {
		opt(pool)
	}

	if pool.totalWorkers < minWorkers {
		pool.totalWorkers = runtime.NumCPU()
	}

	taskChan := make(chan Runnable, pool.maxQueueLength)
	if pool.maxQueueLength <= minQueueLength {
		// The user has set the max queue length to the minQueueLength (0), which means we should use a blocking channel (non-buffered).
		taskChan = make(chan Runnable)
	}
	pool.tasks = taskChan

	if !pool.delayedStart {
		pool.start()
	}

	return pool
}

// TotalWorkers gets the total number of workers in the worker pool.
func (p *WorkerPool) TotalWorkers() int {
	return p.totalWorkers
}

// PendingTasks gets the number of pending tasks in the worker pool.
func (p *WorkerPool) PendingTasks() int {
	return len(p.tasks)
}

// start starts the worker pool by deploying the workers to the pool.
func (p *WorkerPool) start() {
	if p.started {
		return
	}

	p.started = true
	idleWorkers.Set(float64(p.totalWorkers))
	for i := 0; i < p.totalWorkers; i++ {
		p.wg.Add(1)
		go p.deployWorker()
	}
}

// deployWorker deploys a worker to the worker pool. It listens for tasks to execute on the receiving tasks channel.
// It decrements the wait group when the worker is done executing the tasks.
func (p *WorkerPool) deployWorker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.done:
			return
		case task, ok := <-p.tasks:
			if !ok {
				return
			}

			p.runTask(task)
		}
	}
}

// runTask runs the task. It increments the active workers and decrements the idle workers when the task is running.
// It also updates the pending tasks metric.
func (p *WorkerPool) runTask(task Runnable) {
	idleWorkers.Dec()
	activeWorkers.Inc()
	defer func() {
		idleWorkers.Inc()
		activeWorkers.Dec()
	}()

	pendingTasks.Set(float64(len(p.tasks)))
	defer pendingTasks.Set(float64(len(p.tasks)))

	t := prometheus.NewTimer(taskDuration)
	defer t.ObserveDuration()

	task.Run()
}

// Stop stops the worker pool.
//
// Note: This is a blocking operation. It will wait for all the workers to finish executing the tasks.
func (p *WorkerPool) Stop() {
	defer func() {
		idleWorkers.Set(0)
		activeWorkers.Set(0)
		pendingTasks.Set(0)

		p.started = false
	}()

	// Clear down the tasks channel while waiting for the workers to finish
	go func() {
		for range p.tasks {
			// Drain the tasks channel to avoid deadlock
		}
	}()

	close(p.done)
	p.wg.Wait()
}

// StopAsync stops the worker pool asynchronously. It will not wait for all the workers to finish executing the tasks.
//
// Note: This is a non-blocking operation. It will return immediately and stop the worker pool in the background.
func (p *WorkerPool) StopAsync() {
	go p.Stop()
}

// isDone returns true if the worker pool is stopped.
func (p *WorkerPool) isDone() bool {
	select {
	case <-p.done:
		return true
	default:
		return false
	}
}

// MustSchedule schedules a task to be executed by the worker pool.
//
// Note: This is a non-blocking operation. If the worker pool is stopped, it will panic.
func (p *WorkerPool) MustSchedule(task Runnable) {
	if err := p.Schedule(task); err != nil {
		panic(err)
	}
}

// Schedule schedules a task to be executed by the worker pool.
//
// Note: This is a non-blocking operation. If the worker pool is stopped, it will return an error.
func (p *WorkerPool) Schedule(task Runnable) error {
	if p.isDone() {
		return ErrWorkerPoolStopped
	}

	if p.delayedStart && !p.started {
		p.start()
	}

	select {
	case <-p.done:
		return ErrWorkerPoolStopped
	case p.tasks <- task:
		return nil
	default:
		return ErrQueueFull
	}
}

// BlockingSchedule schedules a task to be executed by the worker pool.
//
// Note: This is a blocking operation. If the worker pool is stopped, it will return an error.
func (p *WorkerPool) BlockingSchedule(task Runnable) error {
	if p.isDone() {
		return ErrWorkerPoolStopped
	}

	if p.delayedStart && !p.started {
		p.start()
	}

	select {
	case <-p.done:
		return ErrWorkerPoolStopped
	case p.tasks <- task:
		return nil
	}
}
