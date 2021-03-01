package pngimage

import (
	"fmt"
	"image"
	"math"
)

const (
	x0 = 100 // The coordinate of the origin point for drawing a star (by X).
	y0 = 100 // The coordinate of the origin point for drawing a star (by Y).
)

// The type for passing a method to a function starLine.
type drawingStraightLine func(point1, point2 image.Point, img Image, rgb RGB)

// Creates a 12-pointed star line using the specified method.
func starLine(method drawingStraightLine, img Image, rgb RGB) {
	for i := 0; i < 12; i++ {
		alpha := (float64(2*i) * math.Pi) / 13
		x := int(100 + 95*math.Cos(alpha))
		y := int(100 + 95*math.Sin(alpha))
		method(image.Point{x0, y0}, image.Point{x, y}, img, rgb)
	}
}

// Wrapper function for calling the Line method.
func BresenhamMethod(point1, point2 image.Point, img Image, rgb RGB) {
	img.Line(point1, point2, rgb)
}

// Example of creating a Bresenham method image.
func ExampleBresenhamMethod() {
	var W = 200 // Image Width.
	var H = 200 // Image Height.
	var img = WhiteImage(W, H)
	var rgb = NewRGB(255, 0, 0)
	starLine(BresenhamMethod, *img, *rgb)
	img.Save("pictures/bresenham_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}
