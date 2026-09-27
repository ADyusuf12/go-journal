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
- Method Set rules (Pointer receivers *T vs Value types T for interfaces)
- Type assertions (val, ok := interface.(ConcreteType))
- Type switches (switch v := interface.(type))
- The empty interface (any / interface{})

**Key Insight:**
Interfaces give Go duck typing at compile-time. Values inside interfaces are two-word memory structures (itab pointer + data pointer). Always use the val, ok idiom for type assertions to prevent runtime panics.

---

# Current Position

**Current Phase:** Phase 1
**Current Week:** Week 1
**Current Day:** Day 5

### Next Topic

Concurrency (Goroutines, Channels, Select, and Synchronization)

### Future Topics

- PostgreSQL
- APIs
- Authentication
- Testing
- Messaging
- Kubernetes

---

# Mentor Assessment

**Current Go Level:** Advanced Beginner

### Strengths

- Backend engineering background
- Architecture experience
- PostgreSQL knowledge
- Business domain modelling
- Static typing familiarity

### Likely Challenges

- Thinking in Go instead of Rails
- Advanced Concurrency / Channels
- Distributed systems
- Cloud-native patterns

### Not Likely to Struggle With

- Variables
- Functions
- Basic typing
- Structs & Interfaces
- APIs
- Authentication concepts
