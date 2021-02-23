package obj

import (
	"computer_graphics/obj/scanner"
	"fmt"
	"io"
	"os"
	"strings"
)

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

// Converts a element type constant to its String representation.
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

// Converts a element type constant to its String representation.
func (elementType ElementType) String() string {
	return elementsMap[elementType]
}

// Sets the match between the first word in the line in .obj file and the type of the element that is written in this line.
var elementDeclarationsMap = map[string]ElementType{
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

// One of the possible states of a fluentParser.
type stateType uint8

const (
	start stateType = iota // The initial state and the state of successful completion of parsing.
	err   stateType = iota // The state of the error found during the parsing process.
	warn  stateType = iota // The state of successful completion of parsing with the need to report a warning.
)

// Interface of the parser of a specific element from a .obj file.
// The logic of the parser is based on a finite state machine.
//
// Parser sequentially gives the action method input tokens from the .obj file,
// starting with a space after the name of the element format supported by the fluentParser.
// fluentParser must perform the necessary actions with the input token and return its next state.
//
// The first three state values are reserved:
//
// 0 - start:
//
// 	In the start state, a space is always passed to the fluentParser.
// 	This event should be used for initialization.
// 	Also, the start state is used as the state of successful completion of parsing.
// 	fluentParser must go to the start state if the data about its model element is successfully read.
// 	A fluentParser must necessarily read the string to the end,
// 	do not go to the start state if you have not received the scanner.EOL token or the scanner.EOF token.
//
// 1 - err:
//
// 	fluentParser should go into an err state if an invalid token is received.
// 	If the error method is called, the fluentParser should return an error message.
// 	In this case, you don't need to worry about reaching the end of the line.
//
// 2 - warn:
//
// 	fluentParser should go to the warn state if the data is read correctly, but a warning must be reported.
// 	The warn state like the start state is used to indicate the successful completion of parsing.
// 	fluentParser must only enter the warn state when it reaches the end of the line.
// 	The error method is also used to get the warning message.
//
// The start and warn states are the states of successful completion of parsing.
// After transition to one of these states, the result method must return the read element.
//
// The action method takes the type of the token obtained from the. obj file,
// the previous state of the fluentParser, and the string representation of the token.
// The action method should process the received token and return the next state of the fluentParser.
//
// The implementation of the new fluentParser must be registered in the parsersRegistry.
type fluentParser interface {
	action(tokenType scanner.TokenType, state stateType, token string) stateType
	error() string
	result() interface{}
}

// The base parser implements a method for getting an error message.
// Use it to store an message when going to the err state or the warn state.
// See implementation of the vertexParser.
type baseParser string

// Implements the fluentParser interface.
func (parser baseParser) error() string {
	return string(parser)
}

// Allows you to call the Next method sequentially to get the elements from the .obj file.
// Output information about problems that occur during parsing.
// You can disable the output by using the IgnoreWarnings and IgnoreErrors fields.
// You can also specify os.Writer to output this information to.
type Parser struct {
	scanner        *scanner.Scanner // A scanner that splits the input file into tokens.
	Output         io.Writer        // Recipient of error and warning messages.
	IgnoreWarnings bool             // If true, no error messages will be output to the Output.
	IgnoreErrors   bool             // If true, no warning messages will be output to the Output.
}

// Creates a new .obj file parser.
// By default, it outputs all errors and warnings in os.Stderr.
// This can be changed by using the Parser.Output, Parser.IgnoreWarnings, Parser.IgnoreErrors fields.
func NewParser(file *os.File) *Parser {
	var s = scanner.NewScanner(file)
	s.Error = func(err error) {
		_, err = fmt.Fprintf(
			os.Stderr,
			"[ERROR] line: %d, column: %d, message: %s\n",
			s.Line(),
			s.Column(),
			err,
		)
		if err != nil {
			panic(err)
		}
	}
	return &Parser{scanner: s, Output: os.Stderr}
}

// Outputs an error message in Parser.Output in the format:
// [ERROR] line: {line number}, column: {column number}, token: '{token string}', message: {error message}
// After that, it outputs the line where the error occurred, highlighting the error token.
// Note that the method skips a line and adds information about it to the error message.
func (parser *Parser) error(msg, token string) {
	if parser.IgnoreErrors {
		parser.scanner.SkipLine()
	} else {
		var column = parser.scanner.Column() - len(token)
		fmt.Fprintf(
			parser.Output,
			"[ERROR] line: %d, column: %d, token: '%s', message: %s\n",
			parser.scanner.Line(),
			column,
			token,
			msg+", the line will be skipped",
		)
		fmt.Fprintln(
			parser.Output,
			strings.Repeat(" ", 7),
			"-> ",
			parser.scanner.SkipLine(),
			"\n",
			strings.Repeat(" ", column+10),
			strings.Repeat("^", len(token)),
		)
	}
}

// Outputs an warning message in Parser.Output in the format:
// [WARNING] line: {line number}, message: {warning message}
// Note that the method does not skip the line.
func (parser *Parser) warning(msg string) {
	if parser.IgnoreWarnings {
		parser.scanner.SkipLine()
	} else {
		fmt.Fprintf(
			parser.Output,
			"[WARNING] line: %d, message: %s\n",
			parser.scanner.Line(),
			msg,
		)
	}
}

// Returns the next element read from the file.
// Lines of unsupported format and lines containing an error are skipped and searched for matches further.
// Ensures that the returned object can be safely cast to the type defined by the constant constants.ElementType.
// When the end of the file is reached, it always returns (constants.EndOfFile, nil).
func (parser *Parser) Next() (ElementType, interface{}) {
	// Skipping empty lines.
	var tokenType, token = parser.scanner.Next()
	for tokenType == scanner.EOL {
		tokenType, token = parser.scanner.Next()
	}
	// When the end of the file is reached, it always returns (constants.EndOfFile, nil).
	if tokenType == scanner.EOF {
		return EndOfFile, nil
	}
	// If the first token in the String is found in the registry of possible formats for describing the model element,
	// the String is processed by a parser from the registry.
	if elementType, ok := elementDeclarationsMap[token]; tokenType == scanner.Word && ok {
		var elementParser = parsersRegistry[elementType]
		// If the parser from the registry is nil, then the format is not supported.
		if elementParser != nil {
			var state stateType // Contains the parser state of a specific element.
			for {
				tokenType, token = parser.scanner.Next()
				state = elementParser.action(tokenType, state, token)
				switch state {
				// The transition to the start state means the successful completion of the parser.
				case start:
					return elementType, elementParser.result()
				// The transition to the error state means an erroneous entry of the element.
				// The erroneous line must be skipped and the next element must be searched for.
				case err:
					parser.error(elementParser.error(), token)
					return parser.Next()
				case warn:
					// The transition to the warning state causes the warning to be displayed in Parser.Output.
					// Since the warning state is not a read error, but is the final state, the parser returns the result.
					parser.warning(elementParser.error())
					return elementType, elementParser.result()
				}
			}
		} else {
			parser.warning("unsupported element format - " + elementType.String() + ", the line will be skipped")
			parser.scanner.SkipLine()
		}
	} else {
		parser.error("error in the name of the element type", token)
	}
	// If the line was not read, it means that the parser was not found in the registry
	// or an error occurred during parsing, need to search for the next element.
	return parser.Next()
}
