# piklisp.go
Program in Go with Lisp syntax

This is an experimental project where I'm making a basic Lisp
that "compiles" to valid Go source code, enabling use of Lisp
syntactic macros with minimally-changed Go code. It currently supports
some very basic Go programs.

# Wishlist
* Stable, complete, accurate support for all the Go features it uses
* Translation of Piklisp into itself
* Lisp macros (with quasi-quoting)
* More useful `GoString()` with line numbers of errors
* Prettier output formatting
* Comment preservation
* An Emacs mode file for Piklisp
* Automatic conversion of an entire source tree's Piklisp files in one fell swoop
* Automatic conversion of Go into Piklisp syntax.

# Installation
* [Install and configure Go](https://golang.org/doc/install)
* Run `go get github.com/refola/piklisp_go/piklisp`
* Run `go install github.com/refola/piklisp_go/piklisp`
* Run `piklisp file.{plgo|gol}` to convert `file.{plgo|gol}` into `{plgo|gol}_file.go`
* Run `cd $GOPATH/src/github.com/refola/piklisp_go/piklisp` and `./test.sh` to run current tests.

# More info
* [Credits](credit.md)
* [Motivation/history](motivation.md)

# License
This is licensed as GPLv3 because that's the most restrictive
license GitHub offers by default. I know that this is an inappropriate
license for something resembling a programming language. If for some
reason you want to use experimental learning code for your
GPLv3-incompatible project, then please make a bug report describing
your project and why you want to use piklisp.go for it in a
GPLv3-incompatible way. I'll gladly change the license in exchange for
having a public record of someone wanting to use my project. Depending
on demand, I'm potentially willing to go as far as public domain or
the "Unlicense". But I need a good reason to relicense.
