# Day 6 Notes: Slices, Maps, and Context

## 1. Slice Mechanics & Memory Headers

### The 24-Byte Header

In 64-bit Go, a slice variable is NOT the entire array—it's a lightweight **24-byte header** containing:

1. **Data Pointer** (`unsafe.Pointer`): Address of the first element in the backing array.
2. **Length** (`len`): Number of elements currently accessible.
3. **Capacity** (`cap`): Total contiguous elements allocated before reallocation is required.

### Dynamic Growth (`append`)

- **`len < cap`**: Appending simply writes to index `len` and increments `len` by 1. Zero memory reallocation occurs.
- **`len == cap`**: The backing array is full. Go allocates a new, larger backing array (typically doubling capacity), copies all existing elements byte-for-byte, updates the header pointer, and drops the old array for Garbage Collection (GC).

> **Backend Perf Tip:** Unplanned reallocations under high request loads create GC pressure. Always pre-allocate capacity when the size is known: `make([]T, 0, expectedCount)`.

### Shared Backing Arrays & Sub-slicing

Taking a sub-slice (`slice[low:high]`) creates a new 24-byte header pointing to the **same backing array**.

- Modifying `subSlice[0]` mutates the original parent slice!
- To isolate data safely without mutating the parent, use `copy()`:

```go
dst := make([]T, len(src[low:high]))
copy(dst, src[low:high])
```

## 2. Maps & Safe Lookup Patterns

### The hmap Pointer & Initialization

A Go map is a pointer to an hmap struct holding dynamic hash table buckets.

### nil Map Trap

Reading a nil map is safe (returns zero-values), but writing to a nil map causes a runtime panic.

Always initialize maps dynamically with `make(map[K]V, optionalCap)` or directly with literal syntax.

### Zero-Values vs. Missing Keys (The Comma-Ok Idiom)

Unlike Ruby hashes (which return nil on missing keys), Go map lookups return the key type's zero-value when a key is missing.

To distinguish between a key that actually stores 0/`""`/false versus a missing key, use the comma-ok idiom:

```go
val, ok := myMap["key_name"]
if !ok {
    // Key does NOT exist in the map
}
```

### Map Deletion

Use the built-in `delete(map, key)` function. Deleting a non-existent key is a safe no-op.

## 3. The context.Context Package

### Practical Analogy: The Restaurant Order Ticket

Think of context.Context as an order ticket attached to a meal request. As it travels down the pipeline (handler -> service -> db), it carries:

- Timer/SLA Boundaries ("Cancel this if it takes longer than 2 seconds")
- Customer Metadata (Request ID, User ID)

### Timeouts & Cancellation

In microservices, downstream database or API calls shouldn't hang indefinitely.

```go
// Create a context budget of 200ms
ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
defer cancel() // Always cleanup context resources on function exit!

// Downstream worker selects on ctx.Done()
select {
case <-time.After(500 * time.Millisecond):
    // Finished successfully
case <-ctx.Done():
    // Aborted! ctx.Err() returns context.DeadlineExceeded
}
```

### Request-Scoped Values (WithValue)

Pass request metadata (e.g., Auth Tokens, Trace IDs) safely across layer boundaries without polluting function signatures.

**Collision Prevention:** Always use custom unexported types for context keys instead of raw strings.

```go
type contextKey string

const userIDKey contextKey = "userID"

// Set value
ctx := context.WithValue(parentCtx, userIDKey, "USR-1001")

// Read value with type assertion
if userID, ok := ctx.Value(userIDKey).(string); ok {
    // Use userID safely
}
```
