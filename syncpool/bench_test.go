package syncpool_test

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/ymzuiku/async/syncpool"
)

func Struct(read Read) {
	json.Unmarshal([]byte{}, &read)
}

func Point(read *Read) {
	json.Unmarshal([]byte{}, read)
}

func BenchmarkStruct(b *testing.B) {
	empty := Read{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := Read{}
		r.Join.Table = "aaa"
		r.Sensitives = []string{"aaa"}
		Struct(r)
		r = empty
	}
}

func BenchmarkStructPointFn(b *testing.B) {
	empty := Read{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := Read{}
		r.Join.Table = "aaa"
		r.Sensitives = []string{"aaa"}
		Point(&r)
		r = empty
	}
}

func BenchmarkStuctPoint(b *testing.B) {
	empty := Read{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := &Read{}
		r.Join.Table = "aaa"
		r.Sensitives = []string{"aaa"}
		Point(r)
		*r = empty
	}
}

func BenchmarkStuctPointPoolBase(b *testing.B) {
	b.ResetTimer()
	var read Read
	pool := sync.Pool{
		New: func() any {
			return &Read{}
		},
	}
	for i := 0; i < b.N; i++ {
		r := pool.Get().(*Read)
		r.Join.Table = "aaa"
		r.Sensitives = []string{"aaa"}
		Point(r)
		*r = read
		pool.Put(r)
	}
}

func BenchmarkStuctPointPool(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := readPool.Get()
		r.Join.Table = "aaa"
		r.Sensitives = []string{"aaa"}
		Point(r)
		readPool.Put(r)
	}
}

func BenchmarkStuctPointPoolNotPut(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := readPool.Get()
		r.Join.Table = "aaa"
		r.Sensitives = []string{"aaa"}
		Point(r)
	}
}

func BenchmarkStuctPointPoolDefer(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := readPool.Get()
		defer func() {
			readPool.Put(r)
		}()
		r.Sensitives = []string{"aaa"}
		Point(r)
	}
}

func BenchmarkMap(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := map[string]any{}
		r["Sensitives"] = []string{"aaa"}
		json.Unmarshal([]byte{}, &r)
	}
}

func BenchmarkMapPool(b *testing.B) {
	b.ResetTimer()
	mapPool := syncpool.Map()
	for i := 0; i < b.N; i++ {
		r := mapPool.Get()
		r["Sensitives"] = []string{"aaa"}
		json.Unmarshal([]byte{}, &r)
		mapPool.Put(r)
	}
}
