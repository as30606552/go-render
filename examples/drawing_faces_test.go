package examples

import (
	"computer_graphics/mathutils"
	"computer_graphics/model"
	"computer_graphics/obj/importer"
	"computer_graphics/pngimage"
	"fmt"
	"math"
	"os"
)

// Draws a triangle on the specified image with the specified color.
func DrawTriangle(face *model.Face, img *pngimage.Image, rgb pngimage.RGB) {
	var (
		v1         = face.Vertex1()
		v2         = face.Vertex2()
		v3         = face.Vertex3()
		xMax       = math.Min(float64(img.Width()), mathutils.Max(v1.X, v2.X, v3.X))
		xMin       = math.Max(0, mathutils.Min(v1.X, v2.X, v3.X))
		yMax       = math.Min(float64(img.Height()), mathutils.Max(v1.Y, v2.Y, v3.Y))
		yMin       = math.Max(0, mathutils.Min(v1.Y, v2.Y, v3.Y))
		l1, l2, l3 float64
	)
	for i := int(math.Ceil(xMin)); float64(i) < xMax; i++ {
		for j := int(math.Ceil(yMin)); float64(j) < yMax; j++ {
			l1, l2, l3 = face.BarycentricCoordinates(i, j)
			if l1 > 0 && l2 > 0 && l3 > 0 {
				img.Set(i, j, rgb)
			}
		}
	}
}

// Draws a triangle that fits completely into the image.
func ExampleDrawTriangle_internal() {
	var (
		face = model.NewFace(
			model.NewVertex(50, 10, 0),
			model.NewVertex(10, 80, 0),
			model.NewVertex(90, 80, 0),
		)
		img = pngimage.WhiteImage(100, 100)
		rgb = pngimage.GreenColor()
	)
	DrawTriangle(face, img, rgb)
	if err := img.Save("testdata/pictures/triangles/internal_triangle.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Draws a triangle that partially extends beyond the edges of the image.
func ExampleDrawTriangle_external() {
	var (
		face = model.NewFace(
			model.NewVertex(50, -10, 0),
			model.NewVertex(-10, 110, 0),
			model.NewVertex(120, 110, 0),
		)
		img = pngimage.WhiteImage(100, 100)
		rgb = pngimage.GreenColor()
	)
	DrawTriangle(face, img, rgb)
	if err := img.Save("testdata/pictures/triangles/external_triangle.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Draws a triangle that completely extends beyond the edges of the image.
func ExampleDrawTriangle_huge() {
	var (
		face = model.NewFace(
			model.NewVertex(50, -1000, 0),
			model.NewVertex(-1000, 1000, 0),
			model.NewVertex(1000, 1000, 0),
		)
		img = pngimage.WhiteImage(100, 100)
		rgb = pngimage.GreenColor()
	)
	DrawTriangle(face, img, rgb)
	if err := img.Save("testdata/pictures/triangles/huge_triangle.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Draws all faces from testdata/rabbit.obj with random colors.
func ExampleDrawTriangle_rabbitRainbow() {
	var input, err = os.Open("testdata/rabbit.obj")
	if err != nil {
		fmt.Println(err)
		return
	}
	var (
		ipt = importer.Importer{}
		m   = ipt.Import(input)
	)
	m.Transform(defaultRabbitTransformation)
	var (
		face *model.Face
		img  = pngimage.WhiteImage(2000, 2000)
	)
	for i := 0; i < m.FacesCount(); i++ {
		face = m.GetFace(i)
		DrawTriangle(face, img, pngimage.RandomColor())
	}
	if err := img.Save("testdata/pictures/rabbit_rainbow.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Draws all faces from testdata/rabbit.obj, darkening the faces that are rotated by a larger angle.
func ExampleDrawTriangle_rabbitBasicLighting() {
	var input, err = os.Open("testdata/rabbit.obj")
	if err != nil {
		fmt.Println(err)
		return
	}
	var (
		ipt = importer.Importer{}
		m   = ipt.Import(input)
	)
	m.Transform(defaultRabbitTransformation)
	var (
		face    *model.Face
		img     = pngimage.BlackImage(2000, 2000)
		x, y, z float64
		cos     float64
	)
	for i := 0; i < m.FacesCount(); i++ {
		face = m.GetFace(i)
		x, y, z = face.Normal()
		cos = z / math.Sqrt(x*x+y*y+z*z)
		if cos < 0 {
			DrawTriangle(face, img, pngimage.RGB{
				R: uint8(-255 * cos),
				G: uint8(-255 * cos),
				B: uint8(-255 * cos),
			})
		}
	}
	if err := img.Save("testdata/pictures/rabbit_basic_lighting.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Draws all faces from testdata/rabbit.obj, darkening the faces that are rotated by a larger angle.
// Uses a model method that takes into account the overlap of the faces.
func ExampleModel_BasicLighting_rabbit() {
	var input, err = os.Open("testdata/rabbit.obj")
	if err != nil {
		fmt.Println(err)
		return
	}
	var (
		ipt = importer.Importer{}
		m   = ipt.Import(input)
	)
	m.Transform(defaultRabbitTransformation)
	var img = pngimage.BlackImage(2000, 2000)
	m.BasicLighting(img, pngimage.WhiteColor())
	if err := img.Save("testdata/pictures/rabbit_z_buffer.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Draws all faces from testdata/fox.obj, darkening the faces that are rotated by a larger angle.
// Uses a model method that takes into account the overlap of the faces.
func ExampleModel_BasicLighting_fox() {
	var input, err = os.Open("testdata/fox.obj")
	if err != nil {
		fmt.Println(err)
		return
	}
	var (
		ipt = importer.Importer{}
		m   = ipt.Import(input)
	)
	m.Transform(defaultFoxTransformation)
	var img = pngimage.BlackImage(1000, 1000)
	m.BasicLighting(img, pngimage.RGB{
		R: 224,
		G: 90,
		B: 0,
	})
	if err := img.Save("testdata/pictures/fox_z_buffer.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}
