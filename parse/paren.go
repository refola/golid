// parse parenthetical Lisp notation into syntax tree

package parse

import (
	"errors"
	"regexp"
)

// Syntactic characters separate tokens.
const syntaxChars = " \t\n\""

// Find shortest sequence of double quote, followed by escaped and unescaped characters, followed by double quote
var stringRegex = regexp.MustCompile(`"([^"\\]|\\.)*"`)

// Find longest sequence of non-syntax characters
var tokenRegex = regexp.MustCompile("[^ \t\n\"]+")

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
func ParenString(s string) (*Tree, error) {
	err := func(s string) (*Tree, error) {
		return nil, errors.New(s)
	}
	root := NewNode()
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
