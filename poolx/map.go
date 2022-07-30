package poolx

import "sync"

type MapPool struct {
	pool sync.Pool
}

// 创建一个类型为 map[string]any 的 sync.Pool
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

// 清空所有key, 并且放回 sync.Pool
func (m *MapPool) Put(obj map[string]any) {
	for k := range obj {
		delete(obj, k)
	}
	m.pool.Put(obj)
}
