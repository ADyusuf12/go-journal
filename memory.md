# Go Journey Memory

## Student Profile

**Name:** Yusuf
**Location:** Nigeria

### Background

- Several years of production Ruby on Rails experience
- ERP systems experience
- Marketplace / job-platform experience
- Strong PostgreSQL experience
- API design experience
- Authentication and authorization experience
- Application architecture experience
- RSpec and Minitest experience

### Goal

Become a production-ready Go backend engineer over 6-12 months while continuing professional Rails work.

### Long-term Positioning

Backend Engineer

### Technologies

- Ruby
- Go
- PostgreSQL
- Docker
- Kubernetes
- AWS

### Focus

- Backend systems
- Distributed systems
- APIs
- Cloud-native engineering

### Avoid

- Tutorial projects
- Todo applications
- Beginner portfolio projects

### Leverage

- ERP domain knowledge
- Marketplace domain knowledge
- Business workflows
- Permissions systems
- Inventory systems
- Invoicing systems

## Teaching Style & Pedagogy

- **Iterative & Hands-On**: Teach concepts incrementally while writing code in small, testable chunks. Avoid long upfront theoretical lectures.
- **Practical Real-World Analogies**: Lead with intuitive, real-world analogies (e.g., kitchen order clips, warehouse inventories) before diving into technical syntax or low-level mechanics.
- **Immediate Feedback Loop & Output Deconstruction**: Run code after structural changes and dissect terminal output line-by-line, mapping every piece of output directly back to code mechanics.
- **Senior-to-Senior Peer Tone**: Relate Go concepts directly to production Ruby/Rails architectures, runtime mechanics, and memory behavior.
- **Embedded Check-Ins**: Use code execution and targeted runtime questions to verify understanding inside the coding process.

---

# Mentor Assessment

**Current Go Level:** Solid Intermediate Core Knowledge

### Strengths

- Backend engineering background
- Architecture experience
- PostgreSQL knowledge
- Business domain modelling
- Static typing familiarity
- Rapid mechanical grasp of Go concurrency primitives and channels

### Likely Challenges

- Thinking in Go instead of Rails
- Context cancellation propagation across microservices
- Distributed systems
- Cloud-native patterns

### Not Likely to Struggle With

- Variables & Functions
- Structs & Interfaces
- Basic & Advanced Concurrency (Goroutines, Channels, Select)
- APIs & Domain Architecture

---

# Master Roadmap

## Phase 0 — Environment Setup

**Status:** Completed

Completed:

- Go installed
- VS Code Go extension installed
- WSL2 verified
- PostgreSQL verified
- Rails environment preserved

---

## Phase 1 — Month 1: Core Go

**Status:** In Progress

### Week 1

Topics:

- Go organization
- Packages
- Modules
- Variables
- Types
- Functions
- Error handling

---

# Completed Lessons

## Day 1

**Topics Covered:**

- package main
- func main()
- imports
- modules
- go.mod
- go run .
- packages
- scope
- variable declaration
- reassignment
- type inference
- formatting verbs

**Key Insight:**
Go cares heavily about packages and modules. Filenames are relatively unimportant.

---

## Day 2

**Topics Covered:**

- functions
- parameters
- return values
- multiple return values
- errors.New()
- nil
- Go-style error handling

**Example Pattern:**

```go
result, err := someFunction()

if err != nil {
    return
}
```

**Key Insight:**
Go prefers explicit error handling over exception-driven control flow.

**Exercise Completed:**
Marketplace commission calculator with validation.

---

## Day 3

**Topics Covered:**

- Struct definition & initialization
- Value receivers vs. Pointer receivers
- Automatic pointer dereferencing & address-taking syntactic sugar
- Encapsulation & visibility (Exported vs. Unexported fields/methods)
- Factory/Constructor pattern (`New<StructName>`)
- Struct Embedding (Composition over Inheritance)
- Field & Method promotion

**Key Insight:**
Go enforces pass-by-value strictly. Pointer receivers operate on caller memory addresses, while value receivers operate on stack copies. Go replaces class inheritance with struct embedding composition.

---

## Day 4

**Topics Covered:**

- Interface declaration & implicit satisfaction
- Dependency Injection pattern (CheckoutService accepting PaymentProcessor)
- Method Set rules (Pointer receivers \*T vs Value types T for interfaces)
- Type assertions (val, ok := interface.(ConcreteType))
- Type switches (switch v := interface.(type))
- The empty interface (any / interface{})

**Key Insight:**
Interfaces give Go duck typing at compile-time. Values inside interfaces are two-word memory structures (itab pointer + data pointer). Always use the val, ok idiom for type assertions to prevent runtime panics.

---

## Day 5

**Topics Covered:**

- Goroutines lightweight thread scheduling (`go func()`)
- Synchronization with `sync.WaitGroup` (Add, Done, Wait)
- Pointer receiver mechanics for WaitGroups (`*sync.WaitGroup`)
- Typed Channels (`make(chan T)`) & Directional Channels (`<-chan T`)
- Unbuffered vs. Buffered Channels (`make(chan T, cap)`)
- Channel closing mechanics (`close()`) & `for range` channel iteration
- Multiplexing with `select` and timeout handling using `time.After()`
- Concurrent Worker Pool pattern with dynamic load balancing

**Key Insight:**
Channels are thread-safe FIFO queues. "Do not communicate by sharing memory; instead, share memory by communicating." Always manage channel lifecycles from the sender side to avoid panics and goroutine leaks.

---

## Day 6

**Topics Covered:**

- Slice headers (24-byte struct: ptr, len, cap) & array reallocation mechanics
- Shared underlying array mutations & copy() isolation
- Map memory structure, nil map panics, and the `delete()` built-in
- Key existence checking via the comma-ok idiom (`val, ok := map[key]`)
- Context package (`context.Context`), timeout cancellations (`WithTimeout`), and metadata propagation (`WithValue`)

**Key Insight:**
Slices are lightweight headers pointing to backing arrays; pre-allocating capacity avoids GC allocation churn. Maps return zero-values on missing keys, requiring the comma-ok idiom for safe lookups. Context acts as a request-scoped ticket carrying timers and metadata across boundaries.

---

## Day 7

**Topics Covered:**

- `io.Reader` bucket mechanics & manual buffer allocation (`make([]byte, N)`)
- `io.EOF` stream termination checking & slice re-slicing (`buffer[:n]`)
- Line-by-line text parsing via `bufio.Scanner`
- `io.Writer` interface, stream piping via `io.Copy`, and multiplexing via `io.MultiWriter`
- Custom `io.Writer` implementation (`RedactingWriter`) for stream data sanitization

**Key Insight:**
Files, network sockets, and HTTP request bodies all satisfy `io.Reader` and `io.Writer`. Piping streams with fixed memory buffers allows Go services to process gigabytes of data on a flat RAM footprint.

---

## Day 8

**Topics Covered:**

- `database/sql` connection pooling & `pgx/v5` driver registration
- Connection pool tuning (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`)
- Network health verification via `db.PingContext(ctx)`
- Parameterized SQL execution (`ExecContext`, `QueryRowContext`, `QueryContext`)
- Safe pointer-based scanning with `rows.Scan()` & `defer rows.Close()` resource management
- ACID transaction workflows with `db.BeginTx`, `defer tx.Rollback()`, and `FOR UPDATE` row locking

**Key Insight:**
`sql.DB` is a concurrency-safe connection pool, not a single connection. Executing queries with `*Context` variants guarantees that database operations respect timeout deadlines, while `rows.Scan()` parses binary row streams straight into Go memory pointers.

---

## Day 9

**Topics Covered:**

- `net/http` server mechanics, `http.Server` timeouts, and goroutine-per-request scheduling
- `http.NewServeMux` HTTP method routing (`GET`, `POST`)
- `http.ResponseWriter` (`io.Writer`) and `r.Body` (`io.Reader`) stream processing
- Zero-intermediate JSON streaming using `json.NewEncoder` and `json.NewDecoder`
- Defensive payload bounds via `http.MaxBytesReader` and strict schema enforcement with `DisallowUnknownFields`
- Higher-order middleware chaining (`http.Handler` wrapper pattern) for execution timing and request logging

**Key Insight:**
A Go web server is a high-concurrency TCP listener where `r.Body` and `ResponseWriter` are streaming I/O interfaces. Wrapping `ServeMux` with middleware functions creates modular execution pipelines without reflection magic or heavy framework abstractions.

---

## Day 10 (Phase 1 Capstone)

**Topics Covered:**

- End-to-end service integration (`net/http` + `pgx` + `database/sql`)
- Application dependency injection patterns via struct receivers (`*Application`)
- Context propagation from HTTP request stream (`r.Context()`) down to database socket layers
- Concurrent transactional race-condition mitigation using `FOR UPDATE` row locks
- OS signal trapping (`SIGINT`/`SIGTERM`) with buffered signal channels
- Non-disruptive server lifecycle termination using `server.Shutdown(ctx)`

**Key Insight:**
A production Go backend service is composed of modular standard library building blocks. Running the HTTP listener in a separate goroutine while monitoring OS signals on the main thread guarantees zero dropped requests or dangling database locks during deployments.

---

# Current Position

**Current Phase:** Phase 2 (Concurrency & Production Systems)
**Current Week:** Week 3
**Current Day:** Day 11

### Next Topic

Goroutines & Channels Deep Dive: Go Scheduler Mechanics (GMP), Memory Stack vs Heap, and CSP Concurrency Patterns

### Future Topics

- HTTP APIs & Middleware
- Authentication & JWTs
- Unit & Integration Testing
- Event Messaging
- Kubernetes & Cloud Deployment
