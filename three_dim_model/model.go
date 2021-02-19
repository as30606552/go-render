package three_dim_model

// Describes a point in three-dimensional space.
// Contains three coordinates of the point: X, Y, Z.
type Point struct {
	X, Y, Z float64
}

// Creates a Point object based on its three coordinates.
func NewPoint(x float64, y float64, z float64) *Point {
	return &Point{X: x, Y: y, Z: z}
}

// Another name for the Point slice.
type Points []Point

// Describes a complete three-dimensional model.
type Model struct {
	points Points
}

// Creates a new three-dimensional model with zero points and reserves memory space for 10 points.
// But you can add more than 10 points to the model.
func NewModel() *Model {
	return &Model{points: make([]Point, 0, 10)}
}

// Adds a point to the model.
func (model Model) AppendPoint(point Point) {
	model.points = append(model.points, point)
}

// Returns a point model for the index.
// Be careful, the points are indexed from 0.
func (model Model) GetPoint(index int) Point {
	return model.points[index]
}
