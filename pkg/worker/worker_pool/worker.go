package worker_pool

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Ensure ShutdownWithTimeout is available (defined in shutdown.go)
// This is a forward declaration for documentation

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
	mu          sync.RWMutex       // mutex for thread safety
	started     bool               // flag to track if the pool is started
	shutdown    bool               // flag to track if the pool is shutdown

	// Enhanced features
	priorityQueue  *PriorityQueueManager[T] // Priority queue (optional)
	metrics        *Metrics                 // Metrics tracking (optional)
	enablePriority bool                     // Enable priority queue
	enableMetrics  bool                     // Enable metrics
}

// WorkerPoolOption represents configuration options for WorkerPool
type WorkerPoolOption func(*WorkerPoolConfig)

// WorkerPoolConfig holds configuration for WorkerPool
type WorkerPoolConfig struct {
	MaxWorkers     int
	Timeout        time.Duration
	StopOnError    bool
	EnablePriority bool // Bật priority queue
	EnableMetrics  bool // Bật metrics tracking
}

// resultTask wraps a Task with individual result and error channels

type ResultTask[T any] struct {
	task     Task[T]
	resultCh chan T
	errorCh  chan error
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

// WithPriorityEnabled enables priority queue
func WithPriorityEnabled(enable bool) WorkerPoolOption {
	return func(c *WorkerPoolConfig) {
		c.EnablePriority = enable
	}
}

// WithMetricsEnabled enables metrics tracking
func WithMetricsEnabled(enable bool) WorkerPoolOption {
	return func(c *WorkerPoolConfig) {
		c.EnableMetrics = enable
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

	// Validation
	if config.MaxWorkers <= 0 {
		config.MaxWorkers = 10 // Default to 10 if invalid
	}
	if config.Timeout <= 0 {
		config.Timeout = 10 * time.Minute // Default to 10 minutes if invalid
	}

	ctx, cancel := context.WithTimeout(ctx, config.Timeout)

	// Calculate optimal buffer sizes
	bufferSize := calculateOptimalBufferSize(config.MaxWorkers)

	// For priority queue, use unbuffered channel
	// For regular queue, use a small buffer to avoid blocking on submitting
	taskQueueSize := 0
	if !config.EnablePriority {
		// Small buffer for regular queue to reduce blocking
		taskQueueSize = config.MaxWorkers
	}

	wp := &WorkerPool[T]{
		maxWorkers:     config.MaxWorkers,
		taskQueue:      make(chan Task[T], taskQueueSize),
		results:        make(chan T, bufferSize),
		errors:         make(chan error, bufferSize),
		timeout:        config.Timeout,
		stopOnError:    config.StopOnError,
		ctx:            ctx,
		cancel:         cancel,
		started:        false,
		shutdown:       false,
		enablePriority: config.EnablePriority,
		enableMetrics:  config.EnableMetrics,
	}

	// Initialize priority queue if enabled
	if config.EnablePriority {
		wp.priorityQueue = NewPriorityQueueManager[T]()
	}

	// Initialize metrics if enabled
	if config.EnableMetrics {
		wp.metrics = NewMetrics()
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

// SubmitAsync submits a task to be processed asynchronously (fire-and-forget)
func (wp *WorkerPool[T]) SubmitAsync(task Task[T]) error {
	return wp.Submit(task)
}

// SubmitAndWait is implemented in submit_and_wait.go
// SubmitAndWait submits a task and waits for its result or error
// See submit_and_wait.go for implementation

// Submit adds a task to the worker_pool pool
func (wp *WorkerPool[T]) Submit(task Task[T]) error {
	return wp.SubmitWithPriority(task, 0)
}

// SubmitWithPriority adds a task with priority to the worker pool
func (wp *WorkerPool[T]) SubmitWithPriority(task Task[T], priority int) error {
	wp.mu.RLock()
	shutdown := wp.shutdown
	started := wp.started
	wp.mu.RUnlock()

	if shutdown {
		return fmt.Errorf("worker pool is shutdown")
	}

	if !started {
		return fmt.Errorf("worker pool not started, call Start() first")
	}

	// Use priority queue if enabled
	if wp.enablePriority && wp.priorityQueue != nil {
		// Use object pool to reduce allocations
		priorityTask := WithPriority(task, priority)
		wp.priorityQueue.Push(priorityTask)

		// Record metrics
		if wp.enableMetrics && wp.metrics != nil {
			wp.metrics.RecordTaskSubmitted()
			wp.metrics.SetQueueLength(int32(wp.priorityQueue.Len()))
		}
		return nil
	}

	// Record metrics
	if wp.enableMetrics && wp.metrics != nil {
		wp.metrics.RecordTaskSubmitted()
	}

	// Try non-blocking send first (if buffer available)
	select {
	case wp.taskQueue <- task:
		return nil
	case <-wp.ctx.Done():
		return fmt.Errorf("worker pool is closed or context cancelled")
	default:
		// Buffer full, try blocking send
		select {
		case wp.taskQueue <- task:
			return nil
		case <-wp.ctx.Done():
			return fmt.Errorf("worker pool is closed or context cancelled")
		}
	}
}

// Start begins processing tasks
func (wp *WorkerPool[T]) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.started {
		return // Already started, ignore
	}

	if wp.shutdown {
		return // Cannot start if already shutdown
	}

	wp.started = true

	// Start priority queue processor if enabled
	if wp.enablePriority && wp.priorityQueue != nil {
		wp.wg.Add(1)
		go wp.priorityQueueProcessor()
	}

	for i := 0; i < wp.maxWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

	// Record metrics
	if wp.enableMetrics && wp.metrics != nil {
		for i := 0; i < wp.maxWorkers; i++ {
			wp.metrics.IncrementActiveWorkers()
		}
	}
}

// priorityQueueProcessor xử lý priority queue và đưa tasks vào taskQueue
func (wp *WorkerPool[T]) priorityQueueProcessor() {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		default:
			// Try non-blocking pop first
			priorityTask, ok := wp.priorityQueue.TryPop()
			if !ok {
				// Queue empty, wait a bit to avoid busy waiting
				select {
				case <-wp.ctx.Done():
					return
				case <-time.After(10 * time.Millisecond):
					continue
				}
			}

			// Đưa task vào taskQueue để worker xử lý
			select {
			case wp.taskQueue <- priorityTask.Task:
				// Update metrics
				if wp.enableMetrics && wp.metrics != nil {
					wp.metrics.SetQueueLength(int32(wp.priorityQueue.Len()))
				}
			case <-wp.ctx.Done():
				return
			}
		}
	}
}

// worker processes tasks
func (wp *WorkerPool[T]) worker() {
	defer wp.wg.Done()

	if wp.enableMetrics && wp.metrics != nil {
		defer wp.metrics.DecrementActiveWorkers()
	}

	for {
		var task Task[T]
		var ok bool

		select {
		case <-wp.ctx.Done():
			return
		case task, ok = <-wp.taskQueue:
			if !ok {
				return
			}
		}

		// Record metrics
		if wp.enableMetrics && wp.metrics != nil {
			wp.metrics.IncrementActiveTasks()
			startTime := time.Now()
			defer func() {
				wp.metrics.DecrementActiveTasks()
				duration := time.Since(startTime)
				// Will be recorded after processing
				_ = duration
			}()
		}

		startTime := time.Now()
		result, err := task.Process(wp.ctx)
		duration := time.Since(startTime)

		// Check if this is a WaitableTask (SubmitAndWait)
		// WaitableTask handles its own result/error channels, so we don't need to send it to pool channels
		_, isWaitableTask := task.(*WaitableTask[T])

		if err != nil {
			// Check if it's a retry task
			if retryTask, ok := task.(*RetryTask[T]); ok {
				if retryTask.GetRetryCount() > 0 {
					if wp.enableMetrics && wp.metrics != nil {
						wp.metrics.RecordTaskRetried()
					}
				}
			}

			// Record metrics
			if wp.enableMetrics && wp.metrics != nil {
				wp.metrics.RecordTaskFailed()
				wp.metrics.RecordTaskCompleted(duration)
			}

			// Only send it to a pool error channel if not a WaitableTask
			//  handles its own error channel
			if !isWaitableTask {
				select {
				case wp.errors <- err:
				case <-wp.ctx.Done():
					return
				}
			}

			if wp.stopOnError {
				wp.cancel()
				return
			}
			continue
		}

		// Record metrics
		if wp.enableMetrics && wp.metrics != nil {
			wp.metrics.RecordTaskCompleted(duration)
		}

		// Only send it to a pool result channel if not a WaitableTask
		//  handles its own result channel
		if !isWaitableTask {
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
// This is equivalent to ShutdownWithTimeout with a very long timeout (1 hour)
func (wp *WorkerPool[T]) Shutdown() {
	// Use a very long timeout (essentially unlimited)
	_ = wp.ShutdownWithTimeout(1 * time.Hour)
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

// GetMetrics trả về metrics snapshot (chỉ khi metrics được enable)
func (wp *WorkerPool[T]) GetMetrics() *Stats {
	if !wp.enableMetrics || wp.metrics == nil {
		return nil
	}
	stats := wp.metrics.GetStats()
	return &stats
}
