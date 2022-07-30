package poolx_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ymzuiku/async/poolx"
)

func TestMemoPool(t *testing.T) {
	var studentPool = poolx.New[Student]()
	var bytePool = poolx.New[[]byte]()

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
	var studentPool = poolx.New[Student]()

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
	Table        string
	Where        map[string]any
	On           []string
	Comparable   [][]any
	AllowColumns []string
	Sensitives   []string
}

type Read struct {
	LoadHardDelete bool
	Limit          int
	Offset         int
	Where          map[string]any
	Id             string
	Ids            []any
	Order          string
	Desc           bool
	Total          bool
	Join           Join
	Join2          *Join
	AllowColumns   []string
	Sensitives     []string
	Comparable     [][]any
}

type Student struct {
	Name string
	Age  int
	List [1000]string
}

var key = uuid.NewString()

var readPool = poolx.New[Read]()
var studentPool = poolx.New[Student]()

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
