package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

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

func primeFinder(done <-chan int, randIntStream <-chan int) <-chan int {
	isPrime := func(randomInt int) bool {
		if randomInt < 2 {
			return false
		}
		for i := 2; i*i < randomInt; i++ {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}
	primes := make(chan int)
	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case randomInt := <-randIntStream:
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
		}
	}()
	return primes
}

func fanIn[T any](done <-chan int, channels ...<-chan T) <-chan T {

	var wg sync.WaitGroup
	fannedIntStream := make(chan T)

	transfer := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case fannedIntStream <- i:
			}
		}
	}
	for _, c := range channels {
		wg.Add(1)
		go transfer(c)
	}

	go func() {
		wg.Wait()
		close(fannedIntStream)
	}()
	return fannedIntStream
}

func main() {

	start := time.Now()
	done := make(chan int)
	defer close(done)

	randNumFetcher := func() int { return rand.Intn(500000000) }
	randIntStream := repeatFunc(done, randNumFetcher)
	//primeStream := primeFinder(done, randIntStream)

	// Fan out
	CPUCount := runtime.NumCPU()
	primeFinderChannels := make([]<-chan int, CPUCount)
	for i := 0; i < CPUCount; i++ {
		primeFinderChannels[i] = primeFinder(done, randIntStream)
	}

	// Fan in
	fannedInStream := fanIn(done, primeFinderChannels...)
	for rando := range take(done, fannedInStream, 20) {
		fmt.Println(rando)
	}
	fmt.Println(time.Since(start))
}
