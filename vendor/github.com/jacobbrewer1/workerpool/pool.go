package workerpool

import (
	"runtime"
	"sync"
)

const (
	minWorkers     = 1
	minQueueLength = 0
)

type Pool interface {
	// MustSchedule schedules a task to be executed
	// by the worker pool. It panics if the worker pool is stopped.
	MustSchedule(task Runnable)

	// Schedule schedules a task to be executed by the worker pool.
	// It returns an error if the worker pool is stopped.
	Schedule(task Runnable) error

	// BlockingSchedule schedules a task to be executed by the
	// worker pool. It blocks until the task is scheduled.
	BlockingSchedule(task Runnable) error

	// Stop stops the worker pool.
	Stop()

	// StopAsync stops the worker pool asynchronously.
	StopAsync()
}

type pool struct {
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

// New creates a new worker pool with the given options.
//
// By default, the worker pool will have the same number of workers as the number of CPUs.
// The maximum queue length is set to the total number of workers multiplied by 1000.
//
// You can customize the worker pool by providing the options.
func New(opts ...WorkerOption) Pool {
	p := &pool{
		totalWorkers:   runtime.NumCPU(),
		maxQueueLength: runtime.NumCPU() * 1000,
		wg:             new(sync.WaitGroup),
		done:           make(chan struct{}),
		delayedStart:   false,
		started:        false,
	}

	for _, opt := range opts {
		opt(p)
	}

	if p.totalWorkers < minWorkers {
		p.totalWorkers = runtime.NumCPU()
	}

	taskChan := make(chan Runnable, p.maxQueueLength)
	if p.maxQueueLength <= minQueueLength {
		// The user has set the max queue length to the minQueueLength (0), which means we should use a blocking channel (non-buffered).
		taskChan = make(chan Runnable)
	}
	p.tasks = taskChan

	if !p.delayedStart {
		p.start()
	}

	return p
}
