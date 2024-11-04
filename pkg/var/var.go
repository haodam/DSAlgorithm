package main

import "fmt"

func main() {
	const i = 42
	fmt.Printf("%v, %T\n", i, i) // Print: 42, int

	const s = "Hello World"
	fmt.Printf("%v, %T\n", s, s) // Print: Hello World , string
}
