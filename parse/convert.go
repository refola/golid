// Convenience functions to convert Piklisp files into Go files

package parse

import (
	"fmt"
	"io/ioutil"
	"os"
)

const plgo = "plgo" // praenthesized version of Piklisp-Go
const gol = "gol"   // less parens extension

// split path into directory, filename, and extension, based on
// right-most '.' and '/' characters, and omitting said separators
func dir_name_ext_split(path string) (string, string, string) {
	dot, slash := 0, 0
	foundDot, foundSlash := false, false
	for i, v := range path {
		switch v {
		case '.':
			dot = i
			foundDot = true
		case '/':
			slash = i
			foundSlash = true
		}
	}
	dir, name, ext := path[:slash], path[slash:dot], path[dot:]
	if foundDot {
		ext = ext[1:]
	}
	if foundSlash {
		name = name[1:]
	}
	return dir, name, ext
}

// Read a file into a Piklisp syntax tree
func ReadPiklisp(plfile string) (Expression, error) {
	_, _, ext := dir_name_ext_split(plfile)
	mode := classic
	switch ext {
	case plgo: // valid case, but we already have its parseFn set
	case gol:
		mode = srfi49
	default:
		return nil, fmt.Errorf("File %s has non-Piklisp extension %s", plfile, ext)
	}
	lispBytes, err := ioutil.ReadFile(plfile)
	if err != nil {
		return nil, err
	}
	lispText := string(lispBytes)
	return parseString(lispText, mode)
}

// Convert a Piklisp file into Go. It uses normal Lisp syntax if the extension is .plgo and SRFI#49 if the extension is .gol.
func Convert(plfile string, write_result bool) error {
	parsed, err := ReadPiklisp(plfile)
	if err != nil {
		return err
	}
	go_text := parsed.GoString()
	dir, name, ext := dir_name_ext_split(plfile)
	if write_result {
		gofile := fmt.Sprintf("%s/%s_%s.go", dir, ext, name)
		err = ioutil.WriteFile(gofile, []byte(go_text), os.ModePerm)
	}
	return err
}
