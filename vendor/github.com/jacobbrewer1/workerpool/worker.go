package workerpool

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// ErrQueueFull is returned when the worker pool queue is full.
	ErrQueueFull = errors.New("worker pool queue is full")

	// ErrWorkerPoolStopped is returned when the worker pool is stopped.
	ErrWorkerPoolStopped = errors.New("worker pool stopped")
)

// start starts the worker pool by deploying the workers to the pool.
func (p *pool) start() {
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
func (p *pool) deployWorker() {
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
func (p *pool) runTask(task Runnable) {
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
func (p *pool) Stop() {
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
func (p *pool) StopAsync() {
	go p.Stop()
}

// isDone returns true if the worker pool is stopped.
func (p *pool) isDone() bool {
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
func (p *pool) MustSchedule(task Runnable) {
	if err := p.Schedule(task); err != nil {
		panic(err)
	}
}

// Schedule schedules a task to be executed by the worker pool.
//
// Note: This is a non-blocking operation. If the worker pool is stopped, it will return an error.
func (p *pool) Schedule(task Runnable) error {
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
func (p *pool) BlockingSchedule(task Runnable) error {
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
