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

// Draws a triangle using the z-buffer to cut off overlapping faces.
func DrawTriangleZBuffer(v1, v2, v3 *model.Vertex, buffer [][]float64, img *pngimage.Image, rgb pngimage.RGB) {
	var (
		xMax       = math.Min(float64(img.Width()), mathutils.Max(v1.X, v2.X, v3.X))
		xMin       = math.Max(0, mathutils.Min(v1.X, v2.X, v3.X))
		yMax       = math.Min(float64(img.Height()), mathutils.Max(v1.Y, v2.Y, v3.Y))
		yMin       = math.Max(0, mathutils.Min(v1.Y, v2.Y, v3.Y))
		l1, l2, l3 float64
		x, y, z    float64
	)
	for i := int(math.Ceil(xMin)); float64(i) < xMax; i++ {
		for j := int(math.Ceil(yMin)); float64(j) < yMax; j++ {
			x = float64(i)
			y = float64(j)
			l1 = ((v2.X-v3.X)*(y-v3.Y) - (v2.Y-v3.Y)*(x-v3.X)) / ((v2.X-v3.X)*(v1.Y-v3.Y) - (v2.Y-v3.Y)*(v1.X-v3.X))
			l2 = ((v3.X-v1.X)*(y-v1.Y) - (v3.Y-v1.Y)*(x-v1.X)) / ((v3.X-v1.X)*(v2.Y-v1.Y) - (v3.Y-v1.Y)*(v2.X-v1.X))
			l3 = ((v1.X-v2.X)*(y-v2.Y) - (v1.Y-v2.Y)*(x-v2.X)) / ((v1.X-v2.X)*(v3.Y-v2.Y) - (v1.Y-v2.Y)*(v3.X-v2.X))
			if l1 > 0 && l2 > 0 && l3 > 0 {
				z = l1*v1.Z + l2*v2.Z + l3*v3.Z
				if z < buffer[i][j] {
					img.Set(i, j, rgb)
					buffer[i][j] = z
				}
			}
		}
	}
}

// Draws all faces from the model, darkening the faces that are rotated by a larger angle.
func BasicLighting(m *model.Model, img *pngimage.Image, rgb pngimage.RGB) {
	var (
		face       *model.Face
		v1, v2, v3 model.Vertex
		x, y, z    float64
		cos        float64
		buffer     = make([][]float64, img.Width())
	)
	for i := 0; i < img.Width(); i++ {
		buffer[i] = make([]float64, img.Height())
		for j := 0; j < img.Height(); j++ {
			buffer[i][j] = math.Inf(+1)
		}
	}
	for i := 0; i < m.FacesCount(); i++ {
		face = m.GetFace(i)
		x, y, z = face.Normal()
		cos = z / math.Sqrt(x*x+y*y+z*z)
		if cos < 0 {
			v1 = face.Vertex1()
			v2 = face.Vertex2()
			v3 = face.Vertex3()
			DrawTriangleZBuffer(
				&v1,
				&v2,
				&v3,
				buffer,
				img,
				pngimage.RGB{
					R: uint8(-float64(rgb.R) * cos),
					G: uint8(-float64(rgb.G) * cos),
					B: uint8(-float64(rgb.B) * cos),
				},
			)
		}
	}
}

// Draws all faces from testdata/rabbit.obj, darkening the faces that are rotated by a larger angle.
func ExampleBasicLighting_rabbit() {
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
	BasicLighting(m, img, pngimage.WhiteColor())
	if err := img.Save("testdata/pictures/rabbit_z_buffer.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}

// Draws all faces from testdata/fox.obj, darkening the faces that are rotated by a larger angle.
func ExampleBasicLighting_fox() {
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
	BasicLighting(m, img, pngimage.RGB{
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
