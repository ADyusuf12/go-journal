# Day 15 Notes: Custom Middleware Chains & Structured Logging (`log/slog`)

## 1. Middleware Architecture

Go middlewares are pure functions wrapping an `http.Handler` and returning an `http.Handler`:

```go
func MyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Pre-processing
        next.ServeHTTP(w, r)
        // Post-processing
    })
}
```

This replaces object-oriented Rack middleware arrays (Rails) with explicit, zero-overhead functional composition.

---

## 2. Structured Logging (`log/slog`)

Introduced in Go 1.21 as part of the standard library.

- Outputs machine-readable key-value pairs (JSON or text) natively, optimizing high-throughput ingestion for observability tools.
- Wrapped `http.ResponseWriter` patterns enable interception of status codes and payload sizes post-execution for comprehensive request logging.
