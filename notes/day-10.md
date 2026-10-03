# Day 10 Notes: Phase 1 Capstone Integration & Graceful Shutdown

## 1. Production Architecture Patterns

- **Resource Ownership via App Struct**: Inject dependencies (`*sql.DB`) into an `Application` struct to share connection pools safely across handlers without global state.
- **Context-Bound Database Operations**: Always pass `r.Context()` to `QueryContext`, `ExecContext`, and `BeginTx`. If an HTTP client drops connection mid-flight, database operations abort immediately.

---

## 2. Safe State Mutations

- Combine **`db.BeginTx`**, **`SELECT ... FOR UPDATE`**, and **`defer tx.Rollback()`** for inventory/financial mutations to guarantee thread safety against race conditions under concurrent load.
- Guard against payload memory attacks with `http.MaxBytesReader`.

---

## 3. OS Signal Interception & Graceful Shutdown

```go
// Listen for termination signals from OS, Docker, or Kubernetes
shutdownChan := make(chan os.Signal, 1)
signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

// Run HTTP server in a background Goroutine
go func() {
    if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
        log.Fatalf("Server error: %v", err)
    }
}()

<-shutdownChan // Block main until signal arrives

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Complete in-flight requests before exiting
server.Shutdown(ctx)

```

---
