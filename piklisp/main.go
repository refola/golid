package main

import "github.com/refola/piklisp_go/parse"

func main() {
	test_program := "(package main)\n(import \"fmt\")\n(func main () () (fmt.Println \"Hello piklisp_go!\"))"
	parsed := parse.ParenString(test_program)
	fmt.Println(parsed)
}
