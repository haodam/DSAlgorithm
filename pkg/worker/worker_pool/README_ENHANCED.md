# Worker Pool - Enhanced Features

Worker Pool với các tính năng nâng cao: **Priority Queue**, **Retry Mechanism**, và **Metrics**.

## 📋 Mục lục

1. [Priority Queue](#priority-queue)
2. [Retry Mechanism](#retry-mechanism)
3. [Metrics & Monitoring](#metrics--monitoring)
4. [Ví dụ sử dụng](#ví-dụ-sử-dụng)

---

## 🎯 Priority Queue

### Tổng quan
Priority queue cho phép tasks có priority cao hơn được xử lý trước tasks có priority thấp hơn.

### Cách sử dụng

```go
// Tạo worker pool với priority enabled
pool := worker_pool.NewWorkerPool[string](ctx,
    worker_pool.WithMaxWorkers(5),
    worker_pool.WithPriorityEnabled(true), // Bật priority queue
    worker_pool.WithTimeout(1*time.Minute),
)

pool.Start()

// Submit tasks với priorities khác nhau
// Priority cao hơn = xử lý trước (ví dụ: 10 > 5 > 1)
pool.SubmitWithPriority(task1, 1)  // Low priority
pool.SubmitWithPriority(task2, 5)  // Medium priority
pool.SubmitWithPriority(task3, 10) // High priority

// Tasks sẽ được xử lý theo thứ tự: task3 -> task2 -> task1
```

### Lưu ý
- Priority **cao hơn** = xử lý **trước**
- Nếu không bật priority, `SubmitWithPriority()` vẫn hoạt động nhưng sẽ bỏ qua priority
- Priority queue sử dụng max heap, đảm bảo O(log n) cho insert và pop

---

## 🔄 Retry Mechanism

### Tổng quan
Retry mechanism tự động retry tasks khi chúng fail, với exponential backoff.

### Cách sử dụng

#### 1. Sử dụng retry config mặc định

```go
task := &MyTask{...}

// Wrap task với retry (mặc định: 3 retries)
retryTask := worker_pool.WithRetry(task)

pool.Submit(retryTask)
```

#### 2. Tùy chỉnh retry config

```go
retryConfig := worker_pool.RetryConfig{
    MaxRetries:       5,                      // Retry tối đa 5 lần
    InitialDelay:     100 * time.Millisecond, // Delay ban đầu
    MaxDelay:         5 * time.Second,        // Delay tối đa
    BackoffMultiplier: 2.0,                   // Hệ số nhân (exponential)
    RetryableErrors:  []error{...},           // Chỉ retry các errors này (optional)
}

retryTask := worker_pool.WithRetryConfig(task, retryConfig)
pool.Submit(retryTask)
```

#### 3. Ví dụ với database insert

```go
type InsertTask struct {
    DB   *sql.DB
    Data string
}

func (t *InsertTask) Process(ctx context.Context) (string, error) {
    _, err := t.DB.ExecContext(ctx, "INSERT INTO users (name) VALUES (?)", t.Data)
    return fmt.Sprintf("Inserted: %s", t.Data), err
}

// Tạo retry config cho database operations
retryConfig := worker_pool.RetryConfig{
    MaxRetries:       3,
    InitialDelay:     200 * time.Millisecond,
    MaxDelay:         2 * time.Second,
    BackoffMultiplier: 2.0,
}

task := &InsertTask{DB: db, Data: "John Doe"}
retryTask := worker_pool.WithRetryConfig(task, retryConfig)

pool.Submit(retryTask)
```

### Exponential Backoff
Delay được tính toán như sau:
```
delay = InitialDelay * (BackoffMultiplier ^ (retryCount - 1))
delay = min(delay, MaxDelay)
```

Ví dụ với `InitialDelay = 100ms`, `BackoffMultiplier = 2.0`:
- Retry 1: 100ms
- Retry 2: 200ms
- Retry 3: 400ms
- Retry 4: 800ms
- ...

---

## 📊 Metrics & Monitoring

### Tổng quan
Metrics tracking cung cấp thống kê chi tiết về hiệu suất của worker pool.

### Cách sử dụng

```go
// Tạo worker pool với metrics enabled
pool := worker_pool.NewWorkerPool[string](ctx,
    worker_pool.WithMaxWorkers(10),
    worker_pool.WithMetricsEnabled(true), // Bật metrics
    worker_pool.WithTimeout(1*time.Minute),
)

pool.Start()

// ... submit tasks ...

// Lấy metrics
if stats := pool.GetMetrics(); stats != nil {
    fmt.Printf("Tasks Submitted: %d\n", stats.TasksSubmitted)
    fmt.Printf("Tasks Completed: %d\n", stats.TasksCompleted)
    fmt.Printf("Tasks Failed: %d\n", stats.TasksFailed)
    fmt.Printf("Tasks Retried: %d\n", stats.TasksRetried)
    fmt.Printf("Success Rate: %.2f%%\n", stats.SuccessRate)
    fmt.Printf("Average Process Time: %v\n", stats.AverageProcessTime)
    fmt.Printf("Min Process Time: %v\n", stats.MinProcessTime)
    fmt.Printf("Max Process Time: %v\n", stats.MaxProcessTime)
    fmt.Printf("Active Workers: %d\n", stats.ActiveWorkers)
    fmt.Printf("Queue Length: %d\n", stats.QueueLength)
    fmt.Printf("Active Tasks: %d\n", stats.ActiveTasks)
}
```

### Metrics được track

| Metric | Mô tả |
|--------|-------|
| `TasksSubmitted` | Tổng số tasks đã submit |
| `TasksCompleted` | Tổng số tasks đã hoàn thành thành công |
| `TasksFailed` | Tổng số tasks thất bại |
| `TasksRetried` | Tổng số tasks đã retry |
| `SuccessRate` | Tỷ lệ thành công (%) |
| `AverageProcessTime` | Thời gian xử lý trung bình |
| `MinProcessTime` | Thời gian xử lý ngắn nhất |
| `MaxProcessTime` | Thời gian xử lý dài nhất |
| `ActiveWorkers` | Số workers đang active |
| `QueueLength` | Số tasks trong queue |
| `ActiveTasks` | Số tasks đang được xử lý |

---

## 🎨 Ví dụ sử dụng

### Ví dụ 1: Priority + Retry + Metrics

```go
ctx := context.Background()

// Tạo pool với tất cả tính năng
pool := worker_pool.NewWorkerPool[string](ctx,
    worker_pool.WithMaxWorkers(5),
    worker_pool.WithPriorityEnabled(true),
    worker_pool.WithMetricsEnabled(true),
    worker_pool.WithTimeout(1*time.Minute),
)

pool.Start()

// Retry config
retryConfig := worker_pool.DefaultRetryConfig()
retryConfig.MaxRetries = 3

// Submit tasks với priority và retry
for i := 0; i < 100; i++ {
    task := &MyTask{ID: i}
    retryTask := worker_pool.WithRetryConfig(task, retryConfig)
    
    priority := i % 10 // Priority từ 0-9
    pool.SubmitWithPriority(retryTask, priority)
}

// Xem metrics realtime
time.Sleep(1 * time.Second)
if stats := pool.GetMetrics(); stats != nil {
    fmt.Printf("Realtime: Active Workers: %d, Queue: %d\n",
        stats.ActiveWorkers, stats.QueueLength)
}

// Shutdown và collect results
pool.Shutdown()
results, errors := pool.CollectResults()

// Xem final metrics
if stats := pool.GetMetrics(); stats != nil {
    fmt.Printf("Final Stats: %+v\n", stats)
}
```

### Ví dụ 2: Database Batch Insert với Retry

```go
type InsertDB struct {
    DB   *sql.DB
    Data string
}

func (t *InsertDB) Process(ctx context.Context) (string, error) {
    _, err := t.DB.ExecContext(ctx, "INSERT INTO users (name) VALUES (?)", t.Data)
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("Inserted: %s", t.Data), nil
}

// Tạo pool
pool := worker_pool.NewWorkerPool[string](ctx,
    worker_pool.WithMaxWorkers(10),
    worker_pool.WithMetricsEnabled(true),
)

pool.Start()

// Retry config cho database
retryConfig := worker_pool.RetryConfig{
    MaxRetries:       3,
    InitialDelay:     200 * time.Millisecond,
    MaxDelay:         2 * time.Second,
    BackoffMultiplier: 2.0,
}

// Submit tasks với retry
for _, name := range names {
    task := &InsertDB{DB: db, Data: name}
    retryTask := worker_pool.WithRetryConfig(task, retryConfig)
    pool.Submit(retryTask)
}

pool.Shutdown()
results, errors := pool.CollectResults()

// Xem metrics
if stats := pool.GetMetrics(); stats != nil {
    fmt.Printf("Inserted: %d, Failed: %d, Retried: %d\n",
        stats.TasksCompleted, stats.TasksFailed, stats.TasksRetried)
}
```

---

## 🔧 Best Practices

### 1. Priority Queue
- Sử dụng priority queue khi có tasks quan trọng cần xử lý trước
- Priority nên có range hợp lý (ví dụ: 0-100) để dễ quản lý
- Tránh thay đổi priority quá thường xuyên

### 2. Retry Mechanism
- Chỉ retry các errors có thể recover (network, temporary DB errors)
- Không retry các errors không thể recover (validation errors, permission denied)
- Điều chỉnh `MaxRetries` và `MaxDelay` dựa trên use case
- Với database operations, thường retry 3-5 lần với delay 200ms-2s

### 3. Metrics
- Enable metrics trong production để monitor hiệu suất
- Xem metrics định kỳ để phát hiện bottlenecks
- Sử dụng metrics để tune `MaxWorkers` và các config khác

---

## ⚠️ Lưu ý

1. **Backward Compatibility**: Tất cả các tính năng mới đều optional và không ảnh hưởng đến code cũ
2. **Performance**: Priority queue và metrics có overhead nhỏ, chỉ enable khi cần
3. **Memory**: Metrics lưu processing times (giới hạn 10,000) để tránh memory leak
4. **Thread Safety**: Tất cả các tính năng đều thread-safe

---

## 📝 Tóm tắt

| Tính năng | File | Status |
|-----------|------|--------|
| Priority Queue | `priority.go` | ✅ Hoàn thành |
| Retry Mechanism | `retry.go` | ✅ Hoàn thành |
| Metrics | `metrics.go` | ✅ Hoàn thành |
| Integration | `worker.go` | ✅ Hoàn thành |
| Examples | `example_enhanced.go` | ✅ Hoàn thành |

Tất cả các tính năng đã sẵn sàng sử dụng! 🚀
