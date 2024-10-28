package workerpool

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// idleWorkers is the number of idle workers in the pool.
	idleWorkers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "workerpool_idle_workers",
		Help: "Number of idle workers in the worker pool",
	})

	// activeWorkers is the number of active workers in the pool.
	activeWorkers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "workerpool_active_workers",
		Help: "Number of active workers in the worker pool",
	})

	// pendingTasks is the number of tasks pending in the worker pool.
	pendingTasks = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "workerpool_pending_tasks",
		Help: "Number of tasks pending in the worker pool",
	})

	// taskDuration is the duration of the task.
	taskDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "workerpool_task_duration_seconds",
		Help: "Duration of the task",
	})
)
