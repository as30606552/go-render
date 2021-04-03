package examples

import (
	"computer_graphics/obj/importer"
	"computer_graphics/pngimage"
	"fmt"
	"os"
)

// Draws all sides of the faces from the testdata/rabbit.obj.
func ExampleModel_WireRender_rabbit() {
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
	m.WireRender(img, pngimage.BlackColor())
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
func ExampleModel_WireRender_fox() {
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
	m.WireRender(img, pngimage.WhiteColor())
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
