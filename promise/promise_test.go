package promise_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ymzuiku/async/promise"
)

func findFactorial(n int) int {
	if n == 1 {
		return 1
	}
	return n * findFactorial(n-1)
}

func TestPromise(t *testing.T) {
	p1 := promise.New(func(resolve func(int), _ func(error)) {
		factorial := findFactorial(10)
		resolve(factorial)
	})

	num, err := p1.Await()
	assert.ErrorIs(t, err, nil)
	assert.Greater(t, num, 0)

	p2 := promise.New(func(resolve func(int), _ func(error)) {
		resolve(findFactorial(20))
	})

	nums, err := promise.All(p1, p2).Then(func(data []int) []int {
		log.Printf("__debug__%v", data)
		return data
	}).Await()

	assert.ErrorIs(t, err, nil)
	assert.Len(t, nums, 2)
	assert.Greater(t, nums[1], 100)
}
