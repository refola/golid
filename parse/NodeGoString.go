// convert a Node into Go code

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
