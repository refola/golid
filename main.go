package main

import "internal/parse"

func main() {
	test_program := "(import \"fmt\")\n(fmt.Println \"Hello world!\")"
	parsed := parse.ParenString(test_program)
	fmt.Println(parsed)
}
