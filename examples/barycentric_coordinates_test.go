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
func DrawTriangle(v1, v2, v3 *model.Vertex, img *pngimage.Image, rgb pngimage.RGB) {
	var (
		xMax       = math.Min(float64(img.Width()), mathutils.Max(v1.X, v2.X, v3.X))
		xMin       = math.Max(0, mathutils.Min(v1.X, v2.X, v3.X))
		yMax       = math.Min(float64(img.Height()), mathutils.Max(v1.Y, v2.Y, v3.Y))
		yMin       = math.Max(0, mathutils.Min(v1.Y, v2.Y, v3.Y))
		l1, l2, l3 float64
		x, y       float64
	)
	for i := int(math.Ceil(xMin)); float64(i) < xMax; i++ {
		for j := int(math.Ceil(yMin)); float64(j) < yMax; j++ {
			x = float64(i)
			y = float64(j)
			l1 = ((v2.X-v3.X)*(y-v3.Y) - (v2.Y-v3.Y)*(x-v3.X)) / ((v2.X-v3.X)*(v1.Y-v3.Y) - (v2.Y-v3.Y)*(v1.X-v3.X))
			l2 = ((v3.X-v1.X)*(y-v1.Y) - (v3.Y-v1.Y)*(x-v1.X)) / ((v3.X-v1.X)*(v2.Y-v1.Y) - (v3.Y-v1.Y)*(v2.X-v1.X))
			l3 = ((v1.X-v2.X)*(y-v2.Y) - (v1.Y-v2.Y)*(x-v2.X)) / ((v1.X-v2.X)*(v3.Y-v2.Y) - (v1.Y-v2.Y)*(v3.X-v2.X))
			if l1 > 0 && l2 > 0 && l3 > 0 {
				img.Set(i, j, rgb)
			}
		}
	}
}

// Draws a triangle that fits completely into the image.
func ExampleDrawTriangle_internal() {
	var (
		v1  = model.NewVertex(50, 10, 0)
		v2  = model.NewVertex(10, 80, 0)
		v3  = model.NewVertex(90, 80, 0)
		img = pngimage.WhiteImage(100, 100)
		rgb = pngimage.GreenColor()
	)
	DrawTriangle(v1, v2, v3, img, rgb)
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
		v1  = model.NewVertex(50, -10, 0)
		v2  = model.NewVertex(-10, 110, 0)
		v3  = model.NewVertex(120, 110, 0)
		img = pngimage.WhiteImage(100, 100)
		rgb = pngimage.GreenColor()
	)
	DrawTriangle(v1, v2, v3, img, rgb)
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
		v1  = model.NewVertex(50, -1000, 0)
		v2  = model.NewVertex(-1000, 1000, 0)
		v3  = model.NewVertex(1000, 1000, 0)
		img = pngimage.WhiteImage(100, 100)
		rgb = pngimage.GreenColor()
	)
	DrawTriangle(v1, v2, v3, img, rgb)
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
		face       *model.Face
		v1, v2, v3 model.Vertex
		img        = pngimage.WhiteImage(2000, 2000)
	)
	for i := 0; i < m.FacesCount(); i++ {
		face = m.GetFace(i)
		v1 = face.Vertex1()
		v2 = face.Vertex2()
		v3 = face.Vertex3()
		DrawTriangle(&v1, &v2, &v3, img, pngimage.RandomColor())
	}
	if err := img.Save("testdata/pictures/rabbit_rainbow.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Draws all faces from testdata/rabbit.obj, darkening the faces that are rotated by a larger angle.
func ExampleDrawTriangle_rabbitBarycentricCoordinates() {
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
		face       *model.Face
		v1, v2, v3 model.Vertex
		img        = pngimage.BlackImage(2000, 2000)
		x, y, z    float64
		cos        float64
	)
	for i := 0; i < m.FacesCount(); i++ {
		face = m.GetFace(i)
		x, y, z = face.CalculateNormal()
		cos = z / math.Sqrt(x*x+y*y+z*z)
		if cos < 0 {
			v1 = face.Vertex1()
			v2 = face.Vertex2()
			v3 = face.Vertex3()
			DrawTriangle(&v1, &v2, &v3, img, pngimage.RGB{
				R: uint8(-255 * cos),
				G: uint8(-255 * cos),
				B: uint8(-255 * cos),
			})
		}
	}
	if err := img.Save("testdata/pictures/rabbit_barycentric_coordinates.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}
