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

func (circ *CircularArray) Empty() {
	circ.Buffer = nil
	circ.Start = 0
	circ.End = 0
}

func (circ *CircularArray) Capacity() uint {
	return circ.Size
}

func (circ *CircularArray) Length() uint {
	return circ.End - circ.Start
}

func (circ *CircularArray) Push(val interface{}) (interface{}, bool) {
	dropped := false
	var obj interface{} = nil

	if circ.Length() == circ.Capacity() {
		dropped = true
		obj = circ.Buffer[circ.Start]
		circ.Start += 1
	}

	if len(circ.Buffer) == int(circ.Size) || len(circ.Buffer) > int(circ.End) {
		circ.Buffer[circ.End%circ.Size] = val
	} else {
		circ.Buffer = append(circ.Buffer, val)
	}

	circ.End += 1

	return obj, dropped
}

func (circ *CircularArray) Pop() (interface{}, error) {
	if circ.Start == circ.End {
		return nil, Empty
	}
	circ.Start += 1
	last := circ.Buffer[(circ.Start-1)%circ.Size]

	if circ.Length() == 0 {
		circ.Empty()
	}

	return last, nil
}

func (circ *CircularArray) PopNewest() (interface{}, error) {
	if circ.Start == circ.End {
		return nil, Empty
	}
	circ.End -= 1
	last := circ.Buffer[circ.End%circ.Size]

	if circ.Length() == 0 {
		circ.Empty()
	}

	return last, nil
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

	return circ.Buffer[(circ.End-1)%circ.Size], nil
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
