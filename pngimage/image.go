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

type Image struct {
	img *image.RGBA
}

func NewImage(width, height uint) *Image {
	return &Image{image.NewRGBA(image.Rect(0, 0, int(width), int(height)))}
}

func (img Image) ColorModel() color.Model {
	return img.ColorModel()
}

func (img Image) Bounds() image.Rectangle {
	return img.Bounds()
}

func (img Image) At(x, y int) color.Color {
	return img.At(x, y)
}

func (img Image) Get(x, y int) RGB {
	var r, g, b, _ = img.At(x, y).RGBA()
	return *NewRGB(uint8(r), uint8(g), uint8(b))
}

func (img Image) Set(x, y int, rgb RGB) {
	img.img.SetRGBA(x, y, rgb.ToRGBA())
}

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

// Create a line by Bresenham algorithm: calculate the y-axis offset relative to the center
// of the pixel at each step and, if the value exceeds 0.5, shift the displayed pixel by one position up/down
func (img Image) Line(x0, y0, x1, y1 int, rgb RGB) Image {
	steep := false
	if math.Abs(float64(x0-x1)) < math.Abs(float64(y0-y1)) {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
		steep = true
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	dx := x1 - x0
	dy := y1 - y0
	derr := math.Abs(float64(dy) / float64(dx))
	err := 0.0
	y := y0

	for x := x0; x <= x1; x++ {
		if steep {
			img.Set(y, x, rgb)
		} else {
			img.Set(x, y, rgb)
		}
		err += derr
		if err > 0.5 {
			if y1 > y0 {
				y += 1
			} else {
				y += -1
			}
			err -= 1.0
		}
	}
	return img
}
