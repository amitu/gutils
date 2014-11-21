package gutils

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	// "time"
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
		t, "CircularArray {size: 10, start: 1, end: 1, buffer: [10]}",
		circ.Dump(),
	)

	circ.Push(20)

	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 1, end: 2, buffer: [10 20]}",
	)

	v, err = circ.PeekNewest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 20)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 1, end: 2, buffer: [10 20]}",
	)

	v, err = circ.PeekOldest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 20)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 1, end: 2, buffer: [10 20]}",
	)

	v, err = circ.Ith(0)
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 20)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 1, end: 2, buffer: [10 20]}",
	)

	v, err = circ.Ith(1)
	assert.Equal(t, err, IndexOutOfBounds)
	assert.Equal(t, v, nil)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 1, end: 2, buffer: [10 20]}",
	)

	v, err = circ.Pop()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 20)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 2, end: 2, buffer: [10 20]}",
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
		t, circ.Dump(),
		"CircularArray {size: 10, start: 3, end: 13, buffer: [38 39 40 31 32 33 34 35 36 37]}",
	)

	v, err = circ.Pop()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 31)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 4, end: 13, buffer: [38 39 40 31 32 33 34 35 36 37]}",
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
		t, circ.Dump(),
		"CircularArray {size: 10, start: 13, end: 13, buffer: [38 39 40 31 32 33 34 35 36 37]}",
	)

	v, err = circ.Pop()
	assert.Equal(t, err, Empty)
	assert.Equal(t, v, nil)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 13, end: 13, buffer: [38 39 40 31 32 33 34 35 36 37]}",
	)

	circ.Push(50)
	circ.Push(51)
	circ.Push(52)
	circ.Push(53)

	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 13, end: 17, buffer: [38 39 40 50 51 52 53 35 36 37]}",
	)

	v, err = circ.PopNewest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 53)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 13, end: 16, buffer: [38 39 40 50 51 52 53 35 36 37]}",
	)

	v, err = circ.PeekNewest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 52)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 13, end: 16, buffer: [38 39 40 50 51 52 53 35 36 37]}",
	)

	v, err = circ.PeekOldest()
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 50)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 13, end: 16, buffer: [38 39 40 50 51 52 53 35 36 37]}",
	)

	v, err = circ.Ith(0)
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 50)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 13, end: 16, buffer: [38 39 40 50 51 52 53 35 36 37]}",
	)

	v, err = circ.Ith(2)
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 52)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 13, end: 16, buffer: [38 39 40 50 51 52 53 35 36 37]}",
	)

	v, err = circ.Ith(3)
	assert.Equal(t, err, IndexOutOfBounds)
	assert.Equal(t, v, nil)
	assert.Equal(
		t, circ.Dump(),
		"CircularArray {size: 10, start: 13, end: 16, buffer: [38 39 40 50 51 52 53 35 36 37]}",
	)

	circ.PopNewest()
	circ.PopNewest()
	circ.PopNewest()
	assert.Equal(t, 0, circ.Length())

	v, err = circ.PopNewest()
	assert.Equal(t, err, Empty)
	assert.Equal(t, v, nil)
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
