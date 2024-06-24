package main

import (
	"fmt"
	"time"
)

func printNumbers(number int, done chan bool) {
	fmt.Println(number)
	//time.Sleep(1 * time.Second)
	done <- true // Gửi tín hiệu rằng goroutine này đã hoàn thành
}

func main() {
	done := make(chan bool)

	go func() {
		for i := 1; i <= 100; i++ {
			<-done // Đợi cho đến khi goroutine trước đó hoàn thành
		}
		close(done) // Đóng channel sau khi tất cả goroutine đã hoàn thành
		fmt.Println("All numbers printed.")
	}()

	for i := 1; i <= 100; i++ {
		time.Sleep(1 * time.Second)
		go printNumbers(i, done)
	}
}
