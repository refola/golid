package main

import (
	"fmt"
	"github.com/refola/piklisp_go/parse"
)

const test_program = // This comment saves formatting from being eaten by gofmt.
`(package main)
(import "fmt")
(func main () ()
	(fmt.Println "Hello piklisp_go!"))`

func main() {
	parsed, _ := parse.ParenString(test_program)
	fmt.Println("Lisp code loaded!")
	fmt.Println("\nHere's the Lisp version, parsed and unparsed.")
	fmt.Println(parsed)
	fmt.Println("\nHere's the automatic ('compiled') Go translation.")
	fmt.Println(parsed.GoString())
}
