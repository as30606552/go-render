package types

// One of the possible direction values.
type DirectionType uint8

const (
	V DirectionType = iota // V direction.
	U                      // U direction.
)

// Specifies a geometric vertex.
type Vertex struct {
	X float64 `name:"X coordinate"`                     // X coordinate of the vertex.
	Y float64 `name:"Y coordinate"`                     // Y coordinate of the vertex.
	Z float64 `name:"Z coordinate"`                     // Z coordinate of the vertex.
	W float64 `name:"weight parameter" optional:"true"` // Weight required for rational curves and surfaces.
}

// Creates a new vertex.
func NewVertex() *Vertex {
	return &Vertex{}
}

// Specifies a geometric vertex normal.
type Normal struct {
	X float64 `name:"X coordinate"`                     // X coordinate of the vertex.
	Y float64 `name:"Y coordinate"`                     // Y coordinate of the vertex.
	Z float64 `name:"Z coordinate"`                     // Z coordinate of the vertex.
}

// Creates a new normal.
func NewNormal() *Normal {
	return &Normal{}
}

// Specifies a face element.
type Face struct {
	// Contains information about all vertexes of the face.
	Vertices []struct {
		Index   int `name:"index"`                   // Reference number for the vertex.
		Texture int `name:"texture" optional:"true"` // Reference number for the texture vertex.
		Normal  int `name:"normal" optional:"true"`  // Reference number for the vertex normal.
	} `name:"vertex" delimiter:"slash" min:"3"`
}

// Creates a new face.
func NewFace() *Face {
	return &Face{}
}
