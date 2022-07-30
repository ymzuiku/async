package poolx

import (
	"reflect"
	"sync"
)

// Pool
//
// sync.Pool 的包裹
type Pool[T any] struct {
	empty T
	pool  sync.Pool
}

// New
//
// 创建一个 sync.Pool 对象
// - 利用了泛型减少了很多模版代码, 如类型断言和获取
// - Put的时候自动清空数据
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

// IsEmpty
//
// 利用了反射, 性能很低, 勿滥用
func IsEmpty(obj any) bool {
	// get nil case out of the way
	if obj == nil {
		return true
	}

	objValue := reflect.ValueOf(obj)

	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	// pointers are empty if nil or if the value they point to is empty
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return IsEmpty(deref)
	// for all other types, compare against the zero value
	// array types are empty when they match their zero-initialized state
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(obj, zero.Interface())
	}
}
