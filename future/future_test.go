package future

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	f := New(func() int {
		time.Sleep(100 * time.Millisecond)
		return 420
	})

	n := f.Await()

	assert.Equal(t, 420, n)
}

func TestReady(t *testing.T) {
	fn := Ready(420)

	n := fn.Await()

	assert.Equal(t, 420, n)
}

func TestAwait(t *testing.T) {
	f := &Future[int]{
		Fn: func() int {
			time.Sleep(100 * time.Millisecond)
			return 420
		},
	}

	n := f.Await()

	assert.Equal(t, 420, n)
}

func TestMap(t *testing.T) {
	fn := Ready[int](420)

	fs := Map(fn, func(a int) string {
		return strconv.Itoa(a)
	})

	s := fs.Await()

	assert.Equal(t, "420", s)
}

func TestAndThen(t *testing.T) {
	fn := Ready[int](420)

	fs := AndThen(fn, func(a int) *Future[string] {
		return Ready(strconv.Itoa(a))
	})

	s := fs.Await()

	assert.Equal(t, "420", s)
}

func TestAll(t *testing.T) {
	nfs := []*Future[int]{
		New(func() int {
			time.Sleep(200 * time.Millisecond)
			return 1
		}),
		New(func() int {
			time.Sleep(100 * time.Millisecond)
			return 2
		}),
		Ready(3),
	}

	fns := All(nfs)

	ns := fns.Await()

	assert.Equal(t, []int{1, 2, 3}, ns)
}
