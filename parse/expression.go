// define syntax tree format and basic operations

package parse

// what type of expression something is; needed for expressions to
// interact correctly
type ExprType int

const (
	// a list of expressions
	List ExprType = iota

	// an identifier, constant expression, keyword, or anything else
	// that doesn't syntactically contain other stuff inside of it
	Atom
)

// An Expression represents a parsed Lisp expression, which is either
// a list of Expressions or an Atom. This interface attempts to unify
// both cases.
type Expression interface {
	Type() ExprType // List or Atom

	// Return to canonical (Lisp) form. It doesn't have to be pretty. It
	// just has to have the parentheses match.
	String() string

	// Convert (one-way?) to Go form. It doesn't have to be pretty. It
	// just has to compile if the input code is valid piklisp-go.
	GoString() string
}

// a tree contains either children or a string
type Tree struct {
	parent, prev, next *Tree
	children           []*Tree
	content            string
}

func (t *Tree) String() string {
	output := ""
	if t.children != nil {
		output += "("
		for i, child := range t.children {
			output += child.String()
			output += "\n"
		}
		output = output[:len(output)-1] + ")"
	} else {
		output = t.content
	}
	return output
}
