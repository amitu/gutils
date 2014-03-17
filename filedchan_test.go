package gutils

import (
	"testing"
	"io/ioutil"
	"fmt"
	"github.com/stretchr/testify/assert"
	"encoding/gob"
)

type MyData struct {
	ID string
	Data []byte
	N int
}

func assertDirLength(t *testing.T, n int, dir string) {
	files, err := ioutil.ReadDir(dir)
	assert.Nil(t, err)

	assert.Equal(t, len(files), n)
}

func TestFoo(t *testing.T) {
	gob.Register(MyData{})
	var filedchan FiledChan

	err := filedchan.Init(5)
	defer filedchan.Quit()

	assert.Nil(t, err)

	assertDirLength(t, 0, filedchan.Dir)

	for i := 0; i <= 7 ; i++  {
		fmt.Println("writing packet", i)
		filedchan.Prod <- MyData{ID: "123", Data: []byte("yo"), N: i}
	}

	assertDirLength(t, 2, filedchan.Dir)

	for i := 0; i <= 7 ; i++  {
		fmt.Println("reading packet", i)
		data := (<- filedchan.Cons).(MyData)
		assert.Equal(t, i, data.N)
		fmt.Println("found", i)
	}

	assertDirLength(t, 0, filedchan.Dir)

	for i := 0; i <= 7 ; i++  {
		fmt.Println("writing packet", i)
		filedchan.Prod <- MyData{ID: "123", Data: []byte("yo"), N: i}
	}

	assertDirLength(t, 2, filedchan.Dir)

	data := (<- filedchan.Cons).(MyData)
	assert.Equal(t, 0, data.N)

	assertDirLength(t, 2, filedchan.Dir)

	filedchan.Prod <- MyData{ID: "123", Data: []byte("yo"), N: 8}

	assertDirLength(t, 2, filedchan.Dir)

	for i := 0; i <= 7 ; i++  {
		fmt.Println("reading packet", i)
		data := (<- filedchan.Cons).(MyData)
		assert.Equal(t, i+1, data.N)
		fmt.Println("found", i)
	}

	assertDirLength(t, 0, filedchan.Dir)
}