package parser

import (
	"fmt"
	"os"
)

// Reads all vertices from a file containing errors and an unsupported format.
// Check the testdata/vertices_output.txt file for information about errors and warnings!
func ExampleParser_Next_vertices() {
	input, err := os.Open("testdata/vertices.obj")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = input.Close(); err != nil {
			panic(err)
		}
	}()
	output, err := os.Create("testdata/vertices_output.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = output.Close(); err != nil {
			panic(err)
		}
	}()
	var parser = NewParser(input)
	parser.Output(output)
	var elementType, element = parser.Next()
	for elementType != EndOfFile {
		if elementType == Vertex {
			fmt.Printf("%s : %v\n", elementType, element)
		} else {
			fmt.Fprintf(output, "[INFO] unnecessary element: %s\n", elementType)
		}
		elementType, element = parser.Next()
	}
	// Output:
	//vertex : {-0.046146 0.050437 0.002961 0}
	//vertex : {-0.045498 0.049687 0.001989 0}
	//vertex : {-0.045306 0.049655 0.002956 3434}
	//vertex : {-0.045935 0.050494 0.003832 0}
	//vertex : {-0.044743 0.048768 0.002943 0}
	//vertex : {-0.044832 0.048663 0.001729 0}
	//vertex : {-0.047369 0.051618 0.004211 0}
	//vertex : {-0.044734 0.04789 0.002286 0}
	//vertex : {-0.045207 0.050247 0.004572 0}
	//vertex : {-0.046589 0.05193 0.006586 0}
	//vertex : {-0.044529 0.047892 0.003273 0}
}

// Reads all faces from a file containing errors and an unsupported format.
// Check the testdata/faces_output.txt file for information about errors and warnings!
func ExampleParser_Next_faces() {
	input, err := os.Open("testdata/faces.obj")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = input.Close(); err != nil {
			panic(err)
		}
	}()
	output, err := os.Create("testdata/faces_output.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = output.Close(); err != nil {
			panic(err)
		}
	}()
	var parser = NewParser(input)
	parser.Output(output)
	var elementType, element = parser.Next()
	for elementType != EndOfFile {
		if elementType == Face {
			fmt.Printf("%s : %v\n", elementType, element)
		}
		elementType, element = parser.Next()
	}
	// Output:
	//face : {[{1 1 1} {2 2 2} {3 3 3} {4 4 4}]}
	//face : {[{4 4 4} {3 3 3} {2 2 2}]}
	//face : {[{4 4 4} {2 2 2} {7 7 7}]}
	//face : {[{5 5 5} {3 3 3} {8 8 8}]}
	//face : {[{1 0 0} {1 0 0} {1 0 0}]}
	//face : {[{5 5 5} {9 9 9} {1 1 1}]}
	//face : {[{8 8 8} {3 3 3} {6 6 6}]}
	//face : {[{6 6 6} {4 4 4} {10 10 10}]}
	//face : {[{7 7 7} {2 2 2} {11 11 11}]}
	//face : {[{10 10 10} {4 4 4} {7 7 7}]}
	//face : {[{12 12 12} {5 5 5} {8 8 8}]}
	//face : {[{9 9 9} {5 5 5} {13 13 13}]}
	//face : {[{6 6 6} {14 14 14} {8 8 8}]}
	//face : {[{16 16 16} {7 7 7} {11 11 11}]}
	//face : {[{10 10 10} {7 7 7} {16 16 16}]}
	//face : {[{12 12 12} {17 17 17} {5 5 5}]}
	//face : {[{18 18 18} {12 12 12} {8 8 8}]}
	//face : {[{5 5 5} {17 17 17} {13 13 13}]}
	//face : {[{9 9 9} {13 13 13} {19 19 19}]}
	//face : {[{6 6 6} {15 15 15} {14 14 14}]}
	//face : {[{18 18 18} {8 8 8} {14 14 14}]}
	//face : {[{15 15 15} {10 10 10} {20 20 20}]}
	//face : {[{26 26 26} {14 14 14} {15 15 15}]}
	//face : {[{26 26 26} {18 18 18} {14 14 14}]}
	//face : {[{20 0 20} {10 0 10} {21 0 21}]}
	//face : {[{21 21 21} {16 16 16} {28 28 28}]}
	//face : {[{12 12 12} {23 23 23} {22 22 22}]}
	//face : {[{17 17 17} {22 22 22} {29 29 29}]}
	//face : {[{23 23 23} {18 18 18} {26 26 26}]}
}
