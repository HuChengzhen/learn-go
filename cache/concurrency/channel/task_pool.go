package channel

import (
	"context"
	"sync"
)

type Task func()

type TaskPool struct {
	tasks chan Task
	//close atomic.Bool
	close     chan struct{}
	closeOnce sync.Once
}

// NewTaskPool numG 是goroutine的数量
// capacity 是缓存的容量
func NewTaskPool(numG int, capacity int) *TaskPool {
	res := &TaskPool{
		tasks: make(chan Task, capacity),
		//close: atomic.Bool{},
		close: make(chan struct{}),
	}
	for i := 0; i < numG; i++ {
		go func() {
			for {
				select {
				case <-res.close:
					return
				case t := <-res.tasks:
					t()
				}
			}

			//for t := range res.tasks {
			//	if res.close.Load() {
			//		return
			//	}
			//	t()
			//}
		}()
	}

	return res
}

func (t *TaskPool) Close() error {
	//t.close.Store(true)
	t.closeOnce.Do(func() {
		close(t.close)
	})
	return nil
}

func (p *TaskPool) Submit(ctx context.Context, t Task) error {
	select {
	case p.tasks <- t:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
