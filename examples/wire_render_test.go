package examples

import (
	"computer_graphics/model"
	"computer_graphics/obj/importer"
	"computer_graphics/pngimage"
	"fmt"
	"os"
)

// Draws the sides of all the faces of the model.
func WireRender(m *model.Model, img *pngimage.Image, rgb pngimage.RGB) {
	var (
		face       *model.Face
		v1, v2, v3 model.Vertex
	)
	for i := 0; i < m.FacesCount(); i++ {
		face = m.GetFace(i)
		v1 = face.Vertex1()
		v2 = face.Vertex2()
		v3 = face.Vertex3()
		img.Line(int(v1.X), int(v1.Y), int(v2.X), int(v2.Y), rgb)
		img.Line(int(v1.X), int(v1.Y), int(v3.X), int(v3.Y), rgb)
		img.Line(int(v2.X), int(v2.Y), int(v3.X), int(v3.Y), rgb)
	}
}

// Draws all sides of the faces from the testdata/rabbit.obj.
func ExampleWireRender_rabbit() {
	var input, err = os.Open("testdata/rabbit.obj")
	if err != nil {
		fmt.Println(err)
		return
	}
	var (
		ipt = importer.Importer{}
		m   = ipt.Import(input)
		img = pngimage.WhiteImage(2000, 2000)
	)
	m.Transform(defaultRabbitTransformation)
	WireRender(m, img, pngimage.BlackColor())
	err = img.Save("testdata/pictures/rabbit_faces_sides.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = input.Close()
	if err == nil {
		fmt.Println("Ok")
	} else {
		fmt.Println(err)
	}
	// Output: Ok
}

// Draws all sides of the faces from the testdata/fox.obj.
func ExampleWireRender_fox() {
	var input, err = os.Open("testdata/fox.obj")
	if err != nil {
		fmt.Println(err)
		return
	}
	var (
		ipt = importer.Importer{}
		m   = ipt.Import(input)
		img = pngimage.BlackImage(1000, 1000)
	)
	m.Transform(defaultFoxTransformation)
	WireRender(m, img, pngimage.WhiteColor())
	err = img.Save("testdata/pictures/fox_faces_sides.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = input.Close()
	if err == nil {
		fmt.Println("Ok")
	} else {
		fmt.Println(err)
	}
	// Output: Ok
}
