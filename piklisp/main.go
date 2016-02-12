package main

import (
	"fmt"
	"os"

	"github.com/refola/piklisp_go/parse"
)

const usage = `Usage: piklisp file1.{plgo|gol} [file2.{plgo|gol} [...]]

Converts Piklisp Go code into Go.

Note that .plgo files use "normal" Lisp syntax while .gol files use
SRFI#49 less-parenthetical syntax (which is isomorphic to "normal"
parenthetical syntax).`

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(usage)
		os.Exit(1)
	}
	for _, file := range args {
		err := parse.Convert(file, true)
		if err != nil {
			fmt.Printf("Error converting %s: %s\n", file, err)
		}
	}
}
