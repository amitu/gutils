package gutils

import (
	"errors"
	"fmt"
)

var (
	Empty            = errors.New("CircularArrayEmpty")
	IndexOutOfBounds = errors.New("CircularArrayIndexOutOfBounds")
)

type CircularArray struct {
	size       uint
	buffer     []interface{}
	start, end uint
}

func NewCircularArray(size uint) *CircularArray {
	return &CircularArray{size: size, start: 0, end: 0}
}

func (circ *CircularArray) Capacity() uint {
	return circ.size
}

func (circ *CircularArray) Length() uint {
	return circ.end - circ.start
}

func (circ *CircularArray) Push(val interface{}) {
	if circ.Length() == circ.Capacity() {
		circ.start += 1
	}

	circ.end += 1

	if len(circ.buffer) == int(circ.size) {
		circ.buffer[circ.end%circ.size-1] = val
	} else {
		circ.buffer = append(circ.buffer, val)
	}
}

func (circ *CircularArray) Pop() (interface{}, error) {
	if circ.start == circ.end {
		return nil, Empty
	}
	circ.start += 1
	return circ.buffer[(circ.start-1)%circ.size], nil
}

func (circ *CircularArray) PopNewest() (interface{}, error) {
	if circ.start == circ.end {
		return nil, Empty
	}
	if circ.end-circ.start == 1 {
		circ.start += 1
		return circ.buffer[(circ.start)%circ.size-1], nil
	} else {
		circ.end -= 1
		return circ.buffer[(circ.end+1)%circ.size-1], nil
	}
}

func (circ *CircularArray) PeekOldest() (interface{}, error) {
	if circ.start == circ.end {
		return nil, Empty
	}

	return circ.buffer[circ.start%circ.size], nil
}

func (circ *CircularArray) PeekNewest() (interface{}, error) {
	if circ.start == circ.end {
		return nil, Empty
	}

	return circ.buffer[circ.end%circ.size-1], nil
}

func (circ *CircularArray) Ith(i uint) (interface{}, error) {
	if circ.start == circ.end {
		return nil, Empty
	}

	ti := i + circ.start

	if ti >= circ.end {
		return nil, IndexOutOfBounds
	}

	return circ.buffer[ti%circ.size], nil
}

func (circ *CircularArray) Dump() string {
	return fmt.Sprintf(
		"CircularArray {size: %d, start: %d, end: %d, buffer: %v}",
		circ.size, circ.start, circ.end, circ.buffer,
	)
}

func (circ *CircularArray) P() {
	fmt.Println(circ.Dump(), circ.Length(), circ.Capacity())
}
