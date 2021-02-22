package constants

// One of the possible types of description of model elements, according to the specification of .obj files.
type ElementType uint8

const (
	Vertex                ElementType = iota // Geometric vertices: v x y z [w].
	VertexTexture         ElementType = iota // Texture vertices: vt u [v w].
	VertexNormal          ElementType = iota // Vertex normals: vn i j k.
	VertexParameter       ElementType = iota // Parameter space vertices: vp u v w.
	CurveSurfaceType      ElementType = iota // Rational or non-rational forms of curve or surface type: cstype [rat] type.
	Degree                ElementType = iota // Degree: deg degu [degv].
	BasisMatrix           ElementType = iota // Basis matrix: bmat u matrix || bmat v matrix.
	Step                  ElementType = iota // Step size: step stepu [stepv].
	Point                 ElementType = iota // Point: p  v1 v2 v3 ...
	Line                  ElementType = iota // Line: l v1/[vt1] v2/[vt2] v3/[vt3] ...
	Face                  ElementType = iota // Face: f v1/[vt1]/[vn1] v2/[vt2]/[vn2] v3/[vt3]/[vn3] ...
	Curve                 ElementType = iota // Curve: curv u0 u1 v1 v2 ...
	Curve2D               ElementType = iota // 2D curve: curv2 vp1 vp2 vp3 ...
	Surface               ElementType = iota // Surface: surf s0 s1 t0 t1 v1/[vt1]/[vn1] v2/[vt2]/[vn2] ...
	Parameter             ElementType = iota // Parameter values: parm u p1 p2 p3 ... || parm u p1 p2 p3 ...
	Trim                  ElementType = iota // Outer trimming loop: trim u0 u1 curv2d u0 u1 curv2d ...
	Hole                  ElementType = iota // Inner trimming loop: hole u0 u1 curv2d u0 u1 curv2d ...
	SpecialCurve          ElementType = iota // Special curve: scrv u0 u1 curv2d u0 u1 curv2d ...
	SpecialPoint          ElementType = iota // Special point: sp vp1 vp2 ...
	End                   ElementType = iota // End statement: end.
	Connect               ElementType = iota // Connect: con surf_1 q0_1 q1_1 curv2d_1 surf_2 q0_2 q1_2 curv2d_2.
	Group                 ElementType = iota // Group name: g group_name1 group_name2 ...
	SmoothingGroup        ElementType = iota // Smoothing group: s group_number || s off.
	MergingGroup          ElementType = iota // Merging group: mg group_number res || mg off.
	Object                ElementType = iota // Object name: o object_name.
	BevelInterpolation    ElementType = iota // Bevel interpolation: bevel on || bevel off.
	ColorInterpolation    ElementType = iota // Color interpolation: c_interp on || c_interp off.
	DissolveInterpolation ElementType = iota // Dissolve interpolation: d_interp on || d_interp off.
	LevelOfDetail         ElementType = iota // Level of detail: lod level.
	MapLibrary            ElementType = iota // Map library: maplib filename1 filename2 ...
	UseMapping            ElementType = iota // Use mapping: usemap map_name || usemap off.
	UseMaterial           ElementType = iota // Material name: usemtl material_name.
	MaterialLibrary       ElementType = iota // Material library: mtllib filename1 filename2 ...
	ShadowObject          ElementType = iota // Shadow casting: shadow_obj filename.
	TraceObject           ElementType = iota // Ray tracing: trace_obj filename.
	CurveApproximation    ElementType = iota // Curve approximation technique: ctech technique resolution.
	SurfaceApproximation  ElementType = iota // Surface approximation technique: stech technique resolution.
	Call                  ElementType = iota // Call: call filename.ext arg1 arg2 ...
	Scmp                  ElementType = iota // Scmp: scmp filename.ext arg1 arg2 ...
	Csh                   ElementType = iota // Csh: csh command || csh -command.
	EndOfFile             ElementType = iota // A special marker that indicates that the parser has reached the end of the file.
)

// Converts a element type constant to its string representation.
var elementsMap = [...]string{
	"vertex",
	"vertex texture",
	"vertex normal",
	"vertex parameter",
	"curve surface type",
	"degree",
	"basis matrix",
	"step",
	"point",
	"line",
	"face",
	"curve",
	"curve 2D",
	"surface",
	"parameter",
	"trim",
	"hole",
	"special curve",
	"special point",
	"end",
	"connect",
	"group",
	"smoothing group",
	"merging group",
	"object",
	"bevel interpolation",
	"color interpolation",
	"dissolve interpolation",
	"level of detail",
	"map library",
	"use mapping",
	"use material",
	"material library",
	"shadow object",
	"trace object",
	"curve approximation technique",
	"surface approximation technique",
	"call command",
	"scmp command",
	"csh command",
	"end of file",
}

// Converts a element type constant to its string representation.
func (elementType ElementType) Name() string {
	return elementsMap[elementType]
}

// Sets the match between the first word in the line in .obj file and the type of the element that is written in this line.
var elementNamesMap = map[string]ElementType{
	"v":          Vertex,
	"vt":         VertexTexture,
	"vn":         VertexNormal,
	"vp":         VertexParameter,
	"cstype":     CurveSurfaceType,
	"deg":        Degree,
	"bmat":       BasisMatrix,
	"step":       Step,
	"p":          Point,
	"l":          Line,
	"f":          Face,
	"curv":       Curve,
	"curv2":      Curve2D,
	"surf":       Surface,
	"parm":       Parameter,
	"trim":       Trim,
	"hole":       Hole,
	"scrv":       SpecialCurve,
	"sp":         SpecialPoint,
	"end":        End,
	"con":        Connect,
	"g":          Group,
	"s":          SmoothingGroup,
	"mg":         MergingGroup,
	"o":          Object,
	"bevel":      BevelInterpolation,
	"c_interp":   ColorInterpolation,
	"d_interp":   DissolveInterpolation,
	"lod":        LevelOfDetail,
	"maplib":     MapLibrary,
	"usemap":     UseMapping,
	"usemtl":     UseMaterial,
	"mtllib":     MaterialLibrary,
	"shadow_obj": ShadowObject,
	"trace_obj":  TraceObject,
	"ctech":      CurveApproximation,
	"stech":      SurfaceApproximation,
	"call":       Call,
	"scmp":       Scmp,
	"csh":        Csh,
}

// Converts the first word in a line in .obj file to the type of the element that is written in this line.
func GetElementType(name string) (ElementType, bool) {
	var elementType, ok = elementNamesMap[name]
	return elementType, ok
}
