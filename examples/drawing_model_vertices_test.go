package examples

import (
	"computer_graphics/obj/parser"
	"computer_graphics/obj/parser/types"
	"computer_graphics/pngimage"
	"fmt"
	"image"
	"os"
)

// Draws all vertices from the testdata/rabbit.obj, using the specified coordinate transformation.
// The message output redirects to the corresponding file from the testdata/output directory.
func DrawVertices(transform func(v *types.Vertex) image.Point, name string) error {
	var (
		input  *os.File
		output *os.File
		err    error
	)
	input, err = os.Open("testdata/rabbit.obj")
	if err != nil {
		return err
	}
	output, err = os.Create(fmt.Sprintf("testdata/output/%s_rabbit_vertices_output.txt", name))
	if err != nil {
		return err
	}
	var p = parser.NewParser(input)
	p.Output(output)
	var (
		img                  = pngimage.WhiteImage(1000, 1000)
		rgb                  = pngimage.BlackColor()
		elementType, element = p.Next()
		point                image.Point
	)
	for elementType != parser.EndOfFile {
		if elementType == parser.Vertex {
			point = transform(element.(*types.Vertex))
			img.Set(point.X, point.Y, rgb)
		} else {
			fmt.Fprintf(output, "[INFO] unnecessary element: %s\n", elementType)
		}
		elementType, element = p.Next()
	}
	err = img.Save(fmt.Sprintf("testdata/pictures/%s_rabbit_vertices.png", name))
	if err != nil {
		return err
	}
	err = output.Close()
	if err != nil {
		return err
	}
	return input.Close()
}

// Drawing all vertexes using the first coordinate transformation.
// Check the testdata/output/first_rabbit_vertices_output.txt file for information about errors and warnings!
func ExampleDrawVertices_first() {
	var err = DrawVertices(
		func(v *types.Vertex) image.Point {
			return image.Point{X: int(50*v.X + 500), Y: int(-50*v.Y + 500)}
		},
		"first",
	)
	if err == nil {
		fmt.Println("Ok")
	} else {
		fmt.Println(err)
	}
	// Output: Ok
}

// Drawing all vertexes using the second coordinate transformation.
// Check the testdata/output/second_rabbit_vertices_output.txt file for information about errors and warnings!
func ExampleDrawVertices_second() {
	var err = DrawVertices(
		func(v *types.Vertex) image.Point {
			return image.Point{X: int(100*v.X + 500), Y: int(-100*v.Y + 500)}
		},
		"second",
	)
	if err == nil {
		fmt.Println("Ok")
	} else {
		fmt.Println(err)
	}
	// Output: Ok
}

// Drawing all vertexes using the third coordinate transformation.
// Check the testdata/output/third_rabbit_vertices_output.txt file for information about errors and warnings!
func ExampleDrawVertices_third() {
	var err = DrawVertices(
		func(v *types.Vertex) image.Point {
			return image.Point{X: int(500*v.X + 500), Y: int(-500*v.Y + 500)}
		},
		"third",
	)
	if err == nil {
		fmt.Println("Ok")
	} else {
		fmt.Println(err)
	}
	// Output: Ok
}

// Drawing all vertexes using the fourth coordinate transformation.
// Check the testdata/output/fourth_rabbit_vertices_output.txt file for information about errors and warnings!
func ExampleDrawVertices_fourth() {
	var err = DrawVertices(
		func(v *types.Vertex) image.Point {
			return image.Point{X: int(4000*v.X + 500), Y: int(-4000*v.Y + 500)}
		},
		"fourth",
	)
	if err == nil {
		fmt.Println("Ok")
	} else {
		fmt.Println(err)
	}
	// Output: Ok
}
