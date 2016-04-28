// ngs_doc.go

// This file documents how the Node.GoString() functions are designed and implemented overall.

/* OVERVIEW

Go has many different keywords and syntactic elements, often with
different semantics depending on how and where things are
combined. However, once the contexts are understood, most things can
be thought of separately, as orthogonal concepts.

*/

/* TYPES OF GO SYNTAX TO HANDLE

===== CONTEXTS =====

Top-level: Anything that's valid outside of, e.g., a function body.
* package declaration
* import
* top-level consts and vars
* functions

Action: Places that require the program to _do_ something.
* bodies of functions
* bodies of control structures

Value: Places that need something that results in a single value.
* function arguments
* right-hand-side of ":=" and friends
* most args for if/for/switch
* array indices
* several other places

SimpleStmt (see golang.org/ref/spec#SimpleStmt): Anything except for control structures.
* values, channel sends, ++/--, assignments, and short declarations
* found at beginnings of control structures


===== SYNTAXES BY CONTEXT =====

Top-level:
* (first args) → first args
** package
* (first args ...)
→ first (args; ...)
** import
* (first (arg1 ...) (arg2 ...) ...)
→ first ( arg1 ...; arg2 ...; ... )
** const
** var
* (first second (args1 ...) (args2 ...) (args3 ...))
→ first second(args1, ...) (args2, ...) { args3; ...; }
** func

Action:
* (first second args ...)
→ second first args ...
** = := /= *= += -=
* (first args ...)
→ first(args ...)
** function calls
* (first (arg1 ...) (arg2 ...) ...)
→ first arg1 { ... } ??? arg2 { ... } ...
** if/for/switch/select
* (first second)
→ first second
** return
** goto
** label?

Value:
* (first second third)
→ (second first third)
** math expressions
** comparisons
* (first args ...)
→ first(args ...)
** function calls


===== STRATEGY FOR HANDLING THESE CASES =====

* start at top level
* have each function for handling these cases track what types of things it's expecting
* use per-context function maps to choose how to process each node

*/

package parse
