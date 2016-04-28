// automatically test converting several files from Golid to Go

package parse

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"
)

// const DEBUG = true
const DEBUG = false

// Wrapper function to avoid panics stopping the test
func convertable(t *testing.T, fn string) (ret bool) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic: %s\n", r)
			// TODO: Re-enable this after implementing the less ad-hoc GoString functions.
			if DEBUG {
				trace := make([]byte, 1e4)
				i := runtime.Stack(trace, false)
				t.Errorf("Stack trace:\n%s\n", trace[:i])
			}
		}
	}()
	lisp, err := ReadGolid(fn)
	if err != nil {
		t.Errorf("Error processing %s:\n%s", fn, err)
		return false
	}
	lisp.GoString()
	return true
}

// check that each Golid file converts successfully to Go, without
// crashing
func TestConversions(t *testing.T) {
	root := "../tests"
	dirs, err := ioutil.ReadDir(root)
	if err != nil {
		t.Fatal("could not open tests folder:", err)
	}
	failures := []string{} // list of cases that failed
	nFiles := 0            // how many files were tested
	for _, d := range dirs {
		dirName := root + "/" + d.Name()
		t.Logf("Testing '%s'.", dirName)
		dir, err := ioutil.ReadDir(dirName)
		if err != nil {
			t.Errorf("could not open tests folder:", err)
			continue
		}
		for _, f := range dir {
			filename := dirName + "/" + f.Name()
			if !strings.HasSuffix(filename, ".gol") {
				continue
			}
			nFiles++
			if !convertable(t, filename) {
				t.Errorf("failed parsing %s.\n", filename)
				failures = append(failures, filename)
			}
		}
	}
	if t.Failed() {
		t.Errorf("failed to parse %d/%d test files in %s.", len(failures), nFiles, root)
		t.Errorf("failed files: %v\n", failures)
	}
}

// Convert filename to parsed Golid string, returning an error
// string that's invalid Golid on failure.
func fileToParseString(filename string) (ret string) {
	defer func() {
		if r := recover(); r != nil {
			ret = "Failed converting " + filename + " to parse string. Result is:\n" + ret
		}
	}()
	parse, err := ReadGolid(filename)
	if err != nil {
		return "Could not parse " + filename + "."
	}
	ret = parse.String()
	return
}

// Check that corresponding Golid files (i.e., in .gol files in the
// same test folder) parse the same.
func TestParseEquality(t *testing.T) {
	root := "../tests"
	dirs, err := ioutil.ReadDir(root)
	if err != nil {
		t.Fatal("Could not open tests folder:", err)
	}
	failures := []string{} // list of folders that failed
	for _, d := range dirs {
		dirName := root + "/" + d.Name()
		dir, err := ioutil.ReadDir(dirName)
		if err != nil {
			t.Errorf("could not open tests folder:", err)
			continue
		}
		parseString := ""
		firstFile := ""
		for _, f := range dir {
			filename := dirName + "/" + f.Name()
			if !strings.HasSuffix(filename, ".gol") {
				continue
			}
			if firstFile == "" {
				firstFile = filename
			}
			s := fileToParseString(filename)
			if parseString == "" {
				parseString = s
			}
			if s != parseString {
				failures = append(failures, dirName)
				t.Errorf("Got multiple parse results for the same program!\n%s parsed to: %s\nand %s parsed to: %s\n", firstFile, parseString, filename, s)
			}
		}
	}
	if len(failures) > 0 {
		t.Errorf("Got inconsistent parses with these %v programs: %v\n", len(failures), failures)
	}
}

func TestDirNameExt(t *testing.T) {
	cases := map[string][]string{
		"/foo/bar.baz":  {"/foo", "bar", "baz"},
		"foo":           {"", "foo", ""},
		"foo/bar":       {"foo", "bar", ""},
		"/foo.bar/baz":  {"/foo.bar", "baz", ""},
		"/foo/bar/.baz": {"/foo/bar", ".baz", ""},
	}
	failed := 0
	for in, want := range cases {
		t.Logf("Test '%s' → '%v'", in, want)
		a, b, c := dirNameExt(in)
		result := []string{a, b, c}
		for i, v := range result {
			if want[i] != v {
				t.Errorf("Error: Got '%v' instead.", result)
				failed++
				break
			}
		}
	}
	t.Logf("Failed %v/%v tests.", failed, len(cases))
}

func TestNodeProcessValue(t *testing.T) {
	cases := map[string]string{
		"(< n 2)":                         "(n < 2)",
		"(- n 1)":                         "(n - 1)",
		"(fib (- n 1))":                   "fib((n - 1))",
		"(+ (fib (- n 1)) (fib (- n 2)))": "(fib((n - 1)) + fib((n - 2)))",
		"5":     "5",
		`"foo"`: `"foo"`,
	}
	log := ""
	defer func() {
		if log != "" {
			t.Fatal(log)
		}
	}()
	failed := 0
	for in, want := range cases {
		log = fmt.Sprintf("%s → %s", in, want)
		expr, err := parseString(in)
		if err != nil {
			t.Errorf("%s:\nCould not parse '%s'.", log, in)
			failed++
			continue
		}
		// ".first" gets rid of top-level parens added by parseString()
		n := expr.(*Node)
		out := nc_value(n)
		if out != want {
			t.Errorf("%s:\nGot '%s' instead.", log, out)
			failed++
		}
	}
	t.Logf("Failed %v/%v NodeProcessValue() tests.", failed, len(cases))
	log = ""
}
