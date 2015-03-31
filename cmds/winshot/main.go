package main

import "github.com/amitu/gutils/winshot"
import "fmt"
import "image"
import "image/jpeg"
import "os"

func main() {
	data, w, h, s := winshot.Raw()
	fmt.Println(len(data), w, h, s)
	fmt.Println("yo")
	m := &image.RGBA{data, s, image.Rect(0, 0, w, h)}
    f, err := os.Create("some.jpg")
    if err != nil {
    	panic(err)
    }
    jpeg.Encode(f, m, nil)
    f.Close()
}