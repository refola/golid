// Convenience functions to convert Golid files into Go files

package parse

import (
	"fmt"
	"io/ioutil"
)

// split path into directory, filename, and extension, based on
// right-most '.' and '/' characters, and omitting said separators
func dirNameExt(path string) (string, string, string) {
	pos := make(map[rune]int)
	// find the last slash and dot in the string
	for i, v := range path {
		switch v {
		case '/', '.':
			pos[v] = i
		}
	}
	// figure out validity and split string
	slash, foundSlash := pos['/']
	dot, foundDot := pos['.']
	if dot <= slash+1 { // if there was a dot here, then it was in the directory or at the beginning of a "dotfile", e.g., "/foo.bar/baz" or "/foo/bar/.baz"
		dot = len(path)
		foundDot = false
	}
	dir, name, ext := path[:slash], path[slash:dot], path[dot:]
	if foundDot { // trim leading dot from extension
		ext = ext[1:]
	}
	if foundSlash { // trim leading slash from name
		name = name[1:]
	}
	return dir, name, ext
}

// Read a file into a Golid syntax tree
func ReadGolid(golfile string) (Expression, error) {
	_, _, ext := dirNameExt(golfile)
	if ext != "gol" {
		return nil, fmt.Errorf("File %s has non-Golid extension %s", golfile, ext)
	}
	lispBytes, err := ioutil.ReadFile(golfile)
	if err != nil {
		return nil, err
	}
	lispText := string(lispBytes)
	return parseString(lispText)
}

// Convert a Golid file into Go.
func Convert(golfile string) error {
	parsed, err := ReadGolid(golfile)
	if err != nil {
		return err
	}
	go_text := parsed.GoString()
	dir, name, ext := dirNameExt(golfile)
	if dir == "" {
		// make sure that following dir with "/" doesn't change semantics
		dir = "."
	}
	gofile := fmt.Sprintf("%s/%s_%s.go", dir, ext, name)
	err = ioutil.WriteFile(gofile, []byte(go_text), 0644)
	return err
}
