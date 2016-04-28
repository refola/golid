// ngs_contexts.go

// This file contains the functions that handle particular contexts of
// Lisp→Go syntax conversion. All calls to these functions must pass
// the Node of interest, _not_ its first child Node.

package parse

// Process a top-level Node
func nc_top(n *Node) string {
	first := n.first
	var f func(*Node) string
	switch first.content {
	case "package":
		f = nkw_package
	case "import":
		f = nkw_import
	case "const", "var":
		f = nkw_var
	case "func":
		f = nkw_func
	default:
		panic("Unknown top-level node type: " + first.content)
	}
	return f(first)
}

// Process an action Node
func nc_action(n *Node) string {
	first := n.first
	var f func(*Node) string
	switch first.content {
	case "=", ":=", "+=", "-=", "*=", "/=", "++", "--":
		f = ns_assign
	case "if":
		f = nkw_if
	case "for":
		f = nkw_for
	case "switch":
		f = nkw_switch
	case "select":
		panic("nodeProcessAction: select not implemented")
	case "break", "continue", "return": // TODO: return should use commas between items
		panic("nodeProcessAction: '" + first.content + "' is not properly implemented")
		f = nu_raw_content_space
	default:
		f = ns_funcall
	}
	return f(first)
}

// Process a value Node
func nc_value(n *Node) string {
	if n.content != "" {
		return n.content
	}
	first := n.first
	var f func(*Node) string
	switch first.content {
	case "+", "-", "*", "/", "==", "!=", ">=", "<=", "<", ">":
		f = ns_math
	default:
		f = ns_funcall
	}
	return f(first)
}
