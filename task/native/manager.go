package antask

import (
	"github.com/alioth-center/infrastructure/task"
	"github.com/alioth-center/infrastructure/utils/concurrency"
	"sync"
	"time"
)

type nativeManager struct {
	tickInterval time.Duration
	tasks        concurrency.Slice[*nativeTask]
	closer       chan struct{}
}

func (m *nativeManager) ListenAndServe() {
	for {
		select {
		case <-m.closer:
			return
		case <-time.After(m.tickInterval):
			ts := time.Now()
			for _, nt := range m.tasks.Items() {
				go func(nt *nativeTask, sg *sync.WaitGroup) {
					defer sg.Done()
					if nt.status == task.StatusCancelled {
						return
					}

					nt.execute()
					switch nt.status {
					case task.StatusRetrying, task.StatusLooped:
						if nt.nextScheduled.After(ts.Add(m.tickInterval)) {
							lazy.Append(nt)
						} else {
							next.Append(nt)
						}
					}
				}(nt, sg)
			}
		}
	}
}

func (m *nativeManager) GracefulClose() {
	//TODO implement me
	panic("implement me")
}

func (m *nativeManager) Execute(fn func() bool, retryInterval task.RetryInterval) (cancel func()) {
	t := &nativeTask{
		mutex:         sync.Mutex{},
		execFunc:      fn,
		retries:       0,
		interval:      retryInterval,
		nextScheduled: time.Time{},
		status:        task.StatusCreated,
	}
	m.tasks.Append(t)

	return t.cancel
}

func (m *nativeManager) ExecuteLoop(fn func(), loopInterval time.Duration) (cancel func()) {
	t := &nativeTask{
		mutex:         sync.Mutex{},
		execFunc:      func() (success bool) { fn(); return true },
		retries:       0,
		interval:      task.LoopInterval(loopInterval),
		nextScheduled: time.Time{},
		status:        task.StatusCreated,
	}
	m.tasks.Append(t)

	return t.cancel
}

func (m *nativeManager) execute(wg *sync.WaitGroup) (next concurrency.Slice[*nativeTask], lazy concurrency.Slice[*nativeTask]) {
	defer wg.Done()

	next, lazy = concurrency.NewSlice[*nativeTask](), concurrency.NewSlice[*nativeTask]()
	ts, sg := time.Now(), &sync.WaitGroup{}
	sg.Add(m.nextExecute.Length())

	return next, lazy
}
