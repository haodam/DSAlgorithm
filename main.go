package main

import "fmt"

func sum(n int) int {
	s := 0
	for i := range n {
		s += i
	}
	return s
}

func main() {
	for i := 0; i < 10; i++ {
		fmt.Println(sum(5)) // Kết quả: 15
	}

}
