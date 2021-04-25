package importer

import (
	"computer_graphics/model"
	"computer_graphics/obj/parser"
	"computer_graphics/obj/parser/types"
	"fmt"
	"io"
)

// Allows you to import a model from a .obj file.
// Display information about problems that occur during importing.
// You can disable the output by using the IgnoreInfos, IgnoreWarnings and IgnoreErrors fields.
// You can also specify io.Writer to output this information to.
type Importer struct {
	Output         io.Writer // Recipient of error and warning messages.
	IgnoreInfos    bool      // If true, no info messages will be output to the Output.
	IgnoreWarnings bool      // If true, no warning messages will be output to the Output.
	IgnoreErrors   bool      // If true, no error messages will be output to the Output.
}

// Reads the full model.Model from io.Reader.
// Handles errors according to the settings in the fields.
func (i *Importer) Import(in io.Reader) *model.Model {
	// Setting up the parser.
	var p = parser.NewParser(in)
	p.Output(i.Output)
	p.IgnoreErrors(i.IgnoreErrors)
	p.IgnoreWarnings(i.IgnoreWarnings)
	// Reading the model.
	var m = model.NewModel()
	i.importVertices(p, m)
	i.importNormals(p, m)
	i.importFaces(p, m)
	return m
}

// Outputs a message in Output in the format:
// [INFO] {msg}
func (i *Importer) info(msg string) {
	if i.Output != nil && !i.IgnoreInfos {
		fmt.Fprintln(i.Output, "[INFO]", msg)
	}
}

// Outputs a message in Output in the format:
// [WARNING] line: {line}, message: {msg}
func (i *Importer) warning(line int, msg string) {
	if i.Output != nil && !i.IgnoreWarnings {
		fmt.Fprintf(i.Output, "[WARNING] line: %d, message: %s\n", line, msg)
	}
}

// Outputs a message in Output in the format:
// [ERROR] line: {line}, message: {msg}
func (i *Importer) error(line int, msg string) {
	if i.Output != nil && !i.IgnoreErrors {
		fmt.Fprintf(i.Output, "[ERROR] line: %d, message: %s\n", line, msg)
	}
}

// Imports a single vertex of the model.
func (i *Importer) importVertex(line int, v *types.Vertex, m *model.Model) {
	if v.W != 0 {
		i.warning(line, "vertex weights are not supported")
	}
	m.AppendVertex(v.X, v.Y, v.Z)
}

// Imports all vertices of the model.
func (i *Importer) importVertices(p parser.Parser, m *model.Model) {
	var (
		elementType parser.ElementType
		element     interface{}
		line        int
	)
	for {
		elementType, element = p.Next()
		line = p.Line()
		switch elementType {
		case parser.Vertex:
			i.importVertex(line, element.(*types.Vertex), m)
		case parser.VertexNormal:
			i.importNormal(line, element.(*types.Normal), m)
		case parser.Face, parser.EndOfFile:
			return
		default:
			i.error(line, fmt.Sprintf("An impossible element was read: %s", elementType))
			return
		}
	}
}

// Imports a single vertex normal of the model.
func (i *Importer) importNormal(line int, v *types.Normal, m *model.Model) {
	m.AppendNormal(v.X, v.Y, v.Z)
}

// Imports all model vertex normals.
func (i *Importer) importNormals(p parser.Parser, m *model.Model) {
	var (
		elementType parser.ElementType
		element     interface{}
		line        int
	)
	for {
		elementType, element = p.Next()
		line = p.Line()
		switch elementType {
		case parser.VertexNormal:
			i.importNormal(line, element.(*types.Normal), m)
		case parser.Face, parser.EndOfFile:
			return
		default:
			i.error(line, fmt.Sprintf("An impossible element was read: %s", elementType))
			return
		}

	}
}

// Imports a single face of the model.
func (i *Importer) importFace(line int, f *types.Face, m *model.Model) {
	if len(f.Vertices) > 3 {
		i.warning(line, "only triangular faces are supported, the first three vertices will be used as a triangle")
	}
	if f.Vertices[0].Texture != 0 {
		i.warning(line, "vertex textures are not supported")
	}
	var err = m.AppendFace(f.Vertices[0].Index, f.Vertices[1].Index, f.Vertices[2].Index, f.Vertices[0].Normal, f.Vertices[1].Normal, f.Vertices[2].Normal)
	if err != nil {
		i.error(line, err.Error())
	}
}

// Imports all faces of the model.
func (i *Importer) importFaces(p parser.Parser, m *model.Model) {
	var (
		elementType parser.ElementType
		element     interface{}
		line        int
	)
	for {
		elementType, element = p.Next()
		line = p.Line()
		switch elementType {
		case parser.Face:
			i.importFace(line, element.(*types.Face), m)
		case parser.Vertex:
			i.error(line, "incorrect order of elements (vertices must be defined before faces), the vertex will be skipped")
		case parser.VertexNormal:
			i.error(line, "incorrect order of elements (normal must be defined before faces), the vertex will be skipped")
		case parser.EndOfFile:
			return
		default:
			i.error(line, fmt.Sprintf("An impossible element was read: %s", elementType))
			return
		}
	}
}
