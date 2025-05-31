# Go Programming Guide - Từ PHP đến Go

## Mục lục

1. [Giới thiệu về Go](#giới-thiệu-về-go)
2. [Cài đặt và thiết lập môi trường](#cài-đặt-và-thiết-lập-môi-trường)
3. [Cú pháp cơ bản](#cú-pháp-cơ-bản)
   - [Variables và Types](#variables-và-types)
   - [Constants](#constants)
   - [Operators](#operators)
4. [Cấu trúc dữ liệu](#cấu-trúc-dữ-liệu)
   - [Arrays](#arrays)
   - [Slices](#slices)
   - [Maps](#maps)
   - [Structs](#structs)
5. [Cấu trúc điều khiển](#cấu-trúc-điều-khiển)
   - [If/Else](#ifelse)
   - [Switch](#switch)
   - [Loops](#loops)
6. [Functions và Methods](#functions-và-methods)
   - [Function Declaration](#function-declaration)
   - [Methods](#methods)
   - [Interfaces](#interfaces)
7. [Error Handling](#error-handling)
8. [Concurrency](#concurrency)
   - [Goroutines](#goroutines)
   - [Channels](#channels)
   - [Select Statement](#select-statement)
9. [Package Management](#package-management)
10. [Testing](#testing)
11. [Web Development](#web-development)
12. [Database Operations](#database-operations)
13. [So sánh PHP vs Go](#so-sánh-php-vs-go)

---

## Giới thiệu về Go

Go (hay Golang) là ngôn ngữ lập trình được phát triển bởi Google vào năm 2009. Go được thiết kế để đơn giản, hiệu quả và dễ học, đặc biệt phù hợp cho việc phát triển các ứng dụng web, microservices và hệ thống phân tán.

### Tại sao chuyển từ PHP sang Go?

- **Performance**: Go nhanh hơn PHP đáng kể
- **Concurrency**: Hỗ trợ concurrency tự nhiên với goroutines
- **Static typing**: Giúp phát hiện lỗi sớm hơn
- **Compilation**: Biên dịch thành binary, dễ deploy
- **Memory management**: Garbage collection tự động

---

## Cài đặt và thiết lập môi trường

### Cài đặt Go

```bash
# macOS (sử dụng Homebrew)
brew install go

# Kiểm tra version
go version
```

### Thiết lập GOPATH và workspace

```bash
# Thêm vào ~/.bashrc hoặc ~/.zshrc
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

### Tạo project đầu tiên

```bash
mkdir hello-world
cd hello-world
go mod init hello-world
```

---

## Cú pháp cơ bản

### Variables và Types

#### So sánh với PHP:

**PHP:**
```php
<?php
$name = "John";
$age = 30;
$price = 99.99;
$isActive = true;
?>
```

**Go:**
```go
package main

import "fmt"

func main() {
    // Khai báo với var
    var name string = "John"
    var age int = 30
    var price float64 = 99.99
    var isActive bool = true
    
    // Khai báo ngắn gọn
    name2 := "Jane"
    age2 := 25
    
    fmt.Println(name, age, price, isActive)
    fmt.Println(name2, age2)
}
```

#### Các kiểu dữ liệu cơ bản:

```go
// Số nguyên
var i8 int8 = 127
var i16 int16 = 32767
var i32 int32 = 2147483647
var i64 int64 = 9223372036854775807
var ui uint = 42

// Số thực
var f32 float32 = 3.14
var f64 float64 = 3.141592653589793

// Chuỗi
var str string = "Hello, World!"

// Boolean
var flag bool = true

// Byte
var b byte = 255

// Rune (Unicode code point)
var r rune = 'A'
```

### Constants

**PHP:**
```php
define('PI', 3.14159);
const GRAVITY = 9.8;
```

**Go:**
```go
const PI = 3.14159
const GRAVITY = 9.8

// Nhóm constants
const (
    StatusOK = 200
    StatusNotFound = 404
    StatusInternalServerError = 500
)
```

### Operators

```go
// Toán học
a := 10
b := 3
fmt.Println(a + b)  // 13
fmt.Println(a - b)  // 7
fmt.Println(a * b)  // 30
fmt.Println(a / b)  // 3
fmt.Println(a % b)  // 1

// So sánh
fmt.Println(a == b)  // false
fmt.Println(a != b)  // true
fmt.Println(a > b)   // true
fmt.Println(a < b)   // false

// Logic
x := true
y := false
fmt.Println(x && y)  // false
fmt.Println(x || y)  // true
fmt.Println(!x)      // false
```

---

## Cấu trúc dữ liệu

### Arrays

**PHP:**
```php
$numbers = [1, 2, 3, 4, 5];
$fruits = array("apple", "banana", "orange");
```

**Go:**
```go
// Array có kích thước cố định
var numbers [5]int = [5]int{1, 2, 3, 4, 5}
fruits := [3]string{"apple", "banana", "orange"}

// Tự động xác định kích thước
auto := [...]int{1, 2, 3, 4}

fmt.Println(numbers[0])  // 1
fmt.Println(len(fruits)) // 3
```

### Slices

Slices trong Go tương tự như arrays trong PHP:

```go
// Tạo slice
numbers := []int{1, 2, 3, 4, 5}
fruits := []string{"apple", "banana", "orange"}

// Thêm phần tử
numbers = append(numbers, 6)
fruits = append(fruits, "grape")

// Slice của slice
subNumbers := numbers[1:4]  // [2, 3, 4]

// Tạo slice với make
slice := make([]int, 5)     // [0, 0, 0, 0, 0]
slice2 := make([]int, 3, 5) // length=3, capacity=5

fmt.Println(len(numbers))   // 6
fmt.Println(cap(numbers))   // capacity
```

### Maps

**PHP:**
```php
$person = [
    "name" => "John",
    "age" => 30,
    "city" => "New York"
];

echo $person["name"];
```

**Go:**
```go
// Tạo map
person := map[string]interface{}{
    "name": "John",
    "age":  30,
    "city": "New York",
}

// Hoặc với kiểu cụ thể
ages := map[string]int{
    "John":  30,
    "Jane":  25,
    "Bob":   35,
}

// Truy cập
fmt.Println(person["name"])  // John
fmt.Println(ages["Jane"])    // 25

// Kiểm tra key tồn tại
age, exists := ages["Alice"]
if exists {
    fmt.Println("Alice's age:", age)
} else {
    fmt.Println("Alice not found")
}

// Thêm/sửa
ages["Alice"] = 28

// Xóa
delete(ages, "Bob")

// Tạo map rỗng
emptyMap := make(map[string]int)
```

### Structs

Struct trong Go tương tự như class trong PHP:

**PHP:**
```php
class Person {
    public $name;
    public $age;
    public $email;
    
    public function __construct($name, $age, $email) {
        $this->name = $name;
        $this->age = $age;
        $this->email = $email;
    }
    
    public function introduce() {
        return "Hi, I'm " . $this->name;
    }
}

$person = new Person("John", 30, "john@example.com");
echo $person->introduce();
```

**Go:**
```go
type Person struct {
    Name  string
    Age   int
    Email string
}

// Constructor function
func NewPerson(name string, age int, email string) *Person {
    return &Person{
        Name:  name,
        Age:   age,
        Email: email,
    }
}

// Method
func (p Person) Introduce() string {
    return fmt.Sprintf("Hi, I'm %s", p.Name)
}

// Method với pointer receiver
func (p *Person) SetAge(age int) {
    p.Age = age
}

func main() {
    // Tạo struct
    person1 := Person{"John", 30, "john@example.com"}
    person2 := Person{
        Name:  "Jane",
        Age:   25,
        Email: "jane@example.com",
    }
    person3 := NewPerson("Bob", 35, "bob@example.com")
    
    fmt.Println(person1.Introduce())
    person3.SetAge(36)
    fmt.Println(person3.Age)
}
```

---

## Cấu trúc điều khiển

### If/Else

**PHP:**
```php
$age = 18;

if ($age >= 18) {
    echo "Adult";
} elseif ($age >= 13) {
    echo "Teenager";
} else {
    echo "Child";
}
```

**Go:**
```go
age := 18

if age >= 18 {
    fmt.Println("Adult")
} else if age >= 13 {
    fmt.Println("Teenager")
} else {
    fmt.Println("Child")
}

// If với statement
if num := 10; num%2 == 0 {
    fmt.Println("Even number")
}
```

### Switch

**PHP:**
```php
$day = "Monday";

switch ($day) {
    case "Monday":
        echo "Start of work week";
        break;
    case "Friday":
        echo "TGIF!";
        break;
    default:
        echo "Regular day";
}
```

**Go:**
```go
day := "Monday"

switch day {
case "Monday":
    fmt.Println("Start of work week")
case "Friday":
    fmt.Println("TGIF!")
case "Saturday", "Sunday":
    fmt.Println("Weekend!")
default:
    fmt.Println("Regular day")
}

// Switch không cần expression
num := 15
switch {
case num < 10:
    fmt.Println("Single digit")
case num < 100:
    fmt.Println("Double digit")
default:
    fmt.Println("Large number")
}
```

### Loops

**PHP:**
```php
// For loop
for ($i = 0; $i < 5; $i++) {
    echo $i . "\n";
}

// Foreach
$fruits = ["apple", "banana", "orange"];
foreach ($fruits as $fruit) {
    echo $fruit . "\n";
}

foreach ($fruits as $index => $fruit) {
    echo $index . ": " . $fruit . "\n";
}

// While
$i = 0;
while ($i < 5) {
    echo $i . "\n";
    $i++;
}
```

**Go:**
```go
// For loop cơ bản
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// For như while
i := 0
for i < 5 {
    fmt.Println(i)
    i++
}

// Vòng lặp vô hạn
for {
    // break để thoát
    break
}

// Range loop (như foreach)
fruits := []string{"apple", "banana", "orange"}

// Chỉ value
for _, fruit := range fruits {
    fmt.Println(fruit)
}

// Index và value
for index, fruit := range fruits {
    fmt.Printf("%d: %s\n", index, fruit)
}

// Range với map
ages := map[string]int{"John": 30, "Jane": 25}
for name, age := range ages {
    fmt.Printf("%s is %d years old\n", name, age)
}
```

---

## Functions và Methods

### Function Declaration

**PHP:**
```php
function add($a, $b) {
    return $a + $b;
}

function greet($name, $greeting = "Hello") {
    return $greeting . ", " . $name . "!";
}

$result = add(5, 3);
echo greet("John");
```

**Go:**
```go
// Function cơ bản
func add(a, b int) int {
    return a + b
}

// Multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// Named return values
func calculate(a, b int) (sum, product int) {
    sum = a + b
    product = a * b
    return // naked return
}

// Variadic function
func sum(numbers ...int) int {
    total := 0
    for _, num := range numbers {
        total += num
    }
    return total
}

func main() {
    result := add(5, 3)
    fmt.Println(result) // 8
    
    quotient, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", quotient)
    }
    
    s, p := calculate(4, 5)
    fmt.Println("Sum:", s, "Product:", p)
    
    total := sum(1, 2, 3, 4, 5)
    fmt.Println("Total:", total)
}
```

### Methods

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// Method với value receiver
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Method với pointer receiver
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    fmt.Println("Area:", rect.Area())
    
    rect.Scale(2)
    fmt.Println("New area:", rect.Area())
}
```

### Interfaces

Interface trong Go tương tự như interface trong PHP nhưng được implement ngầm định:

```go
// Định nghĩa interface
type Shape interface {
    Area() float64
    Perimeter() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * 3.14159 * c.Radius
}

type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Function nhận interface
func printShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
    circle := Circle{Radius: 5}
    rectangle := Rectangle{Width: 10, Height: 5}
    
    printShapeInfo(circle)
    printShapeInfo(rectangle)
}
```

---

## Error Handling

Go không có exceptions như PHP, thay vào đó sử dụng error values:

**PHP:**
```php
try {
    $result = riskyOperation();
    echo $result;
} catch (Exception $e) {
    echo "Error: " . $e->getMessage();
}
```

**Go:**
```go
func riskyOperation() (string, error) {
    // Simulate an error
    return "", fmt.Errorf("something went wrong")
}

func main() {
    result, err := riskyOperation()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}

// Custom error type
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
}

func validateAge(age int) error {
    if age < 0 {
        return ValidationError{
            Field:   "age",
            Message: "age cannot be negative",
        }
    }
    return nil
}
```

---

## Concurrency

### Goroutines

Goroutines là lightweight threads trong Go:

```go
func sayHello(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("Hello, %s! (%d)\n", name, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Chạy đồng bộ
    sayHello("Alice")
    
    // Chạy bất đồng bộ với goroutine
    go sayHello("Bob")
    go sayHello("Charlie")
    
    // Đợi goroutines hoàn thành
    time.Sleep(1 * time.Second)
    fmt.Println("Done!")
}
```

### Channels

Channels được sử dụng để giao tiếp giữa các goroutines:

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

### Select Statement

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
        case <-time.After(3 * time.Second):
            fmt.Println("Timeout!")
        }
    }
}
```

---

## Package Management

### Go Modules

```bash
# Tạo module mới
go mod init myproject

# Thêm dependency
go get github.com/gin-gonic/gin

# Cập nhật dependencies
go mod tidy

# Xem dependencies
go list -m all
```

### Import packages

```go
package main

import (
    "fmt"                    // Standard library
    "net/http"              // Standard library
    
    "github.com/gin-gonic/gin" // External package
    
    "myproject/internal/user"  // Local package
)
```

---

## Testing

### Unit Testing

**math.go:**
```go
package math

func Add(a, b int) int {
    return a + b
}

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

**math_test.go:**
```go
package math

import (
    "testing"
)

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    
    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}

func TestDivide(t *testing.T) {
    // Test normal case
    result, err := Divide(10, 2)
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    if result != 5 {
        t.Errorf("Divide(10, 2) = %f; want 5", result)
    }
    
    // Test division by zero
    _, err = Divide(10, 0)
    if err == nil {
        t.Error("Expected error for division by zero")
    }
}

// Benchmark test
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
```

```bash
# Chạy tests
go test

# Chạy với verbose
go test -v

# Chạy benchmark
go test -bench=.

# Test coverage
go test -cover
```

---

## Web Development

### HTTP Server cơ bản

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

var users = []User{
    {ID: 1, Name: "John", Email: "john@example.com"},
    {ID: 2, Name: "Jane", Email: "jane@example.com"},
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
    // Extract ID from URL (simplified)
    // In real app, use a router like Gorilla Mux or Gin
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users[0])
}

func main() {
    http.HandleFunc("/users", getUsers)
    http.HandleFunc("/user", getUser)
    
    fmt.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Sử dụng Gin Framework

```go
package main

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var users = []User{
    {ID: 1, Name: "John", Email: "john@example.com"},
    {ID: 2, Name: "Jane", Email: "jane@example.com"},
}

func getUsers(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"data": users})
}

func getUser(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }
    
    for _, user := range users {
        if user.ID == id {
            c.JSON(http.StatusOK, gin.H{"data": user})
            return
        }
    }
    
    c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func createUser(c *gin.Context) {
    var newUser User
    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    newUser.ID = len(users) + 1
    users = append(users, newUser)
    
    c.JSON(http.StatusCreated, gin.H{"data": newUser})
}

func main() {
    r := gin.Default()
    
    api := r.Group("/api/v1")
    {
        api.GET("/users", getUsers)
        api.GET("/users/:id", getUser)
        api.POST("/users", createUser)
    }
    
    r.Run(":8080")
}
```

---

## Database Operations

### Sử dụng database/sql với MySQL

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    ID    int
    Name  string
    Email string
}

func main() {
    // Kết nối database
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Test connection
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
    
    // Create user
    createUser(db, "John Doe", "john@example.com")
    
    // Get user
    user, err := getUser(db, 1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User: %+v\n", user)
    
    // Get all users
    users, err := getAllUsers(db)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Users: %+v\n", users)
}

func createUser(db *sql.DB, name, email string) error {
    query := "INSERT INTO users (name, email) VALUES (?, ?)"
    _, err := db.Exec(query, name, email)
    return err
}

func getUser(db *sql.DB, id int) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE id = ?"
    row := db.QueryRow(query, id)
    
    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func getAllUsers(db *sql.DB) ([]User, error) {
    query := "SELECT id, name, email FROM users"
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Email)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, nil
}
```

### Sử dụng GORM (ORM)

```go
package main

import (
    "fmt"
    "log"
    
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"size:100;not null"`
    Email string `gorm:"size:100;uniqueIndex"`
}

func main() {
    // Kết nối database
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    // Auto migrate
    db.AutoMigrate(&User{})
    
    // Create
    user := User{Name: "John Doe", Email: "john@example.com"}
    result := db.Create(&user)
    if result.Error != nil {
        log.Fatal(result.Error)
    }
    fmt.Printf("Created user with ID: %d\n", user.ID)
    
    // Read
    var foundUser User
    db.First(&foundUser, user.ID)
    fmt.Printf("Found user: %+v\n", foundUser)
    
    // Update
    db.Model(&foundUser).Update("Name", "Jane Doe")
    
    // Delete
    db.Delete(&foundUser)
    
    // Find all
    var users []User
    db.Find(&users)
    fmt.Printf("All users: %+v\n", users)
}
```

---

## So sánh PHP vs Go

| Aspect | PHP | Go |
|--------|-----|----|
| **Typing** | Dynamic | Static |
| **Performance** | Interpreted | Compiled |
| **Concurrency** | Limited (threads/processes) | Built-in (goroutines) |
| **Memory Management** | Reference counting + GC | Garbage Collection |
| **Error Handling** | Exceptions | Error values |
| **OOP** | Classes, inheritance | Structs, composition |
| **Package Management** | Composer | Go modules |
| **Deployment** | Requires PHP runtime | Single binary |
| **Learning Curve** | Easy | Moderate |

### Ví dụ so sánh cụ thể:

**PHP - Class và Inheritance:**
```php
class Animal {
    protected $name;
    
    public function __construct($name) {
        $this->name = $name;
    }
    
    public function speak() {
        return "Some sound";
    }
}

class Dog extends Animal {
    public function speak() {
        return $this->name . " says Woof!";
    }
}

$dog = new Dog("Buddy");
echo $dog->speak();
```

**Go - Struct và Interface:**
```go
type Animal interface {
    Speak() string
}

type Dog struct {
    Name string
}

func (d Dog) Speak() string {
    return d.Name + " says Woof!"
}

func main() {
    dog := Dog{Name: "Buddy"}
    fmt.Println(dog.Speak())
    
    // Polymorphism
    var animal Animal = dog
    fmt.Println(animal.Speak())
}
```

---

## Kết luận

Go là một ngôn ngữ mạnh mẽ và hiệu quả, đặc biệt phù hợp cho:

- **Microservices và APIs**
- **Concurrent applications**
- **System programming**
- **Cloud-native applications**
- **DevOps tools**

### Lộ trình học Go cho PHP developers:

1. **Tuần 1-2**: Cú pháp cơ bản, types, functions
2. **Tuần 3-4**: Structs, interfaces, methods
3. **Tuần 5-6**: Error handling, testing
4. **Tuần 7-8**: Concurrency (goroutines, channels)
5. **Tuần 9-10**: Web development với Gin
6. **Tuần 11-12**: Database operations, deployment

### Resources hữu ích:

- [Go Tour](https://tour.golang.org/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Documentation](https://golang.org/doc/)

Chúc bạn học Go thành công! 🚀