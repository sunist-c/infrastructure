package antask

import (
	"github.com/alioth-center/infrastructure/task"
	"sync"
	"time"
)

type nativeTask struct {
	mutex         sync.Mutex
	execFunc      func() (success bool)
	retries       int
	interval      task.RetryInterval
	nextScheduled time.Time
	status        task.Status
	loop          bool
}

func (t *nativeTask) execute() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !t.loop {

		if t.retries == 0 {
			t.status = task.StatusRunning
		}

		if t.execFunc() {
			t.status = task.StatusSucceeded

			return
		}

		t.retries++
		needRetry, nextInterval := t.interval(t.retries)
		if !needRetry {
			t.status = task.StatusFailed

			return
		}

		t.status = task.StatusRetrying
		t.nextScheduled = time.Now().Add(nextInterval)
	} else {
		t.status = task.StatusLooped
		t.execute()
		_, nextInterval := t.interval(t.retries)
		t.nextScheduled = time.Now().Add(nextInterval)
	}
}

func (t *nativeTask) cancel() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.status = task.StatusCancelled
}

func (t *nativeTask) getStatus() task.Status {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.status
}
