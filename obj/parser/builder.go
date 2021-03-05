package parser

import (
	"computer_graphics/obj/scanner"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Error message constants.
const (
	noErrorMassage                = ""
	parserUsedInErrorStateMessage = "parser cannot be used in the error state"
)

// Returns a string with a message about an impossible token received in the start state,
// formatted with the received token.
func impossibleTokenInStartStateError(tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received in the start state - %s", tokenType)
}

// Returns a string with a message about an impossible token received when reading the parameter,
// formatted with the received token and the parameter being read
func impossibleTokenWhenReadingParameterError(param parameter, tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received when reading the %s - %s", param, tokenType)
}

// Returns a string with a message about an impossible token received after reading the parameter,
// formatted with the received token and the parameter being read
func impossibleTokenAfterReadingParameterError(param parameter, tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received after reading the %s - %s", param, tokenType)
}

// Returns a string with a message about an impossible token received after the element description,
// formatted with the received token and the element being read
func impossibleTokenAfterDescribingElementError(elementType ElementType, tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received after describing a %s - %s", elementType, tokenType)
}

// Returns a string with a message about an unexpected token received after reading the parameter,
// formatted with the received token and the parameter being read
func unexpectedTokenAfterReadingParameterError(param parameter, tokenType scanner.TokenType) string {
	return fmt.Sprintf("unexpected token received after reading the %s - %s", param, tokenType)
}

// Returns a string with a message about an unexpected token received after the element description,
// formatted with the received token and the element being read
func unexpectedTokenAfterDescribingElementError(elementType ElementType, tokenType scanner.TokenType) string {
	return fmt.Sprintf("unexpected token received after describing a %s - %s", elementType, tokenType)
}

// Returns a string with a message about parameters not specified in the description,
// formatted with the parameter names, separated by commas, passed to the unreadParamsString.
func parametersNotSpecifiedError(unreadParamsString string, unreadParamsCount int) string {
	if unreadParamsCount == 1 {
		return fmt.Sprintf("parameter %s is not specified", unreadParamsString)
	} else {
		return fmt.Sprintf("parameters %s are not specified", unreadParamsString)
	}
}

// Returns a string with a message that all parameters of the read element are not specified in the description,
// formatted by the read element
func allParametersNotSpecifiedError(elementType ElementType) string {
	return fmt.Sprintf("not specified parameters of the %s", elementType)
}

// Returns a string with a message that the parameter is specified incorrectly,
// formatted with the read parameter, the expected token and the received token
func invalidParameterError(param parameter, expected, received scanner.TokenType) string {
	return fmt.Sprintf("invalid %s, excepted: %s, received: %s", param, expected, received)
}

// Returns a string with a message that the delimiter between the parameters is specified incorrectly,
// formatted with these parameters, the expected token and the received token
func invalidDelimiterBetweenParametersError(first, second parameter, expected, received scanner.TokenType) string {
	return fmt.Sprintf(
		"invalid delimiter format between %s and %s, expected: %s, received: %s",
		first,
		second,
		expected,
		received,
	)
}

// The action performed with the token when the elementParser goes to the next state.
type action func(token string)

var (
	// An action that does nothing.
	// It is executed when the elementParser goes to the state of reading the delimiter between the parameters.
	noAction         action = func(token string) {}
	startStateAction action = func(token string) { panic("the action method is called in the err state") }
	lastStateAction  action = func(token string) {
		panic("the action method cannot be called in a state from which transitions are made only to the start state and the err state")
	}
)

// Contains complete information about a single state of the finite state machine.
// Used to fill the finite state machine.
type machineRow struct {
	matrixRow [scanner.TokensCount]stateType // Information about transitions from this state.
	action    action                         // The action that will be performed when switching to this state.
	errorsRow [scanner.TokensCount]string    // Information about errors when transitioning from this state.
}

// Contains complete information about the finite state machine that implements the elementParser.
// The transition to the next state is performed by extracting it from the state matrix.
type finiteStateMachine struct {
	element reflect.Value                    // A structure containing information about the element being read.
	matrix  [][scanner.TokensCount]stateType // The transition matrix.
	actions []action                         // An array of actions that are performed when transitioning to a certain state.
	errors  [][scanner.TokensCount]string    // Array of error messages returned when transitioning to the err state.
}

// Clears the structure fields to read the new line.
// Used in the start state.
func (m *finiteStateMachine) clearElement() {
	var field reflect.Value
	for i := 0; i < m.element.NumField(); i++ {
		field = m.element.Field(i)
		switch field.Kind() {
		case reflect.Int:
			field.SetInt(0)
		case reflect.Float64:
			field.SetFloat(0)
		case reflect.String:
			field.SetString("")
			// TODO add cleaning of structures and slices
		}
	}
}

// Returns the next state that is not yet used by the finite state machine.
func (m *finiteStateMachine) nextState() stateType {
	return stateType(len(m.matrix) + 1)
}

// Updates the finite state machine by adding a new state to it.
func (m *finiteStateMachine) update(row *machineRow) {
	m.matrix = append(m.matrix, row.matrixRow)
	m.actions = append(m.actions, row.action)
	m.errors = append(m.errors, row.errorsRow)
}

// Implementation of the next method in the elementParser interface.
func (m *finiteStateMachine) next(tokenType scanner.TokenType, state stateType) stateType {
	return m.matrix[state][tokenType]
}

// Implementation of the action method in the elementParser interface.
func (m *finiteStateMachine) action(state stateType, token string) {
	m.actions[state](token)
}

// Implementation of the message method in the elementParser interface.
func (m *finiteStateMachine) message(tokenType scanner.TokenType, state stateType) string {
	return m.errors[state][tokenType]
}

// Implementation of the result method in the elementParser interface.
func (m *finiteStateMachine) result() interface{} {
	return m.element.Interface()
}

// Creates a new object of a finiteStateMachine.
// Initializes fields and allocates memory.
// Adds processing of the first two reserved states - the start state and the err state.
func newFiniteStateMachine(element reflect.Value, elementType ElementType) *finiteStateMachine {
	var m = finiteStateMachine{
		element: element,
		matrix:  make([][scanner.TokensCount]stateType, 0, 10),
		actions: make([]action, 0, 10),
		errors:  make([][scanner.TokensCount]string, 0, 10),
	}
	m.update(&machineRow{
		matrixRow: [...]stateType{err, err, err, err, 2, err, err, err, err},
		action:    func(token string) { m.clearElement() },
		errorsRow: [...]string{
			impossibleTokenInStartStateError(scanner.Word),
			impossibleTokenInStartStateError(scanner.Integer),
			impossibleTokenInStartStateError(scanner.Float),
			impossibleTokenInStartStateError(scanner.Slash),
			noErrorMassage,
			allParametersNotSpecifiedError(elementType),
			allParametersNotSpecifiedError(elementType),
			impossibleTokenInStartStateError(scanner.Unknown),
			impossibleTokenInStartStateError(scanner.Comment),
		},
	})
	m.update(&machineRow{
		matrixRow: [...]stateType{err, err, err, err, err, err, err, err, err},
		action:    startStateAction,
		errorsRow: [...]string{
			parserUsedInErrorStateMessage,
			parserUsedInErrorStateMessage,
			parserUsedInErrorStateMessage,
			parserUsedInErrorStateMessage,
			parserUsedInErrorStateMessage,
			parserUsedInErrorStateMessage,
			parserUsedInErrorStateMessage,
			parserUsedInErrorStateMessage,
			parserUsedInErrorStateMessage,
		},
	})
	return &m
}

// Describes methods for updating a finite state machine based on a field from a structure.
type parameter interface {
	// Updates the finiteStateMachine to read the parameter.
	// Since a particular parameter does not have information about the need to read other parameters,
	// the information about the end-of-line transition must be passed from the outside.
	updateMachine(machine *finiteStateMachine, onEndState stateType, onEndMessage string)
	// Returns the parameter name read from the structure tags.
	String() string
	// Returns the name of the parameter, if it is required.
	// For structures, returns a comma-separated list of required parameter names.
	// For slices, returns the number of parameters to be read or a comma-separated list of required parameters.
	requiredString() string
}

type intParameter struct {
	value    reflect.Value
	name     string
	optional bool
}

func (p *intParameter) updateMachine(machine *finiteStateMachine, onEndState stateType, onEndMessage string) {
	var state = machine.nextState()
	machine.update(&machineRow{
		matrixRow: [...]stateType{err, state, err, err, err, onEndState, onEndState, err, err},
		action: func(token string) {
			var value, err = strconv.ParseInt(token, 10, 64)
			if err != nil {
				panic("failed to convert the token to an integer when reading " + p.String())
			}
			p.set(value)
		},
		errorsRow: [...]string{
			invalidParameterError(p, scanner.Integer, scanner.Word),
			noErrorMassage,
			invalidParameterError(p, scanner.Integer, scanner.Float),
			invalidParameterError(p, scanner.Integer, scanner.Slash),
			impossibleTokenWhenReadingParameterError(p, scanner.Space),
			onEndMessage,
			onEndMessage,
			invalidParameterError(p, scanner.Integer, scanner.Unknown),
			impossibleTokenWhenReadingParameterError(p, scanner.Comment),
		},
	})
}

func (p *intParameter) String() string {
	return p.name
}

func (p *intParameter) requiredString() string {
	if p.optional {
		return ""
	} else {
		return p.name
	}
}

func (p intParameter) set(value int64) {
	p.value.SetInt(value)
}

func newStructIntParameter(field reflect.StructField, value reflect.Value) *intParameter {
	var (
		p    = intParameter{value: value}
		tags = field.Tag
	)
	if name, ok := tags.Lookup("name"); ok {
		p.name = name
	} else {
		p.name = field.Name
	}
	if optional, ok := tags.Lookup("optional"); ok {
		switch optional {
		case "true":
			p.optional = true
		case "false":
			p.optional = false
		default:
			panic("the optional tag must take the values 'true' or 'false'")
		}
	} else {
		p.optional = false
	}
	if _, ok := tags.Lookup("delimiter"); ok {
		panic("the delimiter tag cannot be set for an int field")
	}
	if _, ok := tags.Lookup("min"); ok {
		panic("the min tag cannot be set for an int field")
	}
	return &p
}

func newSliceIntParameter(field reflect.StructField, value reflect.Value) *intParameter {
	var (
		p    = intParameter{value: value}
		tags = field.Tag
	)
	if name, ok := tags.Lookup("name"); ok {
		p.name = name
	} else {
		p.name = field.Name
	}
	if _, ok := tags.Lookup("optional"); ok {
		panic("the optional tag cannot be set for a slice of int field")
	}
	if _, ok := tags.Lookup("delimiter"); ok {
		panic("the delimiter tag cannot be set for a slice of int field")
	}
	return &p
}

type floatParameter struct {
	value    reflect.Value
	name     string
	optional bool
}

func (p *floatParameter) updateMachine(machine *finiteStateMachine, onEndState stateType, onEndMessage string) {
	var state = machine.nextState()
	machine.update(&machineRow{
		matrixRow: [...]stateType{err, state, state, err, err, onEndState, onEndState, err, err},
		action: func(token string) {
			var value, err = strconv.ParseFloat(token, 64)
			if err != nil {
				panic("failed to convert the token to a float when reading " + p.String())
			}
			p.set(value)
		},
		errorsRow: [...]string{
			invalidParameterError(p, scanner.Float, scanner.Word),
			noErrorMassage,
			noErrorMassage,
			invalidParameterError(p, scanner.Float, scanner.Slash),
			impossibleTokenWhenReadingParameterError(p, scanner.Space),
			onEndMessage,
			onEndMessage,
			invalidParameterError(p, scanner.Float, scanner.Unknown),
			impossibleTokenWhenReadingParameterError(p, scanner.Comment),
		},
	})
}

func (p *floatParameter) String() string {
	return p.name
}

func (p *floatParameter) requiredString() string {
	if p.optional {
		return ""
	} else {
		return p.name
	}
}

func (p floatParameter) set(value float64) {
	p.value.SetFloat(value)
}

func newStructFloatParameter(field reflect.StructField, value reflect.Value) *floatParameter {
	var (
		p    = floatParameter{value: value}
		tags = field.Tag
	)
	if name, ok := tags.Lookup("name"); ok {
		p.name = name
	} else {
		p.name = field.Name
	}
	if optional, ok := tags.Lookup("optional"); ok {
		switch optional {
		case "true":
			p.optional = true
		case "false":
			p.optional = false
		default:
			panic("the optional tag must take the values 'true' or 'false'")
		}
	} else {
		p.optional = false
	}
	if _, ok := tags.Lookup("delimiter"); ok {
		panic("the delimiter tag cannot be set for a float64 field")
	}
	if _, ok := tags.Lookup("min"); ok {
		panic("the min tag cannot be set for a float64 field")
	}
	return &p
}

func newSliceFloatParameter(field reflect.StructField, value reflect.Value) *floatParameter {
	var (
		p    = floatParameter{value: value}
		tags = field.Tag
	)
	if name, ok := tags.Lookup("name"); ok {
		p.name = name
	} else {
		p.name = field.Name
	}
	if _, ok := tags.Lookup("optional"); ok {
		panic("the optional tag cannot be set for a slice of float64 field")
	}
	if _, ok := tags.Lookup("delimiter"); ok {
		panic("the delimiter tag cannot be set for a slice of float64 field")
	}
	return &p
}

type stringParameter struct {
	value reflect.Value
	name  string
}

func (p *stringParameter) updateMachine(machine *finiteStateMachine, onEndState stateType, onEndMessage string) {
	var state = machine.nextState()
	machine.update(&machineRow{
		matrixRow: [...]stateType{state, err, err, err, err, onEndState, onEndState, err, err},
		action:    func(token string) { p.set(token) },
		errorsRow: [...]string{
			noErrorMassage,
			invalidParameterError(p, scanner.Word, scanner.Integer),
			invalidParameterError(p, scanner.Word, scanner.Float),
			invalidParameterError(p, scanner.Word, scanner.Slash),
			impossibleTokenWhenReadingParameterError(p, scanner.Space),
			onEndMessage,
			onEndMessage,
			invalidParameterError(p, scanner.Word, scanner.Unknown),
			impossibleTokenWhenReadingParameterError(p, scanner.Comment),
		},
	})
}

func (p *stringParameter) String() string {
	return p.name
}

func (p *stringParameter) requiredString() string {
	return p.name
}

func (p stringParameter) set(value string) {
	p.value.SetString(value)
}

func newStructStringParameter(field reflect.StructField, value reflect.Value) *stringParameter {
	var (
		p    = stringParameter{value: value}
		tags = field.Tag
	)
	if name, ok := tags.Lookup("name"); ok {
		p.name = name
	} else {
		p.name = field.Name
	}
	if _, ok := tags.Lookup("optional"); ok {
		panic("the optional tag cannot be set for a string field")
	}
	if _, ok := tags.Lookup("delimiter"); ok {
		panic("the delimiter tag cannot be set for a string field")
	}
	if _, ok := tags.Lookup("min"); ok {
		panic("the min tag cannot be set for a string field")
	}
	return &p
}

func newSliceStringParameter(field reflect.StructField, value reflect.Value) *stringParameter {
	var (
		p    = stringParameter{value: value}
		tags = field.Tag
	)
	if name, ok := tags.Lookup("name"); ok {
		p.name = name
	} else {
		p.name = field.Name
	}
	if _, ok := tags.Lookup("optional"); ok {
		panic("the optional tag cannot be set for a slice of string field")
	}
	if _, ok := tags.Lookup("delimiter"); ok {
		panic("the delimiter tag cannot be set for a slice of string field")
	}
	return &p
}

// TODO optimize the fields for a more convenient update of the finite state machine
type structParameter struct {
	name      string
	params    []parameter
	delimiter scanner.TokenType
}

func (p *structParameter) updateMachine(machine *finiteStateMachine, onEndState stateType, onEndMessage string) {
	// TODO implement updateMachine for structParameter
}

func (p *structParameter) String() string {
	return p.name
}

func (p *structParameter) requiredString() string {
	return p.name
}

func newStructStructParameter(field reflect.StructField, value reflect.Value) *structParameter {
	var t = field.Type
	if t.NumField() < 1 {
		panic("a field of a read object cannot contain a structure without fields")
	}
	var (
		p    = structParameter{}
		tags = field.Tag
	)
	if name, ok := tags.Lookup("name"); ok {
		p.name = name
	} else {
		p.name = field.Name
	}
	if _, ok := tags.Lookup("optional"); ok {
		panic("the optional tag cannot be set for a struct field")
	}
	if delimiter, ok := tags.Lookup("delimiter"); ok {
		switch delimiter {
		case "slash":
			p.delimiter = scanner.Slash
		case "space":
			p.delimiter = scanner.Space
		default:
			panic("the delimiter tag must take the values 'space' or 'slash'")
		}
	} else {
		panic("the structure field must have the delimiter tag specified")
	}
	if _, ok := tags.Lookup("min"); ok {
		panic("the min tag cannot be set for a struct field")
	}
	p.params = make([]parameter, 0, 5)
	var (
		f          reflect.StructField
		optional   = false
		intParam   *intParameter
		floatParam *floatParameter
	)
	for i := 0; i < t.NumField(); i++ {
		f = t.Field(i)
		switch f.Type.Kind() {
		case reflect.Int:
			intParam = newStructIntParameter(f, value.Field(i))
			if i == 0 && intParam.optional {
				panic("the first field of the structure cannot be optional")
			}
			if p.delimiter == scanner.Space && intParam.optional {
				panic("a field of the struct type with a space delimiter cannot contain optional fields")
			}
			if optional && !intParam.optional {
				panic("an optional field cannot be followed by a required field")
			}
			optional = intParam.optional
			p.params = append(p.params, intParam)
		case reflect.Float64:
			floatParam = newStructFloatParameter(f, value.Field(i))
			if i == 0 && floatParam.optional {
				panic("the first field of the structure cannot be optional")
			}
			if p.delimiter == scanner.Space && floatParam.optional {
				panic("a field of the struct type with a space delimiter cannot contain optional fields")
			}
			if optional && !floatParam.optional {
				panic("an optional field cannot be followed by a required field")
			}
			optional = floatParam.optional
			p.params = append(p.params, floatParam)
		default:
			panic("unsupported type of structure field that is a struct field: " + f.Type.Kind().String())
		}
	}
	return &p
}

func newSliceStructParameter(field reflect.StructField, value reflect.Value) *structParameter {
	var t = field.Type.Elem()
	if t.NumField() < 1 {
		panic("a field of a read object cannot contain a slice of structure without fields")
	}
	var (
		p    = structParameter{}
		tags = field.Tag
	)
	if name, ok := tags.Lookup("name"); ok {
		p.name = name
	} else {
		p.name = field.Name
	}
	if _, ok := tags.Lookup("optional"); ok {
		panic("the optional tag cannot be set for a slice of struct field")
	}
	if delimiter, ok := tags.Lookup("delimiter"); ok {
		switch delimiter {
		case "slash":
			p.delimiter = scanner.Slash
		case "space":
			p.delimiter = scanner.Space
		default:
			panic("the delimiter tag must take the values 'space' or 'slash'")
		}
	} else {
		panic("the slice of structure field must have the delimiter tag specified")
	}
	p.params = make([]parameter, 0, 5)
	var (
		f          reflect.StructField
		optional   = false
		intParam   *intParameter
		floatParam *floatParameter
	)
	for i := 0; i < t.NumField(); i++ {
		f = t.Field(i)
		switch f.Type.Kind() {
		case reflect.Int:
			intParam = newStructIntParameter(f, value.Field(i))
			if i == 0 && intParam.optional {
				panic("the first field of the structure cannot be optional")
			}
			if p.delimiter == scanner.Space && intParam.optional {
				panic("a field of the struct type with a space delimiter cannot contain optional fields")
			}
			if optional && !intParam.optional {
				panic("an optional field cannot be followed by a required field")
			}
			optional = intParam.optional
			p.params = append(p.params, intParam)
		case reflect.Float64:
			floatParam = newStructFloatParameter(f, value.Field(i))
			if i == 0 && floatParam.optional {
				panic("the first field of the structure cannot be optional")
			}
			if p.delimiter == scanner.Space && floatParam.optional {
				panic("a field of the struct type with a space delimiter cannot contain optional fields")
			}
			if optional && !floatParam.optional {
				panic("an optional field cannot be followed by a required field")
			}
			optional = floatParam.optional
			p.params = append(p.params, floatParam)
		default:
			panic("unsupported type of structure field that is a slice of struct field: " + f.Type.Kind().String())
		}
	}
	return &p
}

// TODO optimize the fields for a more convenient update of the finite state machine
type sliceParameter struct {
	value reflect.Value
	param parameter
	min   int
}

func (p *sliceParameter) updateMachine(machine *finiteStateMachine, onEndState stateType, onEndMessage string) {
	// TODO implement updateMachine for sliceParameter
}

func (p *sliceParameter) String() string {
	switch p.min {
	case 1:
		return p.param.String()
	case 2:
		return fmt.Sprintf("first %s and second %s", p.param, p.param)
	default:
		return fmt.Sprintf("%d %s values", p.min, p.param)
	}
}

func (p *sliceParameter) requiredString() string {
	return p.String()
}

func (p sliceParameter) new() reflect.Value {
	var newValue = reflect.New(p.value.Type().Elem()).Elem()
	p.value.Set(reflect.Append(p.value, newValue))
	return newValue
}

func newSliceParameter(field reflect.StructField, value reflect.Value) *sliceParameter {
	var (
		tags = field.Tag
		p    = sliceParameter{value: value}
	)
	if min, ok := tags.Lookup("min"); ok {
		minInt, err := strconv.ParseUint(min, 10, 8)
		if err != nil {
			panic("error reading the min tag")
		}
		if minInt < 1 {
			panic("the min tag cannot accept values less than one")
		}
		p.min = int(minInt)
	} else {
		panic("the slice field must have the min tag specified")
	}
	value.Set(reflect.MakeSlice(value.Type(), 0, 5))
	switch field.Type.Elem().Kind() {
	case reflect.Int:
		p.param = newSliceIntParameter(field, p.new())
	case reflect.Float64:
		p.param = newSliceFloatParameter(field, p.new())
	case reflect.String:
		p.param = newSliceStringParameter(field, p.new())
	case reflect.Struct:
		p.param = newSliceStructParameter(field, p.new())
	default:
		panic("unsupported structure field type: slice of " + field.Type.Elem().Kind().String())
	}
	return &p
}

// Converts a structure object to a slice of parameters that parse this object from a line.
func buildParameters(v reflect.Value) []parameter {
	var t = v.Type()
	if t.Kind() != reflect.Struct {
		panic("the element to be read must be a structure object")
	}
	if t.NumField() < 0 {
		panic("the parser cannot be built on a structure without fields")
	}
	var (
		params     = make([]parameter, 0, 5)
		field      reflect.StructField
		optional   = false
		intParam   *intParameter
		floatParam *floatParameter
	)
	for i := 0; i < t.NumField(); i++ {
		field = t.Field(i)
		switch field.Type.Kind() {
		case reflect.Int:
			intParam = newStructIntParameter(field, v.Field(i))
			if i == 0 && intParam.optional {
				panic("the first field of the structure cannot be optional")
			}
			if optional && !intParam.optional {
				panic("an optional field cannot be followed by a required field")
			}
			optional = intParam.optional
			params = append(params, intParam)
		case reflect.Float64:
			floatParam = newStructFloatParameter(field, v.Field(i))
			if i == 0 && floatParam.optional {
				panic("the first field of the structure cannot be optional")
			}
			if optional && !floatParam.optional {
				panic("an optional field cannot be followed by a required field")
			}
			optional = floatParam.optional
			params = append(params, floatParam)
		case reflect.String:
			if optional {
				panic("an optional field cannot be followed by a required field")
			}
			params = append(params, newStructStringParameter(field, v.Field(i)))
		case reflect.Struct:
			if optional {
				panic("an optional field cannot be followed by a required field")
			}
			params = append(params, newStructStructParameter(field, v.Field(i)))
		case reflect.Slice:
			if i != t.NumField()-1 {
				panic("the slice must be the last field of the structure")
			}
			if optional {
				panic("an optional field cannot be followed by a slice")
			}
			params = append(params, newSliceParameter(field, v.Field(i)))
		default:
			panic("unsupported structure field type: " + field.Type.Kind().String())
		}
	}
	return params
}

// Creates an ElementParser that parses the line based on the structure.
// elementType specifies the type of element to be read.
// element specifies the structure on the basis of which fields the line will be read.
// The structure to be read must match the element type.
// The structure fields are extracted from the line in the order in which they are specified in the structure.
//
// The following limitations apply to the structure:
// 	* The element must have the base type struct.
// 	* Only public fields will be parsed.
// 	* Structure fields must have one of the following basic types: int, float64, string, struct, []int, []float64, []string, []struct.
// 	* If the field is of the slice type, it must be the last one in the structure.
// 	* If a field is of the structure type, its fields must be of the base type int or float64.
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
// 	This tag must be specified for the struct, []struct types and cannot be specified for other types.
// 	Used for reading structure field delimiters.
// 	If a structure has a 'space' delimiter, it cannot contain optional fields.
//
// 	min
//
// 	It can only accept integer values that are greater than zero.
// 	This tag must be specified for slices and cannot be specified for other types.
// 	Used to specify the minimum number of slice elements.
func buildParser(elementType ElementType, element interface{}) elementParser {
	var (
		value       = reflect.New(reflect.TypeOf(element)).Elem()
		params      = buildParameters(value)
		paramNames  = make([]string, 0, len(params))
		machine     = newFiniteStateMachine(value, elementType)
		paramString string
	)
	for _, param := range params {
		paramString = param.requiredString()
		if paramString != "" {
			paramNames = append(paramNames, paramString)
		}
	}
	var (
		state             stateType
		unreadParamsCount int
		onEndState        stateType
		onEndMessage      string
	)
	for i, param := range params {
		unreadParamsCount = len(paramNames) - i
		if unreadParamsCount <= 0 {
			onEndState = start
			onEndMessage = noErrorMassage
		} else {
			onEndState = err
			onEndMessage = parametersNotSpecifiedError(strings.Join(paramNames[i:], ", "), unreadParamsCount)
		}
		param.updateMachine(machine, onEndState, onEndMessage)
		state = machine.nextState()
		if i == len(params)-1 {
			machine.update(&machineRow{
				matrixRow: [...]stateType{err, err, err, err, state, start, start, err, err},
				action:    noAction,
				errorsRow: [...]string{
					impossibleTokenAfterReadingParameterError(param, scanner.Word),
					impossibleTokenAfterReadingParameterError(param, scanner.Integer),
					impossibleTokenAfterReadingParameterError(param, scanner.Float),
					unexpectedTokenAfterReadingParameterError(param, scanner.Slash),
					noErrorMassage,
					noErrorMassage,
					noErrorMassage,
					impossibleTokenAfterReadingParameterError(param, scanner.Unknown),
					impossibleTokenAfterReadingParameterError(param, scanner.Comment),
				},
			})
			machine.update(&machineRow{
				matrixRow: [...]stateType{err, err, err, err, err, start, start, err, err},
				action:    lastStateAction,
				errorsRow: [...]string{
					unexpectedTokenAfterDescribingElementError(elementType, scanner.Word),
					unexpectedTokenAfterDescribingElementError(elementType, scanner.Integer),
					unexpectedTokenAfterDescribingElementError(elementType, scanner.Float),
					unexpectedTokenAfterDescribingElementError(elementType, scanner.Slash),
					unexpectedTokenAfterDescribingElementError(elementType, scanner.Space),
					noErrorMassage,
					noErrorMassage,
					unexpectedTokenAfterDescribingElementError(elementType, scanner.Unknown),
					impossibleTokenAfterDescribingElementError(elementType, scanner.Comment),
				},
			})
		} else {
			unreadParamsCount--
			if unreadParamsCount <= 0 {
				onEndState = start
				onEndMessage = noErrorMassage
			} else {
				onEndState = err
				onEndMessage = parametersNotSpecifiedError(strings.Join(paramNames[i+1:], ", "), unreadParamsCount)
			}
			machine.update(&machineRow{
				matrixRow: [...]stateType{err, err, err, err, state, onEndState, onEndState, err, err},
				action:    noAction,
				errorsRow: [...]string{
					impossibleTokenAfterReadingParameterError(param, scanner.Word),
					impossibleTokenAfterReadingParameterError(param, scanner.Integer),
					impossibleTokenAfterReadingParameterError(param, scanner.Float),
					invalidDelimiterBetweenParametersError(param, params[i+1], scanner.Space, scanner.Slash),
					impossibleTokenAfterDescribingElementError(elementType, scanner.Space),
					onEndMessage,
					onEndMessage,
					impossibleTokenAfterReadingParameterError(param, scanner.Unknown),
					impossibleTokenAfterReadingParameterError(param, scanner.Comment),
				},
			})
		}
	}
	return machine
}
