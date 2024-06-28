package main

import (
	"fmt"
	"math/rand"
	"time"
)

func LeakReceive() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.") // (1)
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return

				}
				randStream <- rand.Int()
			}
		}()
		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)
	time.Sleep(1 * time.Second)
}

func Example2() {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c) // (1)
	}()
	fmt.Println("Blocking on read...")
	select {
	case <-c: // (2)
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

func MulChannel() {

	ch1 := make(chan interface{})
	close(ch1)
	ch2 := make(chan interface{})
	close(ch2)

	var ch1Count, ch2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-ch1:
			ch1Count++
		case <-ch2:
			ch2Count++
		}
	}
	fmt.Printf(" ch1Count: %d\n ch2Count: %d\n", ch1Count, ch2Count)
}

func WaitingOther() {
	done := make(chan interface{})
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("done")
		close(done)
	}()
	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}
		fmt.Println(workCounter)
		workCounter++
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("workCounter %v .\n", workCounter)
}

func main() {
	//LeakReceive()
	//Example2()
	MulChannel()
}
