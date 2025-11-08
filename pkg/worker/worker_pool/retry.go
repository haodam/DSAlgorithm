package worker_pool

import (
	"context"
	"fmt"
	"time"
)

// RetryConfig cấu hình retry mechanism
type RetryConfig struct {
	MaxRetries        int           // Số lần retry tối đa
	InitialDelay      time.Duration // Thời gian delay ban đầu
	MaxDelay          time.Duration // Thời gian delay tối đa
	BackoffMultiplier float64       // Hệ số nhân cho exponential backoff
	RetryableErrors   []error       // Danh sách errors có thể retry (optional)
}

// DefaultRetryConfig trả về cấu hình retry mặc định
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:        3,
		InitialDelay:      100 * time.Millisecond,
		MaxDelay:          5 * time.Second,
		BackoffMultiplier: 2.0,
		RetryableErrors:   nil, // Retry tất cả errors
	}
}

// RetryTask là wrapper cho Task với retry mechanism
type RetryTask[T any] struct {
	Task       Task[T]
	Config     RetryConfig
	RetryCount int
	LastError  error
}

// NewRetryTask tạo RetryTask mới
func NewRetryTask[T any](task Task[T], config RetryConfig) *RetryTask[T] {
	return &RetryTask[T]{
		Task:       task,
		Config:     config,
		RetryCount: 0,
	}
}

// Process thực hiện task với retry logic
func (rt *RetryTask[T]) Process(ctx context.Context) (T, error) {
	var zero T
	var result T
	var err error

	for rt.RetryCount <= rt.Config.MaxRetries {
		// Thực hiện task
		result, err = rt.Task.Process(ctx)

		// Nếu thành công, trả về kết quả
		if err == nil {
			return result, nil
		}

		// Kiểm tra xem error có thể retry không
		if !rt.isRetryable(err) {
			return zero, fmt.Errorf("non-retryable error: %w", err)
		}

		// Kiểm tra đã đạt max retries chưa
		if rt.RetryCount >= rt.Config.MaxRetries {
			rt.LastError = err
			return zero, fmt.Errorf("max retries (%d) exceeded, last error: %w", rt.Config.MaxRetries, err)
		}

		rt.RetryCount++
		rt.LastError = err

		// Tính toán delay với exponential backoff
		delay := rt.calculateDelay()

		// Đợi trước khi retry
		select {
		case <-ctx.Done():
			return zero, fmt.Errorf("context cancelled during retry: %w", ctx.Err())
		case <-time.After(delay):
			// Tiếp tục retry
		}
	}

	return zero, fmt.Errorf("unexpected retry loop exit: %w", rt.LastError)
}

// isRetryable kiểm tra xem error có thể retry không
func (rt *RetryTask[T]) isRetryable(err error) bool {
	// Nếu không có danh sách errors cụ thể, retry tất cả
	if len(rt.Config.RetryableErrors) == 0 {
		return true
	}

	// Kiểm tra xem error có trong danh sách không
	for _, retryableErr := range rt.Config.RetryableErrors {
		if err == retryableErr || fmt.Sprintf("%v", err) == fmt.Sprintf("%v", retryableErr) {
			return true
		}
	}

	return false
}

// calculateDelay tính toán delay với exponential backoff
func (rt *RetryTask[T]) calculateDelay() time.Duration {
	delay := float64(rt.Config.InitialDelay) *
		rt.Config.BackoffMultiplier * float64(rt.RetryCount-1)

	if delay > float64(rt.Config.MaxDelay) {
		delay = float64(rt.Config.MaxDelay)
	}

	return time.Duration(delay)
}

// GetRetryCount trả về số lần đã retry
func (rt *RetryTask[T]) GetRetryCount() int {
	return rt.RetryCount
}

// GetLastError trả về error cuối cùng
func (rt *RetryTask[T]) GetLastError() error {
	return rt.LastError
}

// WithRetry tạo RetryTask từ Task thông thường với config mặc định
func WithRetry[T any](task Task[T]) *RetryTask[T] {
	return NewRetryTask(task, DefaultRetryConfig())
}

// WithRetryConfig tạo RetryTask với config tùy chỉnh
func WithRetryConfig[T any](task Task[T], config RetryConfig) *RetryTask[T] {
	return NewRetryTask(task, config)
}
