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
* Run `go install github.com/refola/piklisp_go/cmd/piklisp` or, from the repository root, run `./install.sh`
* Run `piklisp file.gol` to convert `file.gol` into `gol_file.go`

# More info
* [Credits](doc/credit.md)
* [Motivation/history](doc/motivation.md)

# Lisp-iness disclaimer
I am most definitely not a Lisp wizard. I'm just a programmer who like's Go's general simplicity and Lisp's macros. If you are a Lisp wizard, then I would appreciate your advice for how to implement things more easily or more cleanly. *However*, please keep in mind that this is not a "make Lisp like Go" project. This is a "steal Lisp magic and give it to Go" project. As such, as much as possible will come from Go instead of Lisp. For example, this project uses Go's `=` and `:=` for setting variables; there is no `let` or `setf`. If you want a proper Lisp to use with Go, then go [here](https://github.com/glycerine/zygomys).

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
