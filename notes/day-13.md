# Day 13 Notes: Advanced CSP Patterns & Concurrency Pipelines

## 1. Pipeline Architecture Patterns

- **Generator Pattern**: A goroutine that produces values onto a send-only output channel (`<-chan T`) and closes it when complete.
- **Fan-Out (Worker Pool)**: Spawning multiple concurrent worker goroutines reading from a single input channel to distribute heavy processing loads across CPU cores.
- **Fan-In (Multiplexer)**: Merging multiple worker output channels into a single combined output channel using `sync.WaitGroup` and goroutines.

---

## 2. Cancellation Chains (`context.Context`)

- Always select on `<-ctx.Done()` inside channel send and receive operations:

  ```go
  select {
  case <-ctx.Done():
      return // Abort execution cleanly on cancellation/timeout
  case out <- data:
  }

  ```

- Cascades cancellation signals instantly across nested pipeline stages, preventing go routine leaks when processing is aborted

---
