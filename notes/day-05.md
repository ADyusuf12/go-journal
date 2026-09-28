# Day 5: Concurrency, Goroutines, Channels & Select

## Key Concepts Covered

1. **Goroutines**:
   - Spawns lightweight green threads managed by the Go runtime using the `go` keyword.
   - The `main` Goroutine exits immediately upon function completion unless synchronized.

2. **Synchronization (`sync.WaitGroup`)**:
   - `Add(delta)` increments job count.
   - `Done()` (typically deferred) decrements count.
   - `Wait()` blocks until count reaches 0.
   - Must be passed by pointer (`*sync.WaitGroup`) to avoid copying value state.

3. **Channels (`chan T`)**:
   - Thread-safe FIFO data pipelines.
   - **Unbuffered (`make(chan T)`)**: Requires sender and receiver to be ready simultaneously (synchronous handoff).
   - **Buffered (`make(chan T, capacity)`)**: Holds items without blocking until capacity is full.

4. **Channel Lifecycle & Ranging**:
   - `close(ch)` notifies receivers no further data will be sent. Always close from the sender side.
   - `for item := range ch` drains values until closed.

5. **Multiplexing with `select` & Timeouts**:
   - `select` evaluates multiple channel operations simultaneously, executing whichever is ready first.
   - `time.After(duration)` provides clean timeout handling inside `select`.

---

## Code Example Summary (`concurrency/main.go`)
