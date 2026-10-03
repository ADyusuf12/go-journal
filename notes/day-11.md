# Day 11 Notes: Go Scheduler (GMP) & Compiler Escape Analysis

## 1. The GMP Scheduler Model

- **`G` (Goroutine)**: Lightweight user-space execution thread (~2KB starting stack). Managed dynamically by the Go runtime.
- **`M` (Machine)**: Physical OS kernel thread created and managed by the operating system kernel (~1MB-2MB).
- **`P` (Processor)**: Logical execution context representing resources required to run Go code (`runtime.GOMAXPROCS`). An $M$ must hold a $P$ to execute a $G$.

### Runtime Algorithms

- **Work Stealing**: When a $P$'s local run queue empties, its assigned $M$ steals half the runnable goroutines from another $P$'s run queue to keep CPU cores saturated.
- **Syscall Detachment (Netpoller vs Blocking Syscalls)**:
  - Network I/O is parked on the asynchronous **Netpoller** (`epoll`/`kqueue`) without blocking $M$.
  - Blocking OS syscalls (e.g., file reads) detach $P$ from $M$ so another thread can continue executing remaining $G$'s.

---

## 2. Memory Allocation & Escape Analysis

- **Stack Allocation**: Fast pointer bump on a 2KB contiguous, growing stack frame. Automatically discarded on function exit (zero GC overhead).
- **Heap Allocation**: Managed by the Garbage Collector. Required when variables outlive function execution frames or when interface dynamic bounds obscure lifetime analysis.
- **Escape Analysis Compiler Flag**: Inspect compiler decisions at build time:

  ```bash
  go build -gcflags="-m" .

  ```

- **Escape Triggers**: Returning pointers (\*T), channel sends, and interface boxing (fmt.Println(any...)).

---
