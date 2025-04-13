package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	n := 20
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan int)

	go func(ctx context.Context, n int, ch chan int, cancelFunc context.CancelFunc) {
		sum := 0
		count := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context canceled")
				fmt.Printf("Đã nhận %d số Fibonacci\n", count)
				return
			case num := <-ch:
				fmt.Println(num)
				sum += num
				count++
				time.Sleep(100 * time.Millisecond)
				if count >= n {
					fmt.Println("sum fibonacci :", sum)
					cancelFunc()
					return
				}
			}
		}
	}(ctx, n, ch, cancel)

	fina(ctx, ch)
}

func fina(ctx context.Context, c chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-ctx.Done():
			fmt.Println("fina stopped")
			return
		}
	}
}
