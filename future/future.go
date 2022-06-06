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

func Join[A any, B any](fa *Future[A], fb *Future[B]) (A, B) {
	var wg sync.WaitGroup
	wg.Add(2)

	var a A
	var b B

	go func() {
		a = fa.Await()
		defer wg.Done()
	}()
	go func() {
		b = fb.Await()
		defer wg.Done()
	}()
	wg.Wait()

	return a, b
}

func Join3[A any, B any, C any](
	fa *Future[A],
	fb *Future[B],
	fc *Future[C],
) (A, B, C) {
	var wg sync.WaitGroup
	wg.Add(2)

	var a A
	var b B
	var c C

	go func() {
		a, b = Join(fa, fb)
		defer wg.Done()
	}()
	go func() {
		c = fc.Await()
		defer wg.Done()
	}()
	wg.Wait()

	return a, b, c
}

func Join4[A any, B any, C any, D any](
	fa *Future[A],
	fb *Future[B],
	fc *Future[C],
	fd *Future[D],
) (A, B, C, D) {
	var wg sync.WaitGroup
	wg.Add(2)

	var a A
	var b B
	var c C
	var d D

	go func() {
		a, b = Join(fa, fb)
		defer wg.Done()
	}()
	go func() {
		c, d = Join(fc, fd)
		defer wg.Done()
	}()
	wg.Wait()

	return a, b, c, d
}

func Join5[A any, B any, C any, D any, E any](
	fa *Future[A],
	fb *Future[B],
	fc *Future[C],
	fd *Future[D],
	fe *Future[E],
) (A, B, C, D, E) {
	var wg sync.WaitGroup
	wg.Add(2)

	var a A
	var b B
	var c C
	var d D
	var e E

	go func() {
		a, b, c = Join3(fa, fb, fc)
		defer wg.Done()
	}()
	go func() {
		d, e = Join(fd, fe)
		defer wg.Done()
	}()
	wg.Wait()

	return a, b, c, d, e
}
