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

type Join struct {
	Table        string         `json:"table" validate:"required"`
	Where        map[string]any `json:"where"`
	On           []string       `json:"on" validate:"required"`
	Comparable   [][]any        `json:"comparable"`
	AllowColumns []string       `json:"allowColumns"`
	Sensitives   []string       `json:"sensitives"`
}

type Read struct {
	LoadHardDelete bool           `json:"loadHardDelete"`
	Limit          int            `json:"limit"`
	Offset         int            `json:"offset"`
	Where          map[string]any `json:"where"`
	Id             string         `json:"id"`
	Ids            []any          `json:"ids"`
	Order          string         `json:"order"`
	Desc           bool           `json:"desc"`
	Total          bool           `json:"total"`
	Join           Join           `json:"join"`
	Join2          *Join          `json:"join2"`
	AllowColumns   []string       `json:"allowColumns"`
	Sensitives     []string       `json:"sensitives"`
	Comparable     [][]any        `json:"comparable"` // 可比较的
}

var readPool = syncpool.New[Read]()
var studentPool = syncpool.New[Student]()

func TestMemoIsEmpty(t *testing.T) {
	s := studentPool.Get()
	s.Name = "dog"
	assert.NotEmpty(t, s)
	studentPool.Empty(s)
	assert.Empty(t, s)

	r := readPool.Get()
	r.Join = Join{
		Table: "aaaaaa",
	}
	r.Where = map[string]any{"aaa": 2}
	r.Join2 = &Join{
		Table: "aaaaaa",
	}
	r.Sensitives = []string{"aaa"}

	assert.NotEmpty(t, r)
	readPool.Empty(r)
	assert.Empty(t, r)
}
