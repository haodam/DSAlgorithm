# Worker Pool - SubmitAndWait

Tài liệu về hàm `SubmitAndWait` - submit task và đợi kết quả ngay lập tức.

## 📋 Mục lục

1. [Tổng quan](#tổng-quan)
2. [SubmitAndWait](#submitandwait)
3. [SubmitAndWaitWithTimeout](#submitandwaitwithtimeout)
4. [SubmitAndWaitWithContext](#submitandwaitwithcontext)
5. [So sánh với Submit](#so-sánh-với-submit)
6. [Ví dụ sử dụng](#ví-dụ-sử-dụng)
7. [Best Practices](#best-practices)

---

## 🎯 Tổng quan

`SubmitAndWait` là một phương thức **synchronous** - submit task và đợi kết quả trả về ngay lập tức. Khác với `Submit` (asynchronous), `SubmitAndWait` sẽ block cho đến khi task hoàn thành.

### Khi nào sử dụng?

- ✅ Cần kết quả ngay lập tức
- ✅ Cần xử lý tuần tự (sequential processing)
- ✅ Cần đảm bảo task hoàn thành trước khi tiếp tục
- ✅ Cần error handling ngay lập tức

### Khi nào KHÔNG nên sử dụng?

- ❌ Xử lý hàng loạt tasks (dùng `Submit` + `CollectResults`)
- ❌ Không cần đợi kết quả (dùng `Submit`)
- ❌ Performance-critical với nhiều tasks (dùng `Submit`)

---

## 📝 SubmitAndWait

### Signature

```go
func (wp *WorkerPool[T]) SubmitAndWait(task Task[T]) (T, error)
```

### Mô tả

Submit task và đợi kết quả trả về. Sử dụng timeout của pool.

### Ví dụ

```go
pool := worker_pool.NewWorkerPool[string](ctx,
    worker_pool.WithMaxWorkers(10),
    worker_pool.WithTimeout(10*time.Minute),
)

pool.Start()
defer pool.Shutdown()

// Submit task và đợi kết quả
task := &MyTask{ID: 1}
result, err := pool.SubmitAndWait(task)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Result: %s\n", result)
```

---

## ⏱️ SubmitAndWaitWithTimeout

### Signature

```go
func (wp *WorkerPool[T]) SubmitAndWaitWithTimeout(task Task[T], timeout time.Duration) (T, error)
```

### Mô tả

Submit task và đợi kết quả với timeout cụ thể. Nếu timeout = 0, sử dụng timeout của pool.

### Ví dụ

```go
pool := worker_pool.NewWorkerPool[string](ctx,
    worker_pool.WithMaxWorkers(10),
    worker_pool.WithTimeout(10*time.Minute),
)

pool.Start()
defer pool.Shutdown()

// Submit task với timeout 2 giây
task := &MyTask{ID: 1, Duration: 5 * time.Second}
result, err := pool.SubmitAndWaitWithTimeout(task, 2*time.Second)
if err != nil {
    if strings.Contains(err.Error(), "timeout") {
        fmt.Println("Task timeout!")
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return
}

fmt.Printf("Result: %s\n", result)
```

---

## 🔄 SubmitAndWaitWithContext

### Signature

```go
func (wp *WorkerPool[T]) SubmitAndWaitWithContext(ctx context.Context, task Task[T]) (T, error)
```

### Mô tả

Submit task và đợi kết quả với context. Cho phép cancel hoặc timeout linh hoạt hơn.

### Ví dụ

```go
pool := worker_pool.NewWorkerPool[string](ctx,
    worker_pool.WithMaxWorkers(10),
    worker_pool.WithTimeout(10*time.Minute),
)

pool.Start()
defer pool.Shutdown()

// Tạo context với timeout
waitCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

// Submit task với context
task := &MyTask{ID: 1}
result, err := pool.SubmitAndWaitWithContext(waitCtx, task)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        fmt.Println("Context timeout!")
    } else if errors.Is(err, context.Canceled) {
        fmt.Println("Context cancelled!")
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return
}

fmt.Printf("Result: %s\n", result)
```

---

## 🔀 So sánh với Submit

### Submit (Asynchronous)

```go
// Submit task (fire and forget)
err := pool.Submit(task)
if err != nil {
    fmt.Printf("Submit error: %v\n", err)
}

// Collect results later
results, errors := pool.CollectResults()
```

**Đặc điểm:**
- ✅ Non-blocking
- ✅ High throughput
- ✅ Phù hợp cho batch processing
- ❌ Không biết kết quả ngay

### SubmitAndWait (Synchronous)

```go
// Submit task và đợi kết quả
result, err := pool.SubmitAndWait(task)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Result: %s\n", result)
```

**Đặc điểm:**
- ✅ Biết kết quả ngay
- ✅ Error handling ngay lập tức
- ✅ Phù hợp cho sequential processing
- ❌ Blocking
- ❌ Lower throughput

### Khi nào dùng cái nào?

| Use Case | Method | Lý do |
|----------|--------|-------|
| Batch processing 1000 tasks | `Submit` | High throughput, không cần đợi |
| Process 1 task và cần kết quả | `SubmitAndWait` | Cần kết quả ngay |
| Sequential processing | `SubmitAndWait` | Đảm bảo thứ tự |
| Parallel processing | `Submit` | Tăng tốc độ |

---

## 💡 Ví dụ sử dụng

### Ví dụ 1: Basic Usage

```go
pool := worker_pool.NewWorkerPool[string](ctx,
    worker_pool.WithMaxWorkers(10),
    worker_pool.WithTimeout(10*time.Minute),
)

pool.Start()
defer pool.Shutdown()

task := &MyTask{ID: 1}
result, err := pool.SubmitAndWait(task)
if err != nil {
    log.Fatalf("Task failed: %v", err)
}

fmt.Printf("Task completed: %s\n", result)
```

### Ví dụ 2: With Timeout

```go
pool.Start()
defer pool.Shutdown()

task := &MyTask{ID: 1, Duration: 5 * time.Second}

// Wait với timeout 2 giây
result, err := pool.SubmitAndWaitWithTimeout(task, 2*time.Second)
if err != nil {
    if strings.Contains(err.Error(), "timeout") {
        fmt.Println("Task took too long!")
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return
}

fmt.Printf("Result: %s\n", result)
```

### Ví dụ 3: With Context

```go
pool.Start()
defer pool.Shutdown()

// Tạo context với timeout
waitCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

task := &MyTask{ID: 1}
result, err := pool.SubmitAndWaitWithContext(waitCtx, task)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        fmt.Println("Timeout!")
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return
}

fmt.Printf("Result: %s\n", result)
```

### Ví dụ 4: Error Handling

```go
pool.Start()
defer pool.Shutdown()

task := &ErrorTask{
    ID:  1,
    Msg: "simulated error",
}

result, err := pool.SubmitAndWait(task)
if err != nil {
    fmt.Printf("Task failed: %v\n", err)
    fmt.Printf("Result is zero value: %q\n", result)
    // Handle error...
    return
}

fmt.Printf("Result: %s\n", result)
```

### Ví dụ 5: Multiple Tasks Sequential

```go
pool.Start()
defer pool.Shutdown()

// Process tasks sequentially
for i := 0; i < 10; i++ {
    task := &MyTask{ID: i}
    
    result, err := pool.SubmitAndWait(task)
    if err != nil {
        fmt.Printf("Task %d failed: %v\n", i, err)
        continue
    }
    
    fmt.Printf("Task %d completed: %s\n", i, result)
}
```

### Ví dụ 6: Multiple Tasks Concurrent

```go
pool.Start()
defer pool.Shutdown()

// Process tasks concurrently
var wg sync.WaitGroup
results := make([]string, 10)
errors := make([]error, 10)

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        
        task := &MyTask{ID: id}
        result, err := pool.SubmitAndWait(task)
        
        if err != nil {
            errors[id] = err
        } else {
            results[id] = result
        }
    }(i)
}

wg.Wait()

// Process results
for i, result := range results {
    if errors[i] != nil {
        fmt.Printf("Task %d failed: %v\n", i, errors[i])
    } else {
        fmt.Printf("Task %d completed: %s\n", i, result)
    }
}
```

### Ví dụ 7: Database Operation

```go
type InsertTask struct {
    DB   *sql.DB
    Data string
}

func (t *InsertTask) Process(ctx context.Context) (string, error) {
    _, err := t.DB.ExecContext(ctx, "INSERT INTO users (name) VALUES (?)", t.Data)
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("Inserted: %s", t.Data), nil
}

pool.Start()
defer pool.Shutdown()

task := &InsertTask{
    DB:   db,
    Data: "John Doe",
}

result, err := pool.SubmitAndWait(task)
if err != nil {
    fmt.Printf("Insert failed: %v\n", err)
    return
}

fmt.Printf("Success: %s\n", result)
```

---

## 🎯 Best Practices

### 1. Timeout

Luôn đặt timeout hợp lý:

```go
// ✅ Good: Có timeout
result, err := pool.SubmitAndWaitWithTimeout(task, 5*time.Second)

// ❌ Bad: Không có timeout (có thể đợi vô hạn)
result, err := pool.SubmitAndWait(task)
```

### 2. Error Handling

Luôn kiểm tra error:

```go
result, err := pool.SubmitAndWait(task)
if err != nil {
    // Handle error
    log.Printf("Task failed: %v", err)
    return
}

// Use result
fmt.Printf("Result: %s\n", result)
```

### 3. Context Usage

Sử dụng context cho timeout và cancellation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := pool.SubmitAndWaitWithContext(ctx, task)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // Handle timeout
    } else if errors.Is(err, context.Canceled) {
        // Handle cancellation
    } else {
        // Handle other errors
    }
}
```

### 4. Sequential vs Concurrent

**Sequential** (dùng `SubmitAndWait`):
```go
for i := 0; i < 10; i++ {
    result, err := pool.SubmitAndWait(task)
    // Process result...
}
```

**Concurrent** (dùng `Submit`):
```go
for i := 0; i < 10; i++ {
    pool.Submit(task)
}
results, errors := pool.CollectResults()
```

### 5. Performance

- ❌ **Không dùng** `SubmitAndWait` cho batch processing
- ✅ **Dùng** `SubmitAndWait` khi cần kết quả ngay
- ✅ **Dùng** `Submit` + `CollectResults` cho high throughput

---

## ⚠️ Lưu ý

1. **Blocking**: `SubmitAndWait` sẽ block cho đến khi task hoàn thành
2. **Timeout**: Luôn đặt timeout hợp lý để tránh đợi vô hạn
3. **Performance**: Không dùng cho batch processing (dùng `Submit` thay vào đó)
4. **Context**: Sử dụng context cho timeout và cancellation linh hoạt
5. **Error Handling**: Luôn kiểm tra error trước khi sử dụng result

---

## 📊 So sánh Performance

| Method | Throughput | Latency | Use Case |
|--------|------------|---------|----------|
| `Submit` | High | Low | Batch processing |
| `SubmitAndWait` | Low | High | Single task, need result |
| `Submit` + `CollectResults` | High | Medium | Batch processing with results |

---

## 📝 Tóm tắt

| Method | Description | Timeout | Context |
|--------|-------------|---------|---------|
| `SubmitAndWait` | Submit và đợi kết quả | Pool timeout | Pool context |
| `SubmitAndWaitWithTimeout` | Submit với timeout | Custom timeout | Pool context |
| `SubmitAndWaitWithContext` | Submit với context | Context timeout | Custom context |

Tất cả các methods đã sẵn sàng sử dụng! 🚀
