// ngs_main.go

// This file contains Node.GoString() itself. It may be the "main" ngs
// file, but it's also really small because almost everything's done
// by the other functions.

package parse

// Convert a Node into Go code.
func (n *Node) GoString() string {
	return nu_process_many(n.first, nc_top)
}
