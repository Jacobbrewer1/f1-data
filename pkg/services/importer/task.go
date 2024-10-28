package importer

import (
	"github.com/jacobbrewer1/workerpool"
)

type task struct {
	runnable func()
}

func newTask(runnable func()) workerpool.Runnable {
	return &task{
		runnable: runnable,
	}
}

func (t *task) Run() {
	if t.runnable == nil {
		return
	}

	t.runnable()
}
