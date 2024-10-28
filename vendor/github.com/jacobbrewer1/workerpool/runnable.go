package workerpool

// Runnable is the interface objects wishing to be scheduled by this worker pool must conform to.
type Runnable interface {
	Run()
}
