# Day 7 Notes: File I/O, Readers, Writers, and Data Streams

## 1. Core Streaming Interfaces

Go standardizes data processing around two minimal interfaces in the `io` package:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

- **`io.Reader`** (Bucket): Reads available bytes into caller-provided slice `p`. Returns `io.EOF` when the stream ends.
- **`io.Writer`** (Funnel): Flushes bytes from slice `p` to the underlying target (file, stdout, network connection).

## 2. Low-Level Buffering vs High-Level Scanners

- **Manual Buffer Reads**: Using `file.Read(buffer)` gives low-level control. Always slice the buffer using `buffer[:n]` when reading valid bytes to avoid reading stale data from previous iterations.
- **`bufio.Scanner`**: High-level wrapper for line-by-line parsing (`scanner.Scan()` / `scanner.Text()`).
- **`io.ReadAll`**: Loads an entire stream into memory at once. Use only for small payloads (e.g., small JSON API requests) to avoid RAM spikes.

## 3. Data Piping & Composition

- **`io.Copy(dstWriter, srcReader)`**: Streams data directly from reader to writer using an internal 32KB buffer until `io.EOF`.
- **`io.MultiWriter(w1, w2...)`**: Duplicates incoming writes across multiple writers simultaneously (e.g., writing to disk and `os.Stdout` at the same time).
- **Custom `io.Writer`**: Implementing `Write([]byte) (int, error)` on custom structs enables middleware patterns like log redacting, rate limiting, or payload hashing transparently.
