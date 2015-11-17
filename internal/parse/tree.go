// define syntax tree format and basic operations

package parse

// a tree contains either children or a string
type Tree struct {
	parent, prev, next Tree
	children           []Tree
	content            string
}

func (t Tree) String() string {
	output = ""
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
