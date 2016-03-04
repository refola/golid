// automatically test converting several files from Piklisp Go to Go

package parse

import (
	"io/ioutil"
	"runtime"
	"testing"
)

// check that each Piklisp file converts successfully to Go, without
// crashing
func TestConversions(t *testing.T) {
	// Wrapper function to avoid panics stopping the test
	parsable := func(fn string) (ret bool) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Error processing %s:\n%s\n", fn, r)
				trace := make([]byte, 1e4)
				i := runtime.Stack(trace, false)
				t.Errorf("Stack trace:\n%s\n", trace[:i])
			}
		}()
		err := Convert(fn, false)
		if err != nil {
			t.Errorf("Error processing %s:\n%s", fn, err)
			return false
		}
		return true
	}

	testDirs := []string{"srfi49", "classic"}
	prefix := "../tests"
	failures := []string{}
	for _, dir := range testDirs {
		dir = prefix + "/" + dir
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Fatal("Could not open tests folder:", err)
		}
		failcount := 0
		for _, f := range files {
			filename := dir + "/" + f.Name()
			if !parsable(filename) {
				t.Errorf("Failed parsing %s.\n", filename)
				failcount++
				failures = append(failures, filename)
			}
		}
		t.Errorf("Failed to parse %d/%d test files in %s.", failcount, len(files), dir)
	}
	if len(failures) > 0 {
		t.Errorf("Failed files: %v\n", failures)
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
