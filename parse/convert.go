// Convenience functions to convert Piklisp files into Go files

package parse

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const plgo = ".plgo" // praenthesized version of Piklisp-Go
const gol = ".gol"   // less parens extension

// Convert a Piklisp file into Go. It uses normal Lisp syntax if the extension is .plgo and SRFI#49 if the extension is .gol.
func Convert(plfile string, write_result bool) error {
	ext := ""
	parseFn := ParenString
	switch {
	case strings.HasSuffix(plfile, plgo):
		ext = plgo
	case strings.HasSuffix(plfile, gol):
		ext = gol
		parseFn = Srfi49String
	default:
		return fmt.Errorf("Invalid filename: %s", plfile)
	}

	lisp_bytes, err := ioutil.ReadFile(plfile)
	if err != nil {
		return err
	}
	lisp_text := string(lisp_bytes)
	lisp_parse, err := parseFn(lisp_text)
	if err != nil {
		return err
	}
	go_text := lisp_parse.GoString()
	if write_result {
		gofile := fmt.Sprintf("%s_%s.go", plfile[:len(plfile)-len(ext)], ext[1:])
		err = ioutil.WriteFile(gofile, []byte(go_text), os.ModePerm)
	}
	return err
}
