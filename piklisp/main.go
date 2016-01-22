package main

import "github.com/refola/piklisp_go/parse"

const test_program = // This comment saves formatting from being eaten by gofmt.
`(package main)
(import "fmt")
(func main () ()
	(fmt.Println "Hello piklisp_go!"))`

func main() {
	parsed := parse.ParenString(test_program)
	fmt.Println(parsed)
}
