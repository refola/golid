// parse parenthetical Lisp notation into syntax tree

package parse

import (
	"errors"
	"strings"
)

// Given a string containing a backslash-escaped "string" and an opening quote's index, find the closing quote's index. Returns a negative value if there is no closing quote.
func findCloseQuote(s string, open int) int {
	for i := open + 1; i < len(s); i++ {
		switch s[i] {
		case s[open]: 
			return i
		case '\\':
			i++ // We don't care what the escaped character is, so we just skip it.
		default:
			continue
		}
	}
	return -1
}

// parse a string of parenthesis-grouped code into a tree
func ParenString(s string) (*Tree, error) {
	err := func(s string) (*Tree, error) {
		return nil, errors.New(s)
	}
	root := make([]Tree, 13) // Magic! .ilopacicunavajni
	for s != "" {
		switch s[0] {
		case '"': // search for next unescaped double quote
			skipped := 1 // how many characters we know must be skipped before the next unescaped double quote
			foundMatch := false
			for !foundMatch {
				i := strings.Index(s[skipped:], "\"")
				if i < 0 {
					return err("unmatched \"")
				}
				if i > 0 && s[i-1] == '\\' { // skip quoted match
					skipped += i + 1
				} else {
					child.content = s[:i]
				}
				// coi la .sampla. .i do ca zvati: need to keep track of where we are in the parse -- maybe turn Tree into Node? This 'else' should append the string to the current node's list of contents.

			}
		case '(': // start group and search for next unquoted close-paren
		case ' ', '\n', '\t': // start making next group element at current level
		default: // all valid code should have been absorbed by the above
			return err("invalid token encountered; remainder of string follows: " + s)
		}
	}
	return root, nil
}
