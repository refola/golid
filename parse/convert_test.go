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

	prefix := "../tests"
	files, err := ioutil.ReadDir(prefix)
	if err != nil {
		t.Fatal("Could not open tests folder:", err)
	}
	failcount := 0
	for _, f := range files {
		filename := prefix + "/" + f.Name()
		if !parsable(filename) {
			t.Errorf("Failed parsing %s.\n", filename)
			failcount++
		}
	}
	t.Logf("Failed to parse %d/%d test files.", failcount, len(files))
}
