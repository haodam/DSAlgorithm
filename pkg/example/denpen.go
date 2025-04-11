package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	n := 20
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ch := make(chan int)

	go func(ctx context.Context, n int, ch chan int) {
		for i := 0; i <= n; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("context canceled")
				return
			case num := <-ch:
				fmt.Println(num)
				time.Sleep(500 * time.Millisecond)
			}
		}
		cancel()
	}(ctx, n, ch)
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
