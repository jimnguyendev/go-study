# Go Concurrency - Lập trình đồng thời trong Go

## Mục lục

1. [Giới thiệu về Concurrency](#giới-thiệu-về-concurrency)
2. [Goroutines](#goroutines)
   - [Khái niệm cơ bản](#khái-niệm-cơ-bản)
   - [Tạo và sử dụng Goroutines](#tạo-và-sử-dụng-goroutines)
   - [So sánh với PHP](#so-sánh-với-php)
3. [Channels](#channels)
   - [Unbuffered Channels](#unbuffered-channels)
   - [Buffered Channels](#buffered-channels)
   - [Channel Direction](#channel-direction)
   - [Closing Channels](#closing-channels)
4. [Select Statement](#select-statement)
   - [Cú pháp cơ bản](#cú-pháp-cơ-bản)
   - [Non-blocking Operations](#non-blocking-operations)
   - [Timeout Pattern](#timeout-pattern)
5. [Sync Package](#sync-package)
   - [WaitGroup](#waitgroup)
   - [Mutex](#mutex)
   - [RWMutex](#rwmutex)
   - [Cond](#cond-condition-variable)
6. [Context Package](#context-package)
7. [Patterns và Best Practices](#patterns-và-best-practices)
8. [Ví dụ thực tế](#ví-dụ-thực-tế)

---

## Giới thiệu về Concurrency

Concurrency là khả năng thực hiện nhiều tác vụ cùng một lúc. Go được thiết kế với concurrency là một trong những tính năng cốt lõi, giúp xây dựng các ứng dụng hiệu suất cao và có khả năng mở rộng.

### Tại sao Concurrency quan trọng?

- **Performance**: Tận dụng tối đa tài nguyên hệ thống
- **Responsiveness**: Ứng dụng phản hồi nhanh hơn
- **Scalability**: Xử lý được nhiều request đồng thời
- **Resource Utilization**: Sử dụng hiệu quả CPU và I/O

### Concurrency vs Parallelism

- **Concurrency**: Xử lý nhiều việc cùng lúc (dealing with lots of things at once)
- **Parallelism**: Thực hiện nhiều việc cùng lúc (doing lots of things at once)

---

## Goroutines

### Khái niệm cơ bản

Goroutines là lightweight threads được quản lý bởi Go runtime. Chúng có chi phí rất thấp (chỉ khoảng 2KB stack ban đầu) và có thể tạo hàng nghìn goroutines mà không ảnh hưởng đến hiệu suất.

### Tạo và sử dụng Goroutines

#### Ví dụ cơ bản:

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("Hello, %s! (%d)\n", name, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Chạy đồng bộ (sequential)
    sayHello("Alice")
    
    // Chạy bất đồng bộ với goroutine
    go sayHello("Bob")
    go sayHello("Charlie")
    
    // Đợi goroutines hoàn thành
    time.Sleep(1 * time.Second)
    fmt.Println("Done!")
}
```

#### Anonymous Goroutines:

```go
func main() {
    // Goroutine với anonymous function
    go func() {
        fmt.Println("Anonymous goroutine")
    }()
    
    // Goroutine với parameters
    go func(msg string) {
        fmt.Println("Message:", msg)
    }("Hello from goroutine")
    
    time.Sleep(100 * time.Millisecond)
}
```

### So sánh với PHP

**PHP (Traditional):**
```php
// PHP không có built-in concurrency
// Phải sử dụng external libraries như ReactPHP, Swoole

// ReactPHP example
use React\EventLoop\Factory;

$loop = Factory::create();

$loop->addTimer(1, function() {
    echo "Timer 1\n";
});

$loop->addTimer(2, function() {
    echo "Timer 2\n";
});

$loop->run();
```

**Go:**
```go
// Go có built-in concurrency
func main() {
    go func() {
        time.Sleep(1 * time.Second)
        fmt.Println("Timer 1")
    }()
    
    go func() {
        time.Sleep(2 * time.Second)
        fmt.Println("Timer 2")
    }()
    
    time.Sleep(3 * time.Second)
}
```

---

## Channels

Channels là cách để goroutines giao tiếp với nhau. Chúng cho phép gửi và nhận dữ liệu một cách an toàn giữa các goroutines.

### Unbuffered Channels

Unbuffered channels là synchronous - việc gửi sẽ block cho đến khi có goroutine khác nhận.

```go
func main() {
    ch := make(chan string)
    
    go func() {
        ch <- "Hello from goroutine"
    }()
    
    msg := <-ch
    fmt.Println(msg)
}
```

#### Worker Pattern:

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
        time.Sleep(time.Second)
        results <- job * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)
    
    // Start 3 workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }
    
    // Send 5 jobs
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)
    
    // Collect results
    for r := 1; r <= 5; r++ {
        result := <-results
        fmt.Printf("Result: %d\n", result)
    }
}
```

### Buffered Channels

Buffered channels có thể lưu trữ một số lượng giá trị nhất định mà không cần receiver ngay lập tức.

```go
func main() {
    // Buffered channel với capacity = 3
    ch := make(chan int, 3)
    
    // Có thể gửi 3 giá trị mà không bị block
    ch <- 1
    ch <- 2
    ch <- 3
    
    // Nhận giá trị
    fmt.Println(<-ch) // 1
    fmt.Println(<-ch) // 2
    fmt.Println(<-ch) // 3
}
```

### Channel Direction

Có thể hạn chế channel chỉ để gửi hoặc chỉ để nhận:

```go
// Send-only channel
func sender(ch chan<- string) {
    ch <- "Hello"
}

// Receive-only channel
func receiver(ch <-chan string) {
    msg := <-ch
    fmt.Println(msg)
}

func main() {
    ch := make(chan string)
    
    go sender(ch)
    go receiver(ch)
    
    time.Sleep(100 * time.Millisecond)
}
```

### Closing Channels

```go
func main() {
    ch := make(chan int, 3)
    
    // Gửi dữ liệu
    ch <- 1
    ch <- 2
    ch <- 3
    
    // Đóng channel
    close(ch)
    
    // Nhận dữ liệu cho đến khi channel đóng
    for val := range ch {
        fmt.Println(val)
    }
    
    // Kiểm tra channel đã đóng chưa
    val, ok := <-ch
    if !ok {
        fmt.Println("Channel is closed")
    }
}
```

---

## Select Statement

Select statement cho phép goroutine đợi trên nhiều channel operations.

### Cú pháp cơ bản

```go
func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "Message from ch1"
    }()
    
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "Message from ch2"
    }()
    
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Received:", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received:", msg2)
        }
    }
}
```

### Non-blocking Operations

```go
func main() {
    ch := make(chan string)
    
    select {
    case msg := <-ch:
        fmt.Println("Received:", msg)
    default:
        fmt.Println("No message received")
    }
}
```

### Timeout Pattern

```go
func main() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "Hello"
    }()
    
    select {
    case msg := <-ch:
        fmt.Println("Received:", msg)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }
}
```

---

## Sync Package

Sync package cung cấp các primitives để đồng bộ hóa trong môi trường concurrent. Trong production với high traffic, việc sử dụng đúng sync primitives là critical cho performance và data consistency.

### WaitGroup

WaitGroup được sử dụng để đợi một nhóm goroutines hoàn thành. Trong production, đây là tool quan trọng để coordinate batch processing và graceful shutdown.

#### Cách sử dụng cơ bản:

```go
import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup
    
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }
    
    wg.Wait()
    fmt.Println("All workers done")
}
```

#### Production Pattern với Error Handling:

```go
type WorkerPool struct {
    wg     sync.WaitGroup
    errors chan error
    ctx    context.Context
    cancel context.CancelFunc
}

func NewWorkerPool(ctx context.Context) *WorkerPool {
    ctx, cancel := context.WithCancel(ctx)
    return &WorkerPool{
        errors: make(chan error, 100), // Buffered để tránh blocking
        ctx:    ctx,
        cancel: cancel,
    }
}

func (wp *WorkerPool) ProcessBatch(tasks []Task) error {
    // Limit concurrent workers để tránh resource exhaustion
    semaphore := make(chan struct{}, 10)
    
    for _, task := range tasks {
        wp.wg.Add(1)
        
        go func(t Task) {
            defer wp.wg.Done()
            
            // Acquire semaphore
            // Semaphore là một cơ chế đồng bộ hóa (synchronization mechanism) được sử dụng để kiểm soát số lượng goroutine có thể truy cập vào một tài nguyên hoặc thực hiện một tác vụ cùng lúc.
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            select {
            case <-wp.ctx.Done():
                return // Early termination
            default:
                if err := wp.processTask(t); err != nil {
                    select {
                    case wp.errors <- err:
                    default: // Error channel full, log và continue
                        log.Printf("Error channel full, dropping error: %v", err)
                    }
                }
            }
        }(task)
    }
    
    // Wait với timeout
    done := make(chan struct{})
    go func() {
        wp.wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        return nil
    case <-time.After(30 * time.Second): // Production timeout
        wp.cancel() // Cancel remaining workers
        return fmt.Errorf("batch processing timeout")
    }
}
```

### Mutex

Mutex được sử dụng để bảo vệ shared data. Trong high traffic environment, mutex contention có thể trở thành bottleneck.

#### Cách sử dụng cơ bản:

```go
import (
    "fmt"
    "sync"
)

type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

#### Production Optimizations:

**1. Sharded Counters để giảm contention:**

```go
type ShardedCounter struct {
    shards []struct {
        mu    sync.Mutex
        value int64
        _     [56]byte // Cache line padding để tránh false sharing
    }
    numShards int
}

func NewShardedCounter(numShards int) *ShardedCounter {
    if numShards <= 0 {
        numShards = runtime.NumCPU()
    }
    
    sc := &ShardedCounter{
        shards:    make([]struct{mu sync.Mutex; value int64; _ [56]byte}, numShards),
        numShards: numShards,
    }
    return sc
}

func (sc *ShardedCounter) Increment() {
    // Hash goroutine ID để distribute load
    shard := int(getGoroutineID()) % sc.numShards
    
    sc.shards[shard].mu.Lock()
    sc.shards[shard].value++
    sc.shards[shard].mu.Unlock()
}

func (sc *ShardedCounter) Value() int64 {
    var total int64
    for i := 0; i < sc.numShards; i++ {
        sc.shards[i].mu.Lock()
        total += sc.shards[i].value
        sc.shards[i].mu.Unlock()
    }
    return total
}
```

**2. Lock-free alternatives với atomic:**

```go
type AtomicCounter struct {
    value int64
}

func (ac *AtomicCounter) Increment() {
    atomic.AddInt64(&ac.value, 1)
}

func (ac *AtomicCounter) Value() int64 {
    return atomic.LoadInt64(&ac.value)
}

// Benchmark results (typical):
// BenchmarkMutexCounter-8      	10000000	       150 ns/op
// BenchmarkAtomicCounter-8     	50000000	        30 ns/op
// BenchmarkShardedCounter-8    	30000000	        45 ns/op
```

### RWMutex

RWMutex cho phép multiple readers hoặc single writer. Ideal cho read-heavy workloads trong production.

#### Cách sử dụng cơ bản:

```go
type SafeMap struct {
    mu sync.RWMutex
    data map[string]int
}

func (sm *SafeMap) Read(key string) (int, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    val, ok := sm.data[key]
    return val, ok
}

func (sm *SafeMap) Write(key string, value int) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.data[key] = value
}
```

#### Production Cache Implementation:

```go
type ProductionCache struct {
    mu       sync.RWMutex
    data     map[string]*CacheItem
    maxSize  int
    ttl      time.Duration
    stats    CacheStats
}

type CacheItem struct {
    Value     interface{}
    ExpiresAt time.Time
    AccessCount int64
}

type CacheStats struct {
    Hits   int64
    Misses int64
    Evictions int64
}

func (pc *ProductionCache) Get(key string) (interface{}, bool) {
    pc.mu.RLock()
    
    item, exists := pc.data[key]
    if !exists {
        pc.mu.RUnlock()
        atomic.AddInt64(&pc.stats.Misses, 1)
        return nil, false
    }
    
    // Check expiration
    if time.Now().After(item.ExpiresAt) {
        pc.mu.RUnlock()
        
        // Upgrade to write lock để remove expired item
        pc.mu.Lock()
        delete(pc.data, key)
        pc.mu.Unlock()
        
        atomic.AddInt64(&pc.stats.Misses, 1)
        return nil, false
    }
    
    value := item.Value
    atomic.AddInt64(&item.AccessCount, 1)
    pc.mu.RUnlock()
    
    atomic.AddInt64(&pc.stats.Hits, 1)
    return value, true
}

func (pc *ProductionCache) Set(key string, value interface{}) {
    pc.mu.Lock()
    defer pc.mu.Unlock()
    
    // Eviction policy nếu cache full
    if len(pc.data) >= pc.maxSize {
        pc.evictLRU()
    }
    
    pc.data[key] = &CacheItem{
        Value:     value,
        ExpiresAt: time.Now().Add(pc.ttl),
        AccessCount: 0,
    }
}

func (pc *ProductionCache) evictLRU() {
    var oldestKey string
    var oldestAccess int64 = math.MaxInt64
    
    for key, item := range pc.data {
        if item.AccessCount < oldestAccess {
            oldestAccess = item.AccessCount
            oldestKey = key
        }
    }
    
    if oldestKey != "" {
        delete(pc.data, oldestKey)
        atomic.AddInt64(&pc.stats.Evictions, 1)
    }
}
```

### Production Best Practices cho High Traffic

#### 1. Lock Ordering để tránh Deadlocks:

```go
type BankAccount struct {
    id      int
    balance int64
    mu      sync.Mutex
}

// ĐÚNG: Always lock theo thứ tự ID
func Transfer(from, to *BankAccount, amount int64) error {
    if from.id < to.id {
        from.mu.Lock()
        defer from.mu.Unlock()
        to.mu.Lock()
        defer to.mu.Unlock()
    } else {
        to.mu.Lock()
        defer to.mu.Unlock()
        from.mu.Lock()
        defer from.mu.Unlock()
    }
    
    if from.balance < amount {
        return errors.New("insufficient funds")
    }
    
    from.balance -= amount
    to.balance += amount
    return nil
}
```

#### 2. Timeout và Context cho Production:

```go
func (pc *ProductionCache) GetWithTimeout(ctx context.Context, key string) (interface{}, error) {
    done := make(chan struct{})
    var result interface{}
    var found bool
    
    go func() {
        result, found = pc.Get(key)
        close(done)
    }()
    
    select {
    case <-done:
        if !found {
            return nil, ErrNotFound
        }
        return result, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}
```

#### 3. Monitoring và Metrics:

```go
type SyncMetrics struct {
    MutexContentions  int64
    RWMutexReadLocks  int64
    RWMutexWriteLocks int64
    WaitGroupWaits    int64
}

func (sm *SafeMap) ReadWithMetrics(key string) (int, bool) {
    start := time.Now()
    
    sm.mu.RLock()
    atomic.AddInt64(&globalMetrics.RWMutexReadLocks, 1)
    
    val, ok := sm.data[key]
    
    sm.mu.RUnlock()
    
    // Log slow operations
    if duration := time.Since(start); duration > 10*time.Millisecond {
        log.Printf("Slow read operation: %v for key %s", duration, key)
    }
    
    return val, ok
}
```

#### 4. Graceful Shutdown Pattern:

```go
type Server struct {
    mu       sync.RWMutex
    shutdown bool
    wg       sync.WaitGroup
}

func (s *Server) HandleRequest(req Request) {
    s.mu.RLock()
    if s.shutdown {
        s.mu.RUnlock()
        return // Reject new requests during shutdown
    }
    s.wg.Add(1)
    s.mu.RUnlock()
    
    defer s.wg.Done()
    
    // Process request
    s.processRequest(req)
}

func (s *Server) Shutdown(ctx context.Context) error {
    s.mu.Lock()
    s.shutdown = true
    s.mu.Unlock()
    
    done := make(chan struct{})
    go func() {
        s.wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### Performance Considerations cho High Traffic

1. **Lock Granularity**: Fine-grained locks thay vì coarse-grained
2. **Lock Duration**: Minimize time holding locks
3. **Reader/Writer Ratio**: Sử dụng RWMutex khi read >> write
4. **Atomic Operations**: Prefer atomic cho simple operations
5. **Lock-free Data Structures**: Consider cho high-contention scenarios
6. **Sharding**: Distribute load across multiple locks
7. **Monitoring**: Track lock contention và performance metrics

### Cond (Condition Variable)

Cond là một condition variable cho phép goroutines chờ đợi hoặc thông báo về sự thay đổi của một điều kiện cụ thể. Nó thường được sử dụng khi bạn cần goroutines chờ đợi cho đến khi một điều kiện nào đó được thỏa mãn.

#### Khái niệm cơ bản:

```go
type Cond struct {
    L Locker // Thường là *Mutex hoặc *RWMutex
    // ... các field internal
}

// Các method chính:
// Wait()      - Chờ đợi signal
// Signal()    - Đánh thức một goroutine đang chờ
// Broadcast() - Đánh thức tất cả goroutines đang chờ
```

#### Cách sử dụng cơ bản:

```go
import (
    "fmt"
    "sync"
    "time"
)

type Queue struct {
    mu    sync.Mutex
    cond  *sync.Cond
    items []string
}

func NewQueue() *Queue {
    q := &Queue{
        items: make([]string, 0),
    }
    q.cond = sync.NewCond(&q.mu)
    return q
}

func (q *Queue) Put(item string) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    q.items = append(q.items, item)
    fmt.Printf("Put item: %s\n", item)
    
    // Thông báo cho một goroutine đang chờ
    q.cond.Signal()
}

func (q *Queue) Get() string {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    // Chờ đợi cho đến khi có item
    for len(q.items) == 0 {
        q.cond.Wait() // Giải phóng lock và chờ đợi
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    fmt.Printf("Got item: %s\n", item)
    return item
}
```

#### Ví dụ Producer-Consumer Pattern:

```go
func main() {
    queue := NewQueue()
    
    // Consumer goroutines
    for i := 1; i <= 3; i++ {
        go func(id int) {
            for j := 0; j < 2; j++ {
                item := queue.Get()
                fmt.Printf("Consumer %d got: %s\n", id, item)
                time.Sleep(100 * time.Millisecond)
            }
        }(i)
    }
    
    // Producer
    time.Sleep(500 * time.Millisecond) // Đảm bảo consumers đã sẵn sàng
    
    for i := 1; i <= 6; i++ {
        queue.Put(fmt.Sprintf("item-%d", i))
        time.Sleep(200 * time.Millisecond)
    }
    
    time.Sleep(2 * time.Second)
}
```

#### Production Example - Worker Pool với Condition:

```go
type WorkerPool struct {
    mu          sync.Mutex
    cond        *sync.Cond
    tasks       []Task
    workers     int
    activeWorkers int
    shutdown    bool
}

type Task struct {
    ID   int
    Data string
}

func NewWorkerPool(numWorkers int) *WorkerPool {
    wp := &WorkerPool{
        tasks:   make([]Task, 0),
        workers: numWorkers,
    }
    wp.cond = sync.NewCond(&wp.mu)
    
    // Start workers
    for i := 0; i < numWorkers; i++ {
        go wp.worker(i)
    }
    
    return wp
}

func (wp *WorkerPool) worker(id int) {
    for {
        wp.mu.Lock()
        
        // Chờ đợi task hoặc shutdown signal
        for len(wp.tasks) == 0 && !wp.shutdown {
            wp.cond.Wait()
        }
        
        if wp.shutdown {
            wp.mu.Unlock()
            fmt.Printf("Worker %d shutting down\n", id)
            return
        }
        
        // Lấy task
        task := wp.tasks[0]
        wp.tasks = wp.tasks[1:]
        wp.activeWorkers++
        
        wp.mu.Unlock()
        
        // Xử lý task
        fmt.Printf("Worker %d processing task %d\n", id, task.ID)
        time.Sleep(time.Second) // Simulate work
        
        wp.mu.Lock()
        wp.activeWorkers--
        wp.cond.Broadcast() // Thông báo task hoàn thành
        wp.mu.Unlock()
    }
}

func (wp *WorkerPool) AddTask(task Task) {
    wp.mu.Lock()
    defer wp.mu.Unlock()
    
    if wp.shutdown {
        return
    }
    
    wp.tasks = append(wp.tasks, task)
    wp.cond.Signal() // Đánh thức một worker
}

func (wp *WorkerPool) Shutdown() {
    wp.mu.Lock()
    defer wp.mu.Unlock()
    
    wp.shutdown = true
    wp.cond.Broadcast() // Đánh thức tất cả workers
}

func (wp *WorkerPool) WaitForCompletion() {
    wp.mu.Lock()
    defer wp.mu.Unlock()
    
    // Chờ đợi cho đến khi tất cả tasks hoàn thành
    for len(wp.tasks) > 0 || wp.activeWorkers > 0 {
        wp.cond.Wait()
    }
}
```

#### Khi nào sử dụng Cond:

1. **Producer-Consumer Pattern**: Khi cần đồng bộ hóa giữa producers và consumers
2. **Worker Pool**: Quản lý pool of workers chờ đợi tasks
3. **State Changes**: Khi goroutines cần chờ đợi state thay đổi
4. **Resource Availability**: Chờ đợi tài nguyên có sẵn

#### Best Practices:

1. **Always use with Mutex**: Cond phải được sử dụng với Mutex
2. **Use in loop**: Luôn check condition trong vòng lặp với Wait()
3. **Signal vs Broadcast**: 
   - `Signal()`: Đánh thức một goroutine
   - `Broadcast()`: Đánh thức tất cả goroutines
4. **Avoid Spurious Wakeups**: Luôn check condition sau khi Wake up

#### So sánh với Channels:

| Cond | Channels |
|------|----------|
| Shared state với explicit locking | Message passing |
| Multiple waiters cho cùng condition | Point-to-point communication |
| Broadcast capability | Fan-out cần explicit design |
| Lower-level primitive | Higher-level abstraction |

**Khi nào dùng Cond thay vì Channels:**
- Khi cần nhiều goroutines chờ đợi cùng một condition
- Khi cần broadcast signal đến tất cả waiters
- Khi working với shared state phức tạp
- Khi cần fine-grained control over synchronization

---

## Context Package

Context được sử dụng để truyền deadlines, cancellation signals và request-scoped values:

```go
import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d cancelled\n", id)
            return
        default:
            fmt.Printf("Worker %d working...\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    go worker(ctx, 1)
    go worker(ctx, 2)
    
    time.Sleep(3 * time.Second)
    fmt.Println("Main done")
}
```

---

## Patterns và Best Practices

### 1. Fan-out, Fan-in Pattern

```go
func fanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    
    for i := 0; i < workers; i++ {
        output := make(chan int)
        outputs[i] = output
        
        go func() {
            defer close(output)
            for n := range input {
                output <- n * n // Square the number
            }
        }()
    }
    
    return outputs
}

func fanIn(inputs ...<-chan int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup
    
    for _, input := range inputs {
        wg.Add(1)
        go func(ch <-chan int) {
            defer wg.Done()
            for val := range ch {
                output <- val
            }
        }(input)
    }
    
    go func() {
        wg.Wait()
        close(output)
    }()
    
    return output
}
```

### 2. Pipeline Pattern

```go
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

func main() {
    // Pipeline: generator -> square
    numbers := generator(1, 2, 3, 4, 5)
    squares := square(numbers)
    
    for result := range squares {
        fmt.Println(result)
    }
}
```

### 3. Worker Pool Pattern

```go
type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID int
    Value string
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job.ID)
        time.Sleep(time.Second) // Simulate work
        
        results <- Result{
            JobID: job.ID,
            Value: fmt.Sprintf("Processed: %s", job.Data),
        }
    }
}

func main() {
    jobs := make(chan Job, 100)
    results := make(chan Result, 100)
    
    // Start workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }
    
    // Send jobs
    for j := 1; j <= 5; j++ {
        jobs <- Job{ID: j, Data: fmt.Sprintf("job-%d", j)}
    }
    close(jobs)
    
    // Collect results
    for r := 1; r <= 5; r++ {
        result := <-results
        fmt.Printf("Result: %+v\n", result)
    }
}
```

---

## Ví dụ thực tế

### Web Scraper với Concurrency

```go
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "sync"
    "time"
)

type Result struct {
    URL    string
    Status int
    Size   int
    Error  error
}

func fetchURL(url string, results chan<- Result) {
    start := time.Now()
    
    resp, err := http.Get(url)
    if err != nil {
        results <- Result{URL: url, Error: err}
        return
    }
    defer resp.Body.Close()
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        results <- Result{URL: url, Status: resp.StatusCode, Error: err}
        return
    }
    
    results <- Result{
        URL:    url,
        Status: resp.StatusCode,
        Size:   len(body),
    }
    
    fmt.Printf("Fetched %s in %v\n", url, time.Since(start))
}

func main() {
    urls := []string{
        "https://golang.org",
        "https://github.com",
        "https://stackoverflow.com",
        "https://reddit.com",
        "https://news.ycombinator.com",
    }
    
    results := make(chan Result, len(urls))
    
    // Start goroutines
    for _, url := range urls {
        go fetchURL(url, results)
    }
    
    // Collect results
    for i := 0; i < len(urls); i++ {
        result := <-results
        if result.Error != nil {
            fmt.Printf("Error fetching %s: %v\n", result.URL, result.Error)
        } else {
            fmt.Printf("Success: %s - Status: %d, Size: %d bytes\n", 
                result.URL, result.Status, result.Size)
        }
    }
}
```

### Rate Limiter

```go
package main

import (
    "fmt"
    "time"
)

type RateLimiter struct {
    tokens chan struct{}
    ticker *time.Ticker
}

func NewRateLimiter(rate int) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, rate),
        ticker: time.NewTicker(time.Second / time.Duration(rate)),
    }
    
    // Fill initial tokens
    for i := 0; i < rate; i++ {
        rl.tokens <- struct{}{}
    }
    
    // Refill tokens
    go func() {
        for range rl.ticker.C {
            select {
            case rl.tokens <- struct{}{}:
            default:
                // Token bucket is full
            }
        }
    }()
    
    return rl
}

func (rl *RateLimiter) Wait() {
    <-rl.tokens
}

func (rl *RateLimiter) Stop() {
    rl.ticker.Stop()
}

func main() {
    limiter := NewRateLimiter(2) // 2 requests per second
    defer limiter.Stop()
    
    for i := 1; i <= 10; i++ {
        limiter.Wait()
        go func(id int) {
            fmt.Printf("Request %d at %v\n", id, time.Now().Format("15:04:05"))
        }(i)
    }
    
    time.Sleep(6 * time.Second)
}
```

---

## Kết luận

Concurrency trong Go cung cấp các công cụ mạnh mẽ để xây dựng ứng dụng hiệu suất cao:

- **Goroutines**: Lightweight threads dễ sử dụng
- **Channels**: Giao tiếp an toàn giữa goroutines
- **Select**: Xử lý multiple channel operations
- **Sync package**: Đồng bộ hóa và bảo vệ shared data
- **Context**: Quản lý lifecycle và cancellation

Việc thành thạo các khái niệm này sẽ giúp bạn xây dựng các ứng dụng Go hiệu quả và có khả năng mở rộng cao.