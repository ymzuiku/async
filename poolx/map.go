package poolx

import "sync"

type MapPool struct {
	pool sync.Pool
}

func Map() MapPool {
	return MapPool{
		pool: sync.Pool{
			New: func() any {
				return map[string]any{}
			},
		},
	}
}

func (m *MapPool) Get() map[string]any {
	return m.pool.Get().(map[string]any)
}

func (m *MapPool) Put(obj map[string]any) {
	for k := range obj {
		delete(obj, k)
	}
	m.pool.Put(obj)
}
