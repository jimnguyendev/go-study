# Deadlock, Starvation và Livelock trong Go

## Mục lục

1. [Giới thiệu](#giới-thiệu)
2. [Deadlock](#deadlock)
   - [Khái niệm](#khái-niệm-deadlock)
   - [Điều kiện cần thiết](#điều-kiện-cần-thiết)
   - [Ví dụ Deadlock trong Go](#ví-dụ-deadlock-trong-go)
   - [Cách phòng tránh](#cách-phòng-tránh-deadlock)
   - [Tại sao thứ tự lock quan trọng](#tại-sao-thứ-tự-lock-quan-trọng)
3. [Starvation](#starvation)
   - [Khái niệm](#khái-niệm-starvation)
   - [Nguyên nhân](#nguyên-nhân-starvation)
   - [Ví dụ Starvation trong Go](#ví-dụ-starvation-trong-go)
   - [Cách giải quyết](#cách-giải-quyết-starvation)
4. [Livelock](#livelock)
   - [Khái niệm](#khái-niệm-livelock)
   - [Nguyên nhân](#nguyên-nhân-livelock)
   - [Ví dụ Livelock trong Go](#ví-dụ-livelock-trong-go)
   - [Cách phòng tránh](#cách-phòng-tránh-livelock)
5. [So sánh ba khái niệm](#so-sánh-ba-khái-niệm)
6. [Best Practices](#best-practices)
7. [Overhead trong Concurrency](#overhead-trong-concurrency)
   - [Khái niệm Overhead](#khái-niệm-overhead)
   - [Các loại Overhead](#các-loại-overhead)
   - [Ví dụ Overhead trong Go](#ví-dụ-overhead-trong-go)
   - [Cách đo lường và tối ưu](#cách-đo-lường-và-tối-ưu-overhead)

---

## Giới thiệu

Deadlock, Starvation và Livelock là ba vấn đề quan trọng trong lập trình đồng thời (concurrent programming) mà developers cần hiểu và biết cách xử lý. Các vấn đề này có thể xảy ra khi nhiều processes hoặc goroutines cạnh tranh tài nguyên.

---

## Deadlock

### Khái niệm Deadlock

Deadlock là tình huống mà một tập hợp các processes bị block vô thời hạn vì mỗi process đang giữ một tài nguyên và chờ đợi tài nguyên khác được giữ bởi process khác trong tập hợp đó.

### Điều kiện cần thiết

Để deadlock xảy ra, cần có đủ 4 điều kiện sau:

1. **Mutual Exclusion**: Tài nguyên chỉ có thể được sử dụng bởi một process tại một thời điểm
2. **Hold and Wait**: Process giữ ít nhất một tài nguyên và chờ đợi tài nguyên khác
3. **No Preemption**: Tài nguyên không thể bị tước đoạt từ process
4. **Circular Wait**: Có một chuỗi vòng tròn các processes chờ đợi lẫn nhau

### Ví dụ Deadlock trong Go

#### Ví dụ 1: Classic Deadlock với 2 Mutexes

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// Ví dụ deadlock cơ bản với 2 mutexes
func deadlockExample() {
	var mutex1, mutex2 sync.Mutex

	// Goroutine 1: Lock mutex1 trước, sau đó lock mutex2
	go func() {
		fmt.Println("Goroutine 1: Locking mutex1")
		mutex1.Lock()
		fmt.Println("Goroutine 1: Locked mutex1")

		time.Sleep(100 * time.Millisecond) // Simulate some work

		fmt.Println("Goroutine 1: Trying to lock mutex2")
		mutex2.Lock() // Sẽ bị block ở đây
		fmt.Println("Goroutine 1: Locked mutex2")

		mutex2.Unlock()
		mutex1.Unlock()
		fmt.Println("Goroutine 1: Done")
	}()

	// Goroutine 2: Lock mutex2 trước, sau đó lock mutex1
	go func() {
		fmt.Println("Goroutine 2: Locking mutex2")
		mutex2.Lock()
		fmt.Println("Goroutine 2: Locked mutex2")

		time.Sleep(100 * time.Millisecond) // Simulate some work

		fmt.Println("Goroutine 2: Trying to lock mutex1")
		mutex1.Lock() // Sẽ bị block ở đây
		fmt.Println("Goroutine 2: Locked mutex1")

		mutex1.Unlock()
		mutex2.Unlock()
		fmt.Println("Goroutine 2: Done")
	}()

	// Đợi để thấy deadlock xảy ra
	time.Sleep(2 * time.Second)
	fmt.Println("Main: Deadlock detected!")
}

func main() {
	fmt.Println("=== Deadlock Example ===")
	deadlockExample()
}
```

#### Ví dụ 2: Deadlock với Channels

```go
package main

import (
	"fmt"
	"time"
)

// Deadlock với unbuffered channels
func channelDeadlock() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		fmt.Println("Goroutine 1: Sending to ch1")
		ch1 <- 1 // Sẽ block vì không có receiver
		fmt.Println("Goroutine 1: Waiting for ch2")
		<-ch2
		fmt.Println("Goroutine 1: Done")
	}()

	go func() {
		fmt.Println("Goroutine 2: Sending to ch2")
		ch2 <- 2 // Sẽ block vì không có receiver
		fmt.Println("Goroutine 2: Waiting for ch1")
		<-ch1
		fmt.Println("Goroutine 2: Done")
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("Channel deadlock detected!")
}
```

### Cách phòng tránh Deadlock

#### 1. Lock Ordering

```go
func avoidDeadlockWithOrdering() {
	var mutex1, mutex2 sync.Mutex

	// Luôn lock theo thứ tự: mutex1 trước, mutex2 sau
	go func() {
		fmt.Println("Goroutine 1: Locking mutex1")
		mutex1.Lock()
		fmt.Println("Goroutine 1: Locking mutex2")
		mutex2.Lock()

		// Do work
		time.Sleep(100 * time.Millisecond)

		mutex2.Unlock()
		mutex1.Unlock()
		fmt.Println("Goroutine 1: Done")
	}()

	go func() {
		fmt.Println("Goroutine 2: Locking mutex1")
		mutex1.Lock()
		fmt.Println("Goroutine 2: Locking mutex2")
		mutex2.Lock()

		// Do work
		time.Sleep(100 * time.Millisecond)

		mutex2.Unlock()
		mutex1.Unlock()
		fmt.Println("Goroutine 2: Done")
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("No deadlock - both goroutines completed")
}
```

#### 2. Timeout với Context

```go
import "context"

func avoidDeadlockWithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	ch := make(chan int)

	go func() {
		select {
		case ch <- 42:
			fmt.Println("Data sent successfully")
		case <-ctx.Done():
			fmt.Println("Timeout: Could not send data")
		}
	}()

	select {
	case data := <-ch:
		fmt.Printf("Received: %d\n", data)
	case <-ctx.Done():
		fmt.Println("Timeout: Could not receive data")
	}
}
```

### Tại sao thứ tự lock quan trọng

Thứ tự lock (lock ordering) là một trong những kỹ thuật quan trọng nhất để phòng tránh deadlock. Hiểu rõ tại sao thứ tự lock quan trọng sẽ giúp bạn thiết kế hệ thống concurrent an toàn và hiệu quả.

#### Nguyên lý Circular Wait

Deadlock xảy ra khi có **circular wait** - một chuỗi vòng tròn các processes chờ đợi lẫn nhau. Lock ordering phá vỡ điều kiện này bằng cách đảm bảo tất cả goroutines luôn acquire locks theo cùng một thứ tự.

#### Ví dụ minh họa

**Trường hợp gây Deadlock:**
```go
// Goroutine 1: Lock A -> Lock B
func goroutine1() {
    mutexA.Lock()
    mutexB.Lock()
    // work...
    mutexB.Unlock()
    mutexA.Unlock()
}

// Goroutine 2: Lock B -> Lock A (thứ tự ngược lại!)
func goroutine2() {
    mutexB.Lock()  // Có thể gây deadlock
    mutexA.Lock()  // Có thể gây deadlock
    // work...
    mutexA.Unlock()
    mutexB.Unlock()
}
```

**Trường hợp an toàn:**
```go
// Cả hai goroutines đều lock theo thứ tự: A -> B
func goroutine1() {
    mutexA.Lock()
    mutexB.Lock()
    // work...
    mutexB.Unlock()
    mutexA.Unlock()
}

func goroutine2() {
    mutexA.Lock()  // Cùng thứ tự
    mutexB.Lock()  // Cùng thứ tự
    // work...
    mutexB.Unlock()
    mutexA.Unlock()
}
```

#### Ví dụ thực tế: Bank Transfer

```go
type Account struct {
    ID      int
    Balance int
    mu      sync.Mutex
}

// Unsafe: Có thể gây deadlock
func transferUnsafe(from, to *Account, amount int) {
    from.mu.Lock()
    to.mu.Lock()
    
    if from.Balance >= amount {
        from.Balance -= amount
        to.Balance += amount
        fmt.Printf("Transferred %d from account %d to %d\n", amount, from.ID, to.ID)
    }
    
    to.mu.Unlock()
    from.mu.Unlock()
}

// Safe: Sử dụng lock ordering
func transferSafe(from, to *Account, amount int) {
    // Luôn lock account có ID nhỏ hơn trước
    first, second := from, to
    if from.ID > to.ID {
        first, second = to, from
    }
    
    first.mu.Lock()
    second.mu.Lock()
    
    if from.Balance >= amount {
        from.Balance -= amount
        to.Balance += amount
        fmt.Printf("Transferred %d from account %d to %d\n", amount, from.ID, to.ID)
    }
    
    second.mu.Unlock()
    first.mu.Unlock()
}

// Ví dụ sử dụng
func bankTransferExample() {
    account1 := &Account{ID: 1, Balance: 1000}
    account2 := &Account{ID: 2, Balance: 500}
    
    var wg sync.WaitGroup
    
    // Transfer từ account1 sang account2
    wg.Add(1)
    go func() {
        defer wg.Done()
        transferSafe(account1, account2, 100)
    }()
    
    // Transfer từ account2 sang account1 (ngược lại)
    wg.Add(1)
    go func() {
        defer wg.Done()
        transferSafe(account2, account1, 50)
    }()
    
    wg.Wait()
    fmt.Printf("Final balances: Account1=%d, Account2=%d\n", account1.Balance, account2.Balance)
}
```

#### Các chiến lược Lock Ordering

**1. ID/Address-based Ordering:**
```go
// Sắp xếp theo ID hoặc địa chỉ memory
func lockByID(accounts []*Account) {
    sort.Slice(accounts, func(i, j int) bool {
        return accounts[i].ID < accounts[j].ID
    })
    
    for _, acc := range accounts {
        acc.mu.Lock()
    }
    
    // Do work...
    
    // Unlock theo thứ tự ngược lại
    for i := len(accounts) - 1; i >= 0; i-- {
        accounts[i].mu.Unlock()
    }
}
```

**2. Hierarchical Ordering:**
```go
// Định nghĩa hierarchy levels
const (
    DatabaseLevel = 1
    CacheLevel    = 2
    LogLevel      = 3
)

type Resource struct {
    Level int
    mu    sync.Mutex
}

func lockHierarchical(resources []*Resource) {
    // Sắp xếp theo level
    sort.Slice(resources, func(i, j int) bool {
        return resources[i].Level < resources[j].Level
    })
    
    for _, res := range resources {
        res.mu.Lock()
    }
}
```

**3. Global Lock Ordering:**
```go
var globalLockOrder = map[string]int{
    "user_data":    1,
    "session_data": 2,
    "cache_data":   3,
    "log_data":     4,
}

type NamedMutex struct {
    Name string
    mu   sync.Mutex
}

func lockGlobalOrder(mutexes []*NamedMutex) {
    sort.Slice(mutexes, func(i, j int) bool {
        return globalLockOrder[mutexes[i].Name] < globalLockOrder[mutexes[j].Name]
    })
    
    for _, m := range mutexes {
        m.mu.Lock()
    }
}
```

#### Những điều cần lưu ý

1. **Consistency**: Tất cả code phải tuân theo cùng một thứ tự lock
2. **Performance**: Lock ordering có thể ảnh hưởng đến performance do serialization
3. **Documentation**: Cần document rõ ràng lock ordering strategy
4. **Testing**: Sử dụng race detector và stress testing để phát hiện deadlock

```go
// Chạy với race detector
// go run -race main.go

// Stress testing
func stressTest() {
    for i := 0; i < 1000; i++ {
        go func() {
            // Concurrent operations
        }()
    }
}
```

---

## Starvation

### Khái niệm Starvation

Starvation xảy ra khi một process chờ đợi vô thời hạn vì các processes có độ ưu tiên cao hơn liên tục được thực thi, khiến process có độ ưu tiên thấp không bao giờ được truy cập tài nguyên.

### Nguyên nhân Starvation

1. **Priority Scheduling**: Processes có priority cao luôn được ưu tiên
2. **Resource Utilization**: Tài nguyên luôn được sử dụng bởi processes có priority cao
3. **Unfair Scheduling**: Thuật toán scheduling không công bằng

### Ví dụ Starvation trong Go

#### Ví dụ 1: Reader-Writer Starvation

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type RWResource struct {
	mu   sync.RWMutex
	data int
}

func (r *RWResource) Read(id int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fmt.Printf("Reader %d: Reading data = %d\n", id, r.data)
	time.Sleep(100 * time.Millisecond) // Simulate read time
}

func (r *RWResource) Write(id int, value int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Printf("Writer %d: Writing data = %d\n", id, value)
	r.data = value
	time.Sleep(200 * time.Millisecond) // Simulate write time
}

// Ví dụ về Writer Starvation
func writerStarvationExample() {
	resource := &RWResource{data: 0}

	// Tạo nhiều readers
	for i := 1; i <= 10; i++ {
		go func(id int) {
			for j := 0; j < 5; j++ {
				resource.Read(id)
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}

	// Một writer có thể bị starve
	go func() {
		fmt.Println("Writer: Waiting to write...")
		resource.Write(1, 100)
		fmt.Println("Writer: Finally wrote!")
	}()

	time.Sleep(3 * time.Second)
	fmt.Println("Writer starvation example completed")
}
```

#### Ví dụ 2: Priority-based Starvation

```go
type Task struct {
	ID       int
	Priority int
	Work     func()
}

type PriorityScheduler struct {
	tasks chan Task
	mu    sync.Mutex
}

func NewPriorityScheduler() *PriorityScheduler {
	return &PriorityScheduler{
		tasks: make(chan Task, 100),
	}
}

func (ps *PriorityScheduler) AddTask(task Task) {
	ps.tasks <- task
}

func (ps *PriorityScheduler) Start() {
	go func() {
		for task := range ps.tasks {
			// Simulate priority-based execution
			// High priority tasks (priority > 5) get executed immediately
			if task.Priority > 5 {
				fmt.Printf("Executing high priority task %d (priority: %d)\n", task.ID, task.Priority)
				task.Work()
			} else {
				// Low priority tasks might be delayed
				time.Sleep(100 * time.Millisecond)
				fmt.Printf("Executing low priority task %d (priority: %d)\n", task.ID, task.Priority)
				task.Work()
			}
		}
	}()
}

func priorityStarvationExample() {
	scheduler := NewPriorityScheduler()
	scheduler.Start()

	// Add low priority task
	scheduler.AddTask(Task{
		ID:       1,
		Priority: 2,
		Work: func() {
			fmt.Println("Low priority task completed")
		},
	})

	// Continuously add high priority tasks
	for i := 2; i <= 10; i++ {
		scheduler.AddTask(Task{
			ID:       i,
			Priority: 8,
			Work: func() {
				fmt.Println("High priority task completed")
			},
		})
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)
	fmt.Println("Priority starvation example completed")
}
```

### Cách giải quyết Starvation

#### 1. Fair Scheduling với Round Robin

```go
type FairScheduler struct {
	highPriorityTasks chan Task
	lowPriorityTasks  chan Task
	lastServed        string // "high" or "low"
}

func NewFairScheduler() *FairScheduler {
	return &FairScheduler{
		highPriorityTasks: make(chan Task, 50),
		lowPriorityTasks:  make(chan Task, 50),
		lastServed:        "low", // Start with high priority
	}
}

func (fs *FairScheduler) Start() {
	go func() {
		for {
			if fs.lastServed == "low" {
				// Try to serve high priority first
				select {
				case task := <-fs.highPriorityTasks:
					fmt.Printf("Fair: Executing high priority task %d\n", task.ID)
					task.Work()
					fs.lastServed = "high"
				case task := <-fs.lowPriorityTasks:
					fmt.Printf("Fair: Executing low priority task %d\n", task.ID)
					task.Work()
					fs.lastServed = "low"
				case <-time.After(100 * time.Millisecond):
					// No tasks available
					continue
				}
			} else {
				// Try to serve low priority first
				select {
				case task := <-fs.lowPriorityTasks:
					fmt.Printf("Fair: Executing low priority task %d\n", task.ID)
					task.Work()
					fs.lastServed = "low"
				case task := <-fs.highPriorityTasks:
					fmt.Printf("Fair: Executing high priority task %d\n", task.ID)
					task.Work()
					fs.lastServed = "high"
				case <-time.After(100 * time.Millisecond):
					// No tasks available
					continue
				}
			}
		}
	}()
}
```

---

## Livelock

### Khái niệm Livelock

Livelock xảy ra khi hai hoặc nhiều processes liên tục lặp lại cùng một tương tác để phản hồi với các thay đổi trong processes khác mà không thực hiện được công việc hữu ích nào. Các processes này không ở trạng thái chờ đợi và chúng đang chạy đồng thời.

### Nguyên nhân Livelock

1. **Continuous State Changes**: Processes liên tục thay đổi trạng thái để phản hồi lẫn nhau
2. **Resource Competition**: Cạnh tranh tài nguyên hữu hạn
3. **Retry Logic**: Logic retry không hiệu quả

### Ví dụ Livelock trong Go

#### Ví dụ 1: Classic Livelock - Polite People Problem

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Person struct {
	name string
	mu   sync.Mutex
}

func (p *Person) TryToPass(other *Person, hallway chan bool) {
	for i := 0; i < 10; i++ {
		// Try to enter hallway
		select {
		case hallway <- true:
			fmt.Printf("%s: Entered hallway\n", p.name)
			
			// Check if other person is also trying
			time.Sleep(50 * time.Millisecond)
			
			// Be polite - if other person is waiting, step back
			if len(hallway) > 0 || other.isWaiting() {
				fmt.Printf("%s: Being polite, stepping back\n", p.name)
				<-hallway // Step back
				time.Sleep(100 * time.Millisecond) // Wait a bit
				continue
			}
			
			fmt.Printf("%s: Successfully passed through hallway\n", p.name)
			<-hallway // Exit hallway
			return
			
		default:
			fmt.Printf("%s: Hallway occupied, waiting...\n", p.name)
			time.Sleep(100 * time.Millisecond)
		}
	}
	fmt.Printf("%s: Gave up after too many attempts\n", p.name)
}

func (p *Person) isWaiting() bool {
	// Simulate checking if person is waiting
	return true // Always return true for livelock demonstration
}

func livelockExample() {
	hallway := make(chan bool, 1) // Only one person can pass at a time
	
	person1 := &Person{name: "Alice"}
	person2 := &Person{name: "Bob"}
	
	var wg sync.WaitGroup
	wg.Add(2)
	
	go func() {
		defer wg.Done()
		person1.TryToPass(person2, hallway)
	}()
	
	go func() {
		defer wg.Done()
		person2.TryToPass(person1, hallway)
	}()
	
	wg.Wait()
	fmt.Println("Livelock example completed")
}
```

#### Ví dụ 2: Resource Retry Livelock

```go
type Resource struct {
	id       int
	inUse    bool
	mu       sync.Mutex
	attempts int
}

func (r *Resource) TryAcquire(workerID int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.attempts++
	if r.inUse {
		fmt.Printf("Worker %d: Resource %d is busy (attempt %d)\n", workerID, r.id, r.attempts)
		return false
	}
	
	r.inUse = true
	fmt.Printf("Worker %d: Acquired resource %d\n", workerID, r.id)
	return true
}

func (r *Resource) Release(workerID int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.inUse = false
	fmt.Printf("Worker %d: Released resource %d\n", workerID, r.id)
}

func resourceLivelockExample() {
	resource1 := &Resource{id: 1}
	resource2 := &Resource{id: 2}
	
	// Worker 1 needs both resources
	go func() {
		for i := 0; i < 5; i++ {
			if resource1.TryAcquire(1) {
				time.Sleep(100 * time.Millisecond)
				
				if resource2.TryAcquire(1) {
					fmt.Println("Worker 1: Got both resources, doing work...")
					time.Sleep(200 * time.Millisecond)
					resource2.Release(1)
					resource1.Release(1)
					return
				} else {
					// Can't get resource2, release resource1 and retry
					fmt.Println("Worker 1: Can't get resource 2, releasing resource 1")
					resource1.Release(1)
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
		fmt.Println("Worker 1: Gave up")
	}()
	
	// Worker 2 needs both resources (in reverse order)
	go func() {
		for i := 0; i < 5; i++ {
			if resource2.TryAcquire(2) {
				time.Sleep(100 * time.Millisecond)
				
				if resource1.TryAcquire(2) {
					fmt.Println("Worker 2: Got both resources, doing work...")
					time.Sleep(200 * time.Millisecond)
					resource1.Release(2)
					resource2.Release(2)
					return
				} else {
					// Can't get resource1, release resource2 and retry
					fmt.Println("Worker 2: Can't get resource 1, releasing resource 2")
					resource2.Release(2)
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
		fmt.Println("Worker 2: Gave up")
	}()
	
	time.Sleep(3 * time.Second)
	fmt.Println("Resource livelock example completed")
}
```

### Cách phòng tránh Livelock

#### 1. Random Backoff

```go
import "math/rand"

func avoidLivelockWithRandomBackoff() {
	resource1 := &Resource{id: 1}
	resource2 := &Resource{id: 2}
	
	// Worker with random backoff
	worker := func(id int, first, second *Resource) {
		for i := 0; i < 10; i++ {
			if first.TryAcquire(id) {
				if second.TryAcquire(id) {
					fmt.Printf("Worker %d: Success! Got both resources\n", id)
					time.Sleep(100 * time.Millisecond)
					second.Release(id)
					first.Release(id)
					return
				} else {
					first.Release(id)
					// Random backoff to avoid livelock
					backoff := time.Duration(rand.Intn(100)+50) * time.Millisecond
					fmt.Printf("Worker %d: Backing off for %v\n", id, backoff)
					time.Sleep(backoff)
				}
			} else {
				backoff := time.Duration(rand.Intn(100)+50) * time.Millisecond
				time.Sleep(backoff)
			}
		}
		fmt.Printf("Worker %d: Gave up\n", id)
	}
	
	go worker(1, resource1, resource2)
	go worker(2, resource2, resource1)
	
	time.Sleep(3 * time.Second)
	fmt.Println("Random backoff example completed")
}
```

#### 2. Timeout và Circuit Breaker

```go
func avoidLivelockWithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	resource1 := &Resource{id: 1}
	resource2 := &Resource{id: 2}
	
	worker := func(id int, first, second *Resource) {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d: Timeout, giving up\n", id)
				return
			default:
				if first.TryAcquire(id) {
					if second.TryAcquire(id) {
						fmt.Printf("Worker %d: Success with timeout protection\n", id)
						time.Sleep(100 * time.Millisecond)
						second.Release(id)
						first.Release(id)
						return
					} else {
						first.Release(id)
					}
				}
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
	
	go worker(1, resource1, resource2)
	go worker(2, resource2, resource1)
	
	<-ctx.Done()
	fmt.Println("Timeout protection example completed")
}
```

---

## So sánh ba khái niệm

| Đặc điểm | Deadlock | Starvation | Livelock |
|----------|----------|------------|----------|
| **Trạng thái Process** | Blocked/Waiting | Waiting | Running |
| **CPU Usage** | Không sử dụng CPU | Không sử dụng CPU | Sử dụng CPU |
| **Progress** | Không có tiến triển | Không có tiến triển | Không có tiến triển |
| **Phát hiện** | Dễ phát hiện | Khó phát hiện | Rất khó phát hiện |
| **Nguyên nhân** | Circular wait | Priority inversion | Continuous retry |
| **Giải pháp** | Lock ordering, timeout | Fair scheduling | Random backoff |

---

## Best Practices

### 1. Deadlock Prevention

```go
// Always acquire locks in the same order
type BankAccount struct {
	id      int
	balance int
	mu      sync.Mutex
}

func Transfer(from, to *BankAccount, amount int) error {
	// Always lock the account with smaller ID first
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
		return fmt.Errorf("insufficient funds")
	}
	
	from.balance -= amount
	to.balance += amount
	return nil
}
```

### 2. Starvation Prevention

```go
// Use fair scheduling with aging
type FairQueue struct {
	highPriority []Task
	lowPriority  []Task
	mu           sync.Mutex
	lowPriorityAge int
}

func (fq *FairQueue) GetNextTask() *Task {
	fq.mu.Lock()
	defer fq.mu.Unlock()
	
	// Age low priority tasks
	fq.lowPriorityAge++
	
	// If low priority tasks have been waiting too long, serve them
	if fq.lowPriorityAge > 5 && len(fq.lowPriority) > 0 {
		task := fq.lowPriority[0]
		fq.lowPriority = fq.lowPriority[1:]
		fq.lowPriorityAge = 0
		return &task
	}
	
	// Otherwise serve high priority if available
	if len(fq.highPriority) > 0 {
		task := fq.highPriority[0]
		fq.highPriority = fq.highPriority[1:]
		return &task
	}
	
	// Serve low priority if no high priority
	if len(fq.lowPriority) > 0 {
		task := fq.lowPriority[0]
		fq.lowPriority = fq.lowPriority[1:]
		fq.lowPriorityAge = 0
		return &task
	}
	
	return nil
}
```

### 3. Livelock Prevention

```go
// Use exponential backoff with jitter
func ExponentialBackoffWithJitter(attempt int) time.Duration {
	base := time.Duration(100) * time.Millisecond
	max := time.Duration(5) * time.Second
	
	// Exponential backoff: 100ms, 200ms, 400ms, 800ms, ...
	backoff := base * time.Duration(1<<uint(attempt))
	if backoff > max {
		backoff = max
	}
	
	// Add jitter (random component)
	jitter := time.Duration(rand.Intn(int(backoff/2)))
	return backoff + jitter
}

func RetryWithBackoff(operation func() error, maxAttempts int) error {
	for attempt := 0; attempt < maxAttempts; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}
		
		if attempt < maxAttempts-1 {
			backoff := ExponentialBackoffWithJitter(attempt)
			fmt.Printf("Attempt %d failed, retrying in %v\n", attempt+1, backoff)
			time.Sleep(backoff)
		}
	}
	return fmt.Errorf("operation failed after %d attempts", maxAttempts)
}
```

### 4. Monitoring và Detection

```go
// Deadlock detector
type DeadlockDetector struct {
	locks map[string]time.Time
	mu    sync.Mutex
}

func (dd *DeadlockDetector) AcquireLock(lockName string) {
	dd.mu.Lock()
	defer dd.mu.Unlock()
	dd.locks[lockName] = time.Now()
}

func (dd *DeadlockDetector) ReleaseLock(lockName string) {
	dd.mu.Lock()
	defer dd.mu.Unlock()
	delete(dd.locks, lockName)
}

func (dd *DeadlockDetector) CheckForDeadlocks() {
	dd.mu.Lock()
	defer dd.mu.Unlock()
	
	for lockName, acquireTime := range dd.locks {
		if time.Since(acquireTime) > 10*time.Second {
			fmt.Printf("WARNING: Potential deadlock detected on lock %s\n", lockName)
		}
	}
}
```

---

## Overhead trong Concurrency

### Khái niệm Overhead

**Overhead** trong lập trình đồng thời là chi phí bổ sung (thời gian, bộ nhớ, tài nguyên) mà hệ thống phải trả để thực hiện một tác vụ, ngoài công việc chính cần làm. Overhead là trade-off không thể tránh khỏi khi sử dụng concurrency, nhưng có thể được tối ưu hóa thông qua thiết kế hệ thống thông minh.

### Các loại Overhead

#### 1. Context Switching Overhead
Chi phí chuyển đổi giữa các goroutine/thread:
- Lưu và khôi phục trạng thái CPU
- Thời gian scheduler quyết định goroutine nào chạy tiếp
- Cache misses khi chuyển đổi context

#### 2. Memory Overhead
- Bộ nhớ bổ sung cho stack của mỗi goroutine (2KB initial)
- Metadata của sync primitives (mutex, channel)
- Garbage collection overhead
- Memory alignment và padding

#### 3. Synchronization Overhead
- Chi phí lock/unlock mutex
- Channel operations (send/receive)
- Atomic operations
- Wait group operations

#### 4. Communication Overhead
- Chi phí truyền data qua channels
- Network latency trong distributed systems
- Serialization/deserialization
- Buffer management

### Ví dụ Overhead trong Go

#### Ví dụ 1: So sánh Direct Access vs Concurrent Access

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// Overhead thấp - direct access
func directSum(data []int) int {
	sum := 0
	for _, v := range data {
		sum += v
	}
	return sum
}

// Overhead cao - với synchronization
func concurrentSum(data []int) int {
	var mu sync.Mutex
	sum := 0
	
	var wg sync.WaitGroup
	for _, v := range data {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			mu.Lock()         // Synchronization overhead
			sum += val        // Actual work
			mu.Unlock()       // Synchronization overhead
		}(v)
	}
	wg.Wait()
	return sum
}

// Tối ưu - chia nhỏ công việc
func optimizedConcurrentSum(data []int, numWorkers int) int {
	if len(data) < numWorkers {
		return directSum(data)
	}
	
	chunkSize := len(data) / numWorkers
	results := make(chan int, numWorkers)
	
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = len(data)
		}
		
		go func(chunk []int) {
			sum := 0
			for _, v := range chunk {
				sum += v
			}
			results <- sum
		}(data[start:end])
	}
	
	totalSum := 0
	for i := 0; i < numWorkers; i++ {
		totalSum += <-results
	}
	
	return totalSum
}

func benchmarkOverhead() {
	data := make([]int, 1000000)
	for i := range data {
		data[i] = i + 1
	}
	
	// Benchmark direct access
	start := time.Now()
	directResult := directSum(data)
	directTime := time.Since(start)
	
	// Benchmark concurrent access (high overhead)
	start = time.Now()
	concurrentResult := concurrentSum(data)
	concurrentTime := time.Since(start)
	
	// Benchmark optimized concurrent access
	start = time.Now()
	optimizedResult := optimizedConcurrentSum(data, 4)
	optimizedTime := time.Since(start)
	
	fmt.Printf("Direct: %d (Time: %v)\n", directResult, directTime)
	fmt.Printf("Concurrent: %d (Time: %v, Overhead: %.2fx)\n", 
		concurrentResult, concurrentTime, float64(concurrentTime)/float64(directTime))
	fmt.Printf("Optimized: %d (Time: %v, Overhead: %.2fx)\n", 
		optimizedResult, optimizedTime, float64(optimizedTime)/float64(directTime))
}
```

#### Ví dụ 2: Channel vs Mutex Overhead

```go
// Channel-based counter (higher overhead)
type ChannelCounter struct {
	ch chan int
	value int
}

func NewChannelCounter() *ChannelCounter {
	c := &ChannelCounter{
		ch: make(chan int, 1),
	}
	c.ch <- 0 // Initialize
	return c
}

func (c *ChannelCounter) Increment() {
	val := <-c.ch
	c.ch <- val + 1
}

func (c *ChannelCounter) Value() int {
	val := <-c.ch
	c.ch <- val
	return val
}

// Mutex-based counter (lower overhead)
type MutexCounter struct {
	mu    sync.Mutex
	value int
}

func (c *MutexCounter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *MutexCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Atomic counter (lowest overhead)
type AtomicCounter struct {
	value int64
}

func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}
```

#### Ví dụ 3: Overhead trong Round Robin Scheduler

```go
// Trong FairScheduler từ ví dụ trước
func (fs *FairScheduler) Start() {
	for {
		// Overhead: kiểm tra state
		switch fs.lastServed {
		case "high":
			// Overhead: context switching, channel operations
			select {
			case task := <-fs.lowPriorityQueue:
				task()                    // Actual work
				fs.lastServed = "low"     // State management overhead
			case task := <-fs.highPriorityQueue:
				task()
			}
		case "low":
			// Tương tự overhead cho case này
			select {
			case task := <-fs.highPriorityQueue:
				task()
				fs.lastServed = "high"
			case task := <-fs.lowPriorityQueue:
				task()
			}
		}
	}
}
```

### Cách đo lường và tối ưu Overhead

#### 1. Profiling và Benchmarking

```go
// Sử dụng Go's built-in benchmarking
func BenchmarkDirectSum(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		directSum(data)
	}
}

func BenchmarkConcurrentSum(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		concurrentSum(data)
	}
}

// Chạy: go test -bench=. -benchmem
```

#### 2. Strategies để giảm Overhead

**a) Batch Processing**
```go
// Thay vì xử lý từng item
for _, item := range items {
	go processItem(item) // Overhead cao
}

// Xử lý theo batch
chunkSize := 100
for i := 0; i < len(items); i += chunkSize {
	end := i + chunkSize
	if end > len(items) {
		end = len(items)
	}
	go processBatch(items[i:end]) // Overhead thấp hơn
}
```

**b) Worker Pool Pattern**
```go
type WorkerPool struct {
	tasks   chan Task
	workers int
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	wp := &WorkerPool{
		tasks:   make(chan Task, 100),
		workers: numWorkers,
	}
	
	// Tạo workers một lần, tái sử dụng
	for i := 0; i < numWorkers; i++ {
		go wp.worker()
	}
	
	return wp
}

func (wp *WorkerPool) worker() {
	for task := range wp.tasks {
		task() // Không có overhead tạo goroutine mới
	}
}
```

**c) Lock-free Algorithms**
```go
// Sử dụng atomic operations thay vì mutex
type LockFreeCounter struct {
	value int64
}

func (c *LockFreeCounter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

func (c *LockFreeCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}
```

**d) Buffered Channels**
```go
// Unbuffered channel - blocking overhead
ch := make(chan int)

// Buffered channel - giảm blocking
ch := make(chan int, 100)
```

#### 3. Monitoring Overhead

```go
type OverheadMonitor struct {
	goroutineCount int64
	memoryUsage    int64
	contextSwitches int64
}

func (om *OverheadMonitor) TrackGoroutines() {
	for {
		count := runtime.NumGoroutine()
		atomic.StoreInt64(&om.goroutineCount, int64(count))
		time.Sleep(time.Second)
	}
}

func (om *OverheadMonitor) GetStats() (int64, int64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return atomic.LoadInt64(&om.goroutineCount), int64(m.Alloc)
}
```

#### 4. Best Practices để giảm Overhead

1. **Đo lường trước khi tối ưu**: "Premature optimization is the root of all evil"
2. **Sử dụng đúng tool cho đúng job**:
   - Atomic cho simple counters
   - Mutex cho critical sections
   - Channels cho communication
3. **Worker pools** thay vì tạo goroutines liên tục
4. **Batch processing** cho high-throughput scenarios
5. **Buffered channels** để giảm blocking
6. **Context với timeout** để tránh resource leaks
7. **Profile thường xuyên** để phát hiện bottlenecks

---

## Kết luận

Việc hiểu và xử lý Deadlock, Starvation và Livelock là rất quan trọng trong lập trình đồng thời với Go. Các nguyên tắc chính:

1. **Prevention is better than cure**: Thiết kế hệ thống để tránh các vấn đề này từ đầu
2. **Use timeouts**: Luôn sử dụng timeout để tránh chờ đợi vô thời hạn
3. **Fair scheduling**: Đảm bảo tất cả processes đều có cơ hội được thực thi
4. **Monitor and detect**: Implement monitoring để phát hiện sớm các vấn đề
5. **Test thoroughly**: Test với nhiều scenarios khác nhau để phát hiện race conditions

Bằng cách áp dụng các best practices và patterns này, bạn có thể xây dựng các ứng dụng Go concurrent an toàn và hiệu quả.