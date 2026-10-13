package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Define a custom type for context keys to avoid collisions
type contextKey string

const requestIDKey contextKey = "requestID"

// ResponseWriter wrapper to capture HTTP status codes and response sizes for logging
type statusResponseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

// ---------------------------------------------------------
// Middleware 1: Request ID Injection
// ---------------------------------------------------------
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = "req_sys_998124" // In production, generate a UUID here
		}

		// Inject Request ID into context
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ---------------------------------------------------------
// Middleware 2: Structured Request Logging (`log/slog`)
// ---------------------------------------------------------
func LoggingMiddleware(next http.Handler) http.Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID, _ := r.Context().Value(requestIDKey).(string)

		wrappedWriter := &statusResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK, // default if WriteHeader isn't called explicitly
		}

		// Execute downstream handler
		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)

		// Log request telemetry in structured JSON format
		logger.Info("incoming http request",
			slog.String("request_id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrappedWriter.status),
			slog.Int("bytes_written", wrappedWriter.bytes),
			slog.Duration("duration", duration),
		)
	})
}

// ---------------------------------------------------------
// Core Application Handler
// ---------------------------------------------------------
func handleGetBalance(w http.ResponseWriter, r *http.Request) {
	reqID, _ := r.Context().Value(requestIDKey).(string)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","balance":1540.50,"request_id":"` + reqID + `"}`))
}

func main() {
	mux := http.NewServeMux()

	// Register route
	mux.HandleFunc("GET /api/balance", handleGetBalance)

	// Chain middlewares (Outer middleware runs first)
	var handler http.Handler = mux
	handler = LoggingMiddleware(handler)
	handler = RequestIDMiddleware(handler)

	slog.Info("Starting production server on port :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		slog.Error("server failed", slog.String("error", err.Error()))
	}
}
