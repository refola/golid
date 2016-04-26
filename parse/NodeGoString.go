// convert a Node into Go code

/* TYPES OF GO SYNTAX TO HANDLE

===== CONTEXTS =====

Top-level: Anything that's valid outside of, e.g., a function body.
* package declaration
* import
* top-level consts and vars
* functions

Action: Places that require the program to _do_ something.
* bodies of functions
* bodies of control structures

Value: Places that need something that results in a single value.
* function arguments
* right-hand-side of ":=" and friends
* most args for if/for/switch
* array indices
* several other places

SimpleStmt (see golang.org/ref/spec#SimpleStmt): Anything except for control structures.
* values, channel sends, ++/--, assignments, and short declarations
* found at beginnings of control structures


===== SYNTAXES BY CONTEXT =====

Top-level:
* (first args) → first args
** package
* (first args ...)
→ first (args; ...)
** import
* (first (arg1 ...) (arg2 ...) ...)
→ first ( arg1 ...; arg2 ...; ... )
** const
** var
* (first second (args1 ...) (args2 ...) (args3 ...))
→ first second(args1, ...) (args2, ...) { args3; ...; }
** func

Action:
* (first second args ...)
→ second first args ...
** = := /= *= += -=
* (first args ...)
→ first(args ...)
** function calls
* (first (arg1 ...) (arg2 ...) ...)
→ first arg1 { ... } ??? arg2 { ... } ...
** if/for/switch/select
* (first second)
→ first second
** return
** goto
** label?

Value:
* (first second third)
→ (second first third)
** math expressions
** comparisons
* (first args ...)
→ first(args ...)
** function calls


===== STRATEGY FOR HANDLING THESE CASES =====

* start at top level
* have each function for handling these cases track what types of things it's expecting
* use per-context function maps to choose how to process each node

*/

package parse

import "fmt"

// Process a top-level Node
func nodeProcessTop(n *Node) string {
	first := n.first
	var f func(*Node) string
	switch first.content {
	case "package":
		f = nodeUnparenTwo
	case "import":
		f = nodeImport
	case "const", "var":
		f = nodeConstVar
	case "func":
		f = nodeFunc
	default:
		panic("Unknown top-level node type: " + first.content)
	}
	return f(first)
}

// Process an action Node
func nodeProcessAction(n *Node) string {
	first := n.first
	var f func(*Node) string
	switch first.content {
	case "=", ":=", "+=", "-=", "*=", "/=", "++", "--":
		f = nodeAssign
	case "if", "for", "switch", "select":
		f = nodeControlBlock
	case "return": // TODO: This needs to handle more than one thing following return.
		f = nodeUnparenTwo
	default:
		f = nodeFuncall
	}
	return f(first)
}

// Process a value Node
func nodeProcessValue(n *Node) string {
	if n.content != "" {
		return n.content
	}
	first := n.first
	var f func(*Node) string
	switch first.content {
	case "+", "-", "*", "/", "==", "!=", ">=", "<=", "<", ">":
		f = nodeMath
	default:
		f = nodeFuncall
	}
	return f(first)
}

// Apply the correct node-parsing action for the given node and all
// its same-level successors
func nodeProcess(first *Node, f func(*Node) string) string {
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

// Convert a Node into Go code.
func (n *Node) GoString() string {
	return nodeProcess(n.first, nodeProcessTop)
}

// Return the contents of the given Node and its next Node, separated
// by a space.
func nodeUnparenTwo(first *Node) string {
	// TODO: Should have single "unparen" function
	return nodeProcessValue(first) + " " + nodeProcessValue(first.next)
}

// Convert an import Node into a Go import command.
func nodeImport(first *Node) string {
	out := "import ("
	for n := first.next; n != nil; n = n.next {
		out += n.content + "; "
	}
	out += ")"
	return out
}

// Convert a const or var Node into a top-level Go const or var
// declaration.
func nodeConstVar(first *Node) string {
	out := first.content + "("
	for n := first.next; n != nil; n = n.next {
		out += nodeContents(n.first) + "\n"
	}
	out += ")"
	return out
}

// Generate a space-separated list of node contents.
func nodeContents(first *Node) string {
	out := ""
	for n := first; n != nil; n = n.next {
		out += n.content + " "
	}
	return out
}

// Convert a function Node into a Go function declaration.
func nodeFunc(first *Node) string {
	// "func"
	n := first
	out := n.content
	// function name
	n = n.next
	out += " " + n.content
	// function args
	n = n.next
	out += "(" + nodeContents(n.first) + ")"
	// function return types
	n = n.next
	out += "(" + nodeContents(n.first) + ")"
	// function body
	n = n.next
	out += "{" + nodeProcess(n, nodeProcessAction) + "}"
	return out
}

// Process an assignment, starting from the first Node.
func nodeAssign(first *Node) string {
	// Go LHS and assignment operator
	out := first.next.content + first.content
	// RHS
	// TODO: Properly parse as values.
	out += nodeContents(first.next.next)
	return out
}

// Convert a function call into Go
func nodeFuncall(first *Node) string {
	out := first.content + "( " // The space is a hack to make the "len(out)-1" bit not remove the "(" in a nil-adic function
	for n := first.next; n != nil; n = n.next {
		out += nodeProcessValue(n) + ","
	}
	out = out[:len(out)-1] + ")"
	return out
}

// return text representing an "if condition { stuff() ... }" block
func nodeIfCase(firstCond *Node) string {
	out := ""
	if firstCond.content == "else" {
		out += " {\n"
	} else {
		out += "if " + nodeProcessValue(firstCond) + " {\n"
	}
	for n := firstCond.next; n != nil; n = n.next {
		out += nodeProcessAction(n) + "\n"
	}
	out += "}"
	return out
}

// return text representing a "for pre-statement; condition; post-statement { stuff() ... }" block of any type
func nodeForCase(controlClause *Node) string {
	out := "for "
	n := controlClause
	// get header
	switch {
	case n.content != "": // Golid for loops must paren the control clause.
		panic("Invalid 'for' control clause: \"" + n.String() + "\"!")
	case n.first == nil: // "()" case ('infinite' loop)
		out += "{\n"
	case n.first.content == "range": // "(range ...)" case
		panic("nodeForCase: 'range' isn't implemented!")
	case n.first.content != "": // "(condition)" case ('while' loop)
		out += nodeProcessValue(n.first) + "{\n"
	case n.first.first != nil: // "((pre) (cond) (post))" case ('for' loop)
		out += nodeAssign(n.first) + "; " + nodeProcessValue(n.first.next) + "; " + nodeProcessAction(n.first.next.next) + " {\n"
	default:
		panic("nodeForCase: Unhandled case!")
	}
	// go through body
	for n = n.next; n != nil; n = n.next {
		out += nodeProcessAction(n) + "\n"
	}
	// end brace
	out += "}\n"
	return out
}

// Convert if/for/switch/select statements into Go.
func nodeControlBlock(first *Node) string {
	out := ""
	switch first.content {
	case "if":
		n := first.next
		out += nodeIfCase(n.first)
		for n = n.next; n != nil; n = n.next {
			out += "else " + nodeIfCase(n.first)
		}
	case "for":
		out += nodeForCase(first.next)
	//case "switch":
	default:
		// TODO: implement other cases
		panic("nodeControlBlock is unimplemented for '" + first.content + "'!")
	}
	return out
}

// Convert a Lisp math function call into Go form.
func nodeMath(first *Node) string {
	op := first.content
	n := first.next
	lhs := n
	n = n.next
	rhs := n
	return "(" + nodeProcessValue(lhs) + " " + op + " " + nodeProcessValue(rhs) + ")"
}
