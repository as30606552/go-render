package examples

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// Creates a file with the specified name and places the specified image in it.
func makeFile(img image.Image, filename string) error {
	var file, err = os.Create(filename)
	if err != nil {
		return err
	}
	if err = png.Encode(file, img); err != nil {
		_ = file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	return nil
}

// Creates an all-black png image in the testdata/pictures directory.
func BlackImage() error {
	const (
		w = 600 // Image width.
		h = 400 // Image height.
	)
	var img = image.NewGray(image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			img.SetGray(i, j, color.Gray{Y: 0})
		}
	}
	return makeFile(img, "testdata/pictures/black_image.png")
}

// Creates an all-white png image in the testdata/pictures directory.
func WhiteImage() error {
	const (
		w = 600 // Image width.
		h = 400 // Image height.
	)
	var img = image.NewGray(image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			img.SetGray(i, j, color.Gray{Y: 255})
		}
	}
	return makeFile(img, "testdata/pictures/white_image.png")
}

// Creates an all-red png image with the size W*H in the testdata/pictures directory.
func RedImage() error {
	const (
		w = 600 // Image width.
		h = 400 // Image height.
	)
	var img = image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			img.SetRGBA(i, j, color.RGBA{R: 255, A: 255})
		}
	}
	return makeFile(img, "testdata/pictures/red_image.png")
}

// Creates a gradient png image with the size W*H in the testdata/pictures directory.
func GradientImage() error {
	const (
		w = 600 // Image width.
		h = 400 // Image height.
	)
	var img = image.NewGray(image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			img.SetGray(i, j, color.Gray{Y: uint8((i + j) % 256)})
		}
	}
	return makeFile(img, "testdata/pictures/gradient_image.png")
}

// Example of creating a black image.
func ExampleBlackImage() {
	if err := BlackImage(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Example of creating a white image.
func ExampleWhiteImage() {
	if err := WhiteImage(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Example of creating a red image.
func ExampleRedImage() {
	if err := RedImage(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Example of creating a gradient image.
func ExampleGradientImage() {
	if err := GradientImage(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}
