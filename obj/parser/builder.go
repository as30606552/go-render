package parser

import (
	"computer_graphics/obj/scanner"
	"computer_graphics/obj/types"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const noErrorMessage = ""

func impossibleTokenInStartStateMessage(tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received in the start state - %s", tokenType)
}

func impossibleTokenAfterDescribingElementMessage(elementType ElementType, tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received after describing a %s - %s", elementType, tokenType)
}

func unexpectedTokenAfterDescribingElementMessage(elementType ElementType, tokenType scanner.TokenType) string {
	return fmt.Sprintf("unexpected token received after describing a %s - %s", elementType, tokenType)
}

func impossibleTokenMessage(name string, tokenType scanner.TokenType) string {
	return fmt.Sprintf("impossible token received when reading the %s - %s", name, tokenType)
}

func invalidTokenMessage(name string, expected, received scanner.TokenType) string {
	return fmt.Sprintf("invalid %s, expected: %s, received: %s", name, expected, received)
}

func parametersNotSpecifiedMessage(paramNames []string) string {
	if len(paramNames) == 1 {
		return fmt.Sprintf("parameter %s is not specified", paramNames[0])
	} else {
		return fmt.Sprintf("parameters %s are not specified", strings.Join(paramNames, ", "))
	}
}

func invalidParameterFormatMessage(specific, base fmt.Stringer, expected, received scanner.TokenType) string {
	return fmt.Sprintf(
		"invalid format for description of the %s, it must be the same as the first %s, expected: %s, received: %s",
		specific,
		base,
		expected,
		received,
	)
}

func delimiterBetween(predecessor, successor string) string {
	return fmt.Sprintf("delimiter between %s and %s", predecessor, successor)
}

func tokenAfter(name string) string {
	return fmt.Sprintf("token after %s", name)
}

const initMatrixSize = 10

type action func(token string, element reflect.Value) error

type finiteStateMachine struct {
	element reflect.Value
	matrix  [][scanner.TokensCount]stateType
	actions []action
	errors  [][scanner.TokensCount]string
}

func (m *finiteStateMachine) clear() { m.element = reflect.New(m.element.Type()).Elem() }

func (m *finiteStateMachine) transition(tokenType scanner.TokenType, state stateType) stateType {
	return m.matrix[state][tokenType]
}

func (m *finiteStateMachine) action(state stateType, token string) error {
	return m.actions[state](token, m.element)
}

func (m *finiteStateMachine) message(tokenType scanner.TokenType, state stateType) string {
	return m.errors[state][tokenType]
}

func (m *finiteStateMachine) result() interface{} { return m.element.Interface() }

func newMachine(element reflect.Value, size int) *finiteStateMachine {
	return &finiteStateMachine{
		element: element,
		matrix:  make([][scanner.TokensCount]stateType, size),
		actions: make([]action, size),
		errors:  make([][scanner.TokensCount]string, size),
	}
}

type setter interface {
	set(token string, value reflect.Value) error
	expected() scanner.TokenType
	changeName(name string)
}

type boolSetter struct {
	error error
}

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

func (s *boolSetter) expected() scanner.TokenType { return scanner.Word }

func (s *boolSetter) changeName(name string) {
	s.error = errors.New(fmt.Sprintf("the %s parameter must take the values 'on' or 'off'", name))
}

func newBoolSetter(name string) *boolSetter {
	var s = &boolSetter{}
	s.changeName(name)
	return s
}

type directionTypeSetter struct {
	error error
}

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

func (s *directionTypeSetter) expected() scanner.TokenType { return scanner.Word }

func (s *directionTypeSetter) changeName(name string) {
	s.error = errors.New(fmt.Sprintf("the %s parameter must take the values 'v' or 'u'", name))
}

func newDirectionTypeSetter(name string) *directionTypeSetter {
	var s = &directionTypeSetter{}
	s.changeName(name)
	return s
}

type intSetter struct {
	error error
}

func (s *intSetter) set(token string, value reflect.Value) error {
	var val, err = strconv.ParseInt(token, 10, 64)
	if err != nil {
		return s.error
	}
	value.SetInt(val)
	return nil
}

func (s *intSetter) expected() scanner.TokenType { return scanner.Integer }

func (s *intSetter) changeName(name string) {
	s.error = errors.New(fmt.Sprintf("failed to convert the token to an integer when reading %s", name))
}

func newIntSetter(name string) *intSetter {
	var s = &intSetter{}
	s.changeName(name)
	return s
}

type floatSetter struct {
	error error
}

func (s *floatSetter) set(token string, value reflect.Value) error {
	var val, err = strconv.ParseFloat(token, 64)
	if err != nil {
		return s.error
	}
	value.SetFloat(val)
	return nil
}

func (s *floatSetter) expected() scanner.TokenType { return scanner.Float }

func (s *floatSetter) changeName(name string) {
	s.error = errors.New(fmt.Sprintf("failed to convert the token to a float when reading %s", name))
}

func newFloatSetter(name string) *floatSetter {
	var s = &floatSetter{}
	s.changeName(name)
	return s
}

type stringSetter struct{}

func (s *stringSetter) set(token string, value reflect.Value) error {
	value.SetString(token)
	return nil
}

func (s *stringSetter) expected() scanner.TokenType { return scanner.Float }

func (s *stringSetter) changeName(name string) {}

func newStringSetter() *stringSetter { return &stringSetter{} }

type structSetter struct {
	fieldNumber int
	setter
}

func (s *structSetter) set(token string, value reflect.Value) error {
	return s.setter.set(token, value.Field(s.fieldNumber))
}

func newStructSetter(fieldNumber int, setter setter) *structSetter {
	return &structSetter{
		fieldNumber: fieldNumber,
		setter:      setter,
	}
}

type sliceSetter struct {
	setter
}

func (s *sliceSetter) set(token string, value reflect.Value) error {
	return s.setter.set(token, value.Index(value.Len()-1))
}

func newSliceSetter(setter setter) *sliceSetter {
	return &sliceSetter{
		setter: setter,
	}
}

type sliceAppender struct {
	setter
}

func (s *sliceAppender) set(token string, value reflect.Value) error {
	value.Set(reflect.Append(value, reflect.New(value.Type().Elem()).Elem()))
	return s.setter.set(token, value)
}

func newSliceAppender(setter setter) *sliceAppender {
	return &sliceAppender{
		setter: setter,
	}
}

type parameter interface {
	update(b *builder)
	String() string
}

type parameterName string

func (n parameterName) String() string { return string(n) }

type baseParameter struct {
	parameterName
	setter setter
}

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
	if expected == scanner.Integer || expected == scanner.Float {
		b.onInteger(state, act)
	} else {
		b.onIntegerError(invalidTokenMessage(p.String(), expected, scanner.Integer))
	}
	if expected == scanner.Float {
		b.onFloat(state, nil)
	} else {
		b.onFloatError(invalidTokenMessage(p.String(), expected, scanner.Float))
	}
	if len(unread) > 0 {
		b.onEndError(parametersNotSpecifiedMessage(unread))
	} else {
		b.onEnd()
	}
}

func (p *baseParameter) update(b *builder) {
	p.baseUpdate(b.nextParameterRow(p.String(), p.setter.expected()), b.nextState(), b.getUnread())
}

func newBaseParameter(name string, setter setter) *baseParameter {
	return &baseParameter{
		parameterName: parameterName(name),
		setter:        setter,
	}
}

type structParameter struct {
	parameterName
	delimiter     scanner.TokenType
	params        []*baseParameter
	requiredCount int
}

func (p *structParameter) changeName(name string) {
	p.parameterName = parameterName(name)
}

func (p *structParameter) name(num int) string {
	return fmt.Sprintf("%s of the %s", p.params[num], p)
}

func (p *structParameter) requiredNames() []string {
	var res = make([]string, p.requiredCount)
	for i := 0; i < p.requiredCount; i++ {
		res[i] = p.name(i)
	}
	return res
}

func (p *structParameter) allNames(to int) []string {
	var res = make([]string, to)
	for i := 0; i < to; i++ {
		res[i] = p.name(i)
	}
	return res
}

func (p *structParameter) updateRequired(b *builder, unread []string) {
	var (
		param *baseParameter
		name  string
	)
	for i := 0; i < p.requiredCount; i++ {
		param = p.params[i]
		name = unread[i]
		param.setter.changeName(name)
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

func (p *structParameter) update(b *builder) {
	p.updateRequired(b, append(p.requiredNames(), b.getUnread()[1:]...))
}

func newStructParameter(name string, delimiter scanner.TokenType, size int) *structParameter {
	return &structParameter{
		parameterName: parameterName(name),
		delimiter:     delimiter,
		params:        make([]*baseParameter, 0, size),
		requiredCount: 0,
	}
}

type baseSliceParameter struct {
	parameterName
	min   int
	param *baseParameter
}

func (p *baseSliceParameter) String() string { return fmt.Sprintf("%s parameters", p.parameterName) }

func (p *baseSliceParameter) name(num int) string {
	return fmt.Sprintf("%s number %d", p.parameterName, num+1)
}

func (p *baseSliceParameter) names() []string {
	var res = make([]string, p.min)
	for i := 0; i < p.min; i++ {
		res[i] = p.name(i)
	}
	return res
}

func (p *baseSliceParameter) update(b *builder) {
	var (
		name     string
		names    = p.names()
		expected = p.param.setter.expected()
	)
	for i := 0; i < p.min; i++ {
		name = p.name(i)
		p.param.setter.changeName(name)
		p.param.baseUpdate(b.nextParameterRow(name, expected), b.nextState(), names[i:])
		if i != p.min-1 {
			b.waitSpace(delimiterBetween(name, p.name(i+1)), names[i+1:])
		} else {
			b.waitSpace(tokenAfter(name), []string{})
		}
	}
	var loopState = b.nextState()
	name = fmt.Sprintf("additional %s", p.parameterName)
	p.param.setter.changeName(name)
	p.param.baseUpdate(b.nextParameterRow(name, expected), b.nextState(), []string{})
	name = tokenAfter(name)
	b.nextDelimiterRow(name).
		onSlashError(invalidTokenMessage(name, scanner.Space, scanner.Slash)).
		onSpace(loopState).
		onEnd()
}

func newBaseSliceParameter(name string, min int, param *baseParameter) *baseSliceParameter {
	return &baseSliceParameter{
		parameterName: parameterName(name),
		min:           min,
		param:         param,
	}
}

type structSliceParameter struct {
	parameterName
	min   int
	param *structParameter
}

func (p *structSliceParameter) name(num int) string {
	return fmt.Sprintf("%s number %d", p.parameterName, num+1)
}

func (p *structSliceParameter) names() []string {
	var res = make([]string, p.min)
	for i := 0; i < p.min; i++ {
		res[i] = p.name(i)
	}
	return res
}

func (p *structSliceParameter) String() string { return fmt.Sprintf("%s parameters", p.parameterName) }

func (p *structSliceParameter) updateRecursive(
	update func(b *builder, unread []string),
	b *builder,
	paramNumber int,
	lastSlash bool,
) {
	var (
		sliceNames = p.names()
		name       string
	)
	if paramNumber < len(p.param.params) {
		name = p.param.name(paramNumber)
		var delimiterRow *rowBuilder
		if !lastSlash {
			delimiterRow = b.nextDelimiterRow(delimiterBetween(p.param.name(paramNumber-1), name)).
				onSlash(b.nextState())
			if p.min > 1 {
				delimiterRow.onEndError(parametersNotSpecifiedMessage(sliceNames[1:]))
			} else {
				delimiterRow.onEnd()
			}
		}
		var (
			noParamUpdate = func(b *builder, unread []string) {
				update(b, unread)
				var (
					param = p.param.params[paramNumber]
					rb    *rowBuilder
				)
				if lastSlash {
					var (
						extParamMessage = fmt.Sprintf(
							"the %s is specified for %s, but is not specified for the first %s",
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
				rb.
					onSlash(b.nextState()).
					onSpaceError(fmt.Sprintf(
						"the %s is not specified for %s, but is specified for the first %s",
						param,
						p.param,
						p.parameterName,
					)).
					onEndError(parametersNotSpecifiedMessage(unread[paramNumber:]))
			}
			hasParamUpdate = func(b *builder, unread []string) {
				noParamUpdate(b, unread)
				var (
					name     = p.param.name(paramNumber)
					param    = p.param.params[paramNumber]
					expected = param.setter.expected()
					paramRow = b.nextEmptyRow().
							onSlashError(fmt.Sprintf(
							"the %s is not specified for %s, but is specified for the first %s",
							param,
							p.param,
							p.parameterName,
						)).
						onSpaceError(invalidTokenMessage(name, expected, scanner.Space)).
						onUnknownError(invalidTokenMessage(name, expected, scanner.Unknown)).
						onCommentError(impossibleTokenMessage(name, scanner.Comment))
				)
				param.baseUpdate(paramRow, b.nextState(), unread[paramNumber:])
			}
			param    = p.param.params[paramNumber]
			expected = param.setter.expected()
			paramRow = b.nextEmptyRow().
					onSpaceError(invalidTokenMessage(name, expected, scanner.Space)).
					onUnknownError(invalidTokenMessage(name, expected, scanner.Unknown)).
					onCommentError(impossibleTokenMessage(name, scanner.Comment))
		)
		param.baseUpdate(paramRow, b.nextState(), append([]string{name}, sliceNames[1:]...))
		p.updateRecursive(hasParamUpdate, b, paramNumber+1, false)
		if paramNumber != len(p.param.params)-1 {
			paramRow.onSlash(b.nextState())
			p.updateRecursive(noParamUpdate, b, paramNumber+1, true)
		} else {
			paramRow.onSlashError(invalidTokenMessage(name, expected, scanner.Slash))
		}
		if !lastSlash {
			delimiterRow.onSpace(b.nextState())
		}
	} else {
		b.waitSpace(delimiterBetween(sliceNames[0], sliceNames[1]), sliceNames[1:])
	}
	if !lastSlash {
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

func (p *structSliceParameter) update(b *builder) {
	p.param.changeName(p.name(0))
	p.param.updateRequired(b, p.names())
	p.updateRecursive(
		func(b *builder, unread []string) {
			p.param.updateRequired(b, unread)
		},
		b,
		p.param.requiredCount,
		false,
	)
}

func newStructSliceParameter(name string, min int, param *structParameter) *structSliceParameter {
	return &structSliceParameter{
		parameterName: parameterName(name),
		min:           min,
		param:         param,
	}
}

type stateAction struct {
	state  stateType
	action action
}

type rowBuilder struct {
	stateActionRow [scanner.TokensCount]stateAction
	errorsRow      [scanner.TokensCount]string
}

func (b *rowBuilder) onToken(t scanner.TokenType, s stateType, a action) *rowBuilder {
	b.stateActionRow[t] = stateAction{
		state:  s,
		action: a,
	}
	b.errorsRow[t] = noErrorMessage
	return b
}

func (b *rowBuilder) onWord(s stateType, a action) *rowBuilder { return b.onToken(scanner.Word, s, a) }

func (b *rowBuilder) onInteger(s stateType, a action) *rowBuilder {
	return b.onToken(scanner.Integer, s, a)
}

func (b *rowBuilder) onFloat(s stateType, a action) *rowBuilder {
	return b.onToken(scanner.Float, s, a)
}

func (b *rowBuilder) onSlash(s stateType) *rowBuilder { return b.onToken(scanner.Slash, s, nil) }

func (b *rowBuilder) onSpace(s stateType) *rowBuilder { return b.onToken(scanner.Space, s, nil) }

func (b *rowBuilder) onEnd() *rowBuilder {
	return b.onToken(scanner.EOL, start, nil).onToken(scanner.EOF, start, nil)
}

func (b *rowBuilder) onTokenError(t scanner.TokenType, message string) *rowBuilder {
	b.stateActionRow[t] = stateAction{
		state:  err,
		action: nil,
	}
	b.errorsRow[t] = message
	return b
}

func (b *rowBuilder) onWordError(message string) *rowBuilder {
	return b.onTokenError(scanner.Word, message)
}

func (b *rowBuilder) onIntegerError(message string) *rowBuilder {
	return b.onTokenError(scanner.Integer, message)
}

func (b *rowBuilder) onFloatError(message string) *rowBuilder {
	return b.onTokenError(scanner.Float, message)
}

func (b *rowBuilder) onSlashError(message string) *rowBuilder {
	return b.onTokenError(scanner.Slash, message)
}

func (b *rowBuilder) onSpaceError(message string) *rowBuilder {
	return b.onTokenError(scanner.Space, message)
}

func (b *rowBuilder) onEndError(message string) *rowBuilder {
	return b.onTokenError(scanner.EOL, message).onTokenError(scanner.EOF, message)
}

func (b *rowBuilder) onUnknownError(message string) *rowBuilder {
	return b.onTokenError(scanner.Unknown, message)
}

func (b *rowBuilder) onCommentError(message string) *rowBuilder {
	return b.onTokenError(scanner.Comment, message)
}

func newRowBuilder() *rowBuilder { return &rowBuilder{} }

type builder struct {
	value        reflect.Value
	valueType    ElementType
	params       []parameter
	paramNames   []string
	position     int
	builders     []*rowBuilder
	needFinalize bool
}

func (b *builder) createBoolParameter() {
	b.params = make([]parameter, 1)
	b.paramNames = make([]string, 1)
	var name = fmt.Sprintf("%s parameter", b.valueType)
	b.params[0] = newBaseParameter(name, newBoolSetter(name))
	b.paramNames[0] = name
}

func readName(field *reflect.StructField) string {
	if name, ok := field.Tag.Lookup("name"); ok {
		return name
	} else {
		return field.Name
	}
}

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

func requireNoOptional(tags reflect.StructTag, typeName string) {
	if _, ok := tags.Lookup("optional"); ok {
		panic(fmt.Sprintf("the optional tag cannot be set for a %s field", typeName))
	}
}

func requireNoDelimiter(tags reflect.StructTag, typeName string) {
	if _, ok := tags.Lookup("delimiter"); ok {
		panic(fmt.Sprintf("the delimiter tag cannot be set for a %s field", typeName))
	}
}

func requireNoMin(tags reflect.StructTag, typeName string) {
	if _, ok := tags.Lookup("min"); ok {
		panic(fmt.Sprintf("the min tag cannot be set for a %s field", typeName))
	}
}

func requireWasNotOptional(wasOptional bool) {
	if wasOptional {
		panic("an optional field cannot be followed by a required field")
	}
}

func createNestedStructParameter(
	name string,
	delimiter scanner.TokenType,
	t reflect.Type,
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

func (b *builder) createStructParameters() {
	var t = b.value.Type()
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

func (b *builder) nextState() stateType { return stateType(len(b.builders)) }

func (b *builder) nextEmptyRow() *rowBuilder {
	var rb = newRowBuilder()
	b.builders = append(b.builders, rb)
	return rb
}

func (b *builder) nextParameterRow(name string, expected scanner.TokenType) *rowBuilder {
	return b.nextEmptyRow().
		onSlashError(invalidTokenMessage(name, expected, scanner.Slash)).
		onSpaceError(impossibleTokenMessage(name, scanner.Space)).
		onUnknownError(invalidTokenMessage(name, expected, scanner.Unknown)).
		onCommentError(impossibleTokenMessage(name, scanner.Comment))
}

func (b *builder) nextDelimiterRow(name string) *rowBuilder {
	return b.nextEmptyRow().
		onWordError(impossibleTokenMessage(name, scanner.Word)).
		onIntegerError(impossibleTokenMessage(name, scanner.Integer)).
		onFloatError(impossibleTokenMessage(name, scanner.Float)).
		onUnknownError(impossibleTokenMessage(name, scanner.Unknown)).
		onCommentError(impossibleTokenMessage(name, scanner.Comment))
}

func (b *builder) getUnread() []string {
	if b.position < len(b.paramNames) {
		return b.paramNames[b.position:]
	}
	return []string{}
}

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
	m.actions[first] = func(token string, element reflect.Value) error {
		m.clear()
		return nil
	}
	for i, rb := range b.builders {
		for j, sa := range rb.stateActionRow {
			matrixRow[j] = sa.state
			if m.actions[sa.state] == nil {
				m.actions[sa.state] = sa.action
			} else if sa.action != nil {
				panic(fmt.Sprintf("two actions are specified when transitioning to the same state: %d", sa.state))
			}
		}
		m.matrix[i] = matrixRow
		m.errors[i] = rb.errorsRow
	}
	for i := 0; i < len(m.actions); i++ {
		if m.actions[i] == nil {
			m.actions[i] = func(token string, element reflect.Value) error { return nil }
		}
	}
	return m
}

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

func (b *builder) build() *finiteStateMachine {
	b.initialize()
	var param parameter
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

func newBuilder(elementType ElementType, t reflect.Type) *builder {
	var b = &builder{
		value:        reflect.New(t).Elem(),
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
