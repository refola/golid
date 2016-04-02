package main

import "fmt"

func main() {
	for {
		fmt.Println("'infinite' loop")
		break
	}
	i := 0
	for i < 1 {
		fmt.Println("'while' loop")
		i++
	}
	for x := 0; x < 1; x++ {
		fmt.Println("'for' loop")
	}
	for range []int{1} {
		fmt.Println("'range' loop without vars")
	}
	m := map[string]string{"'range'": "with vars"}
	for i, v := range m {
		fmt.Printf("%s loop %s\n", i, v)
	}
}
