package pngimage

import (
	"fmt"
	"image"
	"math"
)

// Example of creating a Bresenham method image.
func ExampleImage_Line() {
	var img = WhiteImage(200, 200)
	var rgb = NewRGB(255, 0, 0)
	for i := 0; i < 12; i++ {
		alpha := (float64(2*i) * math.Pi) / 13
		x := int(100 + 95*math.Cos(alpha))
		y := int(100 + 95*math.Sin(alpha))
		img.Line(image.Point{X: 100, Y: 100}, image.Point{X: x, Y: y}, *rgb)
	}
	img.Save("bresenham_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}
