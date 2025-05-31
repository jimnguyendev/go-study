# Go Generics - Tài liệu Chi tiết

## Mục lục

1. [Giới thiệu về Generics](#giới-thiệu-về-generics)
2. [Tại sao cần Generics](#tại-sao-cần-generics)
3. [Cú pháp cơ bản](#cú-pháp-cơ-bản)
4. [Type Parameters và Type Constraints](#type-parameters-và-type-constraints)
5. [Generic Functions](#generic-functions)
6. [Generic Types](#generic-types)
7. [Type Inference](#type-inference)
8. [Built-in Constraints](#built-in-constraints)
9. [Custom Constraints](#custom-constraints)
10. [Ví dụ thực tế](#ví-dụ-thực-tế)
11. [Best Practices](#best-practices)
12. [Performance Considerations](#performance-considerations)
13. [Migration từ interface{}](#migration-từ-interface)

---

## Giới thiệu về Generics

Generics (hay Type Parameters) được giới thiệu trong Go 1.18, cho phép viết code có thể hoạt động với nhiều kiểu dữ liệu khác nhau mà vẫn đảm bảo type safety tại compile time.

### Trước Go 1.18 (Không có Generics):

```go
// Phải viết riêng cho từng type
func MaxInt(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func MaxFloat64(a, b float64) float64 {
    if a > b {
        return a
    }
    return b
}

// Hoặc sử dụng interface{} (mất type safety)
func Max(a, b interface{}) interface{} {
    // Cần type assertion, có thể panic
    switch a := a.(type) {
    case int:
        if b, ok := b.(int); ok && a > b {
            return a
        }
    case float64:
        if b, ok := b.(float64); ok && a > b {
            return a
        }
    }
    return b
}
```

### Với Go 1.18+ (Có Generics):

```go
// Một function cho tất cả comparable types
func Max[T comparable](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// Sử dụng
func main() {
    fmt.Println(Max(10, 20))       // 20 (int)
    fmt.Println(Max(3.14, 2.71))   // 3.14 (float64)
    fmt.Println(Max("hello", "world")) // "world" (string)
}
```

---

## Tại sao cần Generics

### 1. **Type Safety**
```go
// Trước Generics - Không type safe
func OldSliceContains(slice []interface{}, item interface{}) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

// Có thể gây lỗi runtime
ints := []interface{}{1, 2, 3}
result := OldSliceContains(ints, "hello") // Compile OK, nhưng logic sai

// Với Generics - Type safe
func SliceContains[T comparable](slice []T, item T) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

// Type safe tại compile time
ints := []int{1, 2, 3}
result := SliceContains(ints, "hello") // Compile ERROR!
```

### 2. **Performance**
```go
// Interface{} cần boxing/unboxing
func OldSum(numbers []interface{}) interface{} {
    var sum interface{}
    for _, num := range numbers {
        // Type assertion overhead
        if i, ok := num.(int); ok {
            if s, ok := sum.(int); ok {
                sum = s + i
            } else {
                sum = i
            }
        }
    }
    return sum
}

// Generics - No boxing, direct operations
func Sum[T ~int | ~float64](numbers []T) T {
    var sum T
    for _, num := range numbers {
        sum += num // Direct operation, no type assertion
    }
    return sum
}
```

### 3. **Code Reusability**
```go
// Trước đây phải duplicate code
type IntStack struct {
    items []int
}

func (s *IntStack) Push(item int) {
    s.items = append(s.items, item)
}

type StringStack struct {
    items []string
}

func (s *StringStack) Push(item string) {
    s.items = append(s.items, item)
}

// Với Generics - Một implementation cho tất cả
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    index := len(s.items) - 1
    item := s.items[index]
    s.items = s.items[:index]
    return item, true
}
```

---

## Cú pháp cơ bản

### Type Parameter Declaration

```go
// Cú pháp: [TypeParam Constraint]
func FunctionName[T Constraint](param T) T {
    // implementation
}

// Multiple type parameters
func FunctionName[T, U Constraint](param1 T, param2 U) (T, U) {
    // implementation
}

// Different constraints
func FunctionName[T Constraint1, U Constraint2](param1 T, param2 U) {
    // implementation
}
```

### Ví dụ cụ thể:

```go
package main

import "fmt"

// Generic function với một type parameter
func Identity[T any](value T) T {
    return value
}

// Generic function với multiple type parameters
func Pair[T, U any](first T, second U) (T, U) {
    return first, second
}

// Generic function với constraints
func Add[T ~int | ~float64](a, b T) T {
    return a + b
}

func main() {
    // Type inference
    fmt.Println(Identity(42))        // int
    fmt.Println(Identity("hello"))   // string
    fmt.Println(Identity(3.14))      // float64
    
    // Explicit type specification
    fmt.Println(Identity[string]("world"))
    
    // Multiple type parameters
    name, age := Pair("Alice", 30)
    fmt.Printf("%s is %d years old\n", name, age)
    
    // Constrained types
    fmt.Println(Add(10, 20))      // 30
    fmt.Println(Add(3.14, 2.86))  // 6.0
}
```

---

## Type Parameters và Type Constraints

### 1. **any Constraint**

```go
// any là alias của interface{}
func Print[T any](value T) {
    fmt.Println(value)
}

// Có thể sử dụng với bất kỳ type nào
Print(42)
Print("hello")
Print([]int{1, 2, 3})
Print(map[string]int{"a": 1})
```

### 2. **comparable Constraint**

```go
// comparable cho phép sử dụng == và !=
func Equal[T comparable](a, b T) bool {
    return a == b
}

// Works with:
fmt.Println(Equal(1, 1))           // true
fmt.Println(Equal("a", "b"))       // false
fmt.Println(Equal(3.14, 3.14))     // true

// Không work với slice, map, function
// Equal([]int{1}, []int{1})  // Compile error!
```

### 3. **Union Constraints**

```go
// Union của multiple types
func Numeric[T int | int32 | int64 | float32 | float64](value T) T {
    return value * 2
}

// Sử dụng type approximation (~)
func NumericApprox[T ~int | ~float64](value T) T {
    return value * 2
}

// Custom types work với ~
type MyInt int
type MyFloat float64

var mi MyInt = 10
var mf MyFloat = 3.14

fmt.Println(NumericApprox(mi))  // 20
fmt.Println(NumericApprox(mf))  // 6.28
```

---

## Generic Functions

### 1. **Utility Functions**

```go
// Map function
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Filter function
func Filter[T any](slice []T, predicate func(T) bool) []T {
    var result []T
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

// Reduce function
func Reduce[T, U any](slice []T, initial U, fn func(U, T) U) U {
    result := initial
    for _, v := range slice {
        result = fn(result, v)
    }
    return result
}

// Sử dụng
func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // Map: int -> string
    strings := Map(numbers, func(n int) string {
        return fmt.Sprintf("num_%d", n)
    })
    fmt.Println(strings) // ["num_1", "num_2", "num_3", "num_4", "num_5"]
    
    // Filter: chỉ số chẵn
    evens := Filter(numbers, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println(evens) // [2, 4]
    
    // Reduce: tính tổng
    sum := Reduce(numbers, 0, func(acc, n int) int {
        return acc + n
    })
    fmt.Println(sum) // 15
}
```

### 2. **Slice Utilities**

```go
// Reverse slice
func Reverse[T any](slice []T) []T {
    result := make([]T, len(slice))
    for i, v := range slice {
        result[len(slice)-1-i] = v
    }
    return result
}

// Find element
func Find[T comparable](slice []T, target T) (int, bool) {
    for i, v := range slice {
        if v == target {
            return i, true
        }
    }
    return -1, false
}

// Remove duplicates
func Unique[T comparable](slice []T) []T {
    seen := make(map[T]bool)
    var result []T
    
    for _, v := range slice {
        if !seen[v] {
            seen[v] = true
            result = append(result, v)
        }
    }
    return result
}

// Chunk slice
func Chunk[T any](slice []T, size int) [][]T {
    if size <= 0 {
        return nil
    }
    
    var chunks [][]T
    for i := 0; i < len(slice); i += size {
        end := i + size
        if end > len(slice) {
            end = len(slice)
        }
        chunks = append(chunks, slice[i:end])
    }
    return chunks
}
```

---

## Generic Types

### 1. **Generic Structs**

```go
// Generic Pair
type Pair[T, U any] struct {
    First  T
    Second U
}

func (p Pair[T, U]) String() string {
    return fmt.Sprintf("(%v, %v)", p.First, p.Second)
}

// Generic Result type (như Rust's Result)
type Result[T any] struct {
    value T
    err   error
}

func NewResult[T any](value T, err error) Result[T] {
    return Result[T]{value: value, err: err}
}

func (r Result[T]) IsOk() bool {
    return r.err == nil
}

func (r Result[T]) IsErr() bool {
    return r.err != nil
}

func (r Result[T]) Unwrap() T {
    if r.err != nil {
        panic(r.err)
    }
    return r.value
}

func (r Result[T]) UnwrapOr(defaultValue T) T {
    if r.err != nil {
        return defaultValue
    }
    return r.value
}
```

### 2. **Generic Collections**

```go
// Generic Stack
type Stack[T any] struct {
    items []T
    mu    sync.RWMutex
}

func NewStack[T any]() *Stack[T] {
    return &Stack[T]{
        items: make([]T, 0),
    }
}

func (s *Stack[T]) Push(item T) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    
    index := len(s.items) - 1
    item := s.items[index]
    s.items = s.items[:index]
    return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    
    return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Size() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return len(s.items)
}

func (s *Stack[T]) IsEmpty() bool {
    return s.Size() == 0
}
```

### 3. **Generic Queue**

```go
type Queue[T any] struct {
    items []T
    mu    sync.RWMutex
}

func NewQueue[T any]() *Queue[T] {
    return &Queue[T]{
        items: make([]T, 0),
    }
}

func (q *Queue[T]) Enqueue(item T) {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    if len(q.items) == 0 {
        var zero T
        return zero, false
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    return item, true
}

func (q *Queue[T]) Front() (T, bool) {
    q.mu.RLock()
    defer q.mu.RUnlock()
    
    if len(q.items) == 0 {
        var zero T
        return zero, false
    }
    
    return q.items[0], true
}
```

---

## Type Inference

Go compiler có thể tự động suy luận type trong nhiều trường hợp:

```go
func Max[T comparable](a, b T) T {
    if a > b {
        return a
    }
    return b
}

func main() {
    // Type inference - compiler tự suy luận
    result1 := Max(10, 20)        // T = int
    result2 := Max(3.14, 2.71)    // T = float64
    result3 := Max("hello", "hi") // T = string
    
    // Explicit type specification
    result4 := Max[int](10, 20)
    result5 := Max[float64](3.14, 2.71)
    
    // Mixed types - cần explicit type
    // result6 := Max(10, 3.14)     // Error: type mismatch
    result6 := Max[float64](10, 3.14) // OK: 10 converted to float64
}
```

### Khi Type Inference không hoạt động:

```go
func MakeSlice[T any](size int) []T {
    return make([]T, size)
}

func main() {
    // Cần explicit type vì compiler không thể suy luận T
    ints := MakeSlice[int](5)       // []int{0, 0, 0, 0, 0}
    strings := MakeSlice[string](3) // []string{"", "", ""}
    
    // Error: cannot infer T
    // slice := MakeSlice(5)
}
```

---

## Built-in Constraints

### 1. **golang.org/x/exp/constraints**

```go
import "golang.org/x/exp/constraints"

// Signed integers
func AbsSigned[T constraints.Signed](x T) T {
    if x < 0 {
        return -x
    }
    return x
}

// Unsigned integers
func MaxUnsigned[T constraints.Unsigned](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// All integers
func IsEven[T constraints.Integer](x T) bool {
    return x%2 == 0
}

// Floating point
func Round[T constraints.Float](x T) T {
    return T(math.Round(float64(x)))
}

// All ordered types (can use <, <=, >, >=)
func Clamp[T constraints.Ordered](value, min, max T) T {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}
```

### 2. **Sử dụng Built-in Constraints**

```go
func main() {
    // Signed integers
    fmt.Println(AbsSigned(-42))    // 42
    fmt.Println(AbsSigned(int8(-10))) // 10
    
    // Unsigned integers
    fmt.Println(MaxUnsigned(uint(10), uint(20))) // 20
    
    // All integers
    fmt.Println(IsEven(4))   // true
    fmt.Println(IsEven(5))   // false
    
    // Floating point
    fmt.Println(Round(3.14)) // 3
    fmt.Println(Round(3.67)) // 4
    
    // Ordered types
    fmt.Println(Clamp(15, 10, 20))   // 15
    fmt.Println(Clamp(5, 10, 20))    // 10
    fmt.Println(Clamp(25, 10, 20))   // 20
    fmt.Println(Clamp("m", "a", "z")) // "m"
}
```

---

## Custom Constraints

### 1. **Interface-based Constraints**

```go
// Constraint cho types có method String()
type Stringer interface {
    String() string
}

func PrintAll[T Stringer](items []T) {
    for _, item := range items {
        fmt.Println(item.String())
    }
}

// Constraint cho numeric operations
type Numeric interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64
}

func Sum[T Numeric](numbers []T) T {
    var sum T
    for _, num := range numbers {
        sum += num
    }
    return sum
}

// Constraint kết hợp
type ComparableStringer interface {
    comparable
    Stringer
}

func FindByString[T ComparableStringer](items []T, target string) (T, bool) {
    for _, item := range items {
        if item.String() == target {
            return item, true
        }
    }
    var zero T
    return zero, false
}
```

### 2. **Method Constraints**

```go
// Constraint cho types có method Close()
type Closer interface {
    Close() error
}

// Generic function để đóng resources
func CloseAll[T Closer](resources []T) []error {
    var errors []error
    for _, resource := range resources {
        if err := resource.Close(); err != nil {
            errors = append(errors, err)
        }
    }
    return errors
}

// Constraint cho serializable types
type Serializable interface {
    Marshal() ([]byte, error)
    Unmarshal([]byte) error
}

func SaveToFile[T Serializable](item T, filename string) error {
    data, err := item.Marshal()
    if err != nil {
        return err
    }
    return os.WriteFile(filename, data, 0644)
}

func LoadFromFile[T Serializable](filename string) (T, error) {
    var item T
    data, err := os.ReadFile(filename)
    if err != nil {
        return item, err
    }
    err = item.Unmarshal(data)
    return item, err
}
```

### 3. **Complex Constraints**

```go
// Constraint cho collection types
type Collection[T any] interface {
    Add(T)
    Remove(T) bool
    Contains(T) bool
    Size() int
    Clear()
    ToSlice() []T
}

// Generic function hoạt động với bất kỳ collection nào
func ProcessCollection[T comparable, C Collection[T]](coll C, items []T) {
    // Add all items
    for _, item := range items {
        coll.Add(item)
    }
    
    fmt.Printf("Collection size: %d\n", coll.Size())
    
    // Check if contains specific items
    for _, item := range items[:min(3, len(items))] {
        fmt.Printf("Contains %v: %t\n", item, coll.Contains(item))
    }
}

// Constraint cho comparable và có zero value
type Zeroable interface {
    comparable
    IsZero() bool
}

func RemoveZeros[T Zeroable](slice []T) []T {
    var result []T
    for _, item := range slice {
        if !item.IsZero() {
            result = append(result, item)
        }
    }
    return result
}
```

---

## Ví dụ thực tế

### 1. **Generic HTTP Client**

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// Generic HTTP Client
type HTTPClient struct {
    client  *http.Client
    baseURL string
    headers map[string]string
}

func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
    return &HTTPClient{
        client: &http.Client{
            Timeout: timeout,
        },
        baseURL: baseURL,
        headers: make(map[string]string),
    }
}

func (c *HTTPClient) SetHeader(key, value string) {
    c.headers[key] = value
}

// Generic GET request
func (c *HTTPClient) Get[T any](endpoint string) (*T, error) {
    req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
    if err != nil {
        return nil, err
    }
    
    // Add headers
    for key, value := range c.headers {
        req.Header.Set(key, value)
    }
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result T
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

// Generic POST request
func (c *HTTPClient) Post[T, U any](endpoint string, payload T) (*U, error) {
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }
    
    req, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    for key, value := range c.headers {
        req.Header.Set(key, value)
    }
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result U
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

// Usage example
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserResponse struct {
    ID      int    `json:"id"`
    Message string `json:"message"`
}

func main() {
    client := NewHTTPClient("https://api.example.com", 30*time.Second)
    client.SetHeader("Authorization", "Bearer token123")
    
    // GET request
    user, err := client.Get[User]("/users/1")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("User: %+v\n", user)
    
    // POST request
    createReq := CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    createResp, err := client.Post[CreateUserRequest, CreateUserResponse]("/users", createReq)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Created: %+v\n", createResp)
}
```

### 2. **Generic Cache System**

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Cache item với expiration
type CacheItem[T any] struct {
    Value     T
    ExpiresAt time.Time
}

func (ci *CacheItem[T]) IsExpired() bool {
    return time.Now().After(ci.ExpiresAt)
}

// Generic Cache
type Cache[K comparable, V any] struct {
    items map[K]*CacheItem[V]
    mu    sync.RWMutex
    ttl   time.Duration
}

func NewCache[K comparable, V any](ttl time.Duration) *Cache[K, V] {
    cache := &Cache[K, V]{
        items: make(map[K]*CacheItem[V]),
        ttl:   ttl,
    }
    
    // Start cleanup goroutine
    go cache.cleanup()
    
    return cache
}

func (c *Cache[K, V]) Set(key K, value V) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.items[key] = &CacheItem[V]{
        Value:     value,
        ExpiresAt: time.Now().Add(c.ttl),
    }
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    item, exists := c.items[key]
    if !exists || item.IsExpired() {
        var zero V
        return zero, false
    }
    
    return item.Value, true
}

func (c *Cache[K, V]) Delete(key K) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    delete(c.items, key)
}

func (c *Cache[K, V]) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.items = make(map[K]*CacheItem[V])
}

func (c *Cache[K, V]) Size() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    return len(c.items)
}

func (c *Cache[K, V]) Keys() []K {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    keys := make([]K, 0, len(c.items))
    for key := range c.items {
        keys = append(keys, key)
    }
    return keys
}

// Cleanup expired items
func (c *Cache[K, V]) cleanup() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        c.mu.Lock()
        for key, item := range c.items {
            if item.IsExpired() {
                delete(c.items, key)
            }
        }
        c.mu.Unlock()
    }
}

// GetOrSet - atomic get or set operation
func (c *Cache[K, V]) GetOrSet(key K, factory func() V) V {
    // Try to get first
    if value, exists := c.Get(key); exists {
        return value
    }
    
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // Double-check after acquiring write lock
    if item, exists := c.items[key]; exists && !item.IsExpired() {
        return item.Value
    }
    
    // Create new value
    value := factory()
    c.items[key] = &CacheItem[V]{
        Value:     value,
        ExpiresAt: time.Now().Add(c.ttl),
    }
    
    return value
}

// Usage example
func main() {
    // String cache
    stringCache := NewCache[string, string](5 * time.Minute)
    stringCache.Set("user:1", "John Doe")
    
    if name, exists := stringCache.Get("user:1"); exists {
        fmt.Printf("User name: %s\n", name)
    }
    
    // Struct cache
    type UserProfile struct {
        ID    int
        Name  string
        Email string
    }
    
    profileCache := NewCache[int, UserProfile](10 * time.Minute)
    profileCache.Set(1, UserProfile{
        ID:    1,
        Name:  "Alice",
        Email: "alice@example.com",
    })
    
    if profile, exists := profileCache.Get(1); exists {
        fmt.Printf("Profile: %+v\n", profile)
    }
    
    // GetOrSet example
    expensiveData := profileCache.GetOrSet(2, func() UserProfile {
        fmt.Println("Computing expensive data...")
        time.Sleep(100 * time.Millisecond) // Simulate expensive operation
        return UserProfile{
            ID:    2,
            Name:  "Bob",
            Email: "bob@example.com",
        }
    })
    fmt.Printf("Expensive data: %+v\n", expensiveData)
    
    // Second call should use cached value
    cachedData := profileCache.GetOrSet(2, func() UserProfile {
        fmt.Println("This should not be called")
        return UserProfile{}
    })
    fmt.Printf("Cached data: %+v\n", cachedData)
}
```

### 3. **Generic Repository Pattern**

```go
package main

import (
    "database/sql"
    "fmt"
    "reflect"
    "strings"
)

// Entity interface - tất cả entities phải implement
type Entity interface {
    GetID() interface{}
    SetID(interface{})
    TableName() string
}

// Generic Repository
type Repository[T Entity] struct {
    db *sql.DB
}

func NewRepository[T Entity](db *sql.DB) *Repository[T] {
    return &Repository[T]{db: db}
}

// Create
func (r *Repository[T]) Create(entity T) error {
    tableName := entity.TableName()
    
    // Use reflection để build INSERT query
    v := reflect.ValueOf(entity).Elem()
    t := reflect.TypeOf(entity).Elem()
    
    var columns []string
    var placeholders []string
    var values []interface{}
    
    for i := 0; i < v.NumField(); i++ {
        field := t.Field(i)
        if field.Name == "ID" {
            continue // Skip ID for auto-increment
        }
        
        tag := field.Tag.Get("db")
        if tag == "" {
            tag = strings.ToLower(field.Name)
        }
        
        columns = append(columns, tag)
        placeholders = append(placeholders, "?")
        values = append(values, v.Field(i).Interface())
    }
    
    query := fmt.Sprintf(
        "INSERT INTO %s (%s) VALUES (%s)",
        tableName,
        strings.Join(columns, ", "),
        strings.Join(placeholders, ", "),
    )
    
    result, err := r.db.Exec(query, values...)
    if err != nil {
        return err
    }
    
    // Set ID if auto-increment
    if id, err := result.LastInsertId(); err == nil {
        entity.SetID(id)
    }
    
    return nil
}

// FindByID
func (r *Repository[T]) FindByID(id interface{}) (*T, error) {
    var entity T
    
    // Create new instance
    entityType := reflect.TypeOf(entity).Elem()
    entityValue := reflect.New(entityType)
    entity = entityValue.Interface().(T)
    
    tableName := entity.TableName()
    
    query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)
    row := r.db.QueryRow(query, id)
    
    // Scan into entity fields
    v := reflect.ValueOf(entity).Elem()
    t := reflect.TypeOf(entity).Elem()
    
    var scanArgs []interface{}
    for i := 0; i < v.NumField(); i++ {
        scanArgs = append(scanArgs, v.Field(i).Addr().Interface())
    }
    
    if err := row.Scan(scanArgs...); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    
    return &entity, nil
}

// FindAll
func (r *Repository[T]) FindAll() ([]T, error) {
    var entity T
    tableName := entity.TableName()
    
    query := fmt.Sprintf("SELECT * FROM %s", tableName)
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var entities []T
    
    for rows.Next() {
        // Create new instance
        entityType := reflect.TypeOf(entity).Elem()
        entityValue := reflect.New(entityType)
        newEntity := entityValue.Interface().(T)
        
        // Scan into entity fields
        v := reflect.ValueOf(newEntity).Elem()
        
        var scanArgs []interface{}
        for i := 0; i < v.NumField(); i++ {
            scanArgs = append(scanArgs, v.Field(i).Addr().Interface())
        }
        
        if err := rows.Scan(scanArgs...); err != nil {
            return nil, err
        }
        
        entities = append(entities, newEntity)
    }
    
    return entities, nil
}

// Update
func (r *Repository[T]) Update(entity T) error {
    tableName := entity.TableName()
    id := entity.GetID()
    
    v := reflect.ValueOf(entity).Elem()
    t := reflect.TypeOf(entity).Elem()
    
    var setParts []string
    var values []interface{}
    
    for i := 0; i < v.NumField(); i++ {
        field := t.Field(i)
        if field.Name == "ID" {
            continue // Skip ID
        }
        
        tag := field.Tag.Get("db")
        if tag == "" {
            tag = strings.ToLower(field.Name)
        }
        
        setParts = append(setParts, fmt.Sprintf("%s = ?", tag))
        values = append(values, v.Field(i).Interface())
    }
    
    values = append(values, id) // Add ID for WHERE clause
    
    query := fmt.Sprintf(
        "UPDATE %s SET %s WHERE id = ?",
        tableName,
        strings.Join(setParts, ", "),
    )
    
    _, err := r.db.Exec(query, values...)
    return err
}

// Delete
func (r *Repository[T]) Delete(id interface{}) error {
    var entity T
    tableName := entity.TableName()
    
    query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
    _, err := r.db.Exec(query, id)
    return err
}

// Example entities
type User struct {
    ID    int64  `db:"id"`
    Name  string `db:"name"`
    Email string `db:"email"`
}

func (u *User) GetID() interface{} {
    return u.ID
}

func (u *User) SetID(id interface{}) {
    u.ID = id.(int64)
}

func (u *User) TableName() string {
    return "users"
}

type Product struct {
    ID    int64   `db:"id"`
    Name  string  `db:"name"`
    Price float64 `db:"price"`
}

func (p *Product) GetID() interface{} {
    return p.ID
}

func (p *Product) SetID(id interface{}) {
    p.ID = id.(int64)
}

func (p *Product) TableName() string {
    return "products"
}

// Usage example
func main() {
    // Assume db is initialized
    var db *sql.DB
    
    // User repository
    userRepo := NewRepository[*User](db)
    
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Create
    if err := userRepo.Create(user); err != nil {
        fmt.Printf("Error creating user: %v\n", err)
    }
    
    // Find by ID
    foundUser, err := userRepo.FindByID(user.ID)
    if err != nil {
        fmt.Printf("Error finding user: %v\n", err)
    } else if foundUser != nil {
        fmt.Printf("Found user: %+v\n", foundUser)
    }
    
    // Product repository
    productRepo := NewRepository[*Product](db)
    
    product := &Product{
        Name:  "Laptop",
        Price: 999.99,
    }
    
    if err := productRepo.Create(product); err != nil {
        fmt.Printf("Error creating product: %v\n", err)
    }
    
    // Find all products
    products, err := productRepo.FindAll()
    if err != nil {
        fmt.Printf("Error finding products: %v\n", err)
    } else {
        fmt.Printf("Found %d products\n", len(products))
    }
}
```

---

## Best Practices

### 1. **Naming Conventions**

```go
// ✅ GOOD: Descriptive type parameter names
func Map[T, U any](slice []T, fn func(T) U) []U { ... }
func Cache[Key comparable, Value any]() { ... }
func Repository[Entity any]() { ... }

// ❌ BAD: Non-descriptive names
func Map[A, B any](slice []A, fn func(A) B) []B { ... }
func Cache[K, V any]() { ... } // OK for short functions
```

### 2. **Constraint Design**

```go
// ✅ GOOD: Specific constraints
type Numeric interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64
}

func Sum[T Numeric](numbers []T) T { ... }

// ❌ BAD: Overly broad constraints
func Sum[T any](numbers []T) T {
    // Cannot use + operator with 'any'
}

// ✅ GOOD: Compose constraints
type ComparableNumeric interface {
    Numeric
    comparable
}

func Max[T ComparableNumeric](a, b T) T { ... }
```

### 3. **Error Handling**

```go
// ✅ GOOD: Generic error handling
type Result[T any] struct {
    Value T
    Error error
}

func (r Result[T]) IsOk() bool {
    return r.Error == nil
}

func (r Result[T]) Unwrap() (T, error) {
    return r.Value, r.Error
}

// ✅ GOOD: Option type
type Option[T any] struct {
    value *T
}

func Some[T any](value T) Option[T] {
    return Option[T]{value: &value}
}

func None[T any]() Option[T] {
    return Option[T]{value: nil}
}

func (o Option[T]) IsSome() bool {
    return o.value != nil
}

func (o Option[T]) Unwrap() T {
    if o.value == nil {
        panic("called Unwrap on None value")
    }
    return *o.value
}

func (o Option[T]) UnwrapOr(defaultValue T) T {
    if o.value == nil {
        return defaultValue
    }
    return *o.value
}
```

### 4. **Performance Considerations**

```go
// ✅ GOOD: Pre-allocate slices when size is known
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice)) // Pre-allocate
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// ❌ BAD: Growing slice unnecessarily
func Map[T, U any](slice []T, fn func(T) U) []U {
    var result []U // Will grow multiple times
    for _, v := range slice {
        result = append(result, fn(v))
    }
    return result
}

// ✅ GOOD: Use type approximation for performance
type FastNumeric interface {
    ~int | ~int64 | ~float64 // Limit to common types
}

// ❌ BAD: Too many type variants
type SlowNumeric interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
    ~float32 | ~float64
}
```

### 5. **API Design**

```go
// ✅ GOOD: Consistent API
type Collection[T any] interface {
    Add(T)
    Remove(T) bool
    Contains(T) bool
    Size() int
    IsEmpty() bool
    Clear()
    ToSlice() []T
}

// ✅ GOOD: Builder pattern with generics
type QueryBuilder[T any] struct {
    table      string
    conditions []string
    values     []interface{}
}

func NewQueryBuilder[T any](table string) *QueryBuilder[T] {
    return &QueryBuilder[T]{table: table}
}

func (qb *QueryBuilder[T]) Where(condition string, value interface{}) *QueryBuilder[T] {
    qb.conditions = append(qb.conditions, condition)
    qb.values = append(qb.values, value)
    return qb
}

func (qb *QueryBuilder[T]) Build() (string, []interface{}) {
    query := fmt.Sprintf("SELECT * FROM %s", qb.table)
    if len(qb.conditions) > 0 {
        query += " WHERE " + strings.Join(qb.conditions, " AND ")
    }
    return query, qb.values
}
```

---

## Performance Considerations

### 1. **Compile Time Impact**

```go
// Generics có thể tăng compile time
// Mỗi type instantiation tạo ra code riêng biệt

// ✅ GOOD: Limit type parameters khi có thể
func ProcessInts(numbers []int) int { ... }        // Specific
func ProcessFloats(numbers []float64) float64 { ... } // Specific

// Thay vì:
func Process[T Numeric](numbers []T) T { ... } // Generic cho mọi numeric type

// ✅ GOOD: Sử dụng generics khi thực sự cần thiết
func Map[T, U any](slice []T, fn func(T) U) []U { ... } // Justified
```

### 2. **Runtime Performance**

```go
// Generics trong Go không có runtime overhead như Java
// Code được generate tại compile time

// Benchmark example
func BenchmarkGenericSum(b *testing.B) {
    numbers := make([]int, 1000)
    for i := range numbers {
        numbers[i] = i
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = Sum(numbers) // Generic function
    }
}

func BenchmarkSpecificSum(b *testing.B) {
    numbers := make([]int, 1000)
    for i := range numbers {
        numbers[i] = i
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = SumInt(numbers) // Specific function
    }
}

// Results thường tương đương nhau
```

### 3. **Memory Usage**

```go
// ✅ GOOD: Efficient generic collections
type RingBuffer[T any] struct {
    data  []T
    head  int
    tail  int
    size  int
    count int
}

func NewRingBuffer[T any](size int) *RingBuffer[T] {
    return &RingBuffer[T]{
        data: make([]T, size),
        size: size,
    }
}

func (rb *RingBuffer[T]) Push(item T) bool {
    if rb.count == rb.size {
        return false // Buffer full
    }
    
    rb.data[rb.tail] = item
    rb.tail = (rb.tail + 1) % rb.size
    rb.count++
    return true
}

func (rb *RingBuffer[T]) Pop() (T, bool) {
    if rb.count == 0 {
        var zero T
        return zero, false
    }
    
    item := rb.data[rb.head]
    rb.head = (rb.head + 1) % rb.size
    rb.count--
    return item, true
}
```

---

## Migration từ interface{}

### Trước Generics:

```go
// Old code using interface{}
type OldContainer struct {
    items []interface{}
}

func (c *OldContainer) Add(item interface{}) {
    c.items = append(c.items, item)
}

func (c *OldContainer) Get(index int) interface{} {
    if index < 0 || index >= len(c.items) {
        return nil
    }
    return c.items[index]
}

func (c *OldContainer) GetString(index int) (string, bool) {
    item := c.Get(index)
    if str, ok := item.(string); ok {
        return str, true
    }
    return "", false
}

func (c *OldContainer) GetInt(index int) (int, bool) {
    item := c.Get(index)
    if i, ok := item.(int); ok {
        return i, true
    }
    return 0, false
}
```

### Sau khi migrate sang Generics:

```go
// New code using generics
type Container[T any] struct {
    items []T
}

func NewContainer[T any]() *Container[T] {
    return &Container[T]{
        items: make([]T, 0),
    }
}

func (c *Container[T]) Add(item T) {
    c.items = append(c.items, item)
}

func (c *Container[T]) Get(index int) (T, bool) {
    if index < 0 || index >= len(c.items) {
        var zero T
        return zero, false
    }
    return c.items[index], true
}

func (c *Container[T]) Size() int {
    return len(c.items)
}

func (c *Container[T]) IsEmpty() bool {
    return len(c.items) == 0
}

// Usage
func main() {
    // Type-safe containers
    stringContainer := NewContainer[string]()
    stringContainer.Add("hello")
    stringContainer.Add("world")
    
    if str, ok := stringContainer.Get(0); ok {
        fmt.Println(str) // No type assertion needed!
    }
    
    intContainer := NewContainer[int]()
    intContainer.Add(42)
    intContainer.Add(100)
    
    if num, ok := intContainer.Get(1); ok {
        fmt.Println(num * 2) // Direct arithmetic operations!
    }
}
```

### Migration Strategy:

1. **Identify interface{} usage**
2. **Determine if generics would help**
3. **Create generic version alongside old version**
4. **Gradually migrate callers**
5. **Remove old version when safe**

```go
// Step 1: Create generic version
type NewAPI[T any] struct {
    // generic implementation
}

// Step 2: Keep old version for compatibility
type OldAPI struct {
    inner *NewAPI[interface{}]
}

func (old *OldAPI) OldMethod(item interface{}) {
    old.inner.NewMethod(item)
}

// Step 3: Provide migration helpers
func MigrateToGeneric[T any](old *OldAPI) *NewAPI[T] {
    // Migration logic
    return NewAPI[T]{}
}
```

---

## Kết luận

Generics trong Go 1.18+ mang lại nhiều lợi ích:

### ✅ **Ưu điểm:**
- **Type Safety**: Compile-time type checking
- **Performance**: Không có boxing/unboxing overhead
- **Code Reusability**: Một implementation cho nhiều types
- **Better APIs**: Cleaner, more expressive interfaces
- **Tooling Support**: Better IDE support và refactoring

### ⚠️ **Cần lưu ý:**
- **Learning Curve**: Syntax mới cần thời gian làm quen
- **Compile Time**: Có thể tăng thời gian compile
- **Complexity**: Có thể làm code phức tạp nếu overuse
- **Debugging**: Stack traces có thể dài hơn

### 🎯 **Khi nào nên sử dụng Generics:**
- Collections và data structures
- Utility functions (map, filter, reduce)
- Type-safe APIs
- Khi cần eliminate type assertions
- Khi có duplicate code cho different types

### 🚫 **Khi nào KHÔNG nên sử dụng:**
- Khi interface{} đã đủ tốt
- Cho simple, one-off functions
- Khi làm code phức tạp không cần thiết
- Khi performance không phải concern

Generics là công cụ mạnh mẽ, nhưng như mọi tool khác, cần sử dụng đúng chỗ đúng lúc để tối đa hóa lợi ích và tránh over-engineering.