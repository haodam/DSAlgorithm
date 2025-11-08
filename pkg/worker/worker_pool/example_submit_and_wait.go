package worker_pool

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ExampleSubmitAndWait demonstrates basic SubmitAndWait usage
func ExampleSubmitAndWait() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(5),
		WithTimeout(10*time.Minute),
	)

	pool.Start()
	defer pool.Shutdown()

	// Submit task and wait for result
	task := &ExampleTask{
		ID:       1,
		Duration: 100 * time.Millisecond,
	}

	result, err := pool.SubmitAndWait(task)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", result)
	}
}

// ExampleSubmitAndWaitWithTimeout demonstrates SubmitAndWait with timeout
func ExampleSubmitAndWaitWithTimeout() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(5),
		WithTimeout(10*time.Minute),
	)

	pool.Start()
	defer pool.Shutdown()

	// Task that takes longer than timeout
	task := &ExampleTask{
		ID:       1,
		Duration: 5 * time.Second,
	}

	// Wait with 2 second timeout
	result, err := pool.SubmitAndWaitWithTimeout(task, 2*time.Second)
	if err != nil {
		fmt.Printf("Timeout or error: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", result)
	}
}

// ExampleSubmitAndWaitWithContext demonstrates SubmitAndWait with context
func ExampleSubmitAndWaitWithContext() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(5),
		WithTimeout(10*time.Minute),
	)

	pool.Start()
	defer pool.Shutdown()

	// Create context with timeout
	waitCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Task
	task := &ExampleTask{
		ID:       1,
		Duration: 100 * time.Millisecond,
	}

	// Submit and wait with context
	result, err := pool.SubmitAndWaitWithContext(waitCtx, task)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", result)
	}
}

// ExampleSubmitAndWaitWithError demonstrates handling errors
func ExampleSubmitAndWaitWithError() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(5),
		WithTimeout(10*time.Minute),
	)

	pool.Start()
	defer pool.Shutdown()

	// Task that will fail
	task := &ExampleTask{
		ID:         1,
		Duration:   100 * time.Millisecond,
		ShouldFail: true,
		FailCount:  1, // Fail once
	}

	result, err := pool.SubmitAndWait(task)
	if err != nil {
		fmt.Printf("Task failed as expected: %v\n", err)
	} else {
		fmt.Printf("Unexpected success: %s\n", result)
	}
}

// ExampleSubmitAndWaitMultiple demonstrates multiple SubmitAndWait calls
func ExampleSubmitAndWaitMultiple() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(10),
		WithTimeout(10*time.Minute),
	)

	pool.Start()
	defer pool.Shutdown()

	// Submit multiple tasks and wait for each
	for i := 0; i < 10; i++ {
		task := &ExampleTask{
			ID:       i,
			Duration: 50 * time.Millisecond,
		}

		result, err := pool.SubmitAndWait(task)
		if err != nil {
			fmt.Printf("Task %d failed: %v\n", i, err)
		} else {
			fmt.Printf("Task %d completed: %s\n", i, result)
		}
	}
}

// ExampleSubmitAndWaitConcurrent demonstrates concurrent SubmitAndWait calls
func ExampleSubmitAndWaitConcurrent() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(10),
		WithTimeout(10*time.Minute),
	)

	pool.Start()
	defer pool.Shutdown()

	// Submit tasks concurrently
	results := make(chan string, 10)
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			task := &ExampleTask{
				ID:       id,
				Duration: 100 * time.Millisecond,
			}

			result, err := pool.SubmitAndWait(task)
			if err != nil {
				errors <- err
			} else {
				results <- result
			}
		}(i)
	}

	// Collect results
	successCount := 0
	errorCount := 0

	for i := 0; i < 10; i++ {
		select {
		case result := <-results:
			fmt.Printf("Success: %s\n", result)
			successCount++
		case err := <-errors:
			fmt.Printf("Error: %v\n", err)
			errorCount++
		case <-time.After(5 * time.Second):
			fmt.Println("Timeout waiting for results")
			break
		}
	}

	fmt.Printf("Completed: %d success, %d errors\n", successCount, errorCount)
}

// ExampleSubmitAndWaitVsSubmit demonstrates difference between SubmitAndWait and Submit
func ExampleSubmitAndWaitVsSubmit() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(5),
		WithTimeout(10*time.Minute),
	)

	pool.Start()
	defer pool.Shutdown()

	// Method 1: Submit (async) - fire and forget
	task1 := &ExampleTask{ID: 1, Duration: 100 * time.Millisecond}
	err := pool.Submit(task1)
	if err != nil {
		fmt.Printf("Submit error: %v\n", err)
	}
	fmt.Println("Task 1 submitted (async)")

	// Method 2: SubmitAndWait (sync) - wait for result
	task2 := &ExampleTask{ID: 2, Duration: 100 * time.Millisecond}
	result, err := pool.SubmitAndWait(task2)
	if err != nil {
		fmt.Printf("SubmitAndWait error: %v\n", err)
	} else {
		fmt.Printf("Task 2 completed: %s\n", result)
	}

	// Method 3: SubmitAndWait with timeout
	task3 := &ExampleTask{ID: 3, Duration: 100 * time.Millisecond}
	result, err = pool.SubmitAndWaitWithTimeout(task3, 2*time.Second)
	if err != nil {
		fmt.Printf("SubmitAndWaitWithTimeout error: %v\n", err)
	} else {
		fmt.Printf("Task 3 completed: %s\n", result)
	}
}

// ErrorTask is a task that always fails
type ErrorTask struct {
	ID  int
	Msg string
}

func (t *ErrorTask) Process(ctx context.Context) (string, error) {
	return "", errors.New(t.Msg)
}

// ExampleSubmitAndWaitErrorHandling demonstrates error handling
func ExampleSubmitAndWaitErrorHandling() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(5),
		WithTimeout(10*time.Minute),
	)

	pool.Start()
	defer pool.Shutdown()

	// Task that always fails
	task := &ErrorTask{
		ID:  1,
		Msg: "simulated error",
	}

	result, err := pool.SubmitAndWait(task)
	if err != nil {
		fmt.Printf("Task failed as expected: %v\n", err)
		fmt.Printf("Result is zero value: %q\n", result) // Should be empty string
	} else {
		fmt.Printf("Unexpected success: %s\n", result)
	}
}
