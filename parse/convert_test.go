// automatically test converting several files from Piklisp Go to Go

package parse

import (
	"io/ioutil"
	"runtime"
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
	lisp, err := ReadPiklisp(fn)
	if err != nil {
		t.Errorf("Error processing %s:\n%s", fn, err)
		return false
	}
	lisp.GoString()
	return true
}

// check that each Piklisp file converts successfully to Go, without
// crashing
func TestConversions(t *testing.T) {
	testdirs := []string{"srfi49", "classic"}
	prefix := "../tests"
	failures := []string{}
	for _, dir := range testdirs {
		dir = prefix + "/" + dir
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Fatal("could not open tests folder:", err)
		}
		failcount := 0
		for _, f := range files {
			filename := dir + "/" + f.Name()
			if !convertable(t, filename) {
				t.Errorf("failed parsing %s.\n", filename)
				failcount++
				failures = append(failures, filename)
			}
		}
		if failcount > 0 {
			t.Errorf("failed to parse %d/%d test files in %s.", failcount, len(files), dir)
		}
	}
	if len(failures) > 0 {
		t.Errorf("failed files: %v\n", failures)
	}
}

// Check that corresponding .plgo (classic Lisp syntax) and .gol
// (SRFI#49 syntax) files parse the same.
func TestSrfiClassicEquality(t *testing.T) {
	prefix := "../tests"
	srfi := prefix + "/srfi49"
	classic := prefix + "/classic"
	failures := []string{}

	srfiFiles, err := ioutil.ReadDir(srfi)
	if err != nil {
		t.Fatal("Could not open tests folder:", err)
	}
	failCount := 0
	for _, f := range srfiFiles {
		name := f.Name()
		srfiFile := srfi + "/" + name
		classicFile := classic + "/" + name[:len(name)-len(".gol")] + ".plgo"
		srfiParse, err1 := ReadPiklisp(srfiFile)
		classicParse, err2 := ReadPiklisp(classicFile)
		if err1 == nil && err2 == nil {
			s1, s2 := srfiParse.String(), classicParse.String()
			if s1 != s2 {
				failCount++
				failures = append(failures, srfiFile)
				t.Errorf("Got different parse results for the same program!\n%s parsed to: %s\nand %s parsed to: %s\n", srfiFile, s1, classicFile, s2)
			}
		}
	}
	if len(failures) > 0 {
		t.Errorf("Got inconsistent parses with these %v files: %v\n", failCount, failures)
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
