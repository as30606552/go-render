package pngimage

import (
	"fmt"
	"image"
	"math"
	"os"
	"testing"
)

// Creates directories for output, if there are none.
func TestMain(m *testing.M) {
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		err = os.Mkdir("testdata", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat("testdata/pictures"); os.IsNotExist(err) {
		err = os.Mkdir("testdata/pictures", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	m.Run()
}

// Example of creating a Bresenham method image.
func ExampleImage_Line() {
	var (
		img = WhiteImage(200, 200)
		rgb = RGB{R: 255}
	)
	for i := 0; i < 12; i++ {
		alpha := (float64(2*i) * math.Pi) / 13
		x := int(100 + 95*math.Cos(alpha))
		y := int(100 + 95*math.Sin(alpha))
		img.Line(image.Point{X: 100, Y: 100}, image.Point{X: x, Y: y}, rgb)
	}
	if err := img.Save("testdata/pictures/bresenham_method_image.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}
