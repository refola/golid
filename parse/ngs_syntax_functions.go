// ngs_syntax_functions.go

// This file contains syntax-specific functions for
// Node.GoString(). Lake the functions in ngs_keyword_functions.go,
// these functions should be passed the Node whose content is started
// by the applicable thing, only in this case it's a syntax instead of
// a keyword.

package parse

// Process an assignment, starting from the first Node.
func ns_assign(first *Node) string {
	// Go LHS and assignment operator
	out := first.next.content + first.content
	// RHS
	// TODO: Properly parse as values.
	out += nc_value(first.next.next)
	return out
}

// Convert a function call into Go
func ns_funcall(first *Node) string {
	out := first.content + "("
	for n := first.next; n != nil; n = n.next {
		out += nc_value(n) + ", "
	}
	offset := len(out) - len(", ")
	if out[offset:] == ", " {
		out = out[:offset]
	}
	out += ")"
	return out
}

// Convert a Lisp math function call into Go form.
func ns_math(first *Node) string {
	op := first.content
	n := first.next
	lhs := n
	n = n.next
	rhs := n
	return "(" + nc_value(lhs) + " " + op + " " + nc_value(rhs) + ")"
}
