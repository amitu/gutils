package gutils

import (
	"testing"
	"io/ioutil"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func assertDirLength(t *testing.T, n int, dir string) {
	files, err := ioutil.ReadDir(dir)
	assert.Nil(t, err)

	assert.Equal(t, len(files), n)
}

func TestFoo(t *testing.T) {
	var filedchan FiledChan

	err := filedchan.Init(5)
	defer filedchan.Quit()

	assert.Nil(t, err)

	assertDirLength(t, 0, filedchan.Dir)

	for i := 0; i <=7 ; i++  {
		fmt.Println("writing packet", i)
		filedchan.Prod <- []byte("hello there")
	}

	assertDirLength(t, 2, filedchan.Dir)

	for i := 0; i <=7 ; i++  {
		fmt.Println("reading packet", i)
		fmt.Println("found", i, <- filedchan.Cons)
	}

	assertDirLength(t, 0, filedchan.Dir)

}