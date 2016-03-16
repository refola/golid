/* Macro processing for Piklisp-Go

Piklisp needs to have macros to be a proper Lisp. However, for Piklisp
to still be "Go" and not merely "Go-integrated", it must process the
macros away such that the "compiled" result is pure Go. Thus the top
priority macro is "go-mac", which defines macros that turn into Go
code.

The tricky part is that everything needed to parse away macros must be
interpreted and ran instead of merely translated into Go.

*/

package parse
