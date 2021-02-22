package obj

import (
	"computer_graphics/obj/constants"
	"computer_graphics/obj/scanner"
	"fmt"
	"io"
	"os"
	"strings"
)

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
// 	In the start state, a space is always passed to the fluent pacer.
// 	This event should be used for initialization.
// 	Also, the start state is used as the state of successful completion of parsing.
// 	fluentParser must go to the start state if the data about its model element is successfully read.
// 	A fluent parser must necessarily read the string to the end,
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
// The implementation of the new parser must be registered in the parsersRegistry.
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
func (parser *Parser) Next() (constants.ElementType, interface{}) {
	var tokenType scanner.TokenType // Contains the type of token currently being processed.
	var token string                // Contains the string of the token currently being processed.
	// Skipping empty lines.
	tokenType, token = parser.scanner.Next()
	for tokenType == scanner.EOL {
		tokenType, token = parser.scanner.Next()
	}
	// When the end of the file is reached, it always returns (constants.EndOfFile, nil).
	if tokenType == scanner.EOF {
		return constants.EndOfFile, nil
	}
	// If the first token in the string is found in the registry of possible formats for describing the model element,
	// the string is processed by a parser from the registry.
	if elementType, ok := constants.GetElementType(token); tokenType == scanner.WORD && ok {
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
			parser.warning("unsupported element format - " + elementType.Name() + ", the line will be skipped")
			parser.scanner.SkipLine()
		}
	} else {
		parser.error("error in the name of the element type", token)
	}
	// If the line was not read, it means that the parser was not found in the registry
	// or an error occurred during parsing, need to search for the next element.
	return parser.Next()
}
