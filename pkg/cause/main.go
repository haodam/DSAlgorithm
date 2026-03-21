package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Println("Processing...", os.Getgid())

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-ctx.Done():
			cause := context.Cause(ctx)
			fmt.Println("cause", cause)
			switch {
			case strings.Contains(cause.Error(), "terminated"):
				fmt.Printf("SIGTERM saving checkpoint at item %d\n", i)

			case strings.Contains(cause.Error(), "interrupt"):
				fmt.Println("Manual stop (Ctrl+C), discard progress")

			default:
				fmt.Println("Unknown cause:", cause)
			}
			return
		case <-ticker.C:
			fmt.Printf("item %d Processing\n", i)
			i++
		}
	}
}
