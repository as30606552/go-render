package pngimage

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

// Wrapper around the image.Image for working with images in RGB format without specifying alpha value.
// All pixels have a maximum alfa value, meaning they are completely opaque.
// Implements the interface image.Image, so that all the functions that work with images can be used.
type Image struct {
	img *image.RGBA
}

// Creates a new Image object with the specified width and height.
func NewImage(width, height uint) *Image {
	return &Image{image.NewRGBA(image.Rect(0, 0, int(width), int(height)))}
}

// Implementation of the ColorModel method in the image.Image interface.
func (img Image) ColorModel() color.Model {
	return img.ColorModel()
}

// Implementation of the Bounds method in the image.Image interface.
func (img Image) Bounds() image.Rectangle {
	return img.Bounds()
}

// Implementation of the At method in the image.Image interface.
func (img Image) At(x, y int) color.Color {
	return img.At(x, y)
}

// Returns the color of the pixel at (x, y).
func (img Image) Get(x, y int) RGB {
	var r, g, b, _ = img.At(x, y).RGBA()
	return *NewRGB(uint8(r), uint8(g), uint8(b))
}

// Sets the color of the pixel at (x, y).
func (img Image) Set(x, y int, rgb RGB) {
	img.img.SetRGBA(x, y, rgb.ToRGBA())
}

// Saves the image in a file named filename.
// The file name must contain the .png postfix.
// If an error occurred in the method, the error object is returned, otherwise nil is returned.
func (img Image) Save(filename string) error {
	if !strings.HasSuffix(filename, ".png") {
		return errors.New("file must be in PNG format")
	}
	var file, err = os.Create(filename)
	if err != nil {
		return err
	}
	if err := png.Encode(file, img); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
