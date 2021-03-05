package parser

import (
	"computer_graphics/obj/scanner"
	"computer_graphics/obj/types"
	"reflect"
	"testing"
)

func TestBuildParameters(t *testing.T) {
	var s = &struct {
		A           int `name:"integer" optional:"false"`
		Field2      float64
		Field3      string `name:"string"`
		StructField struct {
			Sf1 int `name:"inner struct field 1"`
			Sf2 int `name:"inner struct field 2"`
		} `delimiter:"space" name:"inner struct"`
		SliceField []struct {
			Ssf1 float64 `name:"inner slice field 1"`
			Ssf2 int     `name:"inner slice field 2" optional:"true"`
			Ssf3 float64 `name:"inner slice field 3" optional:"true"`
		} `delimiter:"slash" min:"3" name:"slice of struct"`
	}{}
	var params = buildParameters(reflect.ValueOf(s).Elem())
	if len(params) != 5 {
		t.Fatalf("Invalid number of parameters, got: %d, want: %d", len(params), 5)
	}
	var intParam = params[0].(*intParameter)
	if intParam.name != "integer" {
		t.Errorf("The parameter name was parsed incorrectly, got: %s, want: %s", intParam.name, "integer")
	}
	if intParam.optional != false {
		t.Errorf("The parameter tag optional was parsed incorrectly, got: %t, want: %t", intParam.optional, false)
	}
	if intParam.value.Kind() != reflect.Int {
		t.Errorf("The parameter refers to an invalid value type, got: %s, want: %s", intParam.value.Kind(), reflect.Int)
	}
	if !intParam.value.CanSet() {
		t.Errorf("The parameter %s cannot change its field", intParam.name)
	}
	var floatParam = params[1].(*floatParameter)
	if floatParam.name != "Field2" {
		t.Errorf("The parameter name was parsed incorrectly, got: %s, want: %s", floatParam.name, "Field2")
	}
	if floatParam.optional != false {
		t.Errorf("The parameter tag optional was parsed incorrectly, got: %t, want: %t", floatParam.optional, false)
	}
	if floatParam.value.Kind() != reflect.Float64 {
		t.Errorf("The parameter refers to an invalid value type, got: %s, want: %s", floatParam.value.Kind(), reflect.Float64)
	}
	if !floatParam.value.CanSet() {
		t.Errorf("The parameter %s cannot change its field", floatParam.String())
	}
	var stringParam = params[2].(*stringParameter)
	if stringParam.name != "string" {
		t.Errorf("The parameter name was parsed incorrectly, got: %s, want: %s", stringParam.name, "string")
	}
	if stringParam.value.Kind() != reflect.String {
		t.Errorf("The parameter refers to an invalid value type, got: %s, want: %s", stringParam.value.Kind(), reflect.String)
	}
	if !stringParam.value.CanSet() {
		t.Errorf("The parameter %s cannot change its field", stringParam.String())
	}
	// TODO add validation of the structure and slice parameters builder
}

func TestBuildParser_vertex(t *testing.T) {
	var (
		parser = buildParser(Vertex, types.Vertex{})
		got    = parser.(*finiteStateMachine).matrix
		want   = [][scanner.TokensCount]stateType{
			{1, 1, 1, 1, 2, 1, 1, 1, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 3, 3, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 4, 1, 1, 1, 1},
			{1, 5, 5, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 6, 1, 1, 1, 1},
			{1, 7, 7, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 8, 0, 0, 1, 1},
			{1, 9, 9, 1, 1, 0, 0, 1, 1},
			{1, 1, 1, 1, 10, 0, 0, 1, 1},
			{1, 1, 1, 1, 1, 0, 0, 1, 1},
		}
		gotDim  = len(got)
		wantDim = len(want)
	)
	if gotDim != wantDim {
		t.Fatalf("Incorrect dimension of the matrix, got: %d, want: %d", gotDim, wantDim)
	}
	var correct = true
	for i := 0; i < gotDim; i++ {
		for j := 0; j < scanner.TokensCount; j++ {
			if got[i][j] != want[i][j] {
				t.Errorf("Invalid matrix element (%d, %d), got: %d, want: %d", i, j, got[i][j], want[i][j])
				correct = false
			}
		}
	}
	if !correct {
		t.Log("got: ", got)
		t.Log("want:", want)
	}
}

// TODO add a test for an element containing a slice of structures (for example, for Face)
