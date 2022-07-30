# async

安全的使用异步库

## sync.pool

```go
type Dog struct {
  Name string
  Age int
}
var pool = syncpool.New[Dog]()

func main(){
  dog := pool.Get()
  defer func(){
    // pool.Put 会还原 dog 的值为初始值
    pool.Put(dog)
  }()
  dog.Name = "the name"
}
```

## promise

fork 自: https://github.com/chebyrash/promise

差异: Promise.Then 改为 : \*promise.Then

```go
func findFactorial(n int) int {
	if n == 1 {
		return 1
	}
	return n * findFactorial(n-1)
}

func main(){
  p := promise.New(func(resolve func(int), _ func(error)) {
  	factorial := findFactorial(10)
  	resolve(factorial)
  })

  num, err := p.Await()
  assert.ErrorIs(t, err, nil)
  assert.Greater(t, num, 0)
}

```
