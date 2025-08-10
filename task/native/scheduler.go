package antask

import (
	"sync"
	"time"

	"github.com/alioth-center/infrastructure/task"
)

type nativeTask struct {
	mutex         sync.Mutex
	execFunc      func() (success bool)
	retries       int
	interval      task.RetryInterval
	nextScheduled time.Time
	status        task.Status
	loop          bool
	triggers      map[task.Status]func()
}

func (t *nativeTask) execute() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !t.loop {
		if t.retries == 0 {
			t.status = task.StatusRunning
			if trigger, exist := t.triggers[task.StatusRunning]; exist {
				go trigger()
			}
		}

		if t.execFunc() {
			t.status = task.StatusSucceeded
			if trigger, exist := t.triggers[task.StatusSucceeded]; exist {
				go trigger()
			}

			return
		}

		t.retries++
		needRetry, nextInterval := t.interval(t.retries)
		if !needRetry {
			t.status = task.StatusFailed
			if trigger, exist := t.triggers[task.StatusFailed]; exist {
				trigger()
			}

			return
		}

		t.status = task.StatusRetrying
		if trigger, exist := t.triggers[task.StatusRetrying]; exist {
			go trigger()
		}

		t.nextScheduled = time.Now().Add(nextInterval)
	} else {
		t.execFunc()
		t.status = task.StatusLooped
		if trigger, exist := t.triggers[task.StatusLooped]; exist {
			go trigger()
		}

		_, nextInterval := t.interval(t.retries)
		t.nextScheduled = time.Now().Add(nextInterval)
	}
}

func (t *nativeTask) cancel() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.status = task.StatusCancelled
	if trigger, exist := t.triggers[task.StatusCancelled]; exist {
		go trigger()
	}
}

func (t *nativeTask) getStatus() task.Status {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.status
}
