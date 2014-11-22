package gutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCircularArray(t *testing.T) {
	circ := NewCircularArray(10)

	assert.Equal(
		t, circ.Dump(), "CircularArray {size: 10, start: 0, end: 0, buffer: []}",
	)

	assert.Equal(t, circ.Capacity(), 10)
	assert.Equal(t, circ.Length(), 0)

	v, err := circ.Pop()
	assert.Equal(t, err, Empty)
	assert.Equal(t, v, nil)

	v, err = circ.PopNewest()
	assert.Equal(t, err, Empty)
	assert.Equal(t, v, nil)

	v, err = circ.PeekNewest()
	assert.Equal(t, err, Empty)
	assert.Equal(t, v, nil)

	v, err = circ.PeekOldest()
	assert.Equal(t, err, Empty)
	assert.Equal(t, v, nil)

	v, err = circ.Ith(0)
	assert.Equal(t, err, Empty)
	assert.Equal(t, v, nil)

	// one element in the array
	circ.Push(10)

	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 0, end: 1, buffer: [10]}",
	)

	v, err = circ.PeekNewest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 10)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 0, end: 1, buffer: [10]}",
	)

	v, err = circ.PeekOldest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 10)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 0, end: 1, buffer: [10]}",
	)

	v, err = circ.Ith(0)
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 10)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 0, end: 1, buffer: [10]}",
	)

	v, err = circ.Ith(1)
	assert.Equal(t, err, IndexOutOfBounds)
	assert.Equal(t, v, nil)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 0, end: 1, buffer: [10]}",
	)

	v, err = circ.Pop()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 10)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 1, end: 1, buffer: [10]}",
	)

	circ = NewCircularArray(10)
	circ.Push(10)

	v, err = circ.PopNewest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 10)
	assert.Equal(
		t, "CircularArray {size: 10, start: 0, end: 0, buffer: [10]}",
		circ.Dump(),
	)

	circ.Push(20)

	assert.Equal(
		t, "CircularArray {size: 10, start: 0, end: 1, buffer: [20]}",
		circ.Dump(),
	)

	v, err = circ.PeekNewest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 20)
	assert.Equal(
		t, "CircularArray {size: 10, start: 0, end: 1, buffer: [20]}",
		circ.Dump(),
	)

	v, err = circ.PeekOldest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 20)
	assert.Equal(
		t, "CircularArray {size: 10, start: 0, end: 1, buffer: [20]}",
		circ.Dump(),
	)

	v, err = circ.Ith(0)
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 20)
	assert.Equal(
		t, "CircularArray {size: 10, start: 0, end: 1, buffer: [20]}",
		circ.Dump(),
	)

	v, err = circ.Ith(1)
	assert.Equal(t, err, IndexOutOfBounds)
	assert.Equal(t, v, nil)
	assert.Equal(
		t, "CircularArray {size: 10, start: 0, end: 1, buffer: [20]}",
		circ.Dump(),
	)

	v, err = circ.Pop()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 20)
	assert.Equal(
		t, "CircularArray {size: 10, start: 1, end: 1, buffer: [20]}",
		circ.Dump(),
	)

	circ.Push(30)
	circ.Push(31)
	circ.Push(32)
	circ.Push(33)
	circ.Push(34)
	circ.Push(35)
	circ.Push(36)
	circ.Push(37)
	circ.Push(38)
	circ.Push(39)
	circ.Push(40)

	assert.Equal(
		t,
		"CircularArray {size: 10, start: 2, end: 12, buffer: [39 40 31 32 33 34 35 36 37 38]}",
		circ.Dump(),
	)

	v, err = circ.Pop()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 31)
	assert.Equal(
		t,
		"CircularArray {size: 10, start: 3, end: 12, buffer: [39 40 31 32 33 34 35 36 37 38]}",
		circ.Dump(),
	)

	circ.Pop() // 32
	circ.Pop()
	circ.Pop()
	circ.Pop()
	circ.Pop()
	circ.Pop()
	circ.Pop()
	circ.Pop()

	v, err = circ.Pop()

	assert.Equal(t, err, nil)
	assert.Equal(t, v, 40)
	assert.Equal(
		t,
		"CircularArray {size: 10, start: 12, end: 12, buffer: [39 40 31 32 33 34 35 36 37 38]}",
		circ.Dump(),
	)

	v, err = circ.Pop()
	assert.Equal(t, err, Empty)
	assert.Equal(t, v, nil)
	assert.Equal(
		t,
		"CircularArray {size: 10, start: 12, end: 12, buffer: [39 40 31 32 33 34 35 36 37 38]}",
		circ.Dump(),
	)

	circ.Push(50)
	circ.Push(51)
	circ.Push(52)
	circ.Push(53)

	assert.Equal(
		t,
		"CircularArray {size: 10, start: 12, end: 16, buffer: [39 40 50 51 52 53 35 36 37 38]}",
		circ.Dump(),
	)

	v, err = circ.PopNewest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 53)

	v, err = circ.PeekNewest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 52)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 12, end: 15, buffer: [39 40 50 51 52 53 35 36 37 38]}",
		circ.Dump(),
	)

	v, err = circ.PeekOldest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 50)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 12, end: 15, buffer: [39 40 50 51 52 53 35 36 37 38]}",
		circ.Dump(),
	)

	v, err = circ.Ith(0)
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 50)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 12, end: 15, buffer: [39 40 50 51 52 53 35 36 37 38]}",
		circ.Dump(),
	)

	v, err = circ.Ith(2)
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 52)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 12, end: 15, buffer: [39 40 50 51 52 53 35 36 37 38]}",
		circ.Dump(),
	)

	v, err = circ.Ith(3)
	assert.Equal(t, err, IndexOutOfBounds)
	assert.Equal(t, v, nil)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 12, end: 15, buffer: [39 40 50 51 52 53 35 36 37 38]}",
		circ.Dump(),
	)

	circ.PopNewest()
	circ.PopNewest()
	circ.PopNewest()
	assert.Equal(t, 0, circ.Length())

	v, err = circ.PopNewest()
	assert.Equal(t, err, Empty)
	assert.Equal(t, v, nil)

	for i := 100; i < 200; i++ {
		circ.Push(i)
	}

	assert.Equal(t, 10, circ.Length())

	for i := 199; i >= 190; i-- {
		v, err = circ.PopNewest()
		assert.Equal(t, nil, err)
		assert.Equal(t, i, v)
	}

	v, err = circ.PopNewest()
	assert.Equal(t, Empty, err)
	assert.Equal(t, nil, v)
}

func TestCircularBufferArray(t *testing.T) {
	circ := NewCircularBufferArray(10)

	assert.Equal(
		t, circ.Dump(), "CircularArray {size: 10, start: 0, end: 0, buffer: []}",
	)

	assert.Equal(t, circ.Capacity(), 10)
	assert.Equal(t, circ.Length(), 0)

	circ.Push([]byte("foo"))
	v, err := circ.Pop()

	assert.Equal(t, []byte("foo"), v)
	assert.Equal(t, nil, err)

	v, err = circ.Pop()

	assert.Equal(t, []byte(nil), v)
	assert.Equal(t, Empty, err)
}
