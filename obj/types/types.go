package types

type DirectionType uint8

const (
	V DirectionType = iota
	U
)

type Vertex struct {
	X float64 `name:"X coordinate"`
	Y float64 `name:"Y coordinate"`
	Z float64 `name:"Z coordinate"`
	W float64 `name:"weight parameter" optional:"true"`
}

func NewVertex() *Vertex {
	return &Vertex{}
}

type Point struct {
	Vertices []int `name:"vertex index" min:"3"`
}

func NewPoint() *Point {
	return &Point{}
}

type Face struct {
	Vertices []struct {
		Index   int `name:"index"`
		Texture int `name:"texture" optional:"true"`
		Normal  int `name:"normal" optional:"true"`
	} `name:"vertex" delimiter:"slash" min:"3"`
}

func NewFace() *Face {
	return &Face{}
}
