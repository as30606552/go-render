package parser

import "computer_graphics/obj/parser/types"

// A registry of parsers for each type of element in the .obj file.
// To add support for the new model description format, you need to implement a parser for this element
// and put this parser in the registry.
// The parser index in the registry must match the value of the ElementType constant corresponding to the element type.
// Look at the comments on the lines of the registry.
var parsersRegistry = [...]elementParser{
	buildParser(Vertex, types.NewVertex()), // Vertex
	nil,                                    // VertexTexture
	nil,                                    // VertexNormal
	nil,                                    // VertexParameter
	nil,                                    // CurveSurfaceType
	nil,                                    // Degree
	nil,                                    // BasisMatrix
	nil,                                    // Step
	nil,                                    // Point
	nil,                                    // Line
	buildParser(Face, types.NewFace()),     // Face
	nil,                                    // Curve
	nil,                                    // Curve2D
	nil,                                    // Surface
	nil,                                    // Parameter
	nil,                                    // Trim
	nil,                                    // Hole
	nil,                                    // SpecialCurve
	nil,                                    // SpecialPoint
	nil,                                    // End
	nil,                                    // Connect
	nil,                                    // Group
	nil,                                    // SmoothingGroup
	nil,                                    // MergingGroup
	nil,                                    // Object
	nil,                                    // BevelInterpolation
	nil,                                    // ColorInterpolation
	nil,                                    // DissolveInterpolation
	nil,                                    // LevelOfDetail
	nil,                                    // MapLibrary
	nil,                                    // UseMapping
	nil,                                    // UseMaterial
	nil,                                    // MaterialLibrary
	nil,                                    // ShadowObject
	nil,                                    // TraceObject
	nil,                                    // CurveApproximation
	nil,                                    // SurfaceApproximation
	nil,                                    // Call
	nil,                                    // Scmp
	nil,                                    // Csh
}
