/*
Piklisp is an experiment in learning how to make a Lisp that compiles to another programming language, effectively combining Lisp syntax and macros with whatever features the target language has. Package group piklisp/* contains the Go implementation. See piklisp-go/piklisp for the command.

Minimal grammar to parse Hello World

Everything between a pair of double quotes is treated as part of the same string.
  "This is one string even though it has
  a bunch of spaces (and a newline (and parentheses))."
Inside a string, a backslash preceding another character escapes its normal meaning. For example,
  "backslash \\ quote \" end of string"
produces a string containing
  backslash \ quote " end of string
Parentheses group stuff and whitespace separates things in a group, like in pretty much any Lisp.
  (first_item_in_outer_group (inner group which is second item in outer group))
The beginning and end of a file implicitly group everything inside.

Hello World example

Here's what I think the first working Hello World in piklisp_go will look like. Later versions should be more concise.
  (package main)
  (import "fmt")
  (func main () () (
    (fmt.Println "Hello World!")
  )

Wishlist

Here's what I'm planning with this. After the first item, priority has yet to be determined.
  * Make it work in the minimal Hello World case.
  * Make it work with all Go constructs.
  * Implement macros and other Lispy goodness.
  * Implement SRFI49 and a bunch of cleverness to get rid of more parentheses.
  * Handle production-y stuff like performance, language stability, and reliability.

Crazy vision

I want to be able to code "in Go" with Lisp syntax and less parentheses/brackets than languages like Python. Lisp macros (or similar) are essential for some types of problems, but there really are too many parentheses for normal use. I want to see such a language come into existence. But first I need to learn how to build a Lisp.

Crazy example

This is something that is possible with any Lisp, requires everything to be in the same global namespace to work in Python, and is impossible with any language that lacks eval.
  // pkg.gol
  package pkg
  import "fmt"
  macro Print-and-eval (code quoted) ()
    fmt.Println code "="
    fmt.Println (eval code)

  // main.gol
  package main
  import "pkg"
  func main () ()
    pkg.Print-and-eval '(+ 1 1)

This trivial example should output something like this.
  (+ 1 1)
  =
  2

This type of "see what is producing the value and the value it produces" macro is extremely useful for debugging and demonstrating what bits of code produce. The power of Lisp is necessary to do this without code duplication.
*/
package piklisp_go
