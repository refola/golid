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
	output := ""
	if n.first != nil {
		output += "("
		child := n.first
		for child != nil {
			output += child.String()
			output += "\n"
			child = child.next
		}
		output = output[:len(output)-1] + ")"
	} else {
		output = n.content
	}
	return output
}

// Convert Lisp to Go.
func (n *Node) GoString() string {
	return "WARNING: This just wraps the String() function.\n" + n.String() // TODO: Implement this for real.
}
