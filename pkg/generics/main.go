package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

func Add(a int, b int) int {
	return a + b
}

func AddGeneric[T int | float64](a T, b T) T {
	return a + b
}

func AddGeneric2[T constraints.Ordered](a T, b T) T {
	return a + b
}

func MapValues[T constraints.Ordered](values []T, mapFunc func(T) T) []T {
	var newValues []T
	for _, value := range values {
		newValue := mapFunc(value)
		newValues = append(newValues, newValue)
	}
	return newValues

}
func main() {

	result := MapValues([]float64{1.1, 1.2, 1.3}, func(n float64) float64 {
		return n * 2
	})
	fmt.Printf("result: %+v\n", result)

	//result := Add(1, 2)
	//result1 := AddGeneric(1.5, 2.9)
	//
	//fmt.Printf("result: %+v\n", result)
	//fmt.Printf("result1: %+v\n", result1)
}
