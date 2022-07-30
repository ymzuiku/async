package syncpool

import "sync"

type Pool[T any] struct {
	empty T
	pool  sync.Pool
}

func New[T any]() Pool[T] {
	return Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return new(T)
			},
		},
	}
}

func (s *Pool[T]) Get() *T {
	return s.pool.Get().(*T)
}

func (s *Pool[T]) Put(obj *T) {
	*obj = s.empty
	s.pool.Put(obj)
}

func (s *Pool[T]) Empty(obj *T) {
	*obj = s.empty
}
