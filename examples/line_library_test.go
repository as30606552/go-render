package examples

import (
	"computer_graphics/pngimage"
	"fmt"
	"image"
	"image/color"
	"math"
)

const (
	W1 = 200
	H1 = 200
	x0 = 100
	y0 = 100
)

var img = image.NewRGBA(image.Rect(0, 0, W1, H1))
var rgba = color.RGBA{255, 0, 0, 255}
var lineImg = pngimage.NewImage(W1, H1)
var rgb = pngimage.NewRGB(255, 255, 255)

// Creates an all-white png image with the size W1*H1 in the examples/pictures directory.
func Image() {
	for i := 0; i < W1; i++ {
		for j := 0; j < H1; j++ {
			img.SetRGBA(i, j, color.RGBA{255, 255, 255, 255})
		}
	}
}

// Creates an all-white png image with the size W1*H1 in the examples/pictures directory.
func ImageBresenham() {
	for i := 0; i < W1; i++ {
		for j := 0; j < H1; j++ {
			lineImg.Set(i, j, *rgb)
		}
	}
}

// Creates a 12-pointed star line using the specified method
func StarLine(numberMethod int) {
	for i := 0; i < 12; i++ {
		alpha := (float64(2*i) * math.Pi) / 13
		x := int(100 + 95*math.Cos(alpha))
		y := int(100 + 95*math.Sin(alpha))
		switch numberMethod {
		case 1:
			SimplestMethod(x0, y0, x, y, rgba)
		case 2:
			SecondMethod(x0, y0, x, y, rgba)
		case 3:
			ThirdMethod(x0, y0, x, y, rgba)
		case 4:
			lineImg.Line(x0, y0, x, y, *rgb)
		default:
			fmt.Println("can't draw a star")
		}
	}
}

// Create a line by drawing N points on a straight line
func SimplestMethod(x0, y0, x1, y1 int, rgba color.RGBA) {
	for t := 0.0; t < 1.0; t += 0.01 {
		x := int(float64(x0)*(1.0-t) + float64(x1)*t)
		y := int(float64(y0)*(1.0-t) + float64(y1)*t)
		img.Set(x, y, rgba)
	}
}

// Create a line by set the x values in increments of one pixel, and calculate the y values for them
func SecondMethod(x0, y0, x1, y1 int, rgba color.RGBA) {
	for x := x0; x <= x1; x++ {
		t := float64(x-x0) / float64(x1-x0)
		y := int(float64(y0)*(1.0-t) + float64(y1)*t)
		img.Set(x, y, rgba)
	}
}

// Create a line by an improved the third method
func ThirdMethod(x0, y0, x1, y1 int, rgba color.RGBA) {
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
	for x := x0; x <= x1; x++ {
		t := float64(x-x0) / float64(x1-x0)
		y := int(float64(y0)*(1.0-t) + float64(y1)*t)
		if steep {
			img.Set(y, x, rgba)
		} else {
			img.Set(x, y, rgba)
		}
	}
}

// Example of creating a simplest method image
func ExampleSimplestMethod() {
	Image()
	StarLine(1)
	makeFile(img, "pictures/simplest_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}

// Example of creating a second method image
func ExampleSecondMethod() {
	Image()
	StarLine(2)
	makeFile(img, "pictures/second_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}

// Example of creating a third method image
func ExampleThirdMethod() {
	Image()
	StarLine(3)
	makeFile(img, "pictures/third_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}

// Example of creating a Bresenham method image
func ExampleBresenhamMethod() {
	ImageBresenham()
	StarLine(4)
	//makeFile(lineImg, "pictures/bresenham_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}
