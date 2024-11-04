package main

import "fmt"

const (
	red     = 0
	yellow  = 1
	green   = 2
	blue    = 3
	magenta = 4
	cyan    = 5
)

func main() {
	fmt.Printf("%v, %T\n", red, red)
	fmt.Printf("%v, %T\n", yellow, yellow)
	fmt.Printf("%v, %T\n", green, green)
	fmt.Printf("%v, %T\n", blue, blue)
	fmt.Printf("%v, %T\n", magenta, magenta)
	fmt.Printf("%v, %T\n", cyan, cyan)
}
