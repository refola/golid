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

type nodeType int

const (
	topNode nodeType = iota
	actionNode
	valueNode
)

func (t nodeType) String() string {
	return map[nodeType]string{
		topNode:    "topNode",
		actionNode: "actionNode",
		valueNode:  "valueNode",
	}[t]
}

func nodeProcessTop(n *Node) string {
	var f func(*Node) string
	switch n.first.content {
	case "package":
		f = nodeUnparenTwo
	case "import":
		f = nodeImport
	case "const", "var":
		f = nodeConstVar
	case "func":
		f = nodeFunc
	default:
		panic("Unknown top-level node: " + n.first.content)
	}
	return f(n)
}

func nodeProcessAction(n *Node) string {
	var f func(*Node) string
	switch n.first.content {
	case "=", ":=", "+=", "-=", "*=", "/=":
		f = nodeAssign
	case "if", "for", "switch", "select":
		f = nodeControlBlock
	case "return":
		f = nodeUnparenTwo
	default:
		f = nodeFuncall
	}
	return f(n)
}

func nodeProcessValue(n *Node) string {
	panic("unimplemented!")
	var f func(*Node) string
	switch n.first.content {
	case "+", "-", "*", "/", "==", "!=", ">=", "<=", "<", ">":
		f = nodeMath
	default:
		f = nodeFuncall
	}
	return f(n)
}

// figure out and apply the correct node-parsing action
func nodeProcess(first *Node, t nodeType) string {
	var f func(*Node) string
	switch t {
	case topNode:
		f = nodeProcessTop
	case actionNode:
		f = nodeProcessAction
	case valueNode:
		f = nodeProcessValue
	default:
		panic(fmt.Errorf("Unknown type of node: %v", t))
	}
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
			panic(fmt.Errorf("Could not process code as %v: %v\n Got error: %v", t, n, err))
		} else {
			out += result + "\n"
		}
	}
	return out
}

func (n *Node) GoString() string {
	return nodeProcess(n.first, topNode)
}

func nodeUnparenTwo(n *Node) string {
	return n.first.content + " " + n.last.content
}

func nodeImport(n *Node) string {
	out := "import ("
	for n = n.first.next; n != nil; n = n.next {
		out += n.content + "; "
	}
	out += ")"
	return out
}

func nodeConstVar(n *Node) string {
	out := n.first.content + "("
	for n = n.first.next; n != nil; n = n.next {
		out += nodeContents(n) + "\n"
	}
	out += ")"
	return out
}

// return space-separated list of node contents
func nodeContents(n *Node) string {
	out := ""
	for n = n.first; n != nil; n = n.next {
		out += n.content + " "
	}
	return out
}

func nodeFunc(n *Node) string {
	// "func"
	n = n.first
	out := n.content
	// function name
	n = n.next
	out += " " + n.content
	// function args
	n = n.next
	out += "(" + nodeContents(n) + ")"
	// function return types
	n = n.next
	out += "(" + nodeContents(n) + ")"
	// function body
	out += "{" + nodeProcess(n.next.first, actionNode) + "}"
	return out
}

func nodeAssign(n *Node) string {
	// Go LHS and assignment operator
	n = n.first
	out := n.next.content + n.content
	// RHS
	for n = n.next.next; n != nil; n = n.next {
		out += " " + n.content
	}
	return out
}

func nodeFuncall(n *Node) string {
	return n.first.content + "(" + nodeProcess(n.first.next, valueNode) + ")"
}

func nodeControlBlock(n *Node) string {
	panic("unimplemented!")
}

func nodeMath(n *Node) string {
	n = n.first
	first := n.content
	n = n.next
	second := n.content
	n = n.next
	third := n.content
	return "(" + second + " " + first + " " + third + ")"
}
