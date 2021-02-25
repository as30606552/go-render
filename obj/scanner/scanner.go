package scanner

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// One of the possible values that the Scanner.Next method returns.
type TokenType uint8

const (
	Word    TokenType = iota // Can consist of letters, numbers, and underscores. Can not start with a number.
	Int                      // Consists of numbers. Can start with a minus.
	Float                    // Consists of digits with a dot between them. Can start with a minus.
	Slash                    // '/' character.
	Space                    // A sequence of spaces and/or tabs.
	EOL                      // '\n' character.
	EOF                      // Indicates that the end of the sequence of bytes being read has been reached.
	Unknown                  // Unknown type of token.
	Comment                  // Starts with the '#' character and ends with the character before the end of the line.
)

// Converts the state of the finite state machine from which it moved to the initial state to the type of the read token.
// See https://github.com/as30606552/ComputerGraphicsProject/wiki/Scanner.
var tokenTypeMap = [...]TokenType{Unknown, Comment, EOL, Space, Slash, Unknown, Unknown, Int, Float, Word, Unknown}

// Converts a token type constant to its string representation.
var tokenTypeNamesMap = [...]string{"WORD", "INT", "FLOAT", "SLASH", "SPACE", "EOL", "EOF", "UNKNOWN", "COMMENT"}

// Converts a token type constant to its string representation.
func (tokenType TokenType) String() string {
	return tokenTypeNamesMap[tokenType]
}

// One of the possible states of a finite state machine.
// See https://github.com/as30606552/ComputerGraphicsProject/wiki/Scanner.
type stateType uint8

const (
	start      stateType = iota // Initial state.
	skipLine                    // Skipping all characters up to the '\n' character.
	foundEol                    // '\n' character found.
	foundSpace                  // Whitespace character found.
	foundSlash                  // '/' character found.
	foundMinus                  // '-' character was found at the beginning of the token, and a digit is expected.
	foundDot                    // A '.' character is found after an integer, a digit is expected.
	foundInt                    // '\n' character found.
	foundFloat                  // A sequence of characters satisfying the Float token is found, a digit is expected.
	foundWord                   // A sequence of characters satisfying the Word token is found.
	unknown                     // An unknown sequence of characters was found.
)

// One of the possible character types that can be contained in a sequence of bytes to be read.
type symbolType uint8

const (
	eol    symbolType = iota // '\n'
	space                    // ' ' or '\t'
	hash                     // '#'
	slash                    // '/'
	minus                    // '-'
	dot                      // '.'
	digit                    // '0' - '9'
	letter                   // 'a' - 'z' or 'A' - 'Z' or '_'
	other                    // Any other character.
)

// Calculates the character type.
func getSymbolType(symbol byte) symbolType {
	switch symbol {
	case '\n':
		return eol
	case ' ':
		return space
	case '\t':
		return space
	case '#':
		return hash
	case '/':
		return slash
	case '-':
		return minus
	case '.':
		return dot
	case '_':
		return letter
	}
	if '0' <= symbol && symbol <= '9' {
		return digit
	}
	if 'a' <= symbol && symbol <= 'z' || 'A' <= symbol && symbol <= 'Z' {
		return letter
	}
	return other
}

// The finite state machine matrix.
// See https://github.com/as30606552/ComputerGraphicsProject/wiki/Scanner.
var matrix = [9][11]stateType{
	{foundEol, start, start, start, start, start, start, start, start, start, start},
	{foundSpace, skipLine, start, foundSpace, start, start, start, start, start, start, start},
	{skipLine, skipLine, start, start, start, start, start, start, start, start, start},
	{foundSlash, skipLine, start, start, start, start, start, start, start, start, start},
	{foundMinus, skipLine, start, start, start, unknown, unknown, unknown, unknown, unknown, unknown},
	{unknown, skipLine, start, start, start, unknown, unknown, foundDot, unknown, unknown, unknown},
	{foundInt, skipLine, start, start, start, foundInt, foundFloat, foundInt, foundFloat, foundWord, unknown},
	{foundWord, skipLine, start, start, start, unknown, unknown, unknown, unknown, foundWord, unknown},
	{unknown, skipLine, start, start, start, unknown, unknown, unknown, unknown, unknown, unknown},
}

// The size of the buffer in which the Scanner stores the read characters.
const bufsize uint8 = 255

// Allows you to sequentially call the Next method to get tokens from a io.Reader that can occur in .obj files.
type Scanner struct {
	reader io.Reader // The io.Reader from which the tokens will be read.
	init   bool      // Contains true if there has already been an attempt to extract a byte from the buffer.

	buffer  [bufsize]byte // Temporary storage for bytes extracted from the reader but not yet processed.
	bufpos  uint8         // The position of the currently processed byte in the buffer.
	buflast uint8         // The number of bytes contained in the buffer.

	lineStr    []byte // Current processed line string.
	switchLine bool   // true if the scanner has moved to a new line.
	line       int    // The number of the currently processed line.
	position   int    // The position of the currently processed character relative to the beginning of the byte sequence.

	Error        func(error) // The function called in case of an error.
	SkipComments bool        // true if comments should be skipped.
}

// Creates a new Scanner that reads from the reader.
// Sets skipping comments by default.
func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{reader: reader, SkipComments: true}
}

// Delegates the execution of the Error function.
// If the function is not specified, by default it outputs the error data to os.Stderr.
func (scanner *Scanner) error(err error) {
	if scanner.Error != nil {
		scanner.Error(err)
		return
	}
	_, err = fmt.Fprintf(os.Stderr, "Error in position %d: %s\n", scanner.position, err.Error())
	if err != nil {
		panic(err)
	}
}

// Reads new values to the buffer.
// The number of bytes read is stored in the buflast field.
// The current buffer position is reset to 0.
func (scanner *Scanner) refreshBuffer() {
	var n, err = scanner.reader.Read(scanner.buffer[:])
	if err != nil && err != io.EOF {
		scanner.error(err)
	}
	scanner.buflast = uint8(n)
	scanner.bufpos = 0
}

// Moving the scanner to the next line.
func (scanner *Scanner) refreshLine() {
	scanner.lineStr = make([]byte, 0, 100)
	scanner.line++
}

// Returns true if there is a next token.
func (scanner *Scanner) has() bool {
	// The buffer is processed to the end.
	// It is necessary to read the new data to the buffer.
	if scanner.bufpos == scanner.buflast {
		// If the number of elements in the buffer is less than the buffer size,
		// it means that the buffer was not fully filled the previous time when reading it.
		if scanner.buflast < bufsize {
			return false
		} else {
			scanner.refreshBuffer()
		}
	}
	return scanner.bufpos != scanner.buflast
}

// Returns the next character from the reader.
// Panics if it can't get the next character, because this method is only used if the next character is present.
func (scanner *Scanner) get() byte {
	if scanner.has() {
		return scanner.buffer[scanner.bufpos]
	}
	// Impossible situation.
	panic(errors.New("can not get the next byte"))
}

// Moves to the next character.
// Calls the get method without checking the existence of the next character,
// so it must only be called if the next character exists.
func (scanner *Scanner) step() {
	if scanner.switchLine {
		scanner.refreshLine()
		scanner.switchLine = false
	}
	var symbol = scanner.get()
	if symbol == '\n' {
		scanner.switchLine = true
	} else {
		scanner.lineStr = append(scanner.lineStr, symbol)
	}
	scanner.bufpos++
	scanner.position++
}

// Returns the next token read from the reader.
// If all bytes are read from the reader before calling the method, the (EOF, "") is always returned.
func (scanner *Scanner) Next() (TokenType, string) {
	// Initialization of the scanner, if it was not initialized earlier.
	if !scanner.init {
		scanner.refreshBuffer()
		scanner.refreshLine()
		scanner.line = 0
		scanner.init = true
	}
	// If all bytes are read from the reader, the Scanner always returns the (EOF, "").
	if !scanner.has() {
		return EOF, ""
	}
	var (
		state     stateType // Contains the current state of finite state machine.
		symbol    byte      // Contains the character currently being processed.
		tokenType TokenType
		buffer    = make([]byte, 0, 100) // Contains the characters that were read.
	)
	for scanner.has() {
		symbol = scanner.get()
		tokenType = tokenTypeMap[state]
		state = matrix[getSymbolType(symbol)][state] // The next state is contained in the matrix.
		// The transition to the start state means the end of the token.
		if state == start {
			// If the comments are omitted, the next token must be returned.
			if scanner.SkipComments && tokenType == Comment {
				return scanner.Next()
			}
			return tokenType, string(buffer)
		}
		buffer = append(buffer, symbol)
		scanner.step()
	}
	// All bytes are read from the reader.
	return tokenTypeMap[state], string(buffer)
}

// Skips all characters until the beginning of the next line.
func (scanner *Scanner) SkipLine() {
	if scanner.switchLine {
		return
	}
	var symbol byte
	for scanner.has() {
		symbol = scanner.get()
		scanner.step()
		if symbol == '\n' {
			return
		}
	}
}

// Returns the position of the character that was last processed by the Scanner
// relative to the beginning of the sequence of bytes being read.
func (scanner *Scanner) Position() int {
	return scanner.position - 1
}

// Returns the number of the line that was last processed by the Scanner.
func (scanner *Scanner) Line() int {
	return scanner.line
}

// Returns the line fragment that was read by the scanner.
func (scanner *Scanner) LineString() string {
	return string(scanner.lineStr)
}

// Returns the position in the line that was last processed by the Scanner.
func (scanner *Scanner) Column() int {
	if scanner.switchLine || !scanner.has() {
		return len(scanner.lineStr)
	}
	return len(scanner.lineStr) - 1
}
