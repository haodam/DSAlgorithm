package worker_pool

import (
	"sync"
	"time"
)

// ObjectPool quản lý object pool để tái sử dụng objects
var (
	// Pool cho metrics processing times slice
	metricsSlicePool = sync.Pool{
		New: func() interface{} {
			return make([]time.Duration, 0, 100)
		},
	}
)

// calculateOptimalBufferSize tính toán buffer size tối ưu dựa trên số workers
func calculateOptimalBufferSize(workers int) int {
	// Buffer size = workers * 2 để tránh blocking
	// Nhưng không quá lớn để tránh memory waste
	if workers <= 10 {
		return workers * 2
	} else if workers <= 50 {
		return workers * 3
	} else {
		return 150 // Cap at 150 for very large pools
	}
}

// GetMetricsSliceFromPool lấy slice từ pool
func GetMetricsSliceFromPool() []time.Duration {
	return metricsSlicePool.Get().([]time.Duration)
}

// PutMetricsSliceToPool trả slice về pool
func PutMetricsSliceToPool(slice []time.Duration) {
	// Reset slice (keep capacity)
	slice = slice[:0]
	metricsSlicePool.Put(slice)
}

// PerformanceConfig cấu hình hiệu suất
type PerformanceConfig struct {
	BufferMultiplier float64 // Hệ số nhân cho buffer size (default: 2.0)
	EnableObjectPool bool    // Bật object pool (default: true)
}

// DefaultPerformanceConfig trả về config mặc định
func DefaultPerformanceConfig() PerformanceConfig {
	return PerformanceConfig{
		BufferMultiplier: 2.0,
		EnableObjectPool: true,
	}
}
