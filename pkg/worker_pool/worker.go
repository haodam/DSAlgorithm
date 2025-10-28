package worker_pool

import (
	"context"
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
		Timeout:     10 * time.Second,
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
