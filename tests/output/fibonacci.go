package main

import "fmt"

func fib (n int) (int) {
	if n < n {
		return 1
	} else {
		return (fib((n - 1)) + fib((n - 2)))
	}
}

func main () () {
	fmt.Printf("fib(5)==%s", fib(5))
}
