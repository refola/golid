package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/refola/piklisp_go/parse"
)

const usage = `Usage: piklisp file1.plgo [file2.plgo [...]]

Converts Piklisp Go code into Go.`

func e(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

const plgo = ".plgo" // praenthesized version of Piklisp-Go
const gol = ".gol"   // less parens extension

// Convert a Piklisp file into Go. "srfi49" means to use parentheses
// reductions according to Scheme Request For Implementation number
// 49.
func convert(in_name, ext string, srfi49 bool) {
	lisp_bytes, err := ioutil.ReadFile(in_name)
	e(err)
	lisp_text := string(lisp_bytes)
	parseFn := parse.ParenString
	if srfi49 {
		parseFn = parse.Srfi49String
	}
	lisp_parse, err := parseFn(lisp_text)
	e(err)
	go_text := lisp_parse.GoString()
	out_name := in_name[:len(in_name)-len(ext)] + ".go"
	err = ioutil.WriteFile(out_name, []byte(go_text), os.ModePerm)
	e(err)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(usage)
		os.Exit(1)
	}
	for _, input_name := range args {
		switch {
		case strings.HasSuffix(input_name, plgo):
			convert(input_name, plgo, false)
		case strings.HasSuffix(input_name, gol):
			convert(input_name, gol, true)
		default:
			e(errors.New("Invalid filename: " + input_name))
		}
	}
}
