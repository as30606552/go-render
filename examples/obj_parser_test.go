package examples

import (
	"computer_graphics/obj"
	"computer_graphics/obj/constants"
	"fmt"
	"os"
)

// Reads all elements from a file containing errors and an unsupported format.
// Check the output/obj_parser_output.txt file for information about errors and warnings!
func ExampleParser() {
	var err error
	var input *os.File
	var output *os.File
	input, err = os.Open("test_data/test.obj")
	if err != nil {
		panic(err)
	}
	defer input.Close()
	output, err = os.Create("output/obj_parser_output.txt")
	if err != nil {
		panic(err)
	}
	defer output.Close()
	var parser = obj.NewParser(input)
	parser.Output = output
	var elementType constants.ElementType
	var element interface{}
	elementType, element = parser.Next()
	for elementType != constants.EndOfFile {
		fmt.Println(elementType.Name(), " : ", element)
		elementType, element = parser.Next()
	}
	// Output:
	//vertex  :  {-0.046146 0.050437 0.002961}
	//vertex  :  {-0.045498 0.049687 0.001989}
	//vertex  :  {-0.045306 0.049655 0.002956}
	//vertex  :  {-0.045935 0.050494 0.003832}
	//vertex  :  {-0.044743 0.048768 0.002943}
	//vertex  :  {-0.044832 0.048663 0.001729}
	//vertex  :  {-0.047369 0.051618 0.004211}
	//vertex  :  {-0.044734 0.04789 0.002286}
	//vertex  :  {-0.045207 0.050247 0.004572}
	//vertex  :  {-0.046589 0.05193 0.006586}
	//vertex  :  {-0.044529 0.047892 0.003273}
}
