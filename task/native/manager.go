package antask

import (
	"sync"
	"time"

	"github.com/alioth-center/infrastructure/grace"
	"github.com/alioth-center/infrastructure/task"
	"github.com/alioth-center/infrastructure/utils/concurrency"
)

type nativeManager struct {
	tickInterval time.Duration
	tasks        concurrency.Slice[*nativeTask]
	closer       chan struct{}
}

func NewNativeManager(tickInterval time.Duration) task.Manager {
	manager := &nativeManager{
		tickInterval: tickInterval,
		tasks:        concurrency.NewSlice[*nativeTask](),
		closer:       make(chan struct{}, 1),
	}
	grace.RegisterGraceful(manager)

	return manager
}

func NewGracefulNativeManager(tickInterval time.Duration) task.GracefulManager {
	return &nativeManager{
		tickInterval: tickInterval,
		tasks:        concurrency.NewSlice[*nativeTask](),
		closer:       make(chan struct{}, 1),
	}
}

func (m *nativeManager) ListenAndServe() {
	for {
		select {
		case <-m.closer:
			return
		case <-time.After(m.tickInterval):
			ts, wg, next := time.Now(), &sync.WaitGroup{}, concurrency.NewSlice[*nativeTask]()
			wg.Add(m.tasks.Length())
			for _, nt := range m.tasks.Items() {
				go func(nt *nativeTask, wg *sync.WaitGroup) {
					defer wg.Done()
					if !ts.After(nt.nextScheduled) {
						next.Append(nt)

						return
					}

					if nt.status == task.StatusCancelled {
						return
					}

					nt.execute()
					switch nt.status {
					case task.StatusRetrying, task.StatusLooped:
						if nt.nextScheduled.After(ts.Add(m.tickInterval)) {
							next.Append(nt)
						}
					}
				}(nt, wg)
			}

			wg.Wait()
			m.tasks = next
		}
	}
}

func (m *nativeManager) GracefulClose() {
	m.closer <- struct{}{}
	for _, nt := range m.tasks.Items() {
		nt.cancel()
	}
}

func (m *nativeManager) Execute(fn func() bool, retryInterval task.RetryInterval) (cancel func()) {
	return m.ExecuteWithTriggers(fn, retryInterval)
}

func (m *nativeManager) ExecuteLoop(fn func(), loopInterval time.Duration) (cancel func()) {
	return m.ExecuteLoopWithTriggers(fn, loopInterval)
}

func (m *nativeManager) ExecuteWithTriggers(fn func() bool, retryInterval task.RetryInterval, triggers ...task.Trigger) (cancel func()) {
	t := &nativeTask{
		mutex:         sync.Mutex{},
		execFunc:      fn,
		retries:       0,
		interval:      retryInterval,
		nextScheduled: time.Time{},
		status:        task.StatusCreated,
		triggers:      make(map[task.Status]func()),
	}
	for _, trigger := range triggers {
		status, triggerFn := trigger.Trigger()
		if triggerFn != nil {
			t.triggers[status] = triggerFn
		}
	}

	m.tasks.Append(t)

	return t.cancel
}

func (m *nativeManager) ExecuteLoopWithTriggers(fn func(), loopInterval time.Duration, triggers ...task.Trigger) (cancel func()) {
	t := &nativeTask{
		mutex:         sync.Mutex{},
		execFunc:      func() (success bool) { fn(); return true },
		retries:       0,
		interval:      task.LoopInterval(loopInterval),
		nextScheduled: time.Time{},
		status:        task.StatusCreated,
	}
	for _, trigger := range triggers {
		status, triggerFn := trigger.Trigger()
		if triggerFn != nil {
			t.triggers[status] = triggerFn
		}
	}

	m.tasks.Append(t)

	return t.cancel
}
