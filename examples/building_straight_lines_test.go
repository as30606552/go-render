package examples

import (
	"computer_graphics/pngimage"
	"fmt"
	"math"
)

// The type for passing a method to a function starLine.
type drawingStraightLine func(x1, y1, x2, y2 int, img *pngimage.Image, rgb pngimage.RGB)

// Creates a 12-pointed star line using the specified method.
func starLine(method drawingStraightLine, img *pngimage.Image, rgb pngimage.RGB) {
	for i := 0; i < 12; i++ {
		alpha := (float64(2*i) * math.Pi) / 13
		x := int(100 + 95*math.Cos(alpha))
		y := int(100 + 95*math.Sin(alpha))
		method(100, 100, x, y, img, rgb)
	}
}

// Creates a line by drawing N points on a straight line.
func SimplestMethod(x1, y1, x2, y2 int, img *pngimage.Image, rgb pngimage.RGB) {
	for t := 0.0; t < 1.0; t += 0.01 {
		x := int(float64(x1)*(1.0-t) + float64(x2)*t)
		y := int(float64(y1)*(1.0-t) + float64(y2)*t)
		img.Set(x, y, rgb)
	}
}

// Create a line by set the x values in increments of one pixel, and calculate the y values for them.
func SecondMethod(x1, y1, x2, y2 int, img *pngimage.Image, rgb pngimage.RGB) {
	for x := x1; x <= x2; x++ {
		t := float64(x-x1) / float64(x2-x1)
		y := int(float64(y1)*(1.0-t) + float64(y2)*t)
		img.Set(x, y, rgb)
	}
}

// Create a line by an improved the third method.
func ThirdMethod(x1, y1, x2, y2 int, img *pngimage.Image, rgb pngimage.RGB) {
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
	for x := x1; x <= x2; x++ {
		t := float64(x-x1) / float64(x2-x1)
		y := int(float64(y1)*(1.0-t) + float64(y2)*t)
		if steep {
			img.Set(y, x, rgb)
		} else {
			img.Set(x, y, rgb)
		}
	}
}

// Example of creating a simplest method image.
func ExampleSimplestMethod() {
	var (
		img = pngimage.WhiteImage(200, 200)
		rgb = pngimage.RedColor()
	)
	starLine(SimplestMethod, img, rgb)
	if err := img.Save("testdata/pictures/lines/simplest_method_image.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Example of creating a second method image.
func ExampleSecondMethod() {
	var (
		img = pngimage.WhiteImage(200, 200)
		rgb = pngimage.RedColor()
	)
	starLine(SecondMethod, img, rgb)
	if err := img.Save("testdata/pictures/lines/second_method_image.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Example of creating a third method image.
func ExampleThirdMethod() {
	var (
		img = pngimage.WhiteImage(200, 200)
		rgb = pngimage.RedColor()
	)
	starLine(ThirdMethod, img, rgb)
	if err := img.Save("testdata/pictures/lines/third_method_image.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}
