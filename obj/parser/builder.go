package parser

import (
	"computer_graphics/obj/scanner"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	noErrorMassage                = ""
	parserUsedInErrorStateMessage = "parser cannot be used in the error state"
)

func impossibleTokenInStartStateError(tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received in the start state - %s", tokenType)
}

func impossibleTokenWhenReadingParameterError(param parameter, tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received when reading the %s - %s", param.getName(), tokenType)
}

func impossibleTokenAfterReadingParameterError(param parameter, tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received after reading the %s - %s", param.getName(), tokenType)
}

func impossibleTokenAfterDescribingElementError(elementType ElementType, tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received after describing a %s - %s", elementType, tokenType)
}

func unexpectedTokenAfterReadingParameterError(param parameter, tokenType scanner.TokenType) string {
	return fmt.Sprintf("unexpected token received after reading the %s - %s", param.getName(), tokenType)
}

func unexpectedTokenAfterDescribingElementError(elementType ElementType, tokenType scanner.TokenType) string {
	return fmt.Sprintf("unexpected token received after describing a %s - %s", elementType, tokenType)
}

func parametersNotSpecifiedError(unreadParamsString string, unreadParamsCount int) string {
	if unreadParamsCount == 1 {
		return fmt.Sprintf("parameter %s is not specified", unreadParamsString)
	} else {
		return fmt.Sprintf("parameters %s are not specified", unreadParamsString)
	}
}

func notSpecifiedParametersOfElementTypeError(elementType ElementType) string {
	return fmt.Sprintf("not specified parameters of the %s", elementType)
}

func invalidParameterError(param parameter, expected, received scanner.TokenType) string {
	return fmt.Sprintf("invalid %s, excepted: %s, received: %s", param.getName(), expected, received)
}

func invalidDelimiterBetweenParametersError(first, second parameter, expected, received scanner.TokenType) string {
	return fmt.Sprintf(
		"invalid delimiter format between %s and %s, expected: %s, received: %s",
		first.getName(),
		second.getName(),
		expected,
		received,
	)
}

type action func(token string)

var (
	noAction         action = func(token string) {}
	startStateAction        = func(token string) { panic("the action method is called in the err state") }
	lastStateAction         = func(token string) {
		panic("the action method cannot be called in a state from which transitions are made only to the start state and the err state")
	}
)

type machineRow struct {
	matrixRow [scanner.TokensCount]stateType
	action    action
	errorsRow [scanner.TokensCount]string
}

type finiteStateMachine struct {
	element reflect.Value
	matrix  [][scanner.TokensCount]stateType
	actions []action
	errors  [][scanner.TokensCount]string
}

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
			// todo add cleaning of structures and slices
		}
	}
}

func (m *finiteStateMachine) nextState() stateType {
	return stateType(len(m.matrix) + 1)
}

func (m *finiteStateMachine) update(row *machineRow) {
	m.matrix = append(m.matrix, row.matrixRow)
	m.actions = append(m.actions, row.action)
	m.errors = append(m.errors, row.errorsRow)
}

func (m *finiteStateMachine) next(tokenType scanner.TokenType, state stateType) stateType {
	return m.matrix[state][tokenType]
}

func (m *finiteStateMachine) action(state stateType, token string) {
	m.actions[state](token)
}

func (m *finiteStateMachine) message(tokenType scanner.TokenType, state stateType) string {
	return m.errors[state][tokenType]
}

func (m *finiteStateMachine) result() interface{} {
	return m.element.Interface()
}

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
			notSpecifiedParametersOfElementTypeError(elementType),
			notSpecifiedParametersOfElementTypeError(elementType),
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

type parameter interface {
	updateMachine(machine *finiteStateMachine, onEndState stateType, onEndMessage string)
	getName() string
	getNameNotOptional() string
	isOptional() bool
}

type intParameter struct {
	value    reflect.Value
	name     string
	optional bool
	wasRead  bool
}

func (p *intParameter) updateMachine(machine *finiteStateMachine, onEndState stateType, onEndMessage string) {
	var (
		state        = machine.nextState()
		convertError = "failed to convert the token to an integer when reading " + p.getName()
	)
	machine.update(&machineRow{
		matrixRow: [...]stateType{err, state, err, err, err, onEndState, onEndState, err, err},
		action: func(token string) {
			var value, err = strconv.ParseInt(token, 10, 64)
			if err != nil {
				panic(convertError)
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

func (p *intParameter) getName() string {
	return p.name
}

func (p *intParameter) getNameNotOptional() string {
	if p.optional {
		return ""
	} else {
		return p.name
	}
}

func (p *intParameter) isOptional() bool {
	return p.optional
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
	var (
		state        = machine.nextState()
		convertError = "failed to convert the token to a float when reading " + p.getName()
	)
	machine.update(&machineRow{
		matrixRow: [...]stateType{err, state, state, err, err, onEndState, onEndState, err, err},
		action: func(token string) {
			var value, err = strconv.ParseFloat(token, 64)
			if err != nil {
				panic(convertError)
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

func (p *floatParameter) getName() string {
	return p.name
}

func (p *floatParameter) getNameNotOptional() string {
	if p.optional {
		return ""
	} else {
		return p.name
	}
}

func (p *floatParameter) isOptional() bool {
	return p.optional
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
	// todo stringParameter updateMachine
}

func (p *stringParameter) getName() string {
	return p.name
}

func (p *stringParameter) getNameNotOptional() string {
	return p.name
}

func (p *stringParameter) isOptional() bool {
	return false
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

type structParameter struct {
	name      string
	params    []parameter
	delimiter scanner.TokenType
}

func (p *structParameter) updateMachine(machine *finiteStateMachine, onEndState stateType, onEndMessage string) {
	// todo structParameter updateMachine
}

func (p *structParameter) getName() string {
	return p.name
}

func (p *structParameter) getNameNotOptional() string {
	return p.name
}

func (p *structParameter) isOptional() bool {
	for _, param := range p.params {
		if param.isOptional() {
			return true
		}
	}
	return false
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
		f        reflect.StructField
		optional = false
		param    parameter
	)
	for i := 0; i < t.NumField(); i++ {
		f = t.Field(i)
		switch f.Type.Kind() {
		case reflect.Int:
			param = newStructIntParameter(f, value.Field(i))
			if i == 0 && param.isOptional() {
				panic("the first field of the structure cannot be optional")
			}
			if p.delimiter == scanner.Space && param.isOptional() {
				panic("a field of the struct type with a space delimiter cannot contain optional fields")
			}
			if optional && !param.isOptional() {
				panic("an optional field cannot be followed by a required field")
			}
			optional = param.isOptional()
			p.params = append(p.params, param)
		case reflect.Float64:
			param = newStructFloatParameter(f, value.Field(i))
			if i == 0 && param.isOptional() {
				panic("the first field of the structure cannot be optional")
			}
			if p.delimiter == scanner.Space && param.isOptional() {
				panic("a field of the struct type with a space delimiter cannot contain optional fields")
			}
			if optional && !param.isOptional() {
				panic("an optional field cannot be followed by a required field")
			}
			optional = param.isOptional()
			p.params = append(p.params, param)
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
		f        reflect.StructField
		optional = false
		param    parameter
	)
	for i := 0; i < t.NumField(); i++ {
		f = t.Field(i)
		switch f.Type.Kind() {
		case reflect.Int:
			param = newStructIntParameter(f, value.Field(i))
			if i == 0 && param.isOptional() {
				panic("the first field of the structure cannot be optional")
			}
			if p.delimiter == scanner.Space && param.isOptional() {
				panic("a field of the struct type with a space delimiter cannot contain optional fields")
			}
			if optional && !param.isOptional() {
				panic("an optional field cannot be followed by a required field")
			}
			optional = param.isOptional()
			p.params = append(p.params, param)
		case reflect.Float64:
			param = newStructFloatParameter(f, value.Field(i))
			if i == 0 && param.isOptional() {
				panic("the first field of the structure cannot be optional")
			}
			if p.delimiter == scanner.Space && param.isOptional() {
				panic("a field of the struct type with a space delimiter cannot contain optional fields")
			}
			if optional && !param.isOptional() {
				panic("an optional field cannot be followed by a required field")
			}
			optional = param.isOptional()
			p.params = append(p.params, param)
		default:
			panic("unsupported type of structure field that is a slice of struct field: " + f.Type.Kind().String())
		}
	}
	return &p
}

type sliceParameter struct {
	value reflect.Value
	param parameter
	min   int
}

func (p *sliceParameter) updateMachine(machine *finiteStateMachine, onEndState stateType, onEndMessage string) {
	// todo sliceParameter updateMachine
}

func (p *sliceParameter) getName() string {
	var paramName = p.param.getName()
	switch p.min {
	case 1:
		return paramName
	case 2:
		return "first " + paramName + " and second " + paramName
	default:
		return strconv.Itoa(p.min) + " " + paramName + " values"
	}
}

func (p *sliceParameter) getNameNotOptional() string {
	return p.getName()
}

func (p *sliceParameter) isOptional() bool {
	return false
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

func buildParameters(v reflect.Value) []parameter {
	var t = v.Type()
	if t.Kind() != reflect.Struct {
		panic("the element to be read must be a structure object")
	}
	if t.NumField() < 0 {
		panic("the parser cannot be built on a structure without fields")
	}
	var (
		params   = make([]parameter, 0, 5)
		field    reflect.StructField
		optional = false
		param    parameter
	)
	for i := 0; i < t.NumField(); i++ {
		field = t.Field(i)
		switch field.Type.Kind() {
		case reflect.Int:
			param = newStructIntParameter(field, v.Field(i))
			if i == 0 && param.isOptional() {
				panic("the first field of the structure cannot be optional")
			}
			if optional && !param.isOptional() {
				panic("an optional field cannot be followed by a required field")
			}
			optional = param.isOptional()
			params = append(params, param)
		case reflect.Float64:
			param = newStructFloatParameter(field, v.Field(i))
			if i == 0 && param.isOptional() {
				panic("the first field of the structure cannot be optional")
			}
			if optional && !param.isOptional() {
				panic("an optional field cannot be followed by a required field")
			}
			optional = param.isOptional()
			params = append(params, param)
		case reflect.String:
			if optional {
				panic("an optional field cannot be followed by a required field")
			}
			params = append(params, newStructStringParameter(field, v.Field(i)))
		case reflect.Struct:
			if optional {
				panic("an optional field cannot be followed by a required field")
			}
			param = newStructStructParameter(field, v.Field(i))
			if param.isOptional() && i != t.NumField()-1 {
				panic("the structure containing the optional fields must be the last field")
			}
			params = append(params, param)
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

func buildParser(elementType ElementType, element interface{}) elementParser {
	var (
		value       = reflect.New(reflect.TypeOf(element)).Elem()
		params      = buildParameters(value)
		paramNames  = make([]string, 0, len(params))
		machine     = newFiniteStateMachine(value, elementType)
		paramString string
	)
	for _, param := range params {
		paramString = param.getNameNotOptional()
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
		unreadParamsCount = len(paramNames) - i - 1
		if unreadParamsCount <= 0 {
			onEndState = start
			onEndMessage = noErrorMassage
		} else {
			onEndState = err
			onEndMessage = parametersNotSpecifiedError(strings.Join(paramNames[i+1:], ", "), unreadParamsCount)
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
