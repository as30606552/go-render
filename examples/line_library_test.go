package examples

import (
	"computer_graphics/pngimage"
	"fmt"
	"image"
	"math"
)

const (
	W = 200 //Image Width
	H = 200 //Image height
	x0 = 100 //The coordinate of the origin point for drawing a star (by X)
	y0 = 100 //The coordinate of the origin point for drawing a star (by Y)
)

// Creates a 12-pointed star line using the specified method
func StarLine(nameMethod string, img pngimage.Image, rgb pngimage.RGB) {
	for i := 0; i < 12; i++ {
		alpha := (float64(2*i) * math.Pi) / 13
		x := int(100 + 95*math.Cos(alpha))
		y := int(100 + 95*math.Sin(alpha))
		switch nameMethod {
		case "SimplestMethod":
			SimplestMethod(image.Point{x0,y0}, image.Point{x,y}, img, rgb)
		case "SecondMethod":
			SecondMethod(image.Point{x0,y0}, image.Point{x,y}, img, rgb)
		case "ThirdMethod":
			ThirdMethod(image.Point{x0,y0}, image.Point{x,y}, img, rgb)
		case "BresenhamMethod":
			img.Line(image.Point{x0,y0}, image.Point{x,y}, rgb)
		default:
			fmt.Println("can't draw a star")
		}
	}
}

// Create a line by drawing N points on a straight line
func SimplestMethod(point1, point2 image.Point, img pngimage.Image, rgb pngimage.RGB) {
	for t := 0.0; t < 1.0; t += 0.01 {
		x := int(float64(point1.X)*(1.0-t) + float64(point2.X)*t)
		y := int(float64(point1.Y)*(1.0-t) + float64(point2.Y)*t)
		img.Set(x, y, rgb)
	}
}

// Create a line by set the x values in increments of one pixel, and calculate the y values for them
func SecondMethod(point1, point2 image.Point, img pngimage.Image, rgb pngimage.RGB) {
	for x := point1.X; x <= point2.X; x++ {
		t := float64(x-point1.X) / float64(point2.X-point1.X)
		y := int(float64(point1.Y)*(1.0-t) + float64(point2.Y)*t)
		img.Set(x, y, rgb)
	}
}

// Create a line by an improved the third method
func ThirdMethod(point1, point2 image.Point, img pngimage.Image, rgb pngimage.RGB) {
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
	for x := point1.X; x <= point2.X; x++ {
		t := float64(x-point1.X) / float64(point2.X-point1.X)
		y := int(float64(point1.Y)*(1.0-t) + float64(point2.Y)*t)
		if steep {
			img.Set(y, x, rgb)
		} else {
			img.Set(x, y, rgb)
		}
	}
}

// Example of creating a simplest method image
func ExampleSimplestMethod() {
	var img = pngimage.WhiteImage(W, H)
	var rgb = pngimage.NewRGB(255, 0, 0)
	StarLine("SimplestMethod", *img, *rgb)
	img.Save( "pictures/simplest_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}

// Example of creating a second method image
func ExampleSecondMethod() {
	var img = pngimage.WhiteImage(W, H)
	var rgb = pngimage.NewRGB(255, 0, 0)
	StarLine("SecondMethod", *img, *rgb)
	img.Save("pictures/second_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}

// Example of creating a third method image
func ExampleThirdMethod() {
	var img = pngimage.WhiteImage(W, H)
	var rgb = pngimage.NewRGB(255, 0, 0)
	StarLine("ThirdMethod", *img, *rgb)
	img.Save("pictures/third_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}

// Example of creating a Bresenham method image
func ExampleBresenhamMethod() {
	var img = pngimage.WhiteImage(W, H)
	var rgb = pngimage.NewRGB(255, 0, 0)
	StarLine("BresenhamMethod", *img, *rgb)
	img.Save( "pictures/bresenham_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}
