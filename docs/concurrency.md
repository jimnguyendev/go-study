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

### WaitGroup

WaitGroup được sử dụng để đợi một nhóm goroutines hoàn thành:

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

### Mutex

Mutex được sử dụng để bảo vệ shared data:

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

func main() {
    counter := &Counter{}
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }
    
    wg.Wait()
    fmt.Println("Final value:", counter.Value())
}
```

### RWMutex

RWMutex cho phép multiple readers hoặc single writer:

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