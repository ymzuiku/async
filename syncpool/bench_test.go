package syncpool_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/ymzuiku/async/syncpool"
)

type Student struct {
	Name string
	Age  int
	List [1000]string
}

var key = uuid.NewString()

func BenchmarkUnmarshal(b *testing.B) {
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		var buf []byte
		stu := new(Student)
		stu.Name = key
		json.Unmarshal(buf, stu)
	}
}

var pool = syncpool.New[Student]()
var buffer = syncpool.New[[]byte]()

func theFor() {
	stu := pool.Get()
	defer func() {
		pool.Put(stu)
	}()
	var buf = buffer.Get()
	if stu.Name != "" {
		panic("no empty")
	}
	stu.Name = key
	json.Unmarshal(*buf, stu)
}

func BenchmarkPool(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		theFor()
	}
}
