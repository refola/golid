// parse parenthetical Lisp notation into syntax tree

package parse

import (
	"errors"
	"fmt"
	"regexp"
)

// Find shortest sequence of double quote, followed by escaped and unescaped characters, followed by double quote
var stringRegex = regexp.MustCompile(`"([^"\\]|\\.)*"`)

// Find longest sequence of non-syntax characters
var tokenRegex = regexp.MustCompile("[^ \t\n\"()]+")

// Given a string starting with a token, find the end of the token (the index of the first character that follows the token).
// Rules:
// * If the token starts with a double quote ("), the token ends at the next double quote that isn't backslash-escaped.
// * Otherwise the token ends right before the next syntax character.
// Returns a negative value if there is no valid token.
func findTokenEnd(s string) int {
	var re *regexp.Regexp
	if s[0] == '"' {
		re = stringRegex
	} else {
		re = tokenRegex
	}
	loc := re.FindStringIndex(s)
	if loc == nil {
		return -1
	} else {
		return loc[1]
	}
}

// change the node according to how the indentation depth level changed
// NOTE: This assumes that the calling parse function is not at a blank line state.
func indentSrfi49(depthChange int, node *Node) *Node {
	switch {
	case depthChange < 0:
		// decrease depth by as many levels as the change
		for i := depthChange; i < 0; i++ {
			node = node.Parent()
		}
		// Now that we're back at the right level, we need to start a new sibling Node.
		node = node.Parent().MakeChild()
	case depthChange == 0:
		// make new sibling node at same depth
		node = node.Parent().MakeChild()
	case depthChange > 0:
		// increasing depth makes a new child node
		for i := 0; i < depthChange; i++ {
			node = node.MakeChild()
		}
	default:
		panic("Impossible! DepthChange is not less than, equal to, or greater than zero!")
	}
	return node
}

// what type of parsing we're doing
type parseMode int

const (
	classic = parseMode(iota)
	srfi49
)

// Parse a Piklisp string into ints AST representation
func parseString(s string, mode parseMode) (Expression, error) {
	// setup everything
	line := 0 // which line was last reached

	// convenience "exit while showing error and line number" function
	err := func(s string) (Expression, error) {
		return nil, errors.New(fmt.Sprintf("Line %s: %v", line, s))
	}

	implicitParenDepth := 0 // how many layers of parentheses are currently elided
	root := Root()          // the top-level node to return
	node := root            // the active node being parsed
	if mode == srfi49 {
		// here there are no extra parens to start the first subnode
		node = node.MakeChild()
	}

	// keep going until the string's empty
	for s != "" {
		switch s[0] {
		case '(': // go a level deeper
			node = node.MakeChild()
			s = s[1:]
		case ')': // go back up a level
			node = node.Parent()
			s = s[1:]
		case ' ', '\t': // ignored whitespace
			s = s[1:]
		case '\n': // either whitespace or marking srfi49 checks
			// absorb all consecutive newlines so blank lines can be used as
			// separation inside an indented group
			for s != "" && s[0] == '\n' {
				line++
				s = s[1:]
			}
			if mode == classic {
				continue
			} else {
				newDepth := 0
				for s != "" && s[0] == '\t' {
					newDepth++
					s = s[1:]
				}
				if s == "" {
					continue // Skip depth-change analysis when we're at the
					// end. The root Node is now complete.
				}
				node = indentSrfi49(newDepth-implicitParenDepth, node)
				implicitParenDepth = newDepth
			}
		case ';': // skip over comments and any following newlines
			for s != "" && s[0] != '\n' {
				s = s[1:]
			}
			if s != "" {
				for s != "" && s[0] == '\n' {
					line++
					s = s[1:]
				}
			}
		default: // must be a token
			end := findTokenEnd(s)
			if end < 0 {
				return err("could not find end of token: " + s)
			} else {
				node.AddToken(s[0:end])
				s = s[end:]
			}
		}
	}
	// Remove all extra layers
	for root != nil && root.first == root.last && root.content == "" {
		root = root.first
	}
	return root, nil
}
