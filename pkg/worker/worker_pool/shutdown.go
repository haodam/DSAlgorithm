package worker_pool

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ShutdownWithTimeout shuts down the worker pool with a timeout
// Returns error if timeout is exceeded
func (wp *WorkerPool[T]) ShutdownWithTimeout(timeout time.Duration) error {
	wp.mu.Lock()

	if wp.shutdown {
		wp.mu.Unlock()
		return nil // Already shutdown
	}

	wp.shutdown = true

	// Close priority queue if enabled
	if wp.enablePriority && wp.priorityQueue != nil {
		wp.priorityQueue.Close()
	}

	// Close the task queue to signal no more tasks
	close(wp.taskQueue)

	wp.mu.Unlock()

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Wait for workers to finish with timeout
	done := make(chan struct{})
	go func() {
		wp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Workers finished successfully
		wp.mu.Lock()
		defer wp.mu.Unlock()

		// Close result channels
		close(wp.results)
		close(wp.errors)

		// Cancel pool context
		wp.cancel()

		return nil
	case <-ctx.Done():
		// Timeout exceeded - force cancel
		wp.mu.Lock()
		defer wp.mu.Unlock()

		// Cancel pool context to stop workers
		wp.cancel()

		// Try to close channels (may panic if already closed, but that's ok)
		defer func() {
			if r := recover(); r != nil {
				// Channel already closed, ignore
			}
		}()

		close(wp.results)
		close(wp.errors)

		return fmt.Errorf("shutdown timeout exceeded: %v", timeout)
	}
}

// ShutdownGracefully shuts down gracefully, draining remaining tasks first
func (wp *WorkerPool[T]) ShutdownGracefully(drainTimeout time.Duration) error {
	wp.mu.Lock()

	if wp.shutdown {
		wp.mu.Unlock()
		return nil
	}

	wp.shutdown = true

	// Don't close task queue yet - drain first
	wp.mu.Unlock()

	// Drain remaining tasks
	drained := wp.drainTasks(drainTimeout)

	if drained {
		// All tasks drained, proceed with normal shutdown
		return wp.ShutdownWithTimeout(30 * time.Second)
	}

	// Drain timeout exceeded, force shutdown
	return wp.ShutdownWithTimeout(5 * time.Second)
}

// drainTasks drains remaining tasks from the queue
func (wp *WorkerPool[T]) drainTasks(timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	drained := make(chan bool, 1)

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				drained <- false
				return
			case <-ticker.C:
				// Check if queues are empty
				wp.mu.RLock()

				// Check priority queue
				priorityEmpty := true
				if wp.enablePriority && wp.priorityQueue != nil {
					priorityEmpty = wp.priorityQueue.Len() == 0
				}

				// Check regular task queue
				regularEmpty := len(wp.taskQueue) == 0

				wp.mu.RUnlock()

				if priorityEmpty && regularEmpty {
					// Both queues are empty, wait a bit more to ensure no new tasks
					time.Sleep(200 * time.Millisecond)

					// Double check
					wp.mu.RLock()
					priorityEmptyCheck := true
					if wp.enablePriority && wp.priorityQueue != nil {
						priorityEmptyCheck = wp.priorityQueue.Len() == 0
					}
					regularEmptyCheck := len(wp.taskQueue) == 0
					wp.mu.RUnlock()

					if priorityEmptyCheck && regularEmptyCheck {
						drained <- true
						return
					}
				}
			}
		}
	}()

	select {
	case result := <-drained:
		return result
	case <-ctx.Done():
		return false // Timeout
	}
}

// ShutdownOnSignal shuts down gracefully when receiving OS signals (SIGINT, SIGTERM)
func (wp *WorkerPool[T]) ShutdownOnSignal(signals ...os.Signal) error {
	if len(signals) == 0 {
		// Default signals
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)

	// Wait for signal
	<-sigChan

	// Shutdown gracefully
	return wp.ShutdownGracefully(30 * time.Second)
}

// ShutdownOnSignalAsync starts graceful shutdown on signal in a goroutine
func (wp *WorkerPool[T]) ShutdownOnSignalAsync(signals ...os.Signal) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		err := wp.ShutdownOnSignal(signals...)
		errChan <- err
	}()

	return errChan
}

// GetPendingTasksCount returns the number of pending tasks in the queue
func (wp *WorkerPool[T]) GetPendingTasksCount() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if wp.shutdown {
		return 0
	}

	if wp.enablePriority && wp.priorityQueue != nil {
		return wp.priorityQueue.Len()
	}

	// For regular queue, we can only estimate
	return len(wp.taskQueue)
}

// IsShuttingDown returns true if the pool is in the process of shutting down
func (wp *WorkerPool[T]) IsShuttingDown() bool {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.shutdown
}
