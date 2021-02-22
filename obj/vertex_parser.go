package obj

import (
	"computer_graphics/obj/scanner"
	"computer_graphics/three_dim_model"
	"strconv"
)

// Parser of the vertex of the three-dimensional model from the .obj file.
// The format of the input data: v x y z [w].
// x - X coordinate of the vertex.
// y - Y coordinate of the vertex.
// z - Z coordinate of the vertex.
// w - weight required for rational curves and surfaces (not supported).
type vertexParser struct {
	point      three_dim_model.Vertex // The read vertex of the model.
	baseParser                        // Error message.
}

// Implements the fluentParser interface.
func (parser *vertexParser) action(tokenType scanner.TokenType, state stateType, token string) stateType {
	switch state {
	// Initial condition: initialize the vertex of the model.
	case start:
		if tokenType == scanner.SPACE {
			parser.point = three_dim_model.Vertex{}
			return 3
		}
		parser.baseParser = baseParser("impossible token was obtained in the start state - " + tokenType.Name())
		return err
	// Impossible situation.
	case err:
		parser.baseParser = "parser cannot be used in the error state"
		return err
	// Impossible situation.
	case warn:
		parser.baseParser = "parser cannot be used in the warning state"
		return err
	// Parsing X coordinate, expected: scanner.FLOAT.
	case 3:
		if tokenType == scanner.INT || tokenType == scanner.FLOAT {
			var e error
			parser.point.X, e = strconv.ParseFloat(token, 64)
			if e != nil {
				parser.baseParser = "could not get float for X coordinate from its string representation"
				return err
			}
			return 4
		}
		if tokenType == scanner.WORD || tokenType == scanner.SLASH || tokenType == scanner.UNKNOWN {
			parser.baseParser = baseParser("invalid X coordinate, expected: FLOAT, received: " + tokenType.Name())
			return err
		}
		if tokenType == scanner.EOL || tokenType == scanner.EOF {
			parser.baseParser = "X, Y, and Z coordinates are not set"
			return err
		}
		parser.baseParser = baseParser("impossible token was obtained when reading the X coordinate - " + tokenType.Name())
		return err
	// Searching for a delimiter between X and Y coordinates, expected: scanner.SPACE.
	case 4:
		if tokenType == scanner.SPACE {
			return 5
		}
		if tokenType == scanner.SLASH {
			parser.baseParser = "invalid delimiter format between X and Y coordinates, expected: SPACE, received: SLASH"
			return err
		}
		if tokenType == scanner.EOL || tokenType == scanner.EOF {
			parser.baseParser = "Y and Z coordinates are not set"
			return err
		}
		parser.baseParser = baseParser("impossible token was obtained after reading the X coordinate - " + tokenType.Name())
		return err
	// Parsing Y coordinate, expected: scanner.FLOAT.
	case 5:
		if tokenType == scanner.INT || tokenType == scanner.FLOAT {
			var e error
			parser.point.Y, e = strconv.ParseFloat(token, 64)
			if e != nil {
				parser.baseParser = "could not get float for Y coordinate from its string representation"
				return err
			}
			return 6
		}
		if tokenType == scanner.WORD || tokenType == scanner.SLASH || tokenType == scanner.UNKNOWN {
			parser.baseParser = baseParser("invalid Y coordinate, expected: FLOAT, received: " + tokenType.Name())
			return err
		}
		if tokenType == scanner.EOL || tokenType == scanner.EOF {
			parser.baseParser = "Y and Z coordinates are not set"
			return err
		}
		parser.baseParser = baseParser("impossible token was obtained when reading the Y coordinate - " + tokenType.Name())
		return err
	// Searching for a delimiter between Y and Z coordinates, expected: scanner.SPACE.
	case 6:
		if tokenType == scanner.SPACE {
			return 7
		}
		if tokenType == scanner.SLASH {
			parser.baseParser = "invalid delimiter format between Y and Z coordinates, expected: SPACE, received: SLASH"
			return err
		}
		if tokenType == scanner.EOL || tokenType == scanner.EOF {
			parser.baseParser = "Z coordinate is not set"
			return err
		}
		parser.baseParser = baseParser("impossible token was obtained after reading the Y coordinate - " + tokenType.Name())
		return err
	// Parsing Z coordinate, expected: scanner.FLOAT.
	case 7:
		if tokenType == scanner.INT || tokenType == scanner.FLOAT {
			var e error
			parser.point.Z, e = strconv.ParseFloat(token, 64)
			if e != nil {
				parser.baseParser = "could not get float for Z coordinate from its string representation"
				return err
			}
			return 8
		}
		if tokenType == scanner.WORD || tokenType == scanner.SLASH || tokenType == scanner.UNKNOWN {
			parser.baseParser = baseParser("invalid Z coordinate, expected: FLOAT, received: " + tokenType.Name())
			return err
		}
		if tokenType == scanner.EOL || tokenType == scanner.EOF {
			parser.baseParser = "Z coordinate is not set"
			return err
		}
		parser.baseParser = baseParser("impossible token was obtained when reading the Z coordinate - " + tokenType.Name())
		return err
	// Searching for a delimiter between Z coordinate and weight, expected: scanner.SPACE.
	case 8:
		if tokenType == scanner.EOL || tokenType == scanner.EOF {
			return start
		}
		if tokenType == scanner.SPACE {
			return 9
		}
		if tokenType == scanner.SLASH {
			parser.baseParser = "invalid delimiter format between Z coordinate and weight parameter, expected: SPACE, received: SLASH"
			return err
		}
		parser.baseParser = baseParser("impossible token was obtained after reading the Z coordinate - " + tokenType.Name())
		return err
	// Parsing weight, expected: scanner.FLOAT.
	case 9:
		if tokenType == scanner.EOL || tokenType == scanner.EOF {
			return start
		}
		if tokenType == scanner.INT || tokenType == scanner.FLOAT {
			return 10
		}
		if tokenType == scanner.WORD || tokenType == scanner.SLASH || tokenType == scanner.UNKNOWN {
			parser.baseParser = baseParser("invalid weight parameter, expected: FLOAT, received: " + tokenType.Name())
			return err
		}
		parser.baseParser = baseParser("impossible token was obtained when reading the weight parameter - " + tokenType.Name())
		return err
	// Searching for a whitespace after vertex data, expected: scanner.SPACE.
	case 10:
		if tokenType == scanner.EOL || tokenType == scanner.EOF {
			parser.baseParser = "unsupported vertex parameter - weight, the parameter will be ignored"
			return warn
		}
		if tokenType == scanner.SPACE {
			return 11
		}
		if tokenType == scanner.SLASH {
			parser.baseParser = "unexpected token received after describing a vertex - SLASH"
			return err
		}
		parser.baseParser = baseParser("impossible token was obtained after reading the weight parameter - " + tokenType.Name())
		return err
	// Search for extra data in a line, expected: scanner.EOL.
	case 11:
		if tokenType == scanner.EOL || tokenType == scanner.EOF {
			parser.baseParser = "unsupported vertex parameter - weight, the parameter will be ignored"
			return warn
		}
		if tokenType == scanner.SPACE || tokenType == scanner.COMMENT {
			parser.baseParser = "impossible token was obtained after reading the weight parameter - SPACE"
			return err
		}
		parser.baseParser = baseParser("unexpected token received after describing a vertex - " + tokenType.Name())
		return err
	}
	// Impossible situation.
	parser.baseParser = "impossible state received, supported states: 0-11"
	return err
}

// Implements the fluentParser interface.
func (parser *vertexParser) result() interface{} {
	return parser.point
}
