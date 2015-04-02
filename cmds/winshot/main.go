package main

import "github.com/amitu/gutils"
import "image/jpeg"
import "os"

func main() {
	m := gutils.Screenshot()

    f, err := os.Create("some.jpg")
    if err != nil {
    	panic(err)
    }
    defer f.Close()

    jpeg.Encode(f, m, nil)
}