package main

import "github.com/amitu/gutils/winshot"
import "fmt"
import "image"
import "image/jpeg"
import "os"
import "image/color"

type myimg struct {
	image.RGBA
}

func (m *myimg) At(x, y int) color.Color {
	r, g, b, a := m.RGBA.At(x, m.Rect.Max.Y-y-1).RGBA()
	return color.RGBA{uint8(b), uint8(g), uint8(r), uint8(a)}
}

func main() {
	data, w, h, s := winshot.Raw()
	fmt.Println(len(data), w, h, s)
	fmt.Println("yo")
	m := &myimg{image.RGBA{data, s, image.Rect(0, 0, w, h)}}
    f, err := os.Create("some.jpg")
    if err != nil {
    	panic(err)
    }
    jpeg.Encode(f, m, nil)
    f.Close()
}