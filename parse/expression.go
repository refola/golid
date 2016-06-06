// define syntax tree format and basic operations

package parse

import "strings"

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
	parent      *Node // the Node that this one's a child of
	next        *Node // the next Node under this Node's parent
	first, last *Node // the first and last child Nodes of this one
	content     string
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

// Indent every line with a leading tab.
func indent(s string) string {
	ret := ""
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		ret += "\t" + line + "\n"
	}
	ret = ret[:len(ret)-1] // remove trailing "\n"
	return ret
}

// Spit out the code in Lisp form
func (n *Node) String() string {
	ret := ""
	switch {
	case n.first != nil:
		ret += "("
		plainVals := true
		for child := n.first; child != nil; child = child.next {
			if child.first != nil { // if we've reached a (...) group
				plainVals = false
			}
			if plainVals {
				ret += " "
			} else {
				ret += "\n"
			}
			indented := indent(child.String())
			indented = indented[1:] // remove leading "\t"
			ret += indented
		}
		ret = "(" + ret[2:] // convert leading "( " to "("
		ret += ")"
	case n.content == "":
		ret += "()"
	default:
		ret += n.content
	}
	return ret
}
