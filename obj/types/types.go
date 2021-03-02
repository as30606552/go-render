package types

type Vertex struct {
	X float64 `name:"X coordinate"`
	Y float64 `name:"Y coordinate"`
	Z float64 `name:"Z coordinate"`
	W float64 `name:"weight parameter" optional:"true"`
}
