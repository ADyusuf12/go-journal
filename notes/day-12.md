# Day 12 Notes: Data Races, Mutexes, and Channel Internals

## 1. What is a Data Race?

- **Concept**: A data race happens when two or more goroutines access the exact same memory location at the same time, and at least one of those accesses is a write.
- **The Problem**: Operations like `count++` or `balance += amount` are **not atomic**. They break down into three separate CPU steps:
  1. Read value from memory into a CPU register.
  2. Modify the value in the register.
  3. Write the value back to memory.
- Without synchronization, concurrent read/write steps overlap, causing silent data corruption.

---

## 2. Shared Memory Synchronization (`sync.Mutex`)

- **`sync.Mutex`**: A mutual exclusion lock used to protect shared low-level state (counters, maps, in-memory caches).
- **Rule of Operation**:
  - `mu.Lock()`: Acquire exclusive access before reading or writing shared memory. Other goroutines block until released.
  - `defer mu.Unlock()`: Guarantee the lock releases when the function finishes, preventing deadlocks.

---

## 3. The Go Race Detector (`go run -race`)

- **Compilation Flag**: `go run -race .` or `go test -race ./...`
- **Mechanics**: Instruments memory access instructions at compile-time (via ThreadSanitizer).
- **Behavior**: Detects concurrent un-synchronized reads/writes at runtime and exits with code `66`, failing builds/tests instantly.

---

## 4. Channels vs. Mutexes (CSP Philosophy)

- **Shared Memory (`Mutex`)**: Goroutines share a variable in RAM and take turns holding a lock to touch it.
- **Message Passing (`Channels`)**: Data ownership is passed down a conveyor belt from sender to receiver.
- **Channel Internals (`runtime.hchan`)**:
  - A heap-allocated struct with a circular ring buffer protected by an internal lock.
  - Goroutines waiting on full or empty channels are put to sleep (`gopark`) on waiting lists (`sendq`/`recvq`) without wasting CPU cycles.
