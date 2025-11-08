package worker_pool

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics theo dõi hiệu suất và thống kê của worker pool
type Metrics struct {
	// Counters (atomic)
	TasksSubmitted int64 // Tổng số tasks đã submit
	TasksCompleted int64 // Tổng số tasks đã hoàn thành
	TasksFailed    int64 // Tổng số tasks thất bại
	TasksRetried   int64 // Tổng số tasks đã retry

	// Timing (atomic)
	TotalProcessingTime int64 // Tổng thời gian xử lý (nanoseconds)
	MinProcessingTime   int64 // Thời gian xử lý ngắn nhất
	MaxProcessingTime   int64 // Thời gian xử lý dài nhất

	// Current state
	ActiveWorkers int32 // Số workers đang active
	QueueLength   int32 // Số tasks trong queue
	ActiveTasks   int32 // Số tasks đang được xử lý

	// Mutex cho các operations không atomic
	mu sync.RWMutex

	// Histogram (optional, có thể mở rộng)
	processingTimes []time.Duration
}

// NewMetrics tạo Metrics mới
func NewMetrics() *Metrics {
	m := &Metrics{
		MinProcessingTime: int64(time.Hour), // Set giá trị lớn ban đầu
		MaxProcessingTime: 0,
		processingTimes:   make([]time.Duration, 0, 1000), // Pre-allocate
	}
	return m
}

// RecordTaskSubmitted ghi nhận task được submit
func (m *Metrics) RecordTaskSubmitted() {
	atomic.AddInt64(&m.TasksSubmitted, 1)
}

// RecordTaskCompleted ghi nhận task hoàn thành
func (m *Metrics) RecordTaskCompleted(duration time.Duration) {
	atomic.AddInt64(&m.TasksCompleted, 1)

	durationNs := int64(duration)

	// Update total processing time
	atomic.AddInt64(&m.TotalProcessingTime, durationNs)

	// Update min processing time
	for {
		currentMin := atomic.LoadInt64(&m.MinProcessingTime)
		if durationNs >= currentMin {
			break
		}
		if atomic.CompareAndSwapInt64(&m.MinProcessingTime, currentMin, durationNs) {
			break
		}
	}

	// Update max processing time
	for {
		currentMax := atomic.LoadInt64(&m.MaxProcessingTime)
		if durationNs <= currentMax {
			break
		}
		if atomic.CompareAndSwapInt64(&m.MaxProcessingTime, currentMax, durationNs) {
			break
		}
	}

	// Record trong histogram (giới hạn để tránh memory leak)
	m.mu.Lock()
	if len(m.processingTimes) < 10000 {
		m.processingTimes = append(m.processingTimes, duration)
	}
	m.mu.Unlock()
}

// RecordTaskFailed ghi nhận task thất bại
func (m *Metrics) RecordTaskFailed() {
	atomic.AddInt64(&m.TasksFailed, 1)
}

// RecordTaskRetried ghi nhận task được retry
func (m *Metrics) RecordTaskRetried() {
	atomic.AddInt64(&m.TasksRetried, 1)
}

// IncrementActiveWorkers tăng số active workers
func (m *Metrics) IncrementActiveWorkers() {
	atomic.AddInt32(&m.ActiveWorkers, 1)
}

// DecrementActiveWorkers giảm số active workers
func (m *Metrics) DecrementActiveWorkers() {
	atomic.AddInt32(&m.ActiveWorkers, -1)
}

// SetQueueLength set số tasks trong queue
func (m *Metrics) SetQueueLength(length int32) {
	atomic.StoreInt32(&m.QueueLength, length)
}

// IncrementActiveTasks tăng số tasks đang xử lý
func (m *Metrics) IncrementActiveTasks() {
	atomic.AddInt32(&m.ActiveTasks, 1)
}

// DecrementActiveTasks giảm số tasks đang xử lý
func (m *Metrics) DecrementActiveTasks() {
	atomic.AddInt32(&m.ActiveTasks, -1)
}

// GetStats trả về snapshot của metrics
func (m *Metrics) GetStats() Stats {
	submitted := atomic.LoadInt64(&m.TasksSubmitted)
	completed := atomic.LoadInt64(&m.TasksCompleted)
	failed := atomic.LoadInt64(&m.TasksFailed)
	retried := atomic.LoadInt64(&m.TasksRetried)
	totalTime := atomic.LoadInt64(&m.TotalProcessingTime)
	minTime := atomic.LoadInt64(&m.MinProcessingTime)
	maxTime := atomic.LoadInt64(&m.MaxProcessingTime)

	var avgTime time.Duration
	if completed > 0 {
		avgTime = time.Duration(totalTime / completed)
	}

	return Stats{
		TasksSubmitted:     submitted,
		TasksCompleted:     completed,
		TasksFailed:        failed,
		TasksRetried:       retried,
		ActiveWorkers:      atomic.LoadInt32(&m.ActiveWorkers),
		QueueLength:        atomic.LoadInt32(&m.QueueLength),
		ActiveTasks:        atomic.LoadInt32(&m.ActiveTasks),
		AverageProcessTime: avgTime,
		MinProcessTime:     time.Duration(minTime),
		MaxProcessTime:     time.Duration(maxTime),
		SuccessRate:        calculateSuccessRate(completed, failed),
	}
}

// Stats là snapshot của metrics
type Stats struct {
	TasksSubmitted     int64
	TasksCompleted     int64
	TasksFailed        int64
	TasksRetried       int64
	ActiveWorkers      int32
	QueueLength        int32
	ActiveTasks        int32
	AverageProcessTime time.Duration
	MinProcessTime     time.Duration
	MaxProcessTime     time.Duration
	SuccessRate        float64
}

// calculateSuccessRate tính tỷ lệ thành công
func calculateSuccessRate(completed, failed int64) float64 {
	total := completed + failed
	if total == 0 {
		return 0.0
	}
	return float64(completed) / float64(total) * 100.0
}

// Reset reset tất cả metrics về 0
func (m *Metrics) Reset() {
	atomic.StoreInt64(&m.TasksSubmitted, 0)
	atomic.StoreInt64(&m.TasksCompleted, 0)
	atomic.StoreInt64(&m.TasksFailed, 0)
	atomic.StoreInt64(&m.TasksRetried, 0)
	atomic.StoreInt64(&m.TotalProcessingTime, 0)
	atomic.StoreInt64(&m.MinProcessingTime, int64(time.Hour))
	atomic.StoreInt64(&m.MaxProcessingTime, 0)

	m.mu.Lock()
	m.processingTimes = m.processingTimes[:0]
	m.mu.Unlock()
}
