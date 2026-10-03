package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Product struct {
	ID    int     `json:"id"`
	SKU   string  `json:"sku"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Custom Logging Middleware
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Pass control to the next handler in the chain
		next.ServeHTTP(w, r)

		// Code executed after handler finishes
		duration := time.Since(start)
		log.Printf("[HTTP] %s %s | Remote: %s | Duration: %v", r.Method, r.URL.Path, r.RemoteAddr, duration)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /products", getProductsHandler)
	mux.HandleFunc("POST /products", createProductHandler)

	// Wrap our ServeMux with the RequestLogger middleware
	loggedRouter := RequestLogger(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      loggedRouter, // Set the wrapped router as primary server handler
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	fmt.Println("Server starting on http://localhost:8080...")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK - System Operational"))
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	products := []Product{
		{ID: 1, SKU: "PRD-101", Name: "Mechanical Keyboard", Price: 85.50},
		{ID: 2, SKU: "PRD-102", Name: "Ergonomic Chair", Price: 299.99},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var newProd Product
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&newProd); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON payload: %v", err), http.StatusBadRequest)
		return
	}

	newProd.ID = 101

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProd)
}
