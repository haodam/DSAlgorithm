package worker_pool

import (
	"context"
	"fmt"
	"time"
)

// WaitableTask là wrapper cho Task với result và error channels
type WaitableTask[T any] struct {
	task     Task[T]
	resultCh chan T
	errorCh  chan error
	done     chan struct{} // Signal when task is done
}

// Process thực hiện task và gửi kết quả vào channels
func (wt *WaitableTask[T]) Process(ctx context.Context) (T, error) {
	result, err := wt.task.Process(ctx)

	// Send result or error to private channels first
	if err != nil {
		// Send error to private channel
		select {
		case wt.errorCh <- err:
		case <-ctx.Done():
			// Context cancelled, still send error if possible
			select {
			case wt.errorCh <- err:
			default:
			}
			close(wt.done)
			var zero T
			return zero, ctx.Err()
		}
		close(wt.done)
		// Return error so worker knows to handle it
		var zero T
		return zero, err
	}

	// Send result to private channel
	select {
	case wt.resultCh <- result:
	case <-ctx.Done():
		// Context cancelled, still send result if possible
		select {
		case wt.resultCh <- result:
		default:
		}
		close(wt.done)
		return result, ctx.Err()
	}

	close(wt.done)

	// Return result so worker can also process it normally
	// (though for WaitableTask, we mainly care about private channels)
	return result, nil
}

// SubmitAndWait submits a task and waits for its result or error
// Returns the result and error immediately when task completes
func (wp *WorkerPool[T]) SubmitAndWait(task Task[T]) (T, error) {
	return wp.SubmitAndWaitWithTimeout(task, 0) // 0 means no timeout (use pool timeout)
}

// SubmitAndWaitWithTimeout submits a task and waits for its result or error with a timeout
// If timeout is 0, it uses the pool's timeout
// Returns the result and error when task completes, or timeout error if timeout exceeded
func (wp *WorkerPool[T]) SubmitAndWaitWithTimeout(task Task[T], timeout time.Duration) (T, error) {
	var zero T

	wp.mu.RLock()
	shutdown := wp.shutdown
	started := wp.started
	ctx := wp.ctx
	wp.mu.RUnlock()

	if shutdown {
		return zero, fmt.Errorf("worker pool is shutdown")
	}

	if !started {
		return zero, fmt.Errorf("worker pool not started, call Start() first")
	}

	// Create waitable task
	waitableTask := &WaitableTask[T]{
		task:     task,
		resultCh: make(chan T, 1),
		errorCh:  make(chan error, 1),
		done:     make(chan struct{}),
	}

	// Submit task
	var submitErr error
	if wp.enablePriority && wp.priorityQueue != nil {
		// Use priority queue (priority 0 = normal)
		priorityTask := WithPriority(waitableTask, 0)
		wp.priorityQueue.Push(priorityTask)
	} else {
		// Use regular queue
		select {
		case wp.taskQueue <- waitableTask:
			// Success
		case <-ctx.Done():
			return zero, fmt.Errorf("worker pool context cancelled: %w", ctx.Err())
		default:
			// Buffer full, try blocking send
			select {
			case wp.taskQueue <- waitableTask:
				// Success
			case <-ctx.Done():
				return zero, fmt.Errorf("worker pool context cancelled: %w", ctx.Err())
			}
		}
	}

	if submitErr != nil {
		return zero, submitErr
	}

	// Record metrics
	if wp.enableMetrics && wp.metrics != nil {
		wp.metrics.RecordTaskSubmitted()
	}

	// Wait for result with timeout
	var result T
	var err error

	// Determine timeout context
	waitCtx := ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// Wait for result or error
	select {
	case result = <-waitableTask.resultCh:
		// Success
		err = nil
	case err = <-waitableTask.errorCh:
		// Error
		result = zero
	case <-waitableTask.done:
		// Task done but no result/error received (shouldn't happen, but handle it)
		// Try to get result or error one more time
		select {
		case result = <-waitableTask.resultCh:
			err = nil
		case err = <-waitableTask.errorCh:
			result = zero
		default:
			return zero, fmt.Errorf("task completed but no result or error received")
		}
	case <-waitCtx.Done():
		// Timeout or context cancelled
		return zero, fmt.Errorf("submit and wait timeout or context cancelled: %w", waitCtx.Err())
	}

	// Record metrics
	if wp.enableMetrics && wp.metrics != nil {
		if err != nil {
			wp.metrics.RecordTaskFailed()
		} else {
			// Note: We don't know the exact duration here, so we skip recording
			// The worker will record it when processing
		}
	}

	return result, err
}

// SubmitAndWaitWithContext submits a task and waits for its result or error with a context
// This allows more control over cancellation and timeout
func (wp *WorkerPool[T]) SubmitAndWaitWithContext(ctx context.Context, task Task[T]) (T, error) {
	var zero T

	wp.mu.RLock()
	shutdown := wp.shutdown
	started := wp.started
	poolCtx := wp.ctx
	wp.mu.RUnlock()

	if shutdown {
		return zero, fmt.Errorf("worker pool is shutdown")
	}

	if !started {
		return zero, fmt.Errorf("worker pool not started, call Start() first")
	}

	// Create waitable task
	waitableTask := &WaitableTask[T]{
		task:     task,
		resultCh: make(chan T, 1),
		errorCh:  make(chan error, 1),
		done:     make(chan struct{}),
	}

	// Submit task
	if wp.enablePriority && wp.priorityQueue != nil {
		priorityTask := WithPriority(waitableTask, 0)
		wp.priorityQueue.Push(priorityTask)
	} else {
		select {
		case wp.taskQueue <- waitableTask:
		case <-poolCtx.Done():
			return zero, fmt.Errorf("worker pool context cancelled: %w", poolCtx.Err())
		case <-ctx.Done():
			return zero, fmt.Errorf("context cancelled before submit: %w", ctx.Err())
		default:
			select {
			case wp.taskQueue <- waitableTask:
			case <-poolCtx.Done():
				return zero, fmt.Errorf("worker pool context cancelled: %w", poolCtx.Err())
			case <-ctx.Done():
				return zero, fmt.Errorf("context cancelled before submit: %w", ctx.Err())
			}
		}
	}

	// Record metrics
	if wp.enableMetrics && wp.metrics != nil {
		wp.metrics.RecordTaskSubmitted()
	}

	// Wait for result with context
	var result T
	var err error

	select {
	case result = <-waitableTask.resultCh:
		err = nil
	case err = <-waitableTask.errorCh:
		result = zero
	case <-waitableTask.done:
		// Task done, try to get result or error
		select {
		case result = <-waitableTask.resultCh:
			err = nil
		case err = <-waitableTask.errorCh:
			result = zero
		default:
			return zero, fmt.Errorf("task completed but no result or error received")
		}
	case <-ctx.Done():
		return zero, fmt.Errorf("context cancelled while waiting: %w", ctx.Err())
	case <-poolCtx.Done():
		return zero, fmt.Errorf("worker pool context cancelled: %w", poolCtx.Err())
	}

	// Record metrics
	if wp.enableMetrics && wp.metrics != nil {
		if err != nil {
			wp.metrics.RecordTaskFailed()
		}
	}

	return result, err
}
