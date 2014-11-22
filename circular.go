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
	Size       uint
	Buffer     []interface{}
	Start, End uint
}

func NewCircularArray(size uint) *CircularArray {
	return &CircularArray{Size: size}
}

func (circ *CircularArray) Capacity() uint {
	return circ.Size
}

func (circ *CircularArray) Length() uint {
	return circ.End - circ.Start
}

func (circ *CircularArray) Push(val interface{}) {
	if circ.Length() == circ.Capacity() {
		circ.Start += 1
	}

	if len(circ.Buffer) == int(circ.Size) || len(circ.Buffer) > int(circ.End) {
		circ.Buffer[circ.End%circ.Size] = val
	} else {
		circ.Buffer = append(circ.Buffer, val)
	}

	circ.End += 1
}

func (circ *CircularArray) Pop() (interface{}, error) {
	if circ.Start == circ.End {
		return nil, Empty
	}
	circ.Start += 1
	return circ.Buffer[(circ.Start-1)%circ.Size], nil
}

func (circ *CircularArray) PopNewest() (interface{}, error) {
	if circ.Start == circ.End {
		return nil, Empty
	}
	circ.End -= 1
	return circ.Buffer[circ.End%circ.Size], nil
}

func (circ *CircularArray) PeekOldest() (interface{}, error) {
	if circ.Start == circ.End {
		return nil, Empty
	}

	return circ.Buffer[circ.Start%circ.Size], nil
}

func (circ *CircularArray) PeekNewest() (interface{}, error) {
	if circ.Start == circ.End {
		return nil, Empty
	}

	return circ.Buffer[circ.End%circ.Size-1], nil
}

func (circ *CircularArray) Ith(i uint) (interface{}, error) {
	if circ.Start == circ.End {
		return nil, Empty
	}

	ti := i + circ.Start

	if ti >= circ.End {
		return nil, IndexOutOfBounds
	}

	return circ.Buffer[ti%circ.Size], nil
}

func (circ *CircularArray) Dump() string {
	return fmt.Sprintf(
		"CircularArray {size: %d, start: %d, end: %d, buffer: %v}",
		circ.Size, circ.Start, circ.End, circ.Buffer,
	)
}

func (circ *CircularArray) P() {
	fmt.Println(circ.Dump(), circ.Length(), circ.Capacity())
}
