package main

import (
	"fmt"
	"time"
)

func printNumber() {
	for i := 0; i < 100; i++ {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}

func main() {

	go printNumber()

	// cho mot chut de goroutine co thoi gian hoan thanh
	time.Sleep(101 * time.Second)
	fmt.Println("All numbers printed.")
}
