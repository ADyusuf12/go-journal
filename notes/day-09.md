# Day 9 Notes: HTTP Standard Library (`net/http`) & ServeMux

## 1. Web Server Mechanics

- Go binaries are standalone HTTP servers (`server.ListenAndServe()`). No external Rack or app server like Puma is required.
- The server spawns a dedicated goroutine for every incoming TCP socket connection (~2KB memory per request).

---

## 2. Handlers & Streaming JSON

- **`http.Handler` Interface**: `ServeHTTP(w http.ResponseWriter, r *http.Request)`.
- **`http.ResponseWriter` (`io.Writer`)**: Streams HTTP headers and body bytes directly back down the open client TCP socket.
- **`*http.Request` (`r.Body` -> `io.Reader`)**: Incoming stream socket for client payloads.
- **Streaming Encoders**:
  - `json.NewDecoder(r.Body).Decode(&dst)` — Stream JSON from body to struct memory pointer.
  - `json.NewEncoder(w).Encode(src)` — Stream JSON from struct memory pointer to response socket.
- **DoS Protection**: Use `http.MaxBytesReader(w, r.Body, maxBytes)` to limit request body size.

---

## 3. Idiomatic Middleware Pattern

Middleware is a higher-order function wrapping `http.Handler`:

```go
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Pre-processing (auth, request tracking, logging start)
        next.ServeHTTP(w, r)
        // Post-processing (response duration, cleanup)
    })
}
```

**Context Cancellation**: `r.Context()` cancels when the client disconnects, allowing database and upstream operations to abort immediately and save CPU/RAM.
