package future

import "sync"

type status int

const (
	pending status = iota
	fulfilled
)

type Future[T any] struct {
	Fn func() T

	value  T
	status status
}

func New[T any](worker func() T) *Future[T] {
	return &Future[T]{
		Fn: worker,
	}
}

func Ready[T any](value T) *Future[T] {
	return New(func() T {
		return value
	})
}

func (f *Future[T]) Await() T {
	if f.status == fulfilled {
		return f.value
	}

	value := f.Fn()
	f.status = fulfilled
	return value
}

func Map[A any, B any](f *Future[A], mapper func(A) B) *Future[B] {
	return New(func() B {
		a := f.Await()
		return mapper(a)
	})
}

func AndThen[A any, B any](f *Future[A], mapper func(A) *Future[B]) *Future[B] {
	return Map(f, mapper).Await()
}

func All[A any](fs []*Future[A]) *Future[[]A] {
	return New(func() []A {
		var wg sync.WaitGroup
		wg.Add(len(fs))

		as := make([]A, len(fs))
		for i, f := range fs {
			go func(i int, f *Future[A]) {
				as[i] = f.Await()
				wg.Done()
			}(i, f)
		}

		wg.Wait()

		return as
	})
}
