package examples

import (
	"computer_graphics/model"
	"computer_graphics/obj/importer"
	"computer_graphics/pngimage"
	"fmt"
	"image"
	"os"
)

// Draws all faces from the testdata/{name}.obj, using the specified coordinate transformation.
func DrawFaces(transform func(v model.Vertex) image.Point, img *pngimage.Image, rgb pngimage.RGB, name string) error {
	var input, err = os.Open(fmt.Sprintf("testdata/%s.obj", name))
	if err != nil {
		return err
	}
	var (
		ipt        = importer.Importer{}
		m          = ipt.Import(input)
		face       model.Face
		v1, v2, v3 image.Point
	)
	for i := 0; i < m.FacesCount(); i++ {
		face = m.GetFace(i)
		v1 = transform(face.Vertex1())
		v2 = transform(face.Vertex2())
		v3 = transform(face.Vertex3())
		img.Line(v1, v2, rgb)
		img.Line(v1, v3, rgb)
		img.Line(v2, v3, rgb)
	}
	err = img.Save(fmt.Sprintf("testdata/pictures/%s_faces.png", name))
	if err != nil {
		return err
	}
	return input.Close()
}

// Draws all faces from the testdata/rabbit.obj.
func ExampleDrawFaces_rabbit() {
	err := DrawFaces(
		func(v model.Vertex) image.Point {
			return image.Point{X: int(10000*v.X + 1000), Y: int(-10000*v.Y + 1500)}
		},
		pngimage.WhiteImage(2000, 2000),
		pngimage.BlackColor(),
		"rabbit",
	)
	if err == nil {
		fmt.Println("Ok")
	} else {
		fmt.Println(err)
	}
	// Output: Ok
}

// Draws all faces from the testdata/fox.obj.
func ExampleDrawFaces_fox() {
	err := DrawFaces(
		func(v model.Vertex) image.Point {
			return image.Point{X: int(-5*v.Z + 500), Y: int(-5*v.Y + 700)}
		},
		pngimage.BlackImage(1000, 1000),
		pngimage.WhiteColor(),
		"fox",
	)
	if err == nil {
		fmt.Println("Ok")
	} else {
		fmt.Println(err)
	}
	// Output: Ok
}
