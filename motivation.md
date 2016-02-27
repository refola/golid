# Motivation for making Piklisp

When programming, I often need to see the value of some variable or
expression in order to understand what's going on. In go, I can do
this via something like `fmt.Printf("x=%v\n", x)`.

That works, but that's a lot of duplicate code. The variable `x` is
mentioned twice. That means that if I change the name of the variable
then I'll have to change it in a string as well as in the actual
code. That's a lot of opportunity to make mistakes.

It's also inconvenient. This works for a single variable/expression,
but what if I want to see 5 values at once?

What I need is a function (ideally stashed away in some package for
easy reuse everywhere) that takes the names of variables (or even
arbitrary expressions) and transforms them into the appropriate
`fmt.Printf` calls. Then I could simply write something like
`magic.ShowAndEvaluate(v, w, x, y, z)` and have it be equivalent to
writing out `fmt.Printf("v=%v\n", v); fmt.Printf("w=%v\n", w);
fmt.Printf("x=%v\n", x); fmt.Printf("y=%v\n", y); fmt.Printf("z=%v\n",
z)` the long way.

Sounds simple enough, right? All I have to do is go through a list of
variables, print out the names, print the values, and throw in a bit
of formatting. This seems like a very basic function that I should be
able to write just a few hours of experience with Go, or with pretty
much any programming language. It's really simple and straightforward,
except for one little problem. How do you represent unexecuted code in
a way the compiler understands?

The short answer in Go is that you can't. In Go, you can pass to a
function a string representing the name of a variable, but the
function has no easy way of knowing what the string means. If you
write `x:=5; magic.GetTheValue("x")` then there's no way in Go to
write a `GetTheValue` function that knows to print `5`. All it has is
`"x"`, not `x`. This is the frustration of being 90% done while the
remaining 10% is blatently obvious, yet impossible.

This 'titanium wall at 90%' problem isn't unique to Go. Pretty much
any programming language has it. I've tried Python and found that it
has the same wall. Python's `eval` builtin moves the wall to 99%, but
that only gives me that much more time to build up momentum before
crashing into the wall. In Python what happens is the "magic function"
works when it's defined in the same variable scope as it's called
from. But then moving the function into its own file (so it can be
easily called from several different files) breaks the magic because
it's in a different variable scope. This wall has become my number one
annoyance with virtually every programming language¹. It doesn't
matter where the wall is. It shouldn't exist in the first place. I
keep hitting this programming wall of "you can't".

But I really, really don't like being told what I can't do. This is
doubly true in something as arbitrary as programming. So instead of
taking the short answer of "you can't", I tried the next shorter
answer of "use a Lisp".

So I went and learned some Common Lisp. By "learned", I mean I did the
"Lisp Koans"(https://github.com/google/lisp-koans) to quickly go
through the basics of how the language works. Syntactically, I was
quite impressed. It was actually really easy in Common Lisp to write
the desired "magic function". Semantically, however, I wasn't so
impressed. I don't want to trade syntactic overload for mentally
juggling almost half a dozen meanings of "is X equal to Y".

So I moved on to the third-shortest answer: "turn it into Lisp". I
tried that with Python, but I got lost in over specifying it. Also,
though Python has some neat tricks, its syntax is too messy for
me. Now, half a year later, I'm finally doing it with Go. I chose Go
because it's the language I'm most familiar with, it's fast enough
that I don't notice weird slowdowns, and it's simple enough that I
won't get overwhelmed as easily.



¹ The exceptions are Lisps (with syntactic macros, hence this project)
and POSIX-like shell scripting (I use Bash). Shell scripting actually
sorta kinda has this capability between `eval` and "sourcing"
files. But shell scripting is also severely lacking in such basic
programming features as a sane quoting syntax.
