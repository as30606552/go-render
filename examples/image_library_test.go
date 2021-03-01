package examples

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)


// Creates a file with the specified name and places the specified image in it.
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

// Creates an all-black png image with the size W*H in the examples/pictures directory.
func BlackImage() {
	var img = image.NewGray(image.Rect(0, 0, W, H))
	for i := 0; i < W; i++ {
		for j := 0; j < H; j++ {
			img.SetGray(i, j, color.Gray{Y: 0})
		}
	}
	makeFile(img, "pictures/black_image.png")
}

// Creates an all-white png image with the size W*H in the examples/pictures directory.
func WhiteImage() {
	var img = image.NewGray(image.Rect(0, 0, W, H))
	for i := 0; i < W; i++ {
		for j := 0; j < H; j++ {
			img.SetGray(i, j, color.Gray{Y: 255})
		}
	}
	makeFile(img, "pictures/white_image.png")
}

// Creates an all-red png image with the size W*H in the examples/pictures directory.
func RedImage() {
	var img = image.NewRGBA(image.Rect(0, 0, W, H))
	for i := 0; i < W; i++ {
		for j := 0; j < H; j++ {
			img.SetRGBA(i, j, color.RGBA{R: 255, A: 255})
		}
	}
	makeFile(img, "pictures/red_image.png")
}

// Creates a gradient png image with the size W*H in the examples/pictures directory.
func GradientImage() {
	var img = image.NewGray(image.Rect(0, 0, W, H))
	for i := 0; i < W; i++ {
		for j := 0; j < H; j++ {
			img.SetGray(i, j, color.Gray{Y: uint8((i + j) % 256)})
		}
	}
	makeFile(img, "pictures/gradient_image.png")
}

func TestMain(m *testing.M) {
	if _, err := os.Stat("pictures"); os.IsNotExist(err) {
		err = os.Mkdir("pictures", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	m.Run()
}

// Example of creating a black image
func ExampleBlackImage() {
	BlackImage()
	fmt.Println("Ok")
	// Output: Ok
}

// Example of creating a white image
func ExampleWhiteImage() {
	WhiteImage()
	fmt.Println("Ok")
	// Output: Ok
}

// Example of creating a red image
func ExampleRedImage() {
	RedImage()
	fmt.Println("Ok")
	// Output: Ok
}

// Example of creating a gradient image
func ExampleGradientImage() {
	GradientImage()
	fmt.Println("Ok")
	// Output: Ok
}
