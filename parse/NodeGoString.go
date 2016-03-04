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
→ first (args ...)
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

func parsePackage(n *Node) string {
	return "package " + n.last.content + "\n"
}
func parseImport(n *Node) string {
	ret := "import ("
	child := n.first.next
	for child != nil {
		ret += child.content + " "
		child = child.next
	}
	ret += ")"
	return ret
}

// return a space-separated list of node contents
func children(n *Node) string {
	child := n.first
	ret := ""
	for child != nil {
		ret += child.content + " "
		child = child.next
	}
	return ret
}

// output a function call in Go form
func parseFunCall(n *Node) string {
	child := n.first
	ret := child.content + "("
	child = child.next
	for child != nil {
		if child.content == "" {
			ret += parseFunCall(child)
		} else {
			ret += child.content
		}
		if child.next != nil {
			ret += ","
		}
		child = child.next
	}
	ret += ")"
	return ret
}
func parseFunc(n *Node) string {
	// This is really crude and will need recursion eventually, but it'll get through Hello World.
	ret := "func"
	child := n.first.next
	ret += " " + child.content
	child = child.next
	for i := 0; i < 2; i++ { // parameters and return
		ret += "(" + children(child) + ")"
		child = child.next
	}
	ret += "{"
	// function body -- for now just assum eit's a bunch of function calls
	for child != nil {
		ret += parseFunCall(child) + ";"
		child = child.next
	}
	ret += "}"
	return ret
}

var parseCases map[string]func(*Node) string = map[string]func(*Node) string{"package": parsePackage, "import": parseImport, "func": parseFunc}

// Convert Lisp to Go.
func (n *Node) GoString() string {
	// This is very crude, but it should get Hello World working.

	// Assume we're at the top level.
	if n.first.first.content != "package" {
		panic("This only works on top-level nodes that represent complete programs. Also, it doesn't like comments.")
	}

	// just parse it
	ret := ""
	child := n.first
	for child != nil {
		ret += parseCases[child.first.content](child) + "\n"
		child = child.next
	}
	return ret
}
