package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func Sum(wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for i := 0; i < 100e8; i++ {
		sum += i
	}
	fmt.Println(sum)
}

func main() {
	start := time.Now()
	numCPU := runtime.NumCPU()
	fmt.Printf("NumCPU: %d\n", numCPU)
	runtime.GOMAXPROCS(numCPU)
	var wg sync.WaitGroup
	wg.Add(1)
	for i := 1; i < numCPU; i++ {
		go Sum(&wg)
	}
	wg.Wait()
	fmt.Printf("Time elapsed: %s\n", time.Since(start))
}
