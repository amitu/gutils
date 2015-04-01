package gutils

import (
	"image"
	"image/color"
)

// MacBGRA is used to represent the image obtained from OS X screenbuffer.
type MacBGRA struct {
	image.RGBA
}

// At can be used to find the color at given pixel. It has been modified to
// accomodate the difference between image.RGBA way of storing pixel data and
// the OS X way of storing pixel data
func (m *MacBGRA) At(x, y int) color.Color {
	r, g, b, a := m.RGBA.At(x, y).RGBA()
	return color.RGBA{uint8(b), uint8(g), uint8(r), uint8(a)}
}

// NewMacBGRA returns a MacARGB object with given pixel data
func NewMacBGRA(r image.Rectangle, stride int, pix []byte) *MacBGRA {
	return &MacBGRA{RGBA: image.RGBA{Pix: pix, Stride: stride, Rect: r}}
}

// WinBGRA represents the image obtained from windows framebuffer
type WinBGRA struct {
	image.RGBA
}

// At can be used to find the color at given pixel. It has been modified to
// accomodate the difference between image.RGBA way of storing pixel data and
// the Windows way of storing pixel data
func (m *WinBGRA) At(x, y int) color.Color {
	r, g, b, a := m.RGBA.At(x, m.Rect.Max.Y-y-1).RGBA()
	return color.RGBA{uint8(b), uint8(g), uint8(r), uint8(a)}
}

// NewWinBGRA returns a WinBGRA object with given pixel data
func NewWinBGRA(r image.Rectangle, stride int, pix []byte) *WinBGRA {
	return &WinBGRA{RGBA: image.RGBA{Pix: pix, Stride: stride, Rect: r}}
}
