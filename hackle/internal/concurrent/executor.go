package concurrent

import "sync"

type Executor struct {
	wg sync.WaitGroup
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Go(task func()) {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		task()
	}()
}

func (e *Executor) Wait() {
	e.wg.Wait()
}
