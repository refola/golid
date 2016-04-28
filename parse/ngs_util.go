// ngs_util.go

// This file contains utility functions for the other ngs_*.go files.

package parse

import "fmt"

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
					err = fmt.Errorf("Recovered panic: %v", r)
				}
			}()
			out = f(n)
			return
		}()
		if err != nil {
			panic(fmt.Errorf("Could not process code: %v\n Got error: %v", n, err))
		} else {
			out += result + "\n"
		}
	}
	return out
}

// Generate a list of raw Node contents (only node.content, ignoring
// children), separated by given separator string. WARNING: This break
// recursion.
func nu_raw_content(first *Node, sep string) string {
	out := first.content
	for n := first.next; n != nil; n = n.next {
		out += sep + n.content
	}
	return out
}

// Call nodeRawContents with a space for the separator.
func nu_raw_content_space(first *Node) string {
	return nu_raw_content(first, " ")
}
