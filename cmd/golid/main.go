package main

import (
	"fmt"
	"os"

	"github.com/refola/golid/parse"
)

const usage = `Usage: golid [-reparse] file1.gol [file2.gol [...]]

Converts Golid code into Go, producing gol_file1.gol, etc.

If the first argument is "-reparse", then print the Lisp parse tree to
standart output instead.
`

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(usage)
		os.Exit(1)
	} else if args[0] == "-reparse" {
		for _, file := range args[1:] {
			parsed, err := parse.ReadGolid(file)
			if err != nil {
				fmt.Printf("Error converting %s: %s\n", file, err)
			}
			reparsed := parsed.String()
			reparsed = reparsed[1 : len(reparsed)-1] // remove wrapping "()"
			fmt.Println(reparsed)
		}
	} else {
		for _, file := range args {
			err := parse.Convert(file)
			if err != nil {
				fmt.Printf("Error converting %s: %s\n", file, err)
			}
		}
	}
}
