package promise_test

import (
	"sync"
	"testing"

	"github.com/ymzuiku/async/promise"
)

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

func BenchmarkSomePromise(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		read, _ := promise.New(func(resolve func(Read), _ func(error)) {
			resolve(Read{
				Total: true,
				Id:    "cat",
			})
		}).Await()
		if read.Id != "cat" {
			panic("id error")
		}
	}
}

func BenchmarkOnlyWg(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var read Read
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			read = Read{
				Id: "dog",
			}
		}()

		wg.Wait()
		if read.Id != "dog" {
			panic("id error")
		}
	}
}
