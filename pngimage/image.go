package pngimage

import (
	"errors"
	"image"
	"image/color"
	"image/png"
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
