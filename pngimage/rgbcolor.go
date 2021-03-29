package pngimage

import "image/color"

// A structure for storing colors in RGB format without specifying alfa value.
// All pixels have a maximum alfa value, meaning they are completely opaque.
// Implements the interface color.Color, so that all the functions that work with color can be used.
type RGB struct {
	R, G, B uint8
}

// Creates a new RGB object with the specified red, green and blue values.
func NewRGB(r, g, b uint8) RGB {
	return RGB{R: r, G: g, B: b}
}

// Implementation of the RGBA method in the color.Color interface.
func (rgb RGB) RGBA() (r, g, b, a uint32) {
	return uint32(rgb.R), uint32(rgb.G), uint32(rgb.B), 255
}

// Converts an RGB object to an color.RGBA object.
func (rgb RGB) ToRGBA() color.RGBA {
	return color.RGBA{
		R: rgb.R,
		G: rgb.G,
		B: rgb.B,
		A: 255,
	}
}
