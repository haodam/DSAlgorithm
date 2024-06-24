package main

import (
	"fmt"
	"sync"
	"time"
)

func printNumberr(number int, wg *sync.WaitGroup) {
	defer wg.Done() // Đánh dấu rằng goroutine này đã hoàn thành
	fmt.Println(number)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 100; i++ {
		time.Sleep(1 * time.Second)
		wg.Add(1) // Tăng số lượng goroutine cần chờ
		go printNumberr(i, &wg)
	}

	wg.Wait() // Chờ tất cả các goroutine hoàn thành
	fmt.Println("All numbers printed.")
}
