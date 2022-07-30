package syncpool_test

import (
	"encoding/json"
	"testing"

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
