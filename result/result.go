package result

type kind int

const (
	errKind kind = iota
	okKind
)

type Result[T any] struct {
	err error
	ok  T
}

func Ok[T any](v T) *Result[T] {
	return &Result[T]{
		ok: v,
	}
}

func Err[T any](err error) *Result[T] {
	return &Result[T]{
		err: err,
	}
}

func Map[A any, B any](r *Result[A], mapper func(A) B) *Result[B] {
	if r.err != nil {
		return Err[B](r.err)
	}

	return Ok(mapper(r.ok))
}

func Flatten[A any, B any](r *Result[*Result[B]]) *Result[B] {
	if r.err != nil {
		return Err[B](r.err)
	}

	return r.ok
}

func AndThen[A any, B any](r *Result[A], mapper func(A) *Result[B]) *Result[B] {
	return Flatten[A](Map(r, mapper))
}

func (r *Result[T]) Unwrap() (T, error) {
	return r.ok, r.err
}

func (r *Result[T]) UnwrapOr(def T) T {
	if r.err != nil {
		return def
	}

	return r.ok
}
