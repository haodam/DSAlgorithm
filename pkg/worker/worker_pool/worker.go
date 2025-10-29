package worker_pool

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Task represents a unit of work that can be processed
type Task[T any] interface {
	Process(ctx context.Context) (T, error)
}

// WorkerPool manages a pool of workers for concurrent task processing
type WorkerPool[T any] struct {
	maxWorkers  int                // maximum number of concurrent workers
	taskQueue   chan Task[T]       // channel for tasks
	results     chan T             // channel for results
	errors      chan error         // channel for errors
	timeout     time.Duration      // timeout for task processing
	stopOnError bool               // whether to stop processing on the first error
	wg          sync.WaitGroup     // wait group for workers
	ctx         context.Context    // context for pool
	cancel      context.CancelFunc // cancel pool context
}

// WorkerPoolOption represents configuration options for WorkerPool
type WorkerPoolOption func(*WorkerPoolConfig)

// WorkerPoolConfig holds configuration for WorkerPool
type WorkerPoolConfig struct {
	MaxWorkers  int
	Timeout     time.Duration
	StopOnError bool
}

// WithMaxWorkers sets the maximum number of concurrent workers
func WithMaxWorkers(n int) WorkerPoolOption {
	return func(c *WorkerPoolConfig) {
		c.MaxWorkers = n
	}
}

// WithTimeout sets the timeout for task processing
func WithTimeout(d time.Duration) WorkerPoolOption {
	return func(c *WorkerPoolConfig) {
		c.Timeout = d
	}
}

// WithStopOnError sets whether to stop processing on the first error
func WithStopOnError(b bool) WorkerPoolOption {
	return func(c *WorkerPoolConfig) {
		c.StopOnError = b
	}
}

// NewWorkerPool creates a new WorkerPool with the given options
func NewWorkerPool[T any](ctx context.Context, options ...WorkerPoolOption) *WorkerPool[T] {
	config := &WorkerPoolConfig{
		MaxWorkers:  10,
		Timeout:     10 * time.Minute,
		StopOnError: false,
	}
	for _, option := range options {
		option(config)
	}
	ctx, cancel := context.WithTimeout(ctx, config.Timeout)

	wp := &WorkerPool[T]{
		maxWorkers:  config.MaxWorkers,
		taskQueue:   make(chan Task[T]),
		results:     make(chan T, config.MaxWorkers),
		errors:      make(chan error, config.MaxWorkers),
		timeout:     config.Timeout,
		stopOnError: config.StopOnError,
		ctx:         ctx,
		cancel:      cancel,
	}

	return wp
}

// NewDefaultExecutor returns a worker_pool pool sized for mixed workloads.
// Default uses ~2 * GOMAXPROCS goroutines, which is a reasonable balance
// for mixed CPU+IO-bound tasks without starving CPU-bound work.
func NewDefaultExecutor[T any](ctx context.Context, opts ...WorkerPoolOption) *WorkerPool[T] {
	defaultWorkers := 2 * runtime.GOMAXPROCS(0)
	opts = append([]WorkerPoolOption{WithMaxWorkers(defaultWorkers)}, opts...)
	return NewWorkerPool[T](ctx, opts...)
}

// NewIOExecutor returns a worker_pool pool optimized for IO-bound tasks.
// IO-bound tasks spend most of their time waiting, so more goroutines
// can be useful. We choose 4 * GOMAXPROCS by default to provide
// plenty of concurrency without unbounded goroutine growth.
func NewIOExecutor[T any](ctx context.Context, opts ...WorkerPoolOption) *WorkerPool[T] {
	ioWorkers := 4 * runtime.GOMAXPROCS(0)
	opts = append([]WorkerPoolOption{WithMaxWorkers(ioWorkers)}, opts...)
	return NewWorkerPool[T](ctx, opts...)
}

// NewCPUExecutor returns a worker_pool pool optimized for CPU-bound tasks.
// CPU-bound tasks should generally be limited to GOMAXPROCS goroutines
// to avoid excess context switching. We default to GOMAXPROCS.
func NewCPUExecutor[T any](ctx context.Context, opts ...WorkerPoolOption) *WorkerPool[T] {
	cpuWorkers := runtime.GOMAXPROCS(0)
	opts = append([]WorkerPoolOption{WithMaxWorkers(cpuWorkers)}, opts...)
	return NewWorkerPool[T](ctx, opts...)
}

// Submit adds a task to the worker_pool pool
func (wp *WorkerPool[T]) Submit(task Task[T]) error {
	// Blocking sending: wait until a worker_pool receives the task or pool context is done
	select {
	case wp.taskQueue <- task:
		return nil
	case <-wp.ctx.Done():
		return fmt.Errorf("worker_pool pool is closed or context cancelled")
	}
}

// Start begins processing tasks
func (wp *WorkerPool[T]) Start() {
	for i := 0; i < wp.maxWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// worker processes tasks
func (wp *WorkerPool[T]) worker() {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case task, ok := <-wp.taskQueue:
			if !ok {
				return
			}

			result, err := task.Process(wp.ctx)
			if err != nil {
				select {
				case wp.errors <- err:
				case <-wp.ctx.Done():
					return
				}

				if wp.stopOnError {
					wp.cancel()
					return
				}
				continue
			}

			select {
			case wp.results <- result:
			case <-wp.ctx.Done():
				return
			}
		}
	}
}

// Results returns channels for results and errors
func (wp *WorkerPool[T]) Results() (<-chan T, <-chan error) {
	return wp.results, wp.errors
}

// Wait waits for worker_pool goroutines to finish processing current tasks.
func (wp *WorkerPool[T]) Wait() {
	wp.wg.Wait()
}

// Shutdown stops accepting new submits, waits for in-flight Submit calls,
// closes the task queue to signal workers, waits for workers to finish and
// then closes result and error channels.
func (wp *WorkerPool[T]) Shutdown() {
	// Close the task queue to signal no more tasks.
	close(wp.taskQueue)

	// Wait for workers to finish
	wp.wg.Wait()

	// Close result channels
	close(wp.results)
	close(wp.errors)

	// Cancel pool context
	wp.cancel()
}

// Stop is an alias for Shutdown for backward-compatibility
func (wp *WorkerPool[T]) Stop() {
	wp.Shutdown()
}

// CollectResults collects all results and errors until the worker_pool pool is done
func (wp *WorkerPool[T]) CollectResults() ([]T, []error) {
	var results []T
	var errors []error

	resultsC, errorsC := wp.Results()

	for {
		select {
		case result, ok := <-resultsC:
			if !ok {
				return results, errors
			}
			results = append(results, result)
		case err, ok := <-errorsC:
			if !ok {
				return results, errors
			}
			errors = append(errors, err)
		}
	}
}
