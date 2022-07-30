package syncpool_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ymzuiku/async/syncpool"
)

func TestMemoPool(t *testing.T) {
	var studentPool = syncpool.New[Student]()
	var bytePool = syncpool.New[[]byte]()

	for range []int{0, 1, 2, 3, 4, 5} {
		stu := studentPool.Get()
		defer func() {
			studentPool.Put(stu)
		}()
		var buf = bytePool.Get()
		if stu.Name != "" {
			panic("no empty")
		}
		stu.Name = key
		json.Unmarshal(*buf, stu)
	}
}

func TestMemoPoolIsEmpty(t *testing.T) {
	var studentPool = syncpool.New[Student]()

	for i := 0; i < 500; i++ {
		s := studentPool.Get()
		defer func() {
			studentPool.Put(s)
		}()
		assert.Empty(t, s.Name)
		s.Name = "dog"
	}
}

func TestMemoIsEmpty(t *testing.T) {
	var studentPool = syncpool.New[Student]()
	s := studentPool.Get()
	s.Name = "dog"
	assert.NotEmpty(t, s)
	studentPool.Empty(s)
	assert.Empty(t, s)
}
