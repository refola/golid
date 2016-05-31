// ngs_util.go

// This file contains utility functions for the other ngs_*.go files.

package parse

import (
	"fmt"
	"runtime/debug"
)

// Apply the correct nc_* function to each Node starting from first
// and going until the end of the current level. WARNING: You will
// probably get bad results if you try using this with functions other
// than the nc_* ("node context") functions found in ngs_context.go.
func nu_process_many(first *Node, f func(*Node) string) string {
	out := ""
	for n := first; n != nil; n = n.next {
		result, err := func() (out string, err error) {
			defer func() {
				if r := recover(); r != nil {
					stack := string(debug.Stack())
					err = fmt.Errorf("Recovered panic: %v.\n\nHere's the stack:\n%s", r, stack)
				}
			}()
			out = f(n)
			return
		}()
		if err != nil {
			panic(fmt.Errorf("Could not process code:\n%v\n\nGot error:\n%v", n, err))
		} else {
			out += result + "\n"
		}
	}
	return out
}

// Generate a list of raw Node contents (only node.content, ignoring
// children), separated by given separator string. WARNING: This
// breaks recursion.
func nu_raw_content(first *Node, sep string) string {
	out := ""
	n := first
	for n != nil {
		out += n.content + sep
		n = n.next
	}
	if out == "" {
		return out
	} else {
		return out[:len(out)-len(sep)]
	}
}

// Call nodeRawContents with a space for the separator.
func nu_raw_content_space(first *Node) string {
	return nu_raw_content(first, " ")
}
