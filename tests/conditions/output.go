package main

import "fmt"

func foo() bool {
	return true
}

func main() {
	if true {
		fmt.Println("true happens")
	}
	if false {
		panic("false happened")
	} else if true {
		fmt.Println("even if it's after false")
	}
	if foo() {
		fmt.Println("or if it's from a function")
	}
	if false {
		panic("false happened")
	} else if false {
		panic("false happened")
	} else if true {
		fmt.Println("or after multiple falses")
	}
	if true {
		fmt.Printf("or when it's on")
		fmt.Println("multiple lines")
	}
}
