package main

import (
	"fmt"
)

func Add(a int, b int) int {
	return a + b
}

func AddGeneric[T int | float64](a T, b T) T {
	return a + b
}

func main() {

	result := Add(1, 2)
	result1 := AddGeneric(1.5, 2.9)

	fmt.Printf("result: %+v\n", result)
	fmt.Printf("result1: %+v\n", result1)
}
