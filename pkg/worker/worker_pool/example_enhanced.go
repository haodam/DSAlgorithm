package worker_pool

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ExampleTask là một task đơn giản để demo
type ExampleTask struct {
	ID         int
	Duration   time.Duration
	ShouldFail bool
	FailCount  int
}

func (t *ExampleTask) Process(ctx context.Context) (string, error) {
	// Simulate processing time
	select {
	case <-time.After(t.Duration):
	case <-ctx.Done():
		return "", ctx.Err()
	}

	// Simulate failure
	if t.ShouldFail && t.FailCount > 0 {
		t.FailCount--
		return "", errors.New("simulated error")
	}

	return fmt.Sprintf("Task %d completed", t.ID), nil
}

// ExampleUsagePriority demonstrates priority queue usage
func ExampleUsagePriority() {
	ctx := context.Background()

	// Tạo worker pool với priority queue enabled
	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(5),
		WithPriorityEnabled(true),
		WithMetricsEnabled(true),
		WithTimeout(1*time.Minute),
	)

	pool.Start()

	// Submit tasks với priorities khác nhau
	// Priority cao hơn = xử lý trước
	pool.SubmitWithPriority(&ExampleTask{ID: 1, Duration: 100 * time.Millisecond}, 1)  // Low priority
	pool.SubmitWithPriority(&ExampleTask{ID: 2, Duration: 100 * time.Millisecond}, 10) // High priority
	pool.SubmitWithPriority(&ExampleTask{ID: 3, Duration: 100 * time.Millisecond}, 5)  // Medium priority
	pool.SubmitWithPriority(&ExampleTask{ID: 4, Duration: 100 * time.Millisecond}, 15) // Highest priority

	// Shutdown và collect results
	pool.Shutdown()
	results, errors := pool.CollectResults()

	fmt.Printf("Completed: %d, Errors: %d\n", len(results), len(errors))

	// Xem metrics
	if stats := pool.GetMetrics(); stats != nil {
		fmt.Printf("Metrics: %+v\n", stats)
	}
}

// ExampleUsageRetry demonstrates retry mechanism
func ExampleUsageRetry() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(3),
		WithMetricsEnabled(true),
		WithTimeout(1*time.Minute),
	)

	pool.Start()

	// Tạo retry config
	retryConfig := RetryConfig{
		MaxRetries:        3,
		InitialDelay:      50 * time.Millisecond,
		MaxDelay:          500 * time.Millisecond,
		BackoffMultiplier: 2.0,
	}

	// Tạo task có thể fail và retry
	task := &ExampleTask{
		ID:         1,
		Duration:   50 * time.Millisecond,
		ShouldFail: true,
		FailCount:  2, // Sẽ fail 2 lần, sau đó thành công
	}

	// Wrap với retry
	retryTask := WithRetryConfig(task, retryConfig)

	pool.Submit(retryTask)
	pool.Shutdown()

	results, errors := pool.CollectResults()
	fmt.Printf("Results: %v, Errors: %v\n", results, errors)

	// Xem metrics
	if stats := pool.GetMetrics(); stats != nil {
		fmt.Printf("Tasks retried: %d\n", stats.TasksRetried)
	}
}

// ExampleUsageMetrics demonstrates metrics tracking
func ExampleUsageMetrics() {
	ctx := context.Background()

	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(10),
		WithMetricsEnabled(true),
		WithTimeout(1*time.Minute),
	)

	pool.Start()

	// Submit nhiều tasks
	for i := 0; i < 100; i++ {
		task := &ExampleTask{
			ID:       i,
			Duration: time.Duration(i%10) * 10 * time.Millisecond,
		}
		pool.Submit(task)
	}

	pool.Shutdown()
	pool.CollectResults()

	// Xem metrics
	if stats := pool.GetMetrics(); stats != nil {
		fmt.Printf("Tasks Submitted: %d\n", stats.TasksSubmitted)
		fmt.Printf("Tasks Completed: %d\n", stats.TasksCompleted)
		fmt.Printf("Tasks Failed: %d\n", stats.TasksFailed)
		fmt.Printf("Success Rate: %.2f%%\n", stats.SuccessRate)
		fmt.Printf("Avg Process Time: %v\n", stats.AverageProcessTime)
		fmt.Printf("Min Process Time: %v\n", stats.MinProcessTime)
		fmt.Printf("Max Process Time: %v\n", stats.MaxProcessTime)
	}
}

// ExampleUsageCombined demonstrates all features together
func ExampleUsageCombined() {
	ctx := context.Background()

	// Tạo pool với tất cả tính năng
	pool := NewWorkerPool[string](ctx,
		WithMaxWorkers(5),
		WithPriorityEnabled(true),
		WithMetricsEnabled(true),
		WithTimeout(1*time.Minute),
	)

	pool.Start()

	// Retry config
	retryConfig := DefaultRetryConfig()
	retryConfig.MaxRetries = 2

	// Submit tasks với priorities và retry
	for i := 0; i < 20; i++ {
		task := &ExampleTask{
			ID:         i,
			Duration:   50 * time.Millisecond,
			ShouldFail: i%5 == 0, // Một số tasks sẽ fail
			FailCount:  1,
		}

		// Wrap với retry
		retryTask := WithRetryConfig(task, retryConfig)

		// Submit với priority
		priority := i % 10
		pool.SubmitWithPriority(retryTask, priority)
	}

	// Đợi một chút để xem metrics realtime
	time.Sleep(500 * time.Millisecond)
	if stats := pool.GetMetrics(); stats != nil {
		fmt.Printf("Realtime Stats: Active Workers: %d, Queue: %d, Active Tasks: %d\n",
			stats.ActiveWorkers, stats.QueueLength, stats.ActiveTasks)
	}

	pool.Shutdown()
	results, errors := pool.CollectResults()

	fmt.Printf("\nFinal Results:\n")
	fmt.Printf("  Completed: %d\n", len(results))
	fmt.Printf("  Errors: %d\n", len(errors))

	if stats := pool.GetMetrics(); stats != nil {
		fmt.Printf("\nFinal Metrics:\n")
		fmt.Printf("  Submitted: %d\n", stats.TasksSubmitted)
		fmt.Printf("  Completed: %d\n", stats.TasksCompleted)
		fmt.Printf("  Failed: %d\n", stats.TasksFailed)
		fmt.Printf("  Retried: %d\n", stats.TasksRetried)
		fmt.Printf("  Success Rate: %.2f%%\n", stats.SuccessRate)
		fmt.Printf("  Avg Time: %v\n", stats.AverageProcessTime)
	}
}
