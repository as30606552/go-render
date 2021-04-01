package parser

import (
	"computer_graphics/obj/parser/types"
	"computer_graphics/obj/scanner"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	noErrorMessage = "" // A message about the absence of an error.
	initMatrixSize = 10 // The initial number of states of the state machine.
)

// Returns a string with a message about an impossible token received in the start state,
// formatted with the received token.
func impossibleTokenInStartStateMessage(tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received in the start state - %s", tokenType)
}

// Returns a string with a message about an impossible token received after the element description,
// formatted with the received token and the element being read.
func impossibleTokenAfterDescribingElementMessage(elementType ElementType, tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received after describing a %s - %s", elementType, tokenType)
}

// Returns a string with a message about an unexpected token received after the element description,
// formatted with the received token and the element being read.
func unexpectedTokenAfterDescribingElementMessage(elementType ElementType, tokenType scanner.TokenType) string {
	return fmt.Sprintf("unexpected token received after describing a %s - %s", elementType, tokenType)
}

// Returns a string with a message about an impossible token received when reading some type,
// formatted with the received token and the type being read.
func impossibleTokenMessage(name string, tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received when reading the %s - %s", name, tokenType)
}

// Returns a string with a message about an invalid token received when reading some type,
// formatted with the received token and the type being read.
func invalidTokenMessage(name string, expected, received scanner.TokenType) string {
	return fmt.Sprintf("invalid %s, expected: %s, received: %s", name, expected, received)
}

// Returns a string with a message about parameters not specified in the description,
// formatted with the parameter names, separated by commas, passed to the paramNames.
func parametersNotSpecifiedMessage(paramNames []string) string {
	if len(paramNames) == 1 {
		return fmt.Sprintf("parameter %s is not specified", paramNames[0])
	} else {
		return fmt.Sprintf("parameters %s are not specified", strings.Join(paramNames, ", "))
	}
}

// Returns a string with a message that the parameter in the slice of structures is specified incorrectly,
// formatted with the structure field name, structure name, the expected token and the received token
func invalidParameterFormatMessage(specific, base fmt.Stringer, expected, received scanner.TokenType) string {
	return fmt.Sprintf(
		"invalid format for description of the %s, it must be the same as the first %s, expected: %s, received: %s",
		specific,
		base,
		expected,
		received,
	)
}

// Returns the name of the delimiter between the two types.
func delimiterBetween(predecessor, successor string) string {
	return fmt.Sprintf("delimiter between %s and %s", predecessor, successor)
}

// Returns the name of the token following the type.
func tokenAfter(name string) string {
	return fmt.Sprintf("token after %s", name)
}

// The action performed with the token when the elementParser goes to the next state.
type action func(token string, element reflect.Value) error

// Contains complete information about the finite state machine that implements the elementParser.
// The transition to the next state is performed by extracting it from the state table - matrix.
type finiteStateMachine struct {
	element reflect.Value                    // A value containing information about the element being read.
	matrix  [][scanner.TokensCount]stateType // The transition table.
	actions []action                         // An array of actions that are performed when transitioning to a certain state.
	errors  [][scanner.TokensCount]string    // Array of error messages returned when transitioning to the err state.
}

// Clears the element of finiteStateMachine to read the new line.
// Used in the start state.
func (m *finiteStateMachine) clear() { m.element = reflect.New(m.element.Type().Elem()) }

// Implementation of the transition method in the elementParser interface.
func (m *finiteStateMachine) transition(tokenType scanner.TokenType, state stateType) stateType {
	return m.matrix[state][tokenType]
}

// Implementation of the action method in the elementParser interface.
func (m *finiteStateMachine) action(state stateType, token string) error {
	return m.actions[state](token, m.element.Elem())
}

// Implementation of the message method in the elementParser interface.
func (m *finiteStateMachine) message(tokenType scanner.TokenType, state stateType) string {
	return m.errors[state][tokenType]
}

// Implementation of the result method in the elementParser interface.
func (m *finiteStateMachine) result() interface{} { return m.element.Interface() }

// Creates a new finiteStateMachine that reads the specified element and has the specified size of the transition table.
func newMachine(element reflect.Value, size int) *finiteStateMachine {
	return &finiteStateMachine{
		element: element,
		matrix:  make([][scanner.TokensCount]stateType, size),
		actions: make([]action, size),
		errors:  make([][scanner.TokensCount]string, size),
	}
}

// Interface for writing a value to the correct location in reflect.Value.
// Implementations create a system of nested setters that allows you to update the internal fields of the structure.
type setter interface {
	// Writes a token to a value, converting it to the desired type.
	set(token string, value reflect.Value) error
	// Returns the type of the token that can be converted to the required type.
	expected() scanner.TokenType
}

// setter for converting on/off values to bool and writing to reflect.Value.
type boolSetter struct {
	error error // bool parsing error message.
}

// Implementation of the set method in the setter interface.
func (s *boolSetter) set(token string, value reflect.Value) error {
	switch token {
	case "on":
		value.SetBool(true)
	case "off":
		value.SetBool(false)
	default:
		return s.error
	}
	return nil
}

// Implementation of the expected method in the setter interface.
func (s *boolSetter) expected() scanner.TokenType { return scanner.Word }

// Implementation of the changeName method in the setter interface.
func (s *boolSetter) changeName(name string) {
	s.error = fmt.Errorf("the %s parameter must take the values 'on' or 'off'", name)
}

// Creates a new boolSetter by the parameter name.
func newBoolSetter(name string) *boolSetter {
	var s = &boolSetter{}
	s.changeName(name)
	return s
}

// setter for converting v/u values to types.DirectionType and writing to reflect.Value.
type directionTypeSetter struct {
	error error // types.DirectionType parsing error message.
}

// Implementation of the set method in the setter interface.
func (s *directionTypeSetter) set(token string, value reflect.Value) error {
	switch token {
	case "v":
		value.SetUint(uint64(types.V))
	case "u":
		value.SetUint(uint64(types.U))
	default:
		return s.error
	}
	return nil
}

// Implementation of the expected method in the setter interface.
func (s *directionTypeSetter) expected() scanner.TokenType { return scanner.Word }

// Implementation of the changeName method in the setter interface.
func (s *directionTypeSetter) changeName(name string) {
	s.error = fmt.Errorf("the %s parameter must take the values 'v' or 'u'", name)
}

// Creates a new directionTypeSetter by the parameter name.
func newDirectionTypeSetter(name string) *directionTypeSetter {
	var s = &directionTypeSetter{}
	s.changeName(name)
	return s
}

// setter for converting integer values to int and writing to reflect.Value.
type intSetter struct {
	error error // int parsing error message.
}

// Implementation of the set method in the setter interface.
func (s *intSetter) set(token string, value reflect.Value) error {
	var val, err = strconv.ParseInt(token, 10, 64)
	if err != nil {
		return s.error
	}
	value.SetInt(val)
	return nil
}

// Implementation of the expected method in the setter interface.
func (s *intSetter) expected() scanner.TokenType { return scanner.Integer }

// Creates a new intSetter by the parameter name.
func newIntSetter(name string) *intSetter {
	return &intSetter{fmt.Errorf("failed to convert the token to an integer when reading %s", name)}
}

// setter for converting float values to float64 and writing to reflect.Value.
type floatSetter struct {
	error error // float64 parsing error message.
}

// Implementation of the set method in the setter interface.
func (s *floatSetter) set(token string, value reflect.Value) error {
	var val, err = strconv.ParseFloat(token, 64)
	if err != nil {
		return s.error
	}
	value.SetFloat(val)
	return nil
}

// Implementation of the expected method in the setter interface.
func (s *floatSetter) expected() scanner.TokenType { return scanner.Float }

// Creates a new floatSetter by the parameter name.
func newFloatSetter(name string) *floatSetter {
	return &floatSetter{fmt.Errorf("failed to convert the token to a float when reading %s", name)}
}

// setter for writing string values to reflect.Value.
type stringSetter struct{}

// Implementation of the set method in the setter interface.
func (s *stringSetter) set(token string, value reflect.Value) error {
	value.SetString(token)
	return nil
}

// Implementation of the expected method in the setter interface.
func (s *stringSetter) expected() scanner.TokenType { return scanner.Float }

// Creates a new stringSetter.
func newStringSetter() *stringSetter { return &stringSetter{} }

// Wrapper for writing a value to the desired field of the structure.
// Retrieves the desired field and delegates writing to it to the nested setter.
type structSetter struct {
	fieldNumber int // The number of the field to write the value to.
	setter          // Delegate.
}

// Implementation of the set method in the setter interface.
func (s *structSetter) set(token string, value reflect.Value) error {
	return s.setter.set(token, value.Field(s.fieldNumber))
}

// Creates a new structSetter.
func newStructSetter(fieldNumber int, setter setter) *structSetter {
	return &structSetter{
		fieldNumber: fieldNumber,
		setter:      setter,
	}
}

// Wrapper for writing a value to the last element of the slice.
// Retrieves the desired element and delegates writing to it to the nested setter.
type sliceSetter struct {
	setter // Delegate.
}

// Implementation of the set method in the setter interface.
func (s *sliceSetter) set(token string, value reflect.Value) error {
	return s.setter.set(token, value.Index(value.Len()-1))
}

// Creates a new sliceSetter.
func newSliceSetter(setter setter) *sliceSetter {
	return &sliceSetter{
		setter: setter,
	}
}

// Wrapper for creating a new element in the slice.
// Creates the desired element and delegates writing to it to the nested setter.
type sliceAppender struct {
	setter // Delegate.
}

// Implementation of the set method in the setter interface.
func (s *sliceAppender) set(token string, value reflect.Value) error {
	value.Set(reflect.Append(value, reflect.New(value.Type().Elem()).Elem()))
	return s.setter.set(token, value)
}

// Creates a new sliceAppender.
func newSliceAppender(setter setter) *sliceAppender {
	return &sliceAppender{
		setter: setter,
	}
}

// Interface for building states for a single value.
type parameter interface {
	// Creates states in builder for reading a single value.
	update(b *builder)
	// Implements the fmt.Stringer interface for correctly displaying error information.
	String() string
}

// The type of parameter name inherited by the parameters for implementing the fmt.Stringer interface.
type parameterName string

// Implementation of the String method in the fmt.Stringer interface.
func (n parameterName) String() string { return string(n) }

// A parameter for a type that requires only one state of the finite state machine.
type baseParameter struct {
	parameterName        // The name of the baseParameter.
	setter        setter // A setter that writes the value in the way required for the baseParameter.
}

// Updating a single state of the finite state machine.
// state - the state to go to if the expected token is received.
// unread - names of parameters to be read after.
func (p *baseParameter) baseUpdate(b *rowBuilder, state stateType, unread []string) {
	var (
		expected = p.setter.expected()
		act      = p.setter.set
	)
	if expected == scanner.Word {
		b.onWord(state, act)
	} else {
		b.onWordError(invalidTokenMessage(p.String(), expected, scanner.Word))
	}
	// The scanner.Integer token is valid if the parameter expects a scanner.Float token.
	if expected == scanner.Integer || expected == scanner.Float {
		b.onInteger(state, act)
	} else {
		b.onIntegerError(invalidTokenMessage(p.String(), expected, scanner.Integer))
	}
	if expected == scanner.Float {
		// You do not need to specify the action that will be performed when transitioning through the scanner.Float,
		// because the necessary action will be recorded when processing the transition by scanner.Integer.
		b.onFloat(state, nil)
	} else {
		b.onFloatError(invalidTokenMessage(p.String(), expected, scanner.Float))
	}
	// Transition to an err state if this or other parameters are to be read.
	if len(unread) > 0 {
		b.onEndError(parametersNotSpecifiedMessage(unread))
	} else {
		b.onEnd()
	}
}

// Implementation of the update method in the parameter interface.
func (p *baseParameter) update(b *builder) {
	p.baseUpdate(b.nextParameterRow(p.String(), p.setter.expected()), b.nextState(), b.getUnread())
}

// Creates a new baseParameter.
func newBaseParameter(name string, setter setter) *baseParameter {
	return &baseParameter{
		parameterName: parameterName(name),
		setter:        setter,
	}
}

// A parameter that generates states for reading the fields of a nested structure.
type structParameter struct {
	parameterName                   // The name of the structParameter.
	delimiter     scanner.TokenType // The type of delimiter between the fields of the nested structure.
	params        []*baseParameter  // Parameters of the structure fields.
	requiredCount int               // The number of required parameters.
}

// Changes the parameterName of the structParameter to display error messages correctly.
func (p *structParameter) changeName(name string) {
	p.parameterName = parameterName(name)
}

// Returns the name of the structure field parameter by its number.
func (p *structParameter) name(num int) string {
	return fmt.Sprintf("%s of the %s", p.params[num], p)
}

// Returns the names of all required fields in the nested structure.
func (p *structParameter) requiredNames() []string {
	var res = make([]string, p.requiredCount)
	for i := 0; i < p.requiredCount; i++ {
		res[i] = p.name(i)
	}
	return res
}

// Returns the names of all parameters up to the specified one.
func (p *structParameter) allNames(to int) []string {
	var res = make([]string, to)
	for i := 0; i < to; i++ {
		res[i] = p.name(i)
	}
	return res
}

// Complements the builder with states for reading the required fields of the nested structure.
func (p *structParameter) updateRequired(b *builder, unread []string) {
	var (
		param *baseParameter // The current parameter being processed.
		name  string         // The name of the current parameter.
	)
	for i := 0; i < p.requiredCount; i++ {
		param = p.params[i]
		name = unread[i]
		param.baseUpdate(b.nextParameterRow(name, param.setter.expected()), b.nextState(), unread[i:])
		if i != p.requiredCount-1 {
			switch p.delimiter {
			case scanner.Space:
				b.waitSpace(delimiterBetween(name, unread[i+1]), unread[i+1:])
			case scanner.Slash:
				b.waitSlash(delimiterBetween(name, unread[i+1]), unread[i+1:])
			}
		}
	}
}

// Implementation of the update method in the parameter interface.
func (p *structParameter) update(b *builder) {
	// If the method is called, then the structure parameter has a delimiter space
	// and reading required parameters means reading all parameters.
	p.updateRequired(b, append(p.requiredNames(), b.getUnread()[1:]...))
}

// Creates a new structParameter.
func newStructParameter(name string, delimiter scanner.TokenType, size int) *structParameter {
	return &structParameter{
		parameterName: parameterName(name),
		delimiter:     delimiter,
		params:        make([]*baseParameter, 0, size),
		requiredCount: 0,
	}
}

// Contains the common logic of the slice parameters.
type sliceParameter struct {
	parameterName     // The name of the sliceParameter.
	min           int // Minimum number of slice elements.
}

// Implementation of the String method in the parameter interface.
func (p *sliceParameter) String() string { return fmt.Sprintf("%s parameters", p.parameterName) }

// Returns the name of the specified slice element.
func (p *sliceParameter) name(num int) string {
	return fmt.Sprintf("%s number %d", p.parameterName, num+1)
}

// Returns the names of all required elements of the slice.
func (p *sliceParameter) names() []string {
	var res = make([]string, p.min)
	for i := 0; i < p.min; i++ {
		res[i] = p.name(i)
	}
	return res
}

// A parameter that generates states for reading the slice of types that requires only one state of the finite state machine.
type baseSliceParameter struct {
	sliceParameter                // Basic structure.
	param          *baseParameter // A baseParameter of a single slice element.
}

// Implementation of the update method in the parameter interface.
// The implementation assumes that the baseSliceParameter is the last parameter of the builder.
func (p *baseSliceParameter) update(b *builder) {
	var (
		name     string                      // The name of the current parameter.
		names    = p.names()                 // The names of all required elements of the slice.
		expected = p.param.setter.expected() // The type of the slice element.
	)
	for i := 0; i < p.min; i++ {
		name = p.name(i)
		p.param.baseUpdate(b.nextParameterRow(name, expected), b.nextState(), names[i:])
		if i != p.min-1 {
			b.waitSpace(delimiterBetween(name, p.name(i+1)), names[i+1:])
		} else {
			b.waitSpace(tokenAfter(name), []string{})
		}
	}
	// Create additional states for reading an arbitrary number of slice elements after the required ones.
	var loopState = b.nextState()
	name = fmt.Sprintf("additional %s", p.parameterName)
	p.param.baseUpdate(b.nextParameterRow(name, expected), b.nextState(), []string{})
	name = tokenAfter(name)
	b.nextDelimiterRow(name).
		onSlashError(invalidTokenMessage(name, scanner.Space, scanner.Slash)).
		onSpace(loopState).
		onEnd()
}

// Creates a new baseSliceParameter.
func newBaseSliceParameter(name string, min int, param *baseParameter) *baseSliceParameter {
	return &baseSliceParameter{
		sliceParameter: sliceParameter{
			parameterName: parameterName(name),
			min:           min,
		},
		param: param,
	}
}

// A parameter that generates states for reading the slice of structures.
type structSliceParameter struct {
	sliceParameter                  // Basic structure.
	param          *structParameter // A structParameter of a single slice element.
}

// Recursively called function to account for different combinations of optional parameters of a nested structure.
// Splits the processing of the paramNumber field baseParameter into three cases:
// 1. When the baseParameter is present in the string (the token type expected by the baseParameter is received).
// 2. When the baseParameter is missing in the string, but other parameters are present (a slash is received).
// 3. When the baseParameter and all next parameters are missing (a space is received).
func (p *structSliceParameter) updateRecursive(
	update func(b *builder, unread []string), // A function that updates the builder based on the previous parameters.
	b *builder, // Updated builder.
	paramNumber int, // The number of the current parameter.
	lastSlash bool, // true if the last processed token was a slash.
) {
	var (
		sliceNames = p.names() // Names of all required parameters of the slice.
		name       string      // The name of the current parameter to be processed.
	)
	p.param.changeName(p.name(0))
	// Checking for other unprocessed parameters.
	if paramNumber < len(p.param.params) {
		// There are other unprocessed parameters, need to split the situation into three cases
		name = p.param.name(paramNumber)
		var delimiterRow *rowBuilder // The search state of the delimiter between the previous parameter and the current one.
		if !lastSlash {
			// If a non-slash was received last, then it is necessary to read the slash before the current parameter.
			delimiterRow = b.nextDelimiterRow(delimiterBetween(p.param.name(paramNumber-1), name)).
				onSlash(b.nextState())
			if p.min > 1 {
				// If the minimum number of elements of the slice is greater than one,
				// then the end of the line means an error (the remaining elements of the slice are not read).
				delimiterRow.onEndError(parametersNotSpecifiedMessage(sliceNames[1:]))
			} else {
				// If the minimum number of slice elements is one,
				// then the end of the line means that only this single slice element is specified in it.
				delimiterRow.onEnd()
			}
		}
		var (
			// Function for reading parameters if the current one is not specified.
			noParamUpdate = func(b *builder, unread []string) {
				// Building the states required for the fart parameters.
				update(b, unread)
				// Reading the slash before the current parameter.
				var (
					param = p.param.params[paramNumber]
					rb    *rowBuilder
				)
				if lastSlash {
					var (
						extParamMessage = fmt.Sprintf(
							"the %s is specified for the %s, but is not specified for the first %s",
							p.param.params[paramNumber-1],
							p.param,
							p.parameterName,
						)
						onFloatError string
					)
					if param.setter.expected() == scanner.Float {
						onFloatError = extParamMessage
					} else {
						onFloatError = invalidParameterFormatMessage(p.param, p.parameterName, scanner.Slash, scanner.Float)
					}
					rb = b.nextEmptyRow().
						onWordError(invalidParameterFormatMessage(p.param, p.parameterName, scanner.Slash, scanner.Word)).
						onIntegerError(extParamMessage).
						onFloatError(onFloatError).
						onUnknownError(invalidParameterFormatMessage(p.param, p.parameterName, scanner.Slash, scanner.Unknown)).
						onCommentError(impossibleTokenMessage(p.param.name(paramNumber-1), scanner.Comment))
				} else {
					rb = b.nextDelimiterRow(tokenAfter(p.param.name(paramNumber - 1)))
				}
				rb.onSlash(b.nextState()).
					onSpaceError(fmt.Sprintf(
						"the %s is not specified for the %s, but is specified for the first %s",
						param,
						p.param,
						p.parameterName,
					)).
					onEndError(parametersNotSpecifiedMessage(unread[paramNumber:]))
			}
			// Function for reading parameters, if the current one is specified.
			hasParamUpdate = func(b *builder, unread []string) {
				// Building the states required for the previous parameters and slash after them.
				noParamUpdate(b, unread)
				// Reading the current parameter.
				var (
					name     = p.param.name(paramNumber)
					param    = p.param.params[paramNumber]
					expected = param.setter.expected()
					paramRow = b.nextEmptyRow().
							onSpaceError(invalidTokenMessage(name, expected, scanner.Space)).
							onUnknownError(invalidTokenMessage(name, expected, scanner.Unknown)).
							onCommentError(impossibleTokenMessage(name, scanner.Comment))
				)
				if paramNumber != len(p.param.params)-1 {
					paramRow.onSlashError(fmt.Sprintf(
						"the %s is not specified for the %s, but is specified for the first %s",
						param,
						p.param,
						p.parameterName,
					))
				} else {
					paramRow.onSlashError(invalidTokenMessage(name, expected, scanner.Slash))
				}
				param.baseUpdate(paramRow, b.nextState(), unread[paramNumber:])
			}
			// The current parameter being processed.
			param = p.param.params[paramNumber]
			// The token type expected by the current parameter.
			expected = param.setter.expected()
			// The read state of the current parameter.
			paramRow = b.nextEmptyRow().
					onSpaceError(invalidTokenMessage(name, expected, scanner.Space)).
					onUnknownError(invalidTokenMessage(name, expected, scanner.Unknown)).
					onCommentError(impossibleTokenMessage(name, scanner.Comment))
		)
		// Recursive processing of the situation when the current parameter is found in a string.
		param.baseUpdate(paramRow, b.nextState(), append([]string{name}, sliceNames[1:]...))
		p.updateRecursive(hasParamUpdate, b, paramNumber+1, false)
		if paramNumber != len(p.param.params)-1 {
			// If the parameter is not the last one, the situation is recursively processed
			// when the current parameter is not found in the string, but others are found.
			paramRow.onSlash(b.nextState())
			p.updateRecursive(noParamUpdate, b, paramNumber+1, true)
		} else {
			// If the parameter is the last one, then the slash token is invalid.
			paramRow.onSlashError(invalidTokenMessage(name, expected, scanner.Slash))
		}
		// If the previous parameter was not found, then a space token in the current state is possible.
		if !lastSlash {
			delimiterRow.onSpace(b.nextState())
		}
	} else {
		// All parameters processed, exit from recursion.
		b.waitSpace(delimiterBetween(sliceNames[0], sliceNames[1]), sliceNames[1:])
	}
	if !lastSlash {
		// If the last token read was not a slash,
		// then it is necessary to read the remaining required elements of the slice
		// corresponding to the processed combination of fields of the structure.
		for i := 1; i < p.min; i++ {
			name = sliceNames[i]
			p.param.changeName(name)
			update(b, append(p.param.allNames(paramNumber), sliceNames[i:]...))
			if i != p.min-1 {
				b.waitSpace(delimiterBetween(name, sliceNames[i+1]), sliceNames[i+1:])
			} else {
				b.waitSpace(tokenAfter(name), []string{})
			}
		}
		// Create additional states for reading an arbitrary number of slice elements after the required ones.
		var loopState = b.nextState()
		name = fmt.Sprintf("additional %s", p.parameterName)
		p.param.changeName(name)
		update(b, p.param.allNames(paramNumber))
		name = tokenAfter(name)
		b.nextDelimiterRow(name).
			onSlashError(invalidTokenMessage(name, scanner.Space, scanner.Slash)).
			onSpace(loopState).
			onEnd()
	}
}

// Implementation of the update method in the parameter interface.
// The implementation assumes that the structSliceParameter is the last parameter of the builder.
func (p *structSliceParameter) update(b *builder) {
	// Starts the recursive builder update process.
	p.param.changeName(p.name(0))
	p.param.updateRequired(b, append(p.param.requiredNames(), p.names()[1:]...))
	p.updateRecursive(
		func(b *builder, unread []string) {
			p.param.updateRequired(b, unread)
		},
		b,
		p.param.requiredCount,
		false,
	)
}

// Creates a new structSliceParameter.
func newStructSliceParameter(name string, min int, param *structParameter) *structSliceParameter {
	return &structSliceParameter{
		sliceParameter: sliceParameter{
			parameterName: parameterName(name),
			min:           min,
		},
		param: param,
	}
}

// Stores the state and the action to be performed when switching to this state.
type stateAction struct {
	state  stateType
	action action
}

// Stores information about transitions from a single state.
type rowBuilder struct {
	stateActionRow [scanner.TokensCount]stateAction // A row of states and actions.
	errorsRow      [scanner.TokensCount]string      // A row of error messages.
}

// Updates the row of states by transitioning through the token without an error.
func (b *rowBuilder) onToken(t scanner.TokenType, s stateType, a action) *rowBuilder {
	b.stateActionRow[t] = stateAction{
		state:  s,
		action: a,
	}
	b.errorsRow[t] = noErrorMessage
	return b
}

// Updates the row of states by transitioning through the scanner.Word token without an error.
func (b *rowBuilder) onWord(s stateType, a action) *rowBuilder { return b.onToken(scanner.Word, s, a) }

// Updates the row of states by transitioning through the scanner.Integer token without an error.
func (b *rowBuilder) onInteger(s stateType, a action) *rowBuilder {
	return b.onToken(scanner.Integer, s, a)
}

// Updates the row of states by transitioning through the scanner.Float token without an error.
func (b *rowBuilder) onFloat(s stateType, a action) *rowBuilder {
	return b.onToken(scanner.Float, s, a)
}

// Updates the row of states by transitioning through the scanner.Slash token without an error.
func (b *rowBuilder) onSlash(s stateType) *rowBuilder { return b.onToken(scanner.Slash, s, nil) }

// Updates the row of states by transitioning through the scanner.Space token without an error.
func (b *rowBuilder) onSpace(s stateType) *rowBuilder { return b.onToken(scanner.Space, s, nil) }

// Updates the row of states by transitioning through the scanner.EOL and scanner.EOF tokens without an error.
func (b *rowBuilder) onEnd() *rowBuilder {
	return b.onToken(scanner.EOL, start, nil).onToken(scanner.EOF, start, nil)
}

// Updates the row of states by transitioning through the token to the error state.
func (b *rowBuilder) onTokenError(t scanner.TokenType, message string) *rowBuilder {
	b.stateActionRow[t] = stateAction{
		state:  err,
		action: nil,
	}
	b.errorsRow[t] = message
	return b
}

// Updates the row of states by transitioning through the scanner.Word token to the error state.
func (b *rowBuilder) onWordError(message string) *rowBuilder {
	return b.onTokenError(scanner.Word, message)
}

// Updates the row of states by transitioning through the scanner.Integer token to the error state.
func (b *rowBuilder) onIntegerError(message string) *rowBuilder {
	return b.onTokenError(scanner.Integer, message)
}

// Updates the row of states by transitioning through the scanner.Float token to the error state.
func (b *rowBuilder) onFloatError(message string) *rowBuilder {
	return b.onTokenError(scanner.Float, message)
}

// Updates the row of states by transitioning through the scanner.Slash token to the error state.
func (b *rowBuilder) onSlashError(message string) *rowBuilder {
	return b.onTokenError(scanner.Slash, message)
}

// Updates the row of states by transitioning through the scanner.Space token to the error state.
func (b *rowBuilder) onSpaceError(message string) *rowBuilder {
	return b.onTokenError(scanner.Space, message)
}

// Updates the row of states by transitioning through the scanner.EOL and scanner.EOL tokens to the error state.
func (b *rowBuilder) onEndError(message string) *rowBuilder {
	return b.onTokenError(scanner.EOL, message).onTokenError(scanner.EOF, message)
}

// Updates the row of states by transitioning through the scanner.Unknown token to the error state.
func (b *rowBuilder) onUnknownError(message string) *rowBuilder {
	return b.onTokenError(scanner.Unknown, message)
}

// Updates the row of states by transitioning through the scanner.Comment token to the error state.
func (b *rowBuilder) onCommentError(message string) *rowBuilder {
	return b.onTokenError(scanner.Comment, message)
}

// Creates a new rowBuilder.
func newRowBuilder() *rowBuilder { return &rowBuilder{} }

// Contains information about the element to be read.
// Builds a finiteStateMachine based on it, which reads this element.
type builder struct {
	value        reflect.Value // The value returned by the finiteStateMachine.
	valueType    ElementType   // The type of the element to be read.
	params       []parameter   // Element parameters that update the builder with states for reading a single field of the structure.
	paramNames   []string      // Names of required parameters.
	position     int           // The current parameter being processed.
	builders     []*rowBuilder // Preliminary information about the rows of the finiteStateMachine.
	needFinalize bool          // true if need to add end-of-line processing states to the finiteStateMachine.
}

// Creates a single parameter that reads on/off values.
func (b *builder) createBoolParameter() {
	b.params = make([]parameter, 1)
	b.paramNames = make([]string, 1)
	var name = fmt.Sprintf("%s parameter", b.valueType)
	b.params[0] = newBaseParameter(name, newBoolSetter(name))
	b.paramNames[0] = name
}

// Reads the name tag (parameter name).
func readName(field *reflect.StructField) string {
	if name, ok := field.Tag.Lookup("name"); ok {
		return name
	} else {
		return field.Name
	}
}

// Reads the optional tag (whether the parameter is optional).
func readOptional(tags reflect.StructTag, isFirst bool) bool {
	if optional, ok := tags.Lookup("optional"); ok {
		if isFirst {
			panic("the first field of the structure cannot be optional")
		}
		if res, err := strconv.ParseBool(optional); err == nil {
			return res
		} else {
			panic("the optional tag must take the values 'true' or 'false'")
		}
	} else {
		return false
	}
}

// Reads the delimiter tag (delimiter between parameters of the nested structure).
func readDelimiter(tags reflect.StructTag) scanner.TokenType {
	if delimiter, ok := tags.Lookup("delimiter"); ok {
		switch delimiter {
		case "slash":
			return scanner.Slash
		case "space":
			return scanner.Space
		default:
			panic("the delimiter tag must take the values 'space' or 'slash'")
		}
	} else {
		panic("the []struct field must have the delimiter tag specified")
	}
}

// Reads the min tag (minimum number of slice elements).
func readMin(tags reflect.StructTag) int {
	if min, ok := tags.Lookup("min"); ok {
		if res, err := strconv.ParseInt(min, 10, 8); err == nil {
			if res < 1 {
				panic("the min tag cannot accept values less than one")
			} else {
				return int(res)
			}
		} else {
			panic("error reading the min tag")
		}
	} else {
		panic("the slice field must have the min tag specified")
	}
}

// Panics if the optional tag is present among the tags.
func requireNoOptional(tags reflect.StructTag, typeName string) {
	if _, ok := tags.Lookup("optional"); ok {
		panic(fmt.Sprintf("the optional tag cannot be set for a %s field", typeName))
	}
}

// Panics if the delimiter tag is present among the tags.
func requireNoDelimiter(tags reflect.StructTag, typeName string) {
	if _, ok := tags.Lookup("delimiter"); ok {
		panic(fmt.Sprintf("the delimiter tag cannot be set for a %s field", typeName))
	}
}

// Panics if the min tag is present among the tags.
func requireNoMin(tags reflect.StructTag, typeName string) {
	if _, ok := tags.Lookup("min"); ok {
		panic(fmt.Sprintf("the min tag cannot be set for a %s field", typeName))
	}
}

// Panics if wasOptional is true.
// It is necessary for all optional fields to be the last in the structure.
func requireWasNotOptional(wasOptional bool) {
	if wasOptional {
		panic("an optional field cannot be followed by a required field")
	}
}

// Creates a nested structure parameter.
func createNestedStructParameter(
	name string, // Name of the parameter.
	delimiter scanner.TokenType, // Fields delimiter of the nested structure.
	t reflect.Type, // Type of the nested structure.
	// A function that generates the setter that the structure field should be wrapped in.
	wrapper func(fieldNumber int, setter setter) setter,
) *structParameter {
	var (
		res         = newStructParameter(name, scanner.Space, t.NumField())
		field       reflect.StructField
		tags        reflect.StructTag
		nestedName  string
		optional    bool
		hasOptional = false
		param       *baseParameter
	)
	// Creating parameters for each field of the structure.
	for i := 0; i < t.NumField(); i++ {
		field = t.Field(i)
		nestedName = readName(&field)
		tags = field.Tag
		switch delimiter {
		case scanner.Space:
			requireNoOptional(tags, "struct with space delimiter")
		case scanner.Slash:
			optional = readOptional(tags, i == 0)
			if !optional {
				requireWasNotOptional(hasOptional)
			}
			hasOptional = optional
		}
		switch field.Type.Kind() {
		case reflect.Int:
			requireNoDelimiter(tags, "int")
			requireNoMin(tags, "int")
			param = newBaseParameter(nestedName, wrapper(i, newIntSetter(nestedName)))
		case reflect.Float64:
			requireNoDelimiter(tags, "float64")
			requireNoMin(tags, "float64")
			param = newBaseParameter(nestedName, wrapper(i, newFloatSetter(nestedName)))
		default:
			panic(fmt.Sprintf("unsupported nested struct field type: %s", field.Type.Kind()))
		}
		res.params = append(res.params, param)
		if !hasOptional {
			res.requiredCount++
		}
	}
	return res
}

// Updates the builder with a set of structure field parameters.
func (b *builder) createStructParameters() {
	var t = b.value.Elem().Type()
	b.params = make([]parameter, 0, t.NumField())
	b.paramNames = make([]string, 0, t.NumField())
	var (
		field       reflect.StructField
		tags        reflect.StructTag
		name        string
		typeName    string
		optional    bool
		hasOptional = false
		min         int
		param       parameter
	)
	if t.NumField() < 1 {
		panic("the parser cannot be built on a structure without fields")
	}
	// Creating parameters for each field of the structure.
	for i := 0; i < t.NumField(); i++ {
		field = t.Field(i)
		name = readName(&field)
		tags = field.Tag
		switch field.Type.Kind() {
		case reflect.Uint8:
			if typeName != "DirectionType" {
				panic("the field with the base type uint8 must have the type DirectionType")
			}
			if i != 0 {
				panic("the DirectionType field must be the first in the structure")
			}
			requireNoOptional(tags, typeName)
			requireNoDelimiter(tags, typeName)
			requireNoMin(tags, typeName)
			param = newBaseParameter(name, newStructSetter(i, newDirectionTypeSetter(name)))
		case reflect.Int:
			typeName = "int"
			requireNoDelimiter(tags, typeName)
			requireNoMin(tags, typeName)
			optional = readOptional(tags, i == 0)
			if !optional {
				requireWasNotOptional(hasOptional)
			}
			hasOptional = optional
			param = newBaseParameter(name, newStructSetter(i, newIntSetter(name)))
		case reflect.Float64:
			typeName = "float64"
			requireNoDelimiter(tags, typeName)
			requireNoMin(tags, typeName)
			optional = readOptional(tags, i == 0)
			if !optional {
				requireWasNotOptional(hasOptional)
			}
			hasOptional = optional
			param = newBaseParameter(name, newStructSetter(i, newFloatSetter(name)))
		case reflect.String:
			typeName = "string"
			requireNoOptional(tags, typeName)
			requireNoDelimiter(tags, typeName)
			requireNoMin(tags, typeName)
			requireWasNotOptional(hasOptional)
			param = newBaseParameter(name, newStringSetter())
		case reflect.Struct:
			typeName = "nested struct"
			requireNoOptional(tags, typeName)
			requireNoDelimiter(tags, typeName)
			requireNoMin(tags, typeName)
			requireWasNotOptional(hasOptional)
			param = createNestedStructParameter(
				name,
				scanner.Space,
				field.Type,
				func(fieldNumber int, setter setter) setter {
					return newStructSetter(i, newStructSetter(fieldNumber, setter))
				},
			)
		case reflect.Slice:
			if i != t.NumField()-1 {
				panic("the slice must be the last field of the structure")
			}
			// The parameters of the slices themselves fill in the last states of the finite state machine.
			b.needFinalize = false
			requireNoOptional(tags, "slice")
			requireWasNotOptional(hasOptional)
			min = readMin(tags)
			switch field.Type.Elem().Kind() {
			case reflect.Int:
				requireNoDelimiter(tags, "[]int")
				param = newBaseSliceParameter(
					name,
					min,
					newBaseParameter("", newSliceAppender(newSliceSetter(newIntSetter("")))),
				)
			case reflect.Float64:
				requireNoDelimiter(tags, "[]float64")
				param = newBaseSliceParameter(
					name,
					min,
					newBaseParameter("", newSliceAppender(newSliceSetter(newFloatSetter("")))),
				)
			case reflect.String:
				requireNoDelimiter(tags, "[]string")
				param = newBaseSliceParameter(
					name,
					min,
					newBaseParameter(name, newSliceAppender(newSliceSetter(newStringSetter()))),
				)
			case reflect.Struct:
				param = newStructSliceParameter(name, min, createNestedStructParameter(
					name,
					readDelimiter(tags),
					field.Type.Elem(),
					func(fieldNumber int, setter setter) setter {
						if fieldNumber == 0 {
							return newStructSetter(i, newSliceAppender(newSliceSetter(newStructSetter(fieldNumber, setter))))
						} else {
							return newStructSetter(i, newSliceSetter(newStructSetter(fieldNumber, setter)))
						}
					},
				))
			default:
				panic(fmt.Sprintf("unsupported struct field type: %s", field.Type.Kind()))
			}
		default:
			panic(fmt.Sprintf("unsupported struct field type: %s", field.Type.Kind()))
		}
		b.params = append(b.params, param)
		if !hasOptional {
			b.paramNames = append(b.paramNames, param.String())
		}
	}
}

// The next free state of the finite state machine.
func (b *builder) nextState() stateType { return stateType(len(b.builders)) }

// Creates a new empty finite state machine row.
func (b *builder) nextEmptyRow() *rowBuilder {
	var rb = newRowBuilder()
	b.builders = append(b.builders, rb)
	return rb
}

// Creates a new row of the finite state machine for reading the parameter.
// By default, creates transitions through delimiter tokens to the err state.
//
func (b *builder) nextParameterRow(name string, expected scanner.TokenType) *rowBuilder {
	return b.nextEmptyRow().
		onSlashError(invalidTokenMessage(name, expected, scanner.Slash)).
		onSpaceError(impossibleTokenMessage(name, scanner.Space)).
		onUnknownError(invalidTokenMessage(name, expected, scanner.Unknown)).
		onCommentError(impossibleTokenMessage(name, scanner.Comment))
}

// Creates a new line of the state machine to read the delimiter between the parameters.
// By default, it creates transitions for all tokens except for delimiter tokens to the err state.
func (b *builder) nextDelimiterRow(name string) *rowBuilder {
	return b.nextEmptyRow().
		onWordError(impossibleTokenMessage(name, scanner.Word)).
		onIntegerError(impossibleTokenMessage(name, scanner.Integer)).
		onFloatError(impossibleTokenMessage(name, scanner.Float)).
		onUnknownError(impossibleTokenMessage(name, scanner.Unknown)).
		onCommentError(impossibleTokenMessage(name, scanner.Comment))
}

// Returns the names of the required parameters that have not yet been processed.
func (b *builder) getUnread() []string {
	if b.position < len(b.paramNames) {
		return b.paramNames[b.position:]
	}
	return []string{}
}

// Creates and fills in the read state of the space.
func (b *builder) waitSpace(name string, unread []string) {
	var rb = b.nextDelimiterRow(name).
		onSlashError(invalidTokenMessage(name, scanner.Space, scanner.Slash)).
		onSpace(b.nextState())
	if len(unread) > 0 {
		rb.onEndError(parametersNotSpecifiedMessage(unread))
	} else {
		rb.onEnd()
	}
}

// Creates and fills in the read state of the slash.
func (b *builder) waitSlash(name string, unread []string) {
	var rb = b.nextDelimiterRow(name).
		onSpaceError(invalidTokenMessage(name, scanner.Slash, scanner.Space)).
		onSlash(b.nextState())
	if len(unread) > 0 {
		rb.onEndError(parametersNotSpecifiedMessage(unread))
	} else {
		rb.onEnd()
	}
}

// Initializes the builder by processing the start state and err state.
func (b *builder) initialize() {
	b.nextEmptyRow().
		onWordError(impossibleTokenInStartStateMessage(scanner.Word)).
		onIntegerError(impossibleTokenInStartStateMessage(scanner.Integer)).
		onFloatError(impossibleTokenInStartStateMessage(scanner.Integer)).
		onSlashError(impossibleTokenInStartStateMessage(scanner.Slash)).
		onSpace(first).
		onEndError(fmt.Sprintf("all parameters of the %s are not specified", b.valueType)).
		onUnknownError(impossibleTokenInStartStateMessage(scanner.Unknown)).
		onCommentError(impossibleTokenInStartStateMessage(scanner.Comment))
	const parserUsedInErrorStateMessage = "parser cannot be used in the error state"
	b.nextEmptyRow().
		onWordError(parserUsedInErrorStateMessage).
		onIntegerError(parserUsedInErrorStateMessage).
		onFloatError(parserUsedInErrorStateMessage).
		onSlashError(parserUsedInErrorStateMessage).
		onSpaceError(parserUsedInErrorStateMessage).
		onEndError(parserUsedInErrorStateMessage).
		onUnknownError(parserUsedInErrorStateMessage).
		onCommentError(parserUsedInErrorStateMessage)
}

// Creates and fills in the states for reading the space after the element description and the end of the line.
func (b *builder) finalize() {
	b.nextEmptyRow().
		onWordError(impossibleTokenAfterDescribingElementMessage(b.valueType, scanner.Word)).
		onIntegerError(impossibleTokenAfterDescribingElementMessage(b.valueType, scanner.Integer)).
		onFloatError(impossibleTokenAfterDescribingElementMessage(b.valueType, scanner.Float)).
		onSlashError(unexpectedTokenAfterDescribingElementMessage(b.valueType, scanner.Slash)).
		onSpace(b.nextState()).
		onEnd().
		onUnknownError(impossibleTokenAfterDescribingElementMessage(b.valueType, scanner.Unknown)).
		onCommentError(impossibleTokenAfterDescribingElementMessage(b.valueType, scanner.Unknown))
	b.nextEmptyRow().
		onWordError(unexpectedTokenAfterDescribingElementMessage(b.valueType, scanner.Word)).
		onIntegerError(unexpectedTokenAfterDescribingElementMessage(b.valueType, scanner.Integer)).
		onFloatError(unexpectedTokenAfterDescribingElementMessage(b.valueType, scanner.Float)).
		onSlashError(unexpectedTokenAfterDescribingElementMessage(b.valueType, scanner.Slash)).
		onSpaceError(impossibleTokenAfterDescribingElementMessage(b.valueType, scanner.Space)).
		onEnd().
		onUnknownError(unexpectedTokenAfterDescribingElementMessage(b.valueType, scanner.Unknown)).
		onCommentError(impossibleTokenAfterDescribingElementMessage(b.valueType, scanner.Unknown))
}

// Builds a state machine based on the information contained in builder.builders.
func (b *builder) buildMachine() *finiteStateMachine {
	var (
		m         = newMachine(b.value, len(b.builders))
		matrixRow [scanner.TokensCount]stateType
	)
	m.actions[start] = func(token string, element reflect.Value) error {
		return errors.New("the action method is called in the start state")
	}
	m.actions[err] = func(token string, element reflect.Value) error {
		return errors.New("the action method is called in the err state")
	}
	// When transitioning to the first unreserved state,
	// it is necessary to clear the value of the element that was read during the previous use of the finiteStateMachine.
	m.actions[first] = func(token string, element reflect.Value) error {
		m.clear()
		return nil
	}
	// Filling in each row of the transition matrix based on elements from builder.builders.
	for i, rb := range b.builders {
		for j, sa := range rb.stateActionRow {
			matrixRow[j] = sa.state
			if m.actions[sa.state] == nil {
				m.actions[sa.state] = sa.action
			} else if sa.action != nil {
				// The action performed during the transition to the state must be defined unambiguously.
				panic(fmt.Sprintf("two actions are specified when transitioning to the same state: %d", sa.state))
			}
		}
		m.matrix[i] = matrixRow
		m.errors[i] = rb.errorsRow
	}
	// Filling the remaining states with actions that do nothing.
	for i := 0; i < len(m.actions); i++ {
		if m.actions[i] == nil {
			m.actions[i] = func(token string, element reflect.Value) error { return nil }
		}
	}
	return m
}

// Builds a finite state machine based on the parameters of the structure fields or a parameter of the bool type.
func (b *builder) build() *finiteStateMachine {
	b.initialize()
	var param parameter
	// Building preliminary information about the states of a finite state machine using parameters.
	for b.position, param = range b.params {
		param.update(b)
		if b.position != len(b.params)-1 {
			b.position++
			b.waitSpace(delimiterBetween(param.String(), b.params[b.position].String()), b.getUnread())
		}
	}
	if b.needFinalize {
		b.finalize()
	}
	return b.buildMachine()
}

// Creates a builder that builds a finiteStateMachine for reading values of a specified type.
func newBuilder(elementType ElementType, t reflect.Type) *builder {
	var b = &builder{
		value:        reflect.New(t),
		valueType:    elementType,
		position:     0,
		builders:     make([]*rowBuilder, 0, initMatrixSize),
		needFinalize: true,
	}
	switch t.Kind() {
	case reflect.Bool:
		b.createBoolParameter()
	case reflect.Struct:
		b.createStructParameters()
	default:
		panic("the element to be read must be a structure object or a bool type")
	}
	return b
}

// Creates an elementParser that parses the line based on the structure.
// elementType specifies the type of element to be read.
// element specifies a pointer to the structure, based on the fields of which the string will be read,
// or a pointer to bool, if the element consists of an on/off value.
// element must have the base type *struct or *bool.
// The element to be read must match the element type.
//
// The following limitations apply to the structure:
// 	* The structure fields are extracted from the line in the order in which they are specified in the structure.
// 	* Only public fields will be parsed.
// 	* Structure fields must have one of the following basic types: uint8, int, float64, string, struct, []int, []float64, []string, []struct.
// 	* If a field is of the slice type, it must be the last one in the structure.
// 	* If a field is of the struct or []struct type, its fields must be of the base type int or float64.
// 	* If a field is of the uint8 base type, it must be of the type DirectionType.
//
// To specify additional information about the fields, use the following tags:
//
// 	name
//
//	The exact name of the structure field (may contain spaces), used when displaying error information.
//	If this tag is not specified, the name of the structure field will be used.
//
// 	optional
//
//	It can take the values 'true' or 'false'.
//	Used to specify optional fields.
//	Optional fields must be the last fields of the structure.
// 	All fields in the structure cannot be optional.
// 	This tag can only be specified for fields of type int and float64.
// 	If the tag value is not specified, the field is processed as required (like optional="false").
//	These rules also apply to nested structures (fields of the struct type).
//
// 	delimiter
//
//	It can take the values 'slash' or 'space'.
// 	This tag must be specified for the []struct type and cannot be specified for other types.
// 	Used for reading structure field delimiters.
// 	If a structure has a 'space' delimiter, it cannot contain optional fields.
// 	If the field is of the type not a slice of structures, but just a structure,
//	it has a space delimiter, which you do not need to specify.
//
// 	min
//
// 	It can only accept integer values that are greater than zero.
// 	This tag must be specified for slices and cannot be specified for other types.
// 	Used to specify the minimum number of slice elements.
func buildParser(elementType ElementType, element interface{}) elementParser {
	var t = reflect.TypeOf(element)
	if t.Kind() != reflect.Ptr {
		panic("the element must be a pointer to a struct or bool")
	}
	return newBuilder(elementType, t.Elem()).build()
}
