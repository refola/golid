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

// Parse a string of parenthesis-grouped code into a Tree
func ParenString(s string) (Expression, error) {
	err := func(s string) (Expression, error) {
		return nil, errors.New(s)
	}
	root := Root()
	node := root
	for s != "" {
		switch s[0] {
		case '(': // go a level deeper
			node = node.MakeChild()
			s = s[1:]
		case ')': // go back up a level
			node = node.Parent()
			s = s[1:]
		case ' ', '\n', '\t': // skip over whitespace
			s = s[1:]
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
	return root, nil
}

// Parse a string of SRFI#49-formatted Lisp code into a Tree
func Srfi49String(s string) (Expression, error) {
	line := 0 // which line we're on
	err := func(s string) (Expression, error) {
		return nil, errors.New(fmt.Sprintf("Line %v: %s", line, s))
	}
	// how deep have we gone with leading tabs?
	depth := 0
	// immediately go a level deeper for the first line
	root := Root().MakeChild()
	node := root
	for s != "" {
		switch s[0] {
		case '(': // go a level deeper
			node = node.MakeChild()
			s = s[1:]
		case ')': // go back up a level
			node = node.Parent()
			s = s[1:]
		case ' ', '\t': // skip over meaningless whitespace
			s = s[1:]
		case '\n':
			line++
			// absorb and count new line's leading tabs
			new_depth := 0
			s = s[1:]
			for s != "" && s[0] == '\t' {
				new_depth++
				s = s[1:]
			}
			switch {
			case new_depth < depth:
				// decrease depth by as many levels as the change
				for i := 0; i < depth-new_depth; i++ {
					node = node.Parent()
				}
			case new_depth == depth:
				// make new sibling node at same depth
				node = node.Parent().MakeChild()
			case new_depth == depth+1:
				// increasing depth makes a new child node
				node = node.MakeChild()
			case new_depth > depth+1:
				err("Attempt to increase depth by more than one level at a time")
			default:
				err("Parser encountered impossible situation!")
			}
			depth = new_depth
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
	return root, nil
}
