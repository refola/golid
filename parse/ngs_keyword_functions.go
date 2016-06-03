// ngs_keyword_functions.go

// This file contains keyword-specific functions ultimately used in
// Node.GoString(). These functions should all be passed the Node
// whose content is the applicable keyword.

package parse

// Convert Golid "(break)", "(break label)", "(continue)", and
// "(continue label)" statements into Go.
var nkw_break func(*Node) string = nu_raw_content_space

// Convert an import Node into a Go import command.
func nkw_import(keywordNode *Node) string {
	out := "import ("
	for n := keywordNode.next; n != nil; n = n.next {
		out += n.content + "; "
	}
	out += ")"
	return out
}

// Convert a package declaration to Go.
func nkw_package(keywordNode *Node) string {
	return nu_raw_content(keywordNode, " ")
}

// Convert Golid "(myVar value)" and "(myVar type value)" expressions
// (which are to the right of var (and const) expressions) into
// corresponding Go "myVar = value" and "myVar type = vaule"
// expressions.
func nkw_var_post_kw(varNameNode *Node) string {
	n := varNameNode
	out := n.content
	n = n.next
	if n.next == nil { // "myVar value" case
		out += " = " + nc_value(n)
	} else { // "myVar type value" case
		out += " " + n.content + " = " + nc_value(n.next)
	}
	return out
}

// Convert a var Node into a Go var declaration.  TODO: This is
// currently broken. Here's how it should convert things:
// (var myVar value)→"var myVar = value"
// (var myVar type value)→"var myVar type = value"
// (var (myVar1 value) (myVar2 type value))
// → "var (
//        myVar1 = value
//        myVar2 type = value
// )"
func nkw_var(keywordNode *Node) string {
	// "var" (or "const")
	n := keywordNode
	out := n.content
	n = n.next
	// if it's a single-var declaration
	if n.content != "" {
		out += " " + nkw_var_post_kw(n) + "\n"
	} else { // if it's a multi-var declaration
		out += " (\n"
		for n != nil {
			out += nkw_var_post_kw(n.first) + "\n"
			n = n.next
		}
		out += ")"
	}
	return out
}

// Convert a function Node into a Go function declaration.
func nkw_func(keywordNode *Node) string {
	// "func"
	n := keywordNode
	out := n.content
	// function name
	n = n.next
	out += " " + n.content
	// function args
	// TODO: This part's broken for functions that take multiple args.
	n = n.next
	out += "(" + nu_raw_content(n.first, " ") + ")"
	// function return types
	// TODO: This part's broken for functions that name their return values in the signature.
	n = n.next
	out += "(" + nu_raw_content(n.first, " ") + ")"
	// function body
	out += "{\n"
	for n = n.next; n != nil; n = n.next {
		out += nc_action(n) + "\n"
	}
	// closing brace
	out += "}\n"
	return out
}

// return text representing an "if condition { stuff() ... }" block
func nkw_if(keywordNode *Node) string {
	n := keywordNode
	// "if"
	out := n.content + " "
	// first case
	n=n.next
	out += nc_value(n.first) + " {\n"
	out += nu_process_many(n.first.next, nc_action)
	// other cases
	for n = n.next; n != nil; n = n.next {
		if n.first.content=="else"{
			out+="} else {\n"
		}else{
			out += "} else if " + nc_value(n.first) + " {\n"
		}
		out += nu_process_many(n.first.next, nc_action)
	}
	// final closing
	out += "}\n"
	return out
}

// return text representing a "for pre-statement; condition; post-statement { stuff() ... }" block of any type
func nkw_for(keywordNode *Node) string {
	// "for"
	out := keywordNode.content + " "
	n := keywordNode.next
	// get header
	switch {
	case n.content != "": // Golid for loops must paren the control clause.
		panic("Invalid 'for' control clause: \"" + n.String() + "\"!")
	case n.first == nil: // "()" case ('infinite' loop)
		out += "{\n"
	case n.first.content == "range": // "(range ...)" case
		panic("nodeForCase: 'range' isn't implemented!")
	case n.first.content != "": // "(condition)" case ('while' loop)
		out += nc_value(n) + "{\n"
	case n.first.first != nil: // "((pre) (cond) (post))" case ('for' loop)
		out += ns_assign(n.first.first) + "; " + nc_value(n.first.next) + "; " + nc_action(n.first.next.next) + " {\n"
	default:
		panic("nodeForCase: Unhandled case!")
	}
	// go through body
	for n = n.next; n != nil; n = n.next {
		out += nc_action(n) + "\n"
	}
	// end brace
	out += "}\n"
	return out
}

// return text representing a "return [values ...]" statement
func nkw_return(keywordNode *Node) string {
	n := keywordNode
	// "return"
	out := n.content
	// args
	n = n.next
	for n != nil {
		out += " " + nc_value(n) + ","
		n = n.next
	}
	if out != keywordNode.content {
		out = out[:len(out)-1]
	}
	return out
}

// return text representing a "switch var { case value: ... case val1 val2: ... }" block
func nkw_switch(keywordNode *Node) string {
	n:=keywordNode
	// "switch" and value expression to switch on
	out := n.content + " " + nc_value(n.next) + " {\n"
	// loop thru cases
	for n = n.next.next; n != nil; n = n.next {
		// "case" statement
		switch c := n.first.content; c {
		case "":
			out += "case " + nu_raw_content(n.first.first, ", ") + ":\n"
		case "default":
			out += "default:\n"
		default:
			out += "case " + c + ":\n"
		}
		// body of case
		for inner := n.first.next; inner != nil; inner = inner.next {
			out += nc_action(inner) + "\n"
		}
	}
	// end brace
	out += "}\n"
	return out
}
