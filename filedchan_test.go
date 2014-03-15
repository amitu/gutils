package gutils

import (
	"testing"
	"io/ioutil"
	"fmt"
)

func assertDirLength(t *testing.T, n int, dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Error(err)
	}

	if len(files) != n {
		t.Errorf("Wrong number of files, expected %d, found %d.", n, files)
	}
}

func TestFoo(t *testing.T) {
	var filedchan FiledChan

	err := filedchan.Init(5)
	defer filedchan.Quit()

	if err != nil {
		t.Error(err)
		return
	}

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