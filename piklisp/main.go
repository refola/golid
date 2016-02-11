package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/refola/piklisp_go/parse"
)

const usage = `Usage: piklisp file1.plgo [file2.plgo [...]]

Converts Piklisp Go code into Go.`

func e(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(usage)
		os.Exit(1)
	}
	for _, input_name := range args {
		if input_name[len(input_name)-len(".plgo"):] != ".plgo" {
			e(errors.New("Invalid filename: " + "input_name"))
		}
		output_name := input_name[:len(input_name)-len(".plgo")] + ".go"
		lisp_bytes, err := ioutil.ReadFile(input_name)
		lisp_text := string(lisp_bytes)
		e(err)
		lisp_parse, err := parse.ParenString(lisp_text)
		e(err)
		go_text := lisp_parse.GoString()
		err = ioutil.WriteFile(output_name, []byte(go_text), os.ModePerm)
		e(err)
	}
}
