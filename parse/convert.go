// Convenience functions to convert Piklisp files into Go files

package parse

import (
	"fmt"
	"io/ioutil"
	"os"
)

const plgo = "plgo" // praenthesized version of Piklisp-Go
const gol = "gol"   // less parens extension

// split filename into base and extension, returning original name if no '.' found
func split_filename(fn string) (string, string) {
	for i := len(fn) - 1; i >= 0; i-- {
		if fn[i] == '.' {
			return fn[:i-1], fn[i+1:]
		}
	}
	return fn, ""
}

// Read a file into a Piklisp syntax tree
func ReadPiklisp(plfile string) (Expression, error) {
	_, ext := split_filename(plfile)
	parseFn := ParenString
	switch ext {
	case plgo: // valid case, but we already have its parseFn set
	case gol:
		parseFn = Srfi49String
	default:
		return nil, fmt.Errorf("File %s has non-Piklisp extension %s", plfile, ext)
	}
	lisp_bytes, err := ioutil.ReadFile(plfile)
	if err != nil {
		return nil, err
	}
	lisp_text := string(lisp_bytes)
	return parseFn(lisp_text)
}

// Convert a Piklisp file into Go. It uses normal Lisp syntax if the extension is .plgo and SRFI#49 if the extension is .gol.
func Convert(plfile string, write_result bool) error {
	parsed, err := ReadPiklisp(plfile)
	if err != nil {
		return err
	}
	go_text := parsed.GoString()
	name, ext := split_filename(plfile)
	if write_result {
		gofile := fmt.Sprintf("%s_%s.go", ext, name)
		err = ioutil.WriteFile(gofile, []byte(go_text), os.ModePerm)
	}
	return err
}
