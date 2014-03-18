package gutils

import (
	"testing"
	"io/ioutil"
	"fmt"
	"time"
	"github.com/stretchr/testify/assert"
	"encoding/gob"
)

type MyData struct {
	ID string
	Data []byte
	N int
}

func assertDirLength(t *testing.T, n int, dir string) {
	<- time.After(1e8)
	files, err := ioutil.ReadDir(dir)
	assert.Nil(t, err)

	assert.Equal(t, len(files), n)
}

func TestFoo(t *testing.T) {
	gob.Register(MyData{})
	var filedchan FiledChan

	err := filedchan.Init(5, 10)
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

	for i := 0; i < 15; i++ {
		fmt.Println("going to write22", i)
		select {
		case filedchan.Prod <- MyData{ID: "123", Data: []byte("yo"), N: i}:
			fmt.Println("written", i)
		case <- time.After(1e9):
			assert.Fail(t, "This must never block")
			return
		}
	}

	assertDirLength(t, 9, filedchan.Dir)

	for i := 0; i < 15; i++ {
		fmt.Println("reading", i)
		select {
		case data := <- filedchan.Cons:
			assert.Equal(t, i, data.(MyData).N)
			fmt.Println("got", data)
		case <- time.After(1e9):
			assert.Fail(t, "This must never block")
			return
		}
	}

	select {
	case <- filedchan.Cons:
		assert.Fail(t, "This must never unblock")
		return
	case <- time.After(1e9):
	}

	fmt.Println("end")
}
