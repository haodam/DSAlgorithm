package worker_pool

import (
	"container/heap"
	"context"
	"sync"
)

// PriorityTask là wrapper cho Task với priority
type PriorityTask[T any] struct {
	Task     Task[T]
	Priority int // Cao hơn = ưu tiên hơn
	Index    int // Index trong heap
}

// PriorityQueue implements heap.Interface
type PriorityQueue[T any] []*PriorityTask[T]

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	// Priority cao hơn sẽ ở đầu queue (max heap)
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue[T]) Push(x interface{}) {
	n := len(*pq)
	task := x.(*PriorityTask[T])
	task.Index = n
	*pq = append(*pq, task)
}

func (pq *PriorityQueue[T]) Pop() interface{} {
	old := *pq
	n := len(old)
	task := old[n-1]
	old[n-1] = nil
	task.Index = -1
	*pq = old[0 : n-1]
	return task
}

// PriorityQueueManager quản lý priority queue với thread safety
type PriorityQueueManager[T any] struct {
	pq   PriorityQueue[T]
	mu   sync.Mutex
	cond *sync.Cond
}

// NewPriorityQueueManager tạo priority queue manager mới
func NewPriorityQueueManager[T any]() *PriorityQueueManager[T] {
	pqm := &PriorityQueueManager[T]{
		pq: make(PriorityQueue[T], 0),
	}
	pqm.cond = sync.NewCond(&pqm.mu)
	heap.Init(&pqm.pq)
	return pqm
}

// Push thêm task vào queue
func (pqm *PriorityQueueManager[T]) Push(task *PriorityTask[T]) {
	pqm.mu.Lock()
	defer pqm.mu.Unlock()
	heap.Push(&pqm.pq, task)
	pqm.cond.Signal()
}

// Pop lấy task có priority cao nhất (blocking)
func (pqm *PriorityQueueManager[T]) Pop() *PriorityTask[T] {
	pqm.mu.Lock()
	defer pqm.mu.Unlock()

	for pqm.pq.Len() == 0 {
		pqm.cond.Wait()
	}

	return heap.Pop(&pqm.pq).(*PriorityTask[T])
}

// PopWithContext lấy task có priority cao nhất với context (có thể cancel)
func (pqm *PriorityQueueManager[T]) PopWithContext(ctx context.Context) (*PriorityTask[T], error) {
	pqm.mu.Lock()
	defer pqm.mu.Unlock()

	for pqm.pq.Len() == 0 {
		// Wait for signal or context cancellation
		waitChan := make(chan struct{})
		go func() {
			pqm.cond.Wait()
			close(waitChan)
		}()

		select {
		case <-waitChan:
			// Got signal, continue loop
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return heap.Pop(&pqm.pq).(*PriorityTask[T]), nil
}

// TryPop thử pop task không block
func (pqm *PriorityQueueManager[T]) TryPop() (*PriorityTask[T], bool) {
	pqm.mu.Lock()
	defer pqm.mu.Unlock()

	if pqm.pq.Len() == 0 {
		return nil, false
	}

	return heap.Pop(&pqm.pq).(*PriorityTask[T]), true
}

// Len trả về số lượng tasks trong queue
func (pqm *PriorityQueueManager[T]) Len() int {
	pqm.mu.Lock()
	defer pqm.mu.Unlock()
	return pqm.pq.Len()
}

// Close đóng queue và wake up tất cả waiters
func (pqm *PriorityQueueManager[T]) Close() {
	pqm.mu.Lock()
	defer pqm.mu.Unlock()
	pqm.cond.Broadcast()
}

// WithPriority tạo PriorityTask từ Task thông thường
func WithPriority[T any](task Task[T], priority int) *PriorityTask[T] {
	return &PriorityTask[T]{
		Task:     task,
		Priority: priority,
		Index:    -1,
	}
}
