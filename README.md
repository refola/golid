# piklisp.go
Program in Go with Lisp syntax

This is an experimental project where I'm trying to make a basic Lisp that "compiles" to valid Go source code. The core goal is to learn how to make a Lisp. Secondary goals are letting it access all Go features, implementing syntactic macros, and getting rid of parentheses (both syntactically as in [SRFI#49](http://srfi.schemers.org/srfi-49/srfi-49.html) and semantically as in [Arc](http://www.arclanguage.org)). Even vague thinking beyond that will kill the primary goal's motivation, so that's all for now.

TODO:
* Generalize from "Hello world" to "real" programs
* Make improvements based on the experience from the previous item

Here's how to use this as-is:
* Checkout this repository with Git.
* Open a console in the directory you downloaded this to.
* Run `go build ./piklisp/main.go`.
* Run `./main file.{plgo|gol}` to convert `file.{plgo|gol}` into `{plgo|gol}_file.go`.
* Run `go test ./parse` to run current tests.

Try running `./main tests/interesting_looking_file.gol` followed by `go run tests/gol_interesting_looking_file.go` on whichever file looks interesting to see what runs with this so far.

License:
This is licensed as GPLv3 because that's the most restrictive license GitHub offers by default. I know that this is an inappropriate license for something resembling a programming language. If for some reason you want to use experimental learning code for your GPLv3-incompatible project, then please make a bug report describing your project and why you want to use piklisp.go for it in a GPLv3-incompatible way. I'll gladly change the license in exchange for having a public record of someone wanting to use my project. Depending on demand, I'm potentially willing to go as far as public domain or the "Unlicense". But I need a good reason to relicense.
