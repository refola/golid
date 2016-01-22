// define syntax tree format and basic operations

package parse

// An Expression represents a parsed Lisp expression, which is either
// a list of Expressions or an Atom. This interface attempts to unify
// both cases.
type Expression interface {
	// Return to canonical (Lisp) form. It doesn't have to be pretty. It
	// just has to have the parentheses match.
	String() string

	// Convert (one-way?) to Go form. It doesn't have to be pretty. It
	// just has to compile if the input code is valid piklisp-go.
	GoString() string
}

// A Node represents a single thing in parsing a Lisp expression.
type Node struct {
	parent, prev, next, first, last *Node
	content                         string
}

// Make a root node.
func Root() *Node {
	return new(Node)
}

// Make a new Node after the current Node's last child
func (n *Node) MakeChild() *Node {
	child := new(Node)
	if n.first == nil {
		n.first = child
	} else {
		n.last.next = child
		child.prev = n.last
	}
	n.last = child
	child.parent = n
	return child
}

// Append a child Node containing a string
func (n *Node) AddToken(s string) {
	n.MakeChild()
	n.last.content = s
}

// Accessor needed for parser
func (n *Node) Parent() *Node { return n.parent }

// Spit out the code in Lisp form
func (n *Node) String() string {
	ret := ""
	switch {
	case n.first != nil:
		ret += "("
		child := n.first
		for child != nil {
			ret += child.String()
			ret += "\n"
			child = child.next
		}
		ret = ret[:len(ret)-1] + ")"
	case n.content == "":
		ret += "()"
	default:
		ret += n.content
	}
	return ret
}

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
