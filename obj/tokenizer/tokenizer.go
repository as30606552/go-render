package tokenizer

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// Is one of the possible values that the Tokenizer.Next method returns.
type TokenType uint8

const (
	WORD    TokenType = iota // Can consist of letters, numbers, and underscores. Can not start with a number.
	INT     TokenType = iota // Consists of numbers. Can start with a minus.
	FLOAT   TokenType = iota // Consists of digits with a dot between them. Can start with a minus.
	SLASH   TokenType = iota // '/' character.
	SPACE   TokenType = iota // A sequence of spaces and/or tabs.
	EOL     TokenType = iota // '\n' character.
	EOF     TokenType = iota // Indicates that the end of the sequence of bytes being read has been reached.
	UNKNOWN TokenType = iota // Unknown type of token.
	COMMENT TokenType = iota // Starts with the '#' character and ends with the character before the end of the line.
)

// Is one of the possible states of a finite state machine.
// See https://github.com/as30606552/ComputerGraphicsProject/wiki/Tokenizer.
type stateType uint8

const (
	start      stateType = iota // Initial state.
	skipLine   stateType = iota // Skipping all characters up to the '\n' character.
	foundEol   stateType = iota // '\n' character found.
	foundSpace stateType = iota // Whitespace character found.
	foundSlash stateType = iota // '/' character found.
	foundMinus stateType = iota // '-' character was found at the beginning of the token, and a digit is expected.
	foundDot   stateType = iota // A '.' character is found after an integer, a digit is expected.
	foundInt   stateType = iota // '\n' character found.
	foundFloat stateType = iota // A sequence of characters satisfying the FLOAT token is found, a digit is expected.
	foundWord  stateType = iota // A sequence of characters satisfying the WORD token is found.
	unknown    stateType = iota // An unknown sequence of characters was found.
)

// Is one of the possible character types that can be contained in a sequence of bytes to be read.
type symbolType uint8

const (
	eol    symbolType = iota // '\n'
	space  symbolType = iota // ' ' or '\t'
	hash   symbolType = iota // '#'
	slash  symbolType = iota // '/'
	minus  symbolType = iota // '-'
	dot    symbolType = iota // '.'
	digit  symbolType = iota // '0' - '9'
	letter symbolType = iota // 'a' - 'z' or 'A' - 'Z' or '_'
	other  symbolType = iota // Any other character.
)

// The finite state machine table.
// See https://github.com/as30606552/ComputerGraphicsProject/wiki/Tokenizer.
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

// The size of the buffer in which the Tokenizer stores the read characters.
const bufsize uint8 = 255

// Allows you to sequentially call the Next method to get tokens from a io.Reader that can occur in obj files.
type Tokenizer struct {
	reader io.Reader // The io.Reader from which the tokens will be read.

	buffer  [bufsize]byte // Temporary storage for bytes extracted from the reader but not yet processed.
	bufpos  uint8         // The position of the currently processed byte in the buffer.
	buflast uint8         // The number of bytes contained in the buffer.
	init    bool          // Contains true if there has already been an attempt to extract a byte from the buffer.

	line     int       // The number of the currently processed line.
	column   int       // The position of the currently processed character in the line.
	position int       // The position of the currently processed character relative to the beginning of the byte sequence.
	state    stateType // The current state of the Tokenizer.

	Error        func(error) // The function called in case of an error.
	SkipComments bool        // true if comments should be skipped.
}

// Creates a new Tokenizer that reads from the reader.
// Sets skipping comments by default.
func NewTokenizer(reader io.Reader) *Tokenizer {
	return &Tokenizer{reader: reader, SkipComments: true}
}

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

// Converts the state of the finite state machine from which it moved to the initial state to the type of the read token.
// See https://github.com/as30606552/ComputerGraphicsProject/wiki/Tokenizer.
func getTokenType(state stateType) TokenType {
	// Impossible situation.
	if state == start {
		panic(errors.New("it is not possible to get the token type from the start state"))
	}
	switch state {
	case foundEol:
		return EOL
	case foundSpace:
		return SPACE
	case foundSlash:
		return SLASH
	case foundInt:
		return INT
	case foundFloat:
		return FLOAT
	case foundWord:
		return WORD
	case skipLine:
		return COMMENT
	}
	return UNKNOWN
}

// Delegates the execution of the Error function.
// If the function is not specified, by default it outputs the error data to os.Stderr.
func (tokenizer *Tokenizer) error(err error) {
	if tokenizer.Error != nil {
		tokenizer.Error(err)
		return
	}
	_, err = fmt.Fprintf(os.Stderr, "Error in position %d: %s\n", tokenizer.position, err.Error())
	if err != nil {
		panic(err)
	}
}

// Reads new values to the buffer.
// The number of bytes read is stored in the buflast field.
// The current buffer position is reset to 0.
func (tokenizer Tokenizer) refreshBuffer() {
	var n, err = tokenizer.reader.Read(tokenizer.buffer[:])
	if err != nil && err != io.EOF {
		tokenizer.error(err)
	}
	tokenizer.buflast = uint8(n)
	tokenizer.bufpos = 0
}

// Returns true if there is a next token.
func (tokenizer *Tokenizer) has() bool {
	// Initialization of the buffer, if it was not initialized earlier.
	if !tokenizer.init {
		tokenizer.refreshBuffer()
		tokenizer.init = true
		return tokenizer.bufpos != tokenizer.buflast
	}
	// The buffer is processed to the end.
	// It is necessary to read the new data to the buffer.
	if tokenizer.bufpos == tokenizer.buflast {
		// If the number of elements in the buffer is less than the buffer size,
		// it means that the buffer was not fully filled the previous time when reading it.
		if tokenizer.buflast < bufsize {
			return false
		} else {
			tokenizer.refreshBuffer()
		}
	}
	return tokenizer.bufpos != tokenizer.buflast
}

// Returns the next character from the reader.
// Panics if it can't get the next character, because this method is only used if the next character is present.
func (tokenizer *Tokenizer) get() byte {
	if tokenizer.has() {
		return tokenizer.buffer[tokenizer.bufpos]
	}
	// Impossible situation.
	panic(errors.New("can not get the next byte"))
}

// Moves to the next character.
// Calls the get method without checking the existence of the next character,
// so it must only be called if the next character exists.
func (tokenizer *Tokenizer) step() {
	if tokenizer.get() == '\n' {
		tokenizer.line++
		tokenizer.column = 0
	} else {
		tokenizer.column++
	}
	tokenizer.bufpos++
	tokenizer.position++
}

// Returns the next token read from the reader.
// If all bytes are read from the reader before calling the method, the EOF is always returned.
func (tokenizer *Tokenizer) Next() (TokenType, string) {
	// If all bytes are read from the reader, the Tokenizer always returns the EOF.
	if !tokenizer.has() {
		return EOF, ""
	}
	var state stateType // Contains the previous state of the Tokenizer, which can be used to determine the type of token.
	var symbol byte     // Contains the character currently being processed.
	var tokenType TokenType
	var buffer = make([]byte, 0, 100) // Contains the characters that were read.
	for tokenizer.has() {
		symbol = tokenizer.get()
		state = tokenizer.state
		tokenizer.state = matrix[getSymbolType(symbol)][tokenizer.state] // The next state is contained in the matrix.
		// The transition to the start state means the end of the token.
		if tokenizer.state == start {
			tokenType = getTokenType(state)
			// If the comments are omitted, the next token must be returned.
			if tokenizer.SkipComments && tokenType == COMMENT {
				return tokenizer.Next()
			}
			return tokenType, string(buffer)
		}
		buffer = append(buffer, symbol)
		tokenizer.step()
	}
	// All bytes are read from the reader.
	tokenizer.state = start
	return getTokenType(state), string(buffer)
}

// Skips all characters until the beginning of the next line.
func (tokenizer *Tokenizer) SkipLine() {
	tokenizer.state = skipLine
	_, _ = tokenizer.Next()
}

// Returns the position of the character currently being processed by the Tokenizer
// relative to the beginning of the sequence of bytes being read.
func (tokenizer *Tokenizer) Position() int {
	return tokenizer.position
}

// Returns the number of the line currently being processed by the Tokenizer.
func (tokenizer *Tokenizer) Line() int {
	return tokenizer.line
}

// Returns the position in the line currently being processed by the Tokenizer.
func (tokenizer *Tokenizer) Column() int {
	return tokenizer.column
}
