package model

import (
	"errors"
	"fmt"
)

// Describes a vertex in three-dimensional space.
// Contains three coordinates of the vertex: X, Y, Z.
type Vertex struct {
	X, Y, Z float64
}

// Creates a Vertex based on its three coordinates.
func newVertex(x, y, z float64) *Vertex {
	return &Vertex{X: x, Y: y, Z: z}
}

// Describes a triangle in three-dimensional space.
// Contains three vertices of the triangle.
type Face struct {
	vertex1, vertex2, vertex3 *Vertex
}

// Returns the first vertex of the triangle.
func (f Face) Vertex1() Vertex {
	return *f.vertex1
}

// Returns the second vertex of the triangle.
func (f Face) Vertex2() Vertex {
	return *f.vertex2
}

// Returns the third vertex of the triangle.
func (f Face) Vertex3() Vertex {
	return *f.vertex3
}

// Creates a Face based on its three vertices.
func newFace(vertex1, vertex2, vertex3 *Vertex) *Face {
	return &Face{
		vertex1: vertex1,
		vertex2: vertex2,
		vertex3: vertex3,
	}
}

// Describes a complete three-dimensional model.
type Model struct {
	vertices []*Vertex // A list of all the vertices of the model.
	faces    []*Face   // A list of all the faces of the model.
}

// Returns a pointer to a vertex by its index and an error if the index is specified incorrectly.
// Supports negative indexing, the index of the first vertex is 1.
func (model *Model) vertexByIndex(index int) (*Vertex, error) {
	var verticesCount = len(model.vertices)
	if index > 0 {
		if index <= verticesCount {
			return model.vertices[index-1], nil
		} else {
			return nil, fmt.Errorf("unresolved vertex index: %d", index)
		}
	} else if index < 0 {
		if -index <= verticesCount {
			return model.vertices[verticesCount+index], nil
		} else {
			return nil, fmt.Errorf("unresolved vertex index: %d", index)
		}
	} else {
		return nil, errors.New("vertex index cannot be zero")
	}
}

// Adds a vertex to the model based on its three coordinates.
func (model *Model) AppendVertex(x, y, z float64) {
	model.vertices = append(model.vertices, newVertex(x, y, z))
}

// Returns the vertex of the model by index and an error if the index is specified incorrectly.
// Supports negative indexing, the index of the first vertex is 1.
func (model *Model) GetVertex(index int) (Vertex, error) {
	var v, err = model.vertexByIndex(index)
	return *v, err
}

// Returns the number of model vertices.
func (model *Model) VerticesCount() int {
	return len(model.vertices)
}

// Adds a face to the model based on its three vertices.
func (model *Model) AppendFace(v1, v2, v3 int) error {
	var (
		err     error
		vertex1 *Vertex
		vertex2 *Vertex
		vertex3 *Vertex
	)
	if vertex1, err = model.vertexByIndex(v1); err != nil {
		return err
	}
	if vertex2, err = model.vertexByIndex(v2); err != nil {
		return err
	}
	if vertex3, err = model.vertexByIndex(v3); err != nil {
		return err
	}
	model.faces = append(model.faces, newFace(vertex1, vertex2, vertex3))
	return nil
}

// Returns the vertex of the model by index.
func (model *Model) GetFace(index int) Face {
	return *model.faces[index]
}

// Returns the number of model faces.
func (model *Model) FacesCount() int {
	return len(model.faces)
}

// Creates a new three-dimensional model with zero vertices and reserves memory space for 10 vertices and 10 faces.
// But you can add more than 10 vertices and faces to the model.
func NewModel() *Model {
	return &Model{
		vertices: make([]*Vertex, 0, 10),
		faces:    make([]*Face, 0, 10),
	}
}
