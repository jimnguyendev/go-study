package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// =============================================================================
// DEADLOCK EXAMPLES
// =============================================================================

// Ví dụ deadlock cơ bản với 2 mutexes
func deadlockExample() {
	fmt.Println("\n=== DEADLOCK EXAMPLE ===")
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
	fmt.Println("Main: Deadlock detected! Program would hang here.")
}

// Deadlock với unbuffered channels
func channelDeadlock() {
	fmt.Println("\n=== CHANNEL DEADLOCK EXAMPLE ===")
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

// Cách phòng tránh deadlock với lock ordering
func avoidDeadlockWithOrdering() {
	fmt.Println("\n=== AVOID DEADLOCK WITH ORDERING ===")
	var mutex1, mutex2 sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)

	// Luôn lock theo thứ tự: mutex1 trước, mutex2 sau
	go func() {
		defer wg.Done()
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
		defer wg.Done()
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

	wg.Wait()
	fmt.Println("No deadlock - both goroutines completed")
}

// Timeout với Context
func avoidDeadlockWithTimeout() {
	fmt.Println("\n=== AVOID DEADLOCK WITH TIMEOUT ===")
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

// =============================================================================
// STARVATION EXAMPLES
// =============================================================================

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
	fmt.Println("\n=== WRITER STARVATION EXAMPLE ===")
	resource := &RWResource{data: 0}

	// Tạo nhiều readers
	for i := 1; i <= 5; i++ {
		go func(id int) {
			for j := 0; j < 3; j++ {
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
	fmt.Println("\n=== PRIORITY STARVATION EXAMPLE ===")
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
	for i := 2; i <= 6; i++ {
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

// Fair Scheduling với Round Robin
type FairScheduler struct {
	highPriorityTasks chan Task
	lowPriorityTasks  chan Task
	lastServed        string // "high" or "low"
	mu                sync.Mutex
}

func NewFairScheduler() *FairScheduler {
	return &FairScheduler{
		highPriorityTasks: make(chan Task, 50),
		lowPriorityTasks:  make(chan Task, 50),
		lastServed:        "low", // Start with high priority
	}
}

func (fs *FairScheduler) AddHighPriorityTask(task Task) {
	fs.highPriorityTasks <- task
}

func (fs *FairScheduler) AddLowPriorityTask(task Task) {
	fs.lowPriorityTasks <- task
}

func (fs *FairScheduler) Start() {
	go func() {
		for {
			fs.mu.Lock()
			lastServed := fs.lastServed
			fs.mu.Unlock()

			if lastServed == "low" {
				// Try to serve high priority first
				select {
				case task := <-fs.highPriorityTasks:
					fmt.Printf("Fair: Executing high priority task %d\n", task.ID)
					task.Work()
					fs.mu.Lock()
					fs.lastServed = "high"
					fs.mu.Unlock()
				case task := <-fs.lowPriorityTasks:
					fmt.Printf("Fair: Executing low priority task %d\n", task.ID)
					task.Work()
					fs.mu.Lock()
					fs.lastServed = "low"
					fs.mu.Unlock()
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
					fs.mu.Lock()
					fs.lastServed = "low"
					fs.mu.Unlock()
				case task := <-fs.highPriorityTasks:
					fmt.Printf("Fair: Executing high priority task %d\n", task.ID)
					task.Work()
					fs.mu.Lock()
					fs.lastServed = "high"
					fs.mu.Unlock()
				case <-time.After(100 * time.Millisecond):
					// No tasks available
					continue
				}
			}
		}
	}()
}

func fairSchedulingExample() {
	fmt.Println("\n=== FAIR SCHEDULING EXAMPLE ===")
	scheduler := NewFairScheduler()
	scheduler.Start()

	// Add mixed priority tasks
	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			scheduler.AddHighPriorityTask(Task{
				ID:       i,
				Priority: 8,
				Work: func() {
					fmt.Println("High priority task completed")
				},
			})
		} else {
			scheduler.AddLowPriorityTask(Task{
				ID:       i,
				Priority: 2,
				Work: func() {
					fmt.Println("Low priority task completed")
				},
			})
		}
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)
	fmt.Println("Fair scheduling example completed")
}

// =============================================================================
// LIVELOCK EXAMPLES
// =============================================================================

type Person struct {
	name    string
	waiting bool
	mu      sync.Mutex
}

func (p *Person) SetWaiting(waiting bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.waiting = waiting
}

func (p *Person) IsWaiting() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.waiting
}

func (p *Person) TryToPass(other *Person, hallway chan bool) {
	for i := 0; i < 5; i++ {
		p.SetWaiting(true)
		// Try to enter hallway
		select {
		case hallway <- true:
			fmt.Printf("%s: Entered hallway\n", p.name)
			p.SetWaiting(false)

			// Check if other person is also trying
			time.Sleep(50 * time.Millisecond)

			// Be polite - if other person is waiting, step back
			if other.IsWaiting() {
				fmt.Printf("%s: Being polite, stepping back\n", p.name)
				<-hallway                          // Step back
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
	p.SetWaiting(false)
}

func livelockExample() {
	fmt.Println("\n=== LIVELOCK EXAMPLE ===")
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
	fmt.Println("\n=== RESOURCE LIVELOCK EXAMPLE ===")
	resource1 := &Resource{id: 1}
	resource2 := &Resource{id: 2}

	var wg sync.WaitGroup
	wg.Add(2)

	// Worker 1 needs both resources
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
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
		defer wg.Done()
		for i := 0; i < 3; i++ {
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

	wg.Wait()
	fmt.Println("Resource livelock example completed")
}

// Cách phòng tránh livelock với random backoff
func avoidLivelockWithRandomBackoff() {
	fmt.Println("\n=== AVOID LIVELOCK WITH RANDOM BACKOFF ===")
	resource1 := &Resource{id: 1}
	resource2 := &Resource{id: 2}

	var wg sync.WaitGroup
	wg.Add(2)

	// Worker with random backoff
	worker := func(id int, first, second *Resource) {
		defer wg.Done()
		for i := 0; i < 5; i++ {
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

	wg.Wait()
	fmt.Println("Random backoff example completed")
}

// Timeout và Circuit Breaker
func avoidLivelockWithTimeout() {
	fmt.Println("\n=== AVOID LIVELOCK WITH TIMEOUT ===")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resource1 := &Resource{id: 1}
	resource2 := &Resource{id: 2}

	var wg sync.WaitGroup
	wg.Add(2)

	worker := func(id int, first, second *Resource) {
		defer wg.Done()
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

	wg.Wait()
	fmt.Println("Timeout protection example completed")
}

// =============================================================================
// MAIN FUNCTION
// =============================================================================

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Go Concurrency Problems Demo")
	fmt.Println("============================")

	// DEADLOCK EXAMPLES
	fmt.Println("\n🔒 DEADLOCK EXAMPLES:")
	// Note: These will actually cause deadlock, so we run them with timeouts
	go func() {
		deadlockExample()
	}()
	time.Sleep(3 * time.Second)

	go func() {
		channelDeadlock()
	}()
	time.Sleep(3 * time.Second)

	// Solutions
	avoidDeadlockWithOrdering()
	avoidDeadlockWithTimeout()

	// STARVATION EXAMPLES
	fmt.Println("\n🍽️ STARVATION EXAMPLES:")
	writerStarvationExample()
	priorityStarvationExample()
	fairSchedulingExample()

	// LIVELOCK EXAMPLES
	fmt.Println("\n🔄 LIVELOCK EXAMPLES:")
	livelockExample()
	resourceLivelockExample()
	avoidLivelockWithRandomBackoff()
	avoidLivelockWithTimeout()

	fmt.Println("\n✅ All examples completed!")
	fmt.Println("\nNote: Some examples demonstrate problems that would normally")
	fmt.Println("cause the program to hang. They are run with timeouts for demonstration.")
}
