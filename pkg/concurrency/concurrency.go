package main

import (
	"fmt"
	"math/rand"
	"time"
)

func sliceToChannel(nums []int) <-chan int {
	c := make(chan int)
	go func() {
		for _, num := range nums[1:] {
			c <- num
		}
		close(c)
	}()
	return c
}

func sq(in <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for num := range in {
			c <- num * num
		}
		close(c)
	}()
	return c
}

func repeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {

	stream := make(chan T)
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}

func take[T any, K any](done <-chan K, steam <-chan T, n int) <-chan T {
	taken := make(chan T)
	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case taken <- <-steam:
			}
		}
	}()
	return taken
}

func main() {

	// input
	//nums := []int{1, 2, 3, 4, 5}

	// stage 1
	//datachannel := sliceToChannel(nums)

	// stage 2
	//finalChannel := sq(datachannel)

	// stage 3
	// value := range finalChannel {
	//fmt.Println(value)
	//}

	start := time.Now()
	done := make(chan bool)
	defer close(done)

	//nums := runtime.NumCPU()

	randNumFetcher := func() int { return rand.Intn(10000000) }
	randIntStream := repeatFunc(done, randNumFetcher)

	for rando := range take(done, randIntStream, 10) {
		fmt.Println(rando)
	}
	fmt.Println(time.Since(start))
}
