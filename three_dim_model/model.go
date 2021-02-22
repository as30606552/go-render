package three_dim_model

// Describes a point in three-dimensional space.
// Contains three coordinates of the point: X, Y, Z.
type Vertex struct {
	X, Y, Z float64
}

// Creates a Vertex object based on its three coordinates.
func NewVertex(x float64, y float64, z float64) *Vertex {
	return &Vertex{X: x, Y: y, Z: z}
}

// Another name for the Vertex slice.
type Vertices []Vertex

// Describes a complete three-dimensional model.
type Model struct {
	vertices Vertices
}

// Creates a new three-dimensional model with zero vertices and reserves memory space for 10 vertices.
// But you can add more than 10 points to the model.
func NewModel() *Model {
	return &Model{vertices: make([]Vertex, 0, 10)}
}

// Adds a vertex to the model.
func (model *Model) AppendVertex(point Vertex) {
	model.vertices = append(model.vertices, point)
}

// Returns the vertex of the model by index.
// Be careful, the vertices are indexed from 0.
func (model *Model) GetVertex(index int) Vertex {
	return model.vertices[index]
}
