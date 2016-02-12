// automatically test converting several files from Piklisp Go to Go

package parse

import (
	"testing"
	//"os"
	"io/ioutil"
)

// check that each Piklisp file converts successfully to Go, without
// crashing
func TestConversions(t *testing.T) {
	prefix := "../tests"
	files, err := ioutil.ReadDir(prefix)
	if err != nil {
		t.Fatal("Could not open tests folder:", err)
	}
	for _, f := range files {
		filename := prefix + "/" + f.Name()
		t.Logf("Testing %s.\n", filename)
		err := Convert(filename, false)
		if err != nil {
			t.Error(err)
		}
	}
}
