package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

const (
	W = 600
	H = 400
)

func makeFile(img image.Image, filename string) {
	var file, err = os.Create(filename)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(file, img); err != nil {
		_ = file.Close()
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}
}

func blackImage() {
	var img = image.NewGray(image.Rect(0, 0, W, H))
	for i := 0; i < W; i++ {
		for j := 0; j < H; j++ {
			img.SetGray(i, j, color.Gray{Y: 0})
		}
	}
	makeFile(img, "pictures/black_image.png")
}

func whiteImage() {
	var img = image.NewGray(image.Rect(0, 0, W, H))
	for i := 0; i < W; i++ {
		for j := 0; j < H; j++ {
			img.SetGray(i, j, color.Gray{Y: 255})
		}
	}
	makeFile(img, "pictures/white_image.png")
}

func redImage() {
	var img = image.NewRGBA(image.Rect(0, 0, W, H))
	for i := 0; i < W; i++ {
		for j := 0; j < H; j++ {
			img.SetRGBA(i, j, color.RGBA{R: 255, A: 255})
		}
	}
	makeFile(img, "pictures/red_image.png")
}

func gradientImage() {
	var img = image.NewGray(image.Rect(0, 0, W, H))
	for i := 0; i < W; i++ {
		for j := 0; j < H; j++ {
			img.SetGray(i, j, color.Gray{Y: uint8((i + j) % 256)})
		}
	}
	makeFile(img, "pictures/gradient_image.png")
}

func main() {
	blackImage()
	whiteImage()
	redImage()
	gradientImage()
}
