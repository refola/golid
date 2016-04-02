package main

import "fmt"

func main() {
	x := "foo"
	switch x {
	case "foo":
		fmt.Println("first case matches var")
	}
	x = "bar"
	switch x {
	case "foo":
		panic("x should not be 'foo' anymore!")
	case "bar":
		fmt.Println("second case matches var")
	}
	x = "baz"
	switch x {
	case "foo", "bar":
		panic("x should not be 'foo' or 'bar' anymore!")
	default:
		fmt.Println("neither case matches var")
	}
}
