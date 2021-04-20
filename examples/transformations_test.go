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

// Converts the coordinates of a vertex to the coordinates of a pixel in the image.
func projectiveTransformation(v *model.Vertex, img *pngimage.Image, scale float64) (float64, float64) {
	var (
		width   = float64(img.Width())
		height  = float64(img.Height())
		xCenter = width / 2
		yCenter = height / 2
	)
	scale = math.Max(width, height) * scale
	return scale*v.X/v.Z + xCenter, scale*v.Y/v.Z + yCenter
}

// Draws a triangle on the image with the specified color.
func renderTriangle(face *model.Face, buffer [][]float64, img *pngimage.Image, rgb pngimage.RGB, scale float64) {
	var (
		// Vertices.
		v1 = face.Vertex1()
		v2 = face.Vertex2()
		v3 = face.Vertex3()
		// Coordinates of the vertices in the image.
		x1, y1 = projectiveTransformation(&v1, img, scale)
		x2, y2 = projectiveTransformation(&v2, img, scale)
		x3, y3 = projectiveTransformation(&v3, img, scale)
		// The boundaries of the rectangle inside which the face is located.
		xMax = math.Min(float64(img.Width()), mathutils.Max(x1, x2, x3))
		xMin = math.Max(0, mathutils.Min(x1, x2, x3))
		yMax = math.Min(float64(img.Height()), mathutils.Max(y1, y2, y3))
		yMin = math.Max(0, mathutils.Min(y1, y2, y3))
		// Barycentric coordinates.
		l1, l2, l3 float64
		// Coordinates of the current pixel.
		x, y, z float64
	)
	for i := int(xMin); float64(i) < xMax; i++ {
		for j := int(yMin); float64(j) < yMax; j++ {
			x = float64(i)
			y = float64(j)
			// Calculation of barycentric coordinates.
			l1 = ((x2-x3)*(y-y3) - (y2-y3)*(x-x3)) / ((x2-x3)*(y1-y3) - (y2-y3)*(x1-x3))
			l2 = ((x3-x1)*(y-y1) - (y3-y1)*(x-x1)) / ((x3-x1)*(y2-y1) - (y3-y1)*(x2-x1))
			l3 = ((x1-x2)*(y-y2) - (y1-y2)*(x-x2)) / ((x1-x2)*(y3-y2) - (y1-y2)*(x3-x2))
			if l1 > 0 && l2 > 0 && l3 > 0 {
				z = l1*v1.Z + l2*v2.Z + l3*v3.Z
				if z < buffer[i][j] {
					img.Set(i, img.Height()-j, rgb)
					buffer[i][j] = z
				}
			}
		}
	}
}

// Draws a model on an image using a projective coordinate transformation.
func RenderWithProjectiveTransformation(m *model.Model, img *pngimage.Image, scale float64) {
	var (
		face    *model.Face
		x, y, z float64
		cos     float64
		buffer  = make([][]float64, img.Width())
	)
	// Initializing the z-buffer.
	for i := 0; i < img.Width(); i++ {
		buffer[i] = make([]float64, img.Height())
		for j := 0; j < img.Height(); j++ {
			buffer[i][j] = math.Inf(+1)
		}
	}
	// Rendering triangles.
	for i := 0; i < m.FacesCount(); i++ {
		face = m.GetFace(i)
		x, y, z = face.Normal()
		cos = -z / math.Sqrt(x*x+y*y+z*z)
		if cos < 0 {
			renderTriangle(
				face,
				buffer,
				img,
				pngimage.RGB{
					R: uint8(-cos * 255),
					G: uint8(-cos * 255),
					B: uint8(-cos * 255),
				},
				scale,
			)
		}
	}
}

// Draws all faces from testdata/rabbit.obj, darkening the faces that are rotated by a larger angle.
// Uses the projective transformation of the model coordinates to the pixel coordinates of the image.
func ExampleRenderWithProjectiveTransformation_rabbit() {
	var input, err = os.Open("testdata/rabbit.obj")
	if err != nil {
		fmt.Println(err)
		return
	}
	var (
		ipt = importer.Importer{}
		m   = ipt.Import(input)
		img = pngimage.BlackImage(2000, 2000)
	)
	m.Rotate(0, math.Pi*3/2, 0)
	m.Shift(0.005, -0.045, 15)
	RenderWithProjectiveTransformation(m, img, 100)
	if err := img.Save("testdata/pictures/rabbit_transformations.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}
