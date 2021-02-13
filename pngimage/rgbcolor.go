package pngimage

import "image/color"

type RGB struct {
	R, G, B uint8
}

func NewRGB(r, g, b uint8) *RGB {
	return &RGB{R: r, G: g, B: b}
}

func (rgb RGB) RGBA() (r, g, b, a uint32) {
	return uint32(rgb.R), uint32(rgb.G), uint32(rgb.B), 255
}

func (rgb RGB) ToRGBA() color.RGBA {
	return color.RGBA{
		R: rgb.R,
		G: rgb.G,
		B: rgb.B,
		A: 255,
	}
}
