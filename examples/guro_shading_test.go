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
func projectiveTransformationGuro(v *model.Vertex, img *pngimage.Image, scale float64) (float64, float64) {
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
func renderTriangleGuro(face *model.Face, buffer [][]float64, img *pngimage.Image, scale float64) {
	var (
		// Vertices.
		v1 = face.Vertex1()
		v2 = face.Vertex2()
		v3 = face.Vertex3()
		// Coordinates of the vertices in the image.
		x1, y1 = projectiveTransformationGuro(&v1, img, scale)
		x2, y2 = projectiveTransformationGuro(&v2, img, scale)
		x3, y3 = projectiveTransformationGuro(&v3, img, scale)
		// The boundaries of the rectangle inside which the face is located.
		xMax = math.Min(float64(img.Width()), mathutils.Max(x1, x2, x3))
		xMin = math.Max(0, mathutils.Min(x1, x2, x3))
		yMax = math.Min(float64(img.Height()), mathutils.Max(y1, y2, y3))
		yMin = math.Max(0, mathutils.Min(y1, y2, y3))
		// Barycentric coordinates.
		lambda1, lambda2, lambda3 float64
		// Polygon illumination.
		factor float64
		// Coordinates of the current pixel.
		x, y, z float64
		// Polygon vertex normals.
		normal1, normal2, normal3 model.Normal
		// Scalar product calculation coefficients.
		l1, l2, l3 float64
	)
	for i := int(xMin); float64(i) < xMax; i++ {
		for j := int(yMin); float64(j) < yMax; j++ {
			x = float64(i)
			y = float64(j)
			normal1 = face.Normal1()
			normal2 = face.Normal2()
			normal3 = face.Normal3()
			// Calculation polygon illumination using barycentric coordinates.
			l1 = normal1.Z / math.Sqrt(normal1.X*normal1.X + normal1.Y*normal1.Y + normal1.Z*normal1.Z)
			l2 = normal2.Z / math.Sqrt(normal2.X*normal2.X + normal2.Y*normal2.Y + normal2.Z*normal2.Z)
			l3 = normal3.Z / math.Sqrt(normal3.X*normal3.X + normal3.Y*normal3.Y + normal3.Z*normal3.Z)
			// Calculation of barycentric coordinates.
			lambda1 = ((x2-x3)*(y-y3) - (y2-y3)*(x-x3)) / ((x2-x3)*(y1-y3) - (y2-y3)*(x1-x3))
			lambda2 = ((x3-x1)*(y-y1) - (y3-y1)*(x-x1)) / ((x3-x1)*(y2-y1) - (y3-y1)*(x2-x1))
			lambda3 = ((x1-x2)*(y-y2) - (y1-y2)*(x-x2)) / ((x1-x2)*(y3-y2) - (y1-y2)*(x3-x2))
			// Calculation polygon illumination.
			factor = l1*lambda1 + l2*lambda2 +l3*lambda3
			if lambda1 > 0 && lambda2 > 0 && lambda3 > 0 {
				z = lambda1*v1.Z + lambda2*v2.Z + lambda3*v3.Z
				if z < buffer[i][j] {
					img.Set(i, img.Height()-j, pngimage.RGB{
						R: uint8(-factor * 255),
						G: uint8(-factor * 255),
						B: uint8(-factor * 255),
					},)
					buffer[i][j] = z
				}
			}
		}
	}
}

// Draws a model on an image using the Guro shading method.
func TransformationGuro(m *model.Model, img *pngimage.Image, scale float64) {
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
		x, y, z = face.CalculateNormal()
		cos = -z / math.Sqrt(x*x+y*y+z*z)
		if cos < 0 {
			renderTriangleGuro(
				face,
				buffer,
				img,
				scale,
			)
		}
	}
}

// Draws all faces from testdata/rabbit.obj, darkening the faces that are rotated by a larger angle.
// Implements advanced lighting using the Guro shading method.
// Uses the projective transformation of the model coordinates to the pixel coordinates of the image.
func ExampleTransformationGuro_rabbit() {
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
	m.Shift(0.005, -0.045, 15)
	TransformationGuro(m, img, 100)
	if err := img.Save("testdata/pictures/rabbit_guro_shading.png"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Ok")
	}
	// Output: Ok
}



