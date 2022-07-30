package promise

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

type Promise[T any] struct {
	res     T
	err     error
	pending bool
	mu      *sync.Mutex
	wg      *sync.WaitGroup
}

func New[T any](exec func(resolve func(T), reject func(error))) *Promise[T] {
	if exec == nil {
		panic("executor cannot be nil")
	}
	p := &Promise[T]{
		pending: true,
		mu:      &sync.Mutex{},
		wg:      &sync.WaitGroup{},
	}
	p.wg.Add(1)
	go func() {
		defer p.handlePanic()
		exec(p.resolve, p.reject)
	}()

	return p
}

func (p *Promise[T]) resolve(res T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.pending {
		return
	}

	p.res = res
	p.pending = false

	p.wg.Done()
}

func (p *Promise[T]) reject(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.pending {
		return
	}
	p.err = err
	p.pending = false
	p.wg.Done()
}

func (p *Promise[T]) handlePanic() {
	err := recover()
	if validErr, ok := err.(error); ok {
		p.reject(validErr)
	} else {
		p.reject(fmt.Errorf("%+v", err))
	}
}

func (p *Promise[T]) Await() (T, error) {
	p.wg.Wait()
	return p.res, p.err
}

func (p *Promise[T]) Then(resolveA func(data T) T) *Promise[T] {
	return New(func(resolve func(T), reject func(error)) {
		res, err := p.Await()
		if err != nil {
			reject(err)
			return
		}
		resolve(resolveA(res))
	})
}

func (p *Promise[T]) Catch(rejection func(err error) error) *Promise[T] {
	return New(func(resolve func(T), reject func(error)) {
		res, err := p.Await()
		if err != nil {
			reject(rejection(err))
			return
		}
		resolve(res)
	})
}

func Resolve[T any](resolution T) *Promise[T] {
	return &Promise[T]{
		res:     resolution,
		pending: false,
		mu:      &sync.Mutex{},
		wg:      &sync.WaitGroup{},
	}
}

func Reject[T any](err error) *Promise[T] {
	return &Promise[T]{
		err:     err,
		pending: false,
		mu:      &sync.Mutex{},
		wg:      &sync.WaitGroup{},
	}
}

// func Then[T, O any](p *Promise[T], resolveA func(data T) O) *Promise[O] {
// 	return New(func(resolve func(O), reject func(error)) {
// 		res, err := p.Await()
// 		if err != nil {
// 			reject(err)
// 			return
// 		}
// 		resolve(resolveA(res))
// 	})
// }

type tuple[D, I any] struct {
	data  D
	index I
}

func All[T any](promises ...*Promise[T]) *Promise[[]T] {
	if len(promises) == 0 {
		return nil
	}

	return New(func(resolve func([]T), reject func(error)) {
		length := len(promises)
		valsCh := make(chan tuple[T, int], length)
		errsCh := make(chan error, 1)
		for idx, p := range promises {
			idx := idx
			_ = p.Then(func(data T) T {
				valsCh <- tuple[T, int]{data: data, index: idx}
				return data
			})
			_ = p.Catch(func(err error) error {
				errsCh <- err
				return err
			})
		}

		resolutions := make([]T, length)
		for idx := 0; idx < length; idx++ {
			select {
			case val := <-valsCh:
				resolutions[val.index] = val.data
			case err := <-errsCh:
				reject(err)
				return
			}
		}

		resolve(resolutions)
	})
}

func Race[T any](promises ...*Promise[T]) *Promise[T] {
	if len(promises) == 0 {
		return nil
	}

	return New(func(resolve func(T), reject func(error)) {
		valsCh := make(chan T, 1)
		errsCh := make(chan error, 1)
		for _, p := range promises {
			_ = p.Then(func(data T) T {
				valsCh <- data
				return data
			})
			_ = p.Catch(func(err error) error {
				errsCh <- err
				return err
			})
		}

		select {
		case v := <-valsCh:
			resolve(v)
		case err := <-errsCh:
			reject(err)
		}
	})
}

func Any[T any](promises ...*Promise[T]) *Promise[T] {
	if len(promises) == 0 {
		return nil
	}

	return New(func(resolve func(T), reject func(error)) {
		valsCh := make(chan T, 1)
		errsCh := make(chan tuple[error, int], 1)
		for i, p := range promises {
			i := i
			_ = p.Then(func(data T) T {
				valsCh <- data
				return data
			})
			_ = p.Catch(func(err error) error {
				errsCh <- tuple[error, int]{data: err, index: i}
				return err
			})
		}

		errs := make([]error, len(promises))
		for i := 0; i < len(promises); i++ {
			select {
			case v := <-valsCh:
				resolve(v)
			case err := <-errsCh:
				errs[err.index] = err.data
			}
		}

		var err error

		for _, v := range errs[1:] {
			if v != nil {
				err = errors.Wrap(v, err.Error())
			}
		}

		reject(err)
	})
}
