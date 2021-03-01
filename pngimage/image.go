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

// Creates a new Image object with the specified width and height.
func NewImage(width, height uint) *Image {
	return &Image{image.NewRGBA(image.Rect(0, 0, int(width), int(height)))}
}

// Creates an all-white png image with the size W*H in the examples/pictures directory.
func WhiteImage(width, height  int) *Image {
	var img = image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			img.SetRGBA(i, j, color.RGBA{255, 255, 255, 255})
		}
	}
	return &Image{img}
}

// Creates an all-black png image with the size W*H in the examples/pictures directory.
func BlackImage(width, height  int) *Image {
	var img = image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			img.SetRGBA(i, j, color.RGBA{0, 0, 0, 255})
		}
	}
	return &Image{img}
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

//Line drawing method.
//Takes 2 points coordinates (x0, y0), (x1, y1) and line color (rgb) as input.
//Returns the picture from the drawn line.
func (img *Image) Line(point1, point2 image.Point, rgb RGB) Image {
	steep := false
	if math.Abs(float64(point1.X-point2.X)) < math.Abs(float64(point1.Y-point2.Y)) {
		point1.X, point1.Y = point1.Y, point1.X
		point2.X, point2.Y = point2.Y, point2.X
		steep = true
	}
	if point1.X > point2.X {
		point1.X, point2.X = point2.X, point1.X
		point1.Y, point2.Y = point2.Y, point1.Y
	}
	dx := point2.X - point1.X
	dy := point2.Y - point1.Y
	derr := math.Abs(float64(dy) / float64(dx))
	err := 0.0
	y := point1.Y
	//calculate the y-axis offset relative to the center of the pixel at each step
	for x := point1.X; x <= point2.X; x++ {
		if steep {
			img.Set(y, x, rgb)
		} else {
			img.Set(x, y, rgb)
		}
		err += derr
		if err > 0.5 { //if the value exceeds 0.5, shift the displayed pixel by one position up/dow
			if point2.Y > point1.Y {
				y += 1
			} else {
				y -= 1
			}
			err -= 1.0
		}
	}
	return *img
}
