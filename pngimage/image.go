package pngimage

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strings"
)

// Wrapper around the image.Image for working with images in RGB format without specifying alpha value.
// All pixels have a maximum alfa value, meaning they are completely opaque.
// Implements the interface image.Image, so that all the functions that work with images can be used.
type Image struct {
	img *image.RGBA
}

// Creates a new Image with the specified width and height.
func NewImage(width, height uint) *Image {
	return &Image{image.NewRGBA(image.Rect(0, 0, int(width), int(height)))}
}

// Creates an all-white Image with the specified width and height.
func WhiteImage(width, height uint) *Image {
	var (
		img = NewImage(width, height)
		rgb = WhiteColor()
	)
	for i := 0; i < int(width); i++ {
		for j := 0; j < int(height); j++ {
			img.Set(i, j, rgb)
		}
	}
	return img
}

// Creates an all-black Image with the specified width and height.
func BlackImage(width, height uint) *Image {
	var (
		img = NewImage(width, height)
		rgb = BlackColor()
	)
	for i := 0; i < int(width); i++ {
		for j := 0; j < int(height); j++ {
			img.Set(i, j, rgb)
		}
	}
	return img
}

// Implementation of the ColorModel method in the image.Image interface.
func (img *Image) ColorModel() color.Model {
	return img.img.ColorModel()
}

// Implementation of the Bounds method in the image.Image interface.
func (img *Image) Bounds() image.Rectangle {
	return img.img.Bounds()
}

// Implementation of the At method in the image.Image interface.
func (img *Image) At(x, y int) color.Color {
	return img.img.At(x, y)
}

// Returns the color of the pixel at (x, y).
func (img *Image) Get(x, y int) RGB {
	var r, g, b, _ = img.At(x, y).RGBA()
	return RGB{uint8(r), uint8(g), uint8(b)}
}

// Sets the color of the pixel at (x, y).
func (img *Image) Set(x, y int, rgb RGB) {
	img.img.Set(x, y, rgb.ToRGBA())
}

// Returns the width of the image in pixels.
func (img *Image) Width() int {
	return img.img.Rect.Max.X
}

// Returns the height of the image in pixels.
func (img *Image) Height() int {
	return img.img.Rect.Max.Y
}

// Line drawing method.
// Takes 2 points coordinates (x0, y0), (x1, y1) and line color (rgb) as input.
// Draw a line by Bresenham algorithm.
func (img *Image) Line(x1, y1, x2, y2 int, rgb RGB) {
	var steep = false
	if math.Abs(float64(x1-x2)) < math.Abs(float64(y1-y2)) {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
		steep = true
	}
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	var (
		deltaX          = x2 - x1
		deltaY          = y2 - y1
		deltaInaccuracy = math.Abs(float64(deltaY) / float64(deltaX))
		inaccuracy      = 0.0
		y               = y1
	)
	// Calculate the y-axis offset relative to the center of the pixel at each step.
	for x := x1; x <= x2; x++ {
		if steep {
			img.Set(y, x, rgb)
		} else {
			img.Set(x, y, rgb)
		}
		inaccuracy += deltaInaccuracy
		if inaccuracy > 0.5 {
			// If the value exceeds 0.5, shift the displayed pixel by one position up/down.
			if y2 > y1 {
				y += 1
			} else {
				y -= 1
			}
			inaccuracy -= 1.0
		}
	}
}

// Saves the image in a file named filename.
// The file name must contain the .png postfix.
// If an error occurred in the method, the error object is returned, otherwise nil is returned.
func (img *Image) Save(filename string) error {
	if !strings.HasSuffix(filename, ".png") {
		return errors.New("file must be in PNG format")
	}
	var file, err = os.Create(filename)
	if err != nil {
		return err
	}
	if err := png.Encode(file, img.img); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
