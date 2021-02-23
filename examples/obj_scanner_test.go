package examples

import (
	"computer_graphics/obj/scanner"
	"fmt"
	"strings"
)

// Reading the correct data.
func ExampleScanner() {
	var s = scanner.NewScanner(strings.NewReader("word 123/-321 0.01"))
	var tokenType scanner.TokenType
	var token string
	tokenType, token = s.Next()
	for tokenType != scanner.EOF {
		fmt.Printf("%s : '%s'\n", tokenType.Name(), token)
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
func ExampleScanner_second() {
	var s = scanner.NewScanner(strings.NewReader("invalid&word validWord 123-321 0.0.1"))
	var tokenType scanner.TokenType
	var token string
	tokenType, token = s.Next()
	for tokenType != scanner.EOF {
		fmt.Printf("%s : '%s'\n", tokenType.Name(), token)
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
