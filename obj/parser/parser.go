package parser

import (
	"computer_graphics/obj/scanner"
	"fmt"
	"io"
	"os"
	"strings"
)

// One of the possible types of description of model types, according to the specification of .obj files.
type ElementType uint8

const (
	Vertex                ElementType = iota // Geometric vertices: v x y z [w].
	VertexTexture                            // Texture vertices: vt u [v] [w].
	VertexNormal                             // Vertex normals: vn i j k.
	VertexParameter                          // Parameter space vertices: vp u [v] [w].
	CurveSurfaceType                         // Rational or non-rational forms of curve or surface type: cstype [rat] type.
	Degree                                   // Degree: deg degu [degv].
	BasisMatrix                              // Basis matrix: bmat u matrix || bmat v matrix.
	Step                                     // Step size: step stepu [stepv].
	Point                                    // Point: p v1 v2 v3 ...
	Line                                     // Line: l v1/[vt1] v2/[vt2] v3/[vt3] ...
	Face                                     // Face: f v1/[vt1]/[vn1] v2/[vt2]/[vn2] v3/[vt3]/[vn3] ...
	Curve                                    // Curve: curv u0 u1 v1 v2 ...
	Curve2D                                  // 2D curve: curv2 vp1 vp2 vp3 ...
	Surface                                  // Surface: surf s0 s1 t0 t1 v1/[vt1]/[vn1] v2/[vt2]/[vn2] ...
	Parameter                                // Parameter values: parm u p1 p2 p3 ... || parm v p1 p2 p3 ...
	Trim                                     // Outer trimming loop: trim u0 u1 curv2d u0 u1 curv2d ...
	Hole                                     // Inner trimming loop: hole u0 u1 curv2d u0 u1 curv2d ...
	SpecialCurve                             // Special curve: scrv u0 u1 curv2d u0 u1 curv2d ...
	SpecialPoint                             // Special point: sp vp1 vp2 ...
	End                                      // End statement: end.
	Connect                                  // Connect: con surf_1 q0_1 q1_1 curv2d_1 surf_2 q0_2 q1_2 curv2d_2.
	Group                                    // Group name: g group_name1 group_name2 ...
	SmoothingGroup                           // Smoothing group: s group_number || s off.
	MergingGroup                             // Merging group: mg group_number res || mg off.
	Object                                   // Object name: o object_name.
	BevelInterpolation                       // Bevel interpolation: bevel on || bevel off.
	ColorInterpolation                       // Color interpolation: c_interp on || c_interp off.
	DissolveInterpolation                    // Dissolve interpolation: d_interp on || d_interp off.
	LevelOfDetail                            // Level of detail: lod level.
	MapLibrary                               // Map library: maplib filename1 filename2 ...
	UseMapping                               // Use mapping: usemap map_name || usemap off.
	UseMaterial                              // Material name: usemtl material_name.
	MaterialLibrary                          // Material library: mtllib filename1 filename2 ...
	ShadowObject                             // Shadow casting: shadow_obj filename.
	TraceObject                              // Ray tracing: trace_obj filename.
	CurveApproximation                       // Curve approximation technique: ctech technique resolution.
	SurfaceApproximation                     // Surface approximation technique: stech technique resolution.
	Call                                     // Call: call filename.ext arg1 arg2 ...
	Scmp                                     // Scmp: scmp filename.ext arg1 arg2 ...
	Csh                                      // Csh: csh command || csh -command.
	EndOfFile                                // A special marker that indicates that the parser has reached the end of the file.
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

// One of the possible states of a elementParser.
type stateType uint8

const (
	start stateType = iota // The initial state and the state of successful completion of parsing.
	err                    // The state of the error found during the parsing process.
)

// Interface of the parser of a specific element from a .obj file.
// The logic of the parser is based on a finite state machine.
// Implementations must contain the logic of transitions between states of the state machine
// and the logic of actions that must be performed when switching to a certain state.
//
// Parser sequentially gives the next method input tokens from the .obj file,
// starting with a space after the name of the element format supported by the elementParser.
// After that, the received state is checked: if the elementParser has moved to the start state,
// the result method is called to get the read data; if the elementParser has moved to the err state,
// the message method is called to get the error information.
//
// The first two state values are reserved:
//
// 0 - start:
//
// 	In the start state, a space is always passed to the elementParser.
// 	This event should be used for initialization.
// 	Also, the start state is used as the state of successful completion of parsing.
// 	elementParser must go to the start state if the data about its model element is successfully read.
// 	A elementParser must necessarily read the string to the end,
// 	do not go to the start state if you have not received the scanner.EOL token or the scanner.EOF token.
//
// 1 - err:
//
// 	elementParser should go into an err state if an invalid token is received.
// 	In this case, you don't need to worry about reaching the end of the line.
//
// The implementation of the new elementParser must be registered in the parsersRegistry.
// See the parsersRegistry documentation for more information.
type elementParser interface {
	// Returns the next state of the state machine based on the previous state and the received token type.
	next(tokenType scanner.TokenType, state stateType) stateType
	// Performs the necessary actions on the received token when switching to the state.
	action(state stateType, token string)
	// Returns information about the error by the state from which the elementParser went to the err state
	// and the type of token that was received when going to the err state.
	message(tokenType scanner.TokenType, state stateType) string
	// Returns a structure containing the read data from the string.
	// The  elementParser must ensure that the return value can be safely cast to the type of the element
	// that the elementParser processes.
	result() interface{}
}

// Allows you to call the Next method sequentially to get the types from the .obj file.
// Display information about problems that occur during parsing.
// You can disable the output by using the IgnoreWarnings and IgnoreErrors fields.
// You can also specify io.Writer to output this information to.
type Parser struct {
	scanner        scanner.Scanner // A scanner that splits the input file into tokens.
	Output         io.Writer       // Recipient of error and warning messages.
	IgnoreWarnings bool            // If true, no error messages will be output to the Output.
	IgnoreErrors   bool            // If true, no warning messages will be output to the Output.
}

// Creates a new .obj file parser.
// By default, it outputs all errors and warnings in os.Stderr.
// This can be changed by using the Parser.Output, Parser.IgnoreWarnings, Parser.IgnoreErrors fields.
func NewParser(reader io.Reader) *Parser {
	return &Parser{scanner: scanner.NewScanner(reader), Output: os.Stderr}
}

type logType uint8

const (
	e logType = iota
	w
)

func (t logType) String() string {
	switch t {
	case e:
		return "ERROR"
	case w:
		return "WARNING"
	default:
		panic("unknown log type")
	}
}

// Outputs a message in Output in the format:
// [{log type}] line: {line number}, column: {column number}, token: '{token string}', message: {log message}
// After that, it outputs the line where the token occurred, highlighting the token.
// Note that the method skips a line and adds information about it to the msg parameter.
func (parser *Parser) log(msg, token string, t logType) {
	if t == e && parser.IgnoreErrors || t == w && parser.IgnoreWarnings {
		parser.scanner.SkipLine()
	} else {
		var (
			tokenLength   int
			logTypeString = t.String()
		)
		switch token {
		case "\n":
			token = "eol"
			tokenLength = 1
		case "":
			token = "eof"
			tokenLength = 1
		default:
			tokenLength = len(token)
		}
		var column = parser.scanner.Column() - tokenLength + 2
		parser.scanner.SkipLine()
		fmt.Fprintf(
			parser.Output,
			"[%s] line: %d, column: %d, token: '%s', message: %s%s\n",
			logTypeString,
			parser.scanner.Line()+1,
			column,
			token,
			msg,
			", the line will be skipped",
		)
		fmt.Fprintln(
			parser.Output,
			strings.Repeat(" ", len(logTypeString)+2),
			"->",
			parser.scanner.LineString(),
			"\n",
			strings.Repeat(" ", column+len(logTypeString)+3),
			strings.Repeat("^", tokenLength),
		)
	}
}

// Returns the next element read from the reader.
// Lines of unsupported format and lines containing an error are skipped and searched for matches further.
// Ensures that the returned object can be safely cast to the type defined by the constant ElementType.
// When the end of the file is reached, it always returns (EndOfFile, nil).
func (parser *Parser) Next() (ElementType, interface{}) {
	// Skipping empty lines.
	var tokenType, token = parser.scanner.Next()
	for tokenType == scanner.EOL || tokenType == scanner.Space {
		tokenType, token = parser.scanner.Next()
	}
	// When the end of the file is reached, it always returns (EndOfFile, nil).
	if tokenType == scanner.EOF {
		return EndOfFile, nil
	}
	// If the first token in the String is found in the registry of possible formats for describing the model element,
	// the String is processed by a parser from the registry.
	if elementType, ok := elementDeclarationsMap[token]; tokenType == scanner.Word && ok {
		var p = parsersRegistry[elementType]
		// If the parser from the registry is nil, then the format is not supported.
		if p != nil {
			var (
				prevState stateType // Contains the previous state of the parser to get the error message.
				state     stateType // Contains the parser state of a specific element.
			)
			for {
				tokenType, token = parser.scanner.Next()
				prevState = state
				state = p.next(tokenType, state)
				switch state {
				// The transition to the start state means the successful completion of the parser.
				case start:
					return elementType, p.result()
				// The transition to the error state means an erroneous entry of the element.
				// The erroneous line must be skipped and the next element must be searched for.
				case err:
					parser.log(p.message(tokenType, prevState), token, e)
					return parser.Next()
				default:
					p.action(prevState, token)
				}
			}
		} else {
			parser.log("unsupported element format - "+elementType.String(), token, w)
		}
	} else {
		parser.log("error in the name of the element type", token, e)
	}
	// If the line was not read, it means that the parser was not found in the registry,
	// need to search for the next element.
	return parser.Next()
}
