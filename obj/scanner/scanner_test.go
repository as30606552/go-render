package scanner

import (
	"fmt"
	"strings"
)

// Reading the correct data.
func ExampleScanner_Next_correct() {
	var s = NewScanner(strings.NewReader("word 123/-321 0.01"))
	var tokenType, token = s.Next()
	for tokenType != EOF {
		fmt.Printf("%s : '%s'\n", tokenType.String(), token)
		tokenType, token = s.Next()
	}
	// Output:
	//WORD : 'word'
	//SPACE : ' '
	//INT : '123'
	//SLASH : '/'
	//INT : '-321'
	//SPACE : ' '
	//FLOAT : '0.01'
}

// Reading data containing errors.
func ExampleScanner_Next_incorrect() {
	var s = NewScanner(strings.NewReader("invalid&word validWord 123-321 0.0.1"))
	var tokenType, token = s.Next()
	for tokenType != EOF {
		fmt.Printf("%s : '%s'\n", tokenType.String(), token)
		tokenType, token = s.Next()
	}
	// Output:
	//UNKNOWN : 'invalid&word'
	//SPACE : ' '
	//WORD : 'validWord'
	//SPACE : ' '
	//UNKNOWN : '123-321'
	//SPACE : ' '
	//UNKNOWN : '0.0.1'
}
