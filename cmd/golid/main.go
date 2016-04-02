package main

import (
	"fmt"
	"os"

	"github.com/refola/golid/parse"
)

const usage = `Usage: golid file1.gol [file2.gol [...]]

Converts Golid code into Go, producing gol_file1.gol, etc.
`

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(usage)
		os.Exit(1)
	}
	for _, file := range args {
		err := parse.Convert(file)
		if err != nil {
			fmt.Printf("Error converting %s: %s\n", file, err)
		}
	}
}
