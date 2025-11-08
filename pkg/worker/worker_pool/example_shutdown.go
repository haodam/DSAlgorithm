package worker_pool

import (
	"context"
	"fmt"
	"os"
	"time"
)

// ExampleShutdownWithTimeout demonstrates shutdown with timeout
func ExampleShutdownWithTimeout() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(5),
		WithTimeout(10*time.Minute),
	)

	pool.Start()

	// Submit some tasks
	for i := 0; i < 100; i++ {
		task := &ExampleTask{ID: i, Duration: 100 * time.Millisecond}
		pool.Submit(task)
	}

	// Shutdown with 5 second timeout
	err := pool.ShutdownWithTimeout(5 * time.Second)
	if err != nil {
		fmt.Printf("Shutdown error: %v\n", err)
	} else {
		fmt.Println("Shutdown completed successfully")
	}
}

// ExampleShutdownGracefully demonstrates graceful shutdown with draining
func ExampleShutdownGracefully() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(10),
		WithTimeout(10*time.Minute),
	)

	pool.Start()

	// Submit tasks
	for i := 0; i < 1000; i++ {
		task := &ExampleTask{ID: i, Duration: 50 * time.Millisecond}
		pool.Submit(task)
	}

	// Graceful shutdown: drain tasks first, then shutdown
	// Drain timeout: 30 seconds, shutdown timeout: 30 seconds
	err := pool.ShutdownGracefully(30 * time.Second)
	if err != nil {
		fmt.Printf("Graceful shutdown error: %v\n", err)
	} else {
		fmt.Println("Graceful shutdown completed")
	}

	// Get final stats
	if stats := pool.GetMetrics(); stats != nil {
		fmt.Printf("Final stats: Completed: %d, Failed: %d\n",
			stats.TasksCompleted, stats.TasksFailed)
	}
}

// ExampleShutdownOnSignal demonstrates shutdown on OS signals
func ExampleShutdownOnSignal() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(10),
		WithMetricsEnabled(true),
		WithTimeout(10*time.Minute),
	)

	pool.Start()

	// Submit tasks continuously
	go func() {
		for i := 0; ; i++ {
			task := &ExampleTask{ID: i, Duration: 100 * time.Millisecond}
			if err := pool.Submit(task); err != nil {
				break // Pool shutdown
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Shutdown on SIGINT or SIGTERM
	fmt.Println("Pool running. Press Ctrl+C to shutdown gracefully...")
	err := pool.ShutdownOnSignal()
	if err != nil {
		fmt.Printf("Shutdown error: %v\n", err)
	} else {
		fmt.Println("Shutdown completed on signal")
	}
}

// ExampleShutdownOnSignalAsync demonstrates async signal handling
func ExampleShutdownOnSignalAsync() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(10),
		WithTimeout(10*time.Minute),
	)

	pool.Start()

	// Submit tasks
	for i := 0; i < 100; i++ {
		task := &ExampleTask{ID: i, Duration: 100 * time.Millisecond}
		pool.Submit(task)
	}

	// Start async signal handler
	errChan := pool.ShutdownOnSignalAsync()

	// Do other work...
	fmt.Println("Pool is running. Will shutdown on signal...")

	// Wait for shutdown
	err := <-errChan
	if err != nil {
		fmt.Printf("Shutdown error: %v\n", err)
	} else {
		fmt.Println("Shutdown completed")
	}
}

// ExampleGetPendingTasks demonstrates checking pending tasks
func ExampleGetPendingTasks() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(5),
		WithPriorityEnabled(true),
		WithTimeout(10*time.Minute),
	)

	pool.Start()

	// Submit many tasks
	for i := 0; i < 1000; i++ {
		task := &ExampleTask{ID: i, Duration: 50 * time.Millisecond}
		pool.SubmitWithPriority(task, i%10)
	}

	// Check pending tasks periodically
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 10; i++ {
		<-ticker.C
		pending := pool.GetPendingTasksCount()
		fmt.Printf("Pending tasks: %d\n", pending)

		if pending == 0 {
			break
		}
	}

	// Shutdown
	pool.Shutdown()
}

// ExampleProductionShutdown demonstrates production-ready shutdown pattern
func ExampleProductionShutdown() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(10),
		WithPriorityEnabled(true),
		WithMetricsEnabled(true),
		WithTimeout(10*time.Minute),
	)

	pool.Start()

	// Your application logic here...
	// Submit tasks, process results, etc.

	// Production shutdown pattern:
	// 1. Stop accepting new tasks (already done by shutdown flag)
	// 2. Try graceful shutdown first
	// 3. If timeout, force shutdown

	fmt.Println("Starting graceful shutdown...")

	// Try graceful shutdown first (drain + wait)
	err := pool.ShutdownGracefully(30 * time.Second)
	if err != nil {
		fmt.Printf("Graceful shutdown failed: %v\n", err)
		fmt.Println("Forcing shutdown...")

		// Force shutdown with shorter timeout
		err = pool.ShutdownWithTimeout(5 * time.Second)
		if err != nil {
			fmt.Printf("Force shutdown error: %v\n", err)
			// At this point, you may want to log and exit
			os.Exit(1)
		}
	}

	// Get final metrics
	if stats := pool.GetMetrics(); stats != nil {
		fmt.Printf("\nFinal Metrics:\n")
		fmt.Printf("  Submitted: %d\n", stats.TasksSubmitted)
		fmt.Printf("  Completed: %d\n", stats.TasksCompleted)
		fmt.Printf("  Failed: %d\n", stats.TasksFailed)
		fmt.Printf("  Success Rate: %.2f%%\n", stats.SuccessRate)
	}

	fmt.Println("Shutdown completed")
}
