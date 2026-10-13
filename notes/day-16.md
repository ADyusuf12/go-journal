# Day 16 Notes: Advanced Routing, Method Matching, and Strict JSON Payload Validation

## 1. Advanced Routing (`net/http` Method Routing)

Go 1.22+ standard library `ServeMux` natively supports HTTP method matching in route patterns:

```go
mux.HandleFunc("POST /api/orders", handleCreateOrder)
```

This eliminates the need for external routing packages (like Gorilla Mux or Chi) for standard RESTful resource endpoints.

---

## 2. Strict JSON Validation & Security

### Payload Bomb Protection
`http.MaxBytesReader` caps incoming request bodies (e.g., to 1MB) to prevent memory exhaustion/DoS attacks.

### Strict Schema Enforcement
`decoder.DisallowUnknownFields()` forces a 400 Bad Request when clients send unauthorized or unexpected JSON keys, preventing silent mass-assignment bugs common in dynamic frameworks.

### Business Rule Validation
Implementing custom `.Validate()` methods on structs keeps business constraints clean and decoupled from HTTP handling layers.
