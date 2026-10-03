package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// --- Models ---

type Product struct {
	ID        int       `json:"id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderRequest struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type OrderResponse struct {
	OrderID   int       `json:"order_id"`
	SKU       string    `json:"sku"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// --- Application Server Container ---

type Application struct {
	DB *sql.DB
}

func main() {
	// 1. Connection String
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	// 2. Initialize PostgreSQL Connection Pool
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to initialize driver: %v", err)
	}
	defer db.Close()

	// Tune Connection Pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(15 * time.Minute)

	// Verify Database Reachability
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	fmt.Println("Database connection pool established!")

	// 3. Ensure Database Schema Exists
	if err := setupDatabaseSchema(ctx, db); err != nil {
		log.Fatalf("Schema setup failed: %v", err)
	}

	app := &Application{DB: db}

	// 4. Setup Routes & Middleware Pipeline
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/products", app.getProductsHandler)
	mux.HandleFunc("POST /api/v1/orders", app.createOrderHandler)

	loggedRouter := RequestLogger(mux)

	// 5. Configure HTTP Server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      loggedRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// 6. Graceful Shutdown Signal Interception Setup
	// Create a channel listening for OS signals (SIGINT from Ctrl+C, SIGTERM from Docker/K8s)
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Start HTTP server in a separate background Goroutine so it doesn't block main!
	go func() {
		fmt.Println("Capstone API Server listening on http://localhost:8080...")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Block main thread until an OS termination signal arrives!
	sig := <-shutdownChan
	fmt.Printf("\nReceived shutdown signal (%v). Initiating Graceful Shutdown...\n", sig)

	// 7. Execute Graceful Shutdown with a 10-second deadline
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Stop accepting new connections and wait for in-flight HTTP requests to complete cleanly
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown abruptly: %v", err)
	}

	fmt.Println("All in-flight requests completed. Server and Database resources closed cleanly.")
}

// --- Middleware ---

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[HTTP] %s %s | Remote: %s | Duration: %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

// --- Schema Initialization ---

func setupDatabaseSchema(ctx context.Context, db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		sku VARCHAR(50) UNIQUE NOT NULL,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(10,2) NOT NULL,
		stock INT NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		sku VARCHAR(50) NOT NULL,
		quantity INT NOT NULL,
		total_amount NUMERIC(10,2) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	INSERT INTO products (sku, name, price, stock)
	VALUES
		('PRD-101', 'Mechanical Keyboard', 85.50, 20),
		('PRD-102', 'Ergonomic Chair', 299.99, 10)
	ON CONFLICT (sku) DO NOTHING;`

	_, err := db.ExecContext(ctx, schema)
	return err
}

// --- Handlers ---

func (app *Application) getProductsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.QueryContext(r.Context(), "SELECT id, sku, name, price, stock, created_at FROM products ORDER BY id ASC")
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock, &p.CreatedAt); err != nil {
			http.Error(w, "Error scanning products", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Stream processing error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (app *Application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var req OrderRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid order payload: %v", err), http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 || req.SKU == "" {
		http.Error(w, "Invalid SKU or quantity", http.StatusUnprocessableEntity)
		return
	}

	tx, err := app.DB.BeginTx(r.Context(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		http.Error(w, "Transaction initiation failed", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var productID, stock int
	var price float64

	checkSQL := `SELECT id, price, stock FROM products WHERE sku = $1 FOR UPDATE;`
	err = tx.QueryRowContext(r.Context(), checkSQL, req.SKU).Scan(&productID, &price, &stock)
	if err == sql.ErrNoRows {
		http.Error(w, "Product SKU not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if stock < req.Quantity {
		http.Error(w, fmt.Sprintf("Insufficient stock! Available: %d, Requested: %d", stock, req.Quantity), http.StatusConflict)
		return
	}

	updateStockSQL := `UPDATE products SET stock = stock - $1 WHERE id = $2;`
	if _, err := tx.ExecContext(r.Context(), updateStockSQL, req.Quantity, productID); err != nil {
		http.Error(w, "Failed to update stock", http.StatusInternalServerError)
		return
	}

	totalAmount := price * float64(req.Quantity)
	var orderID int
	var orderCreatedAt time.Time

	insertOrderSQL := `
	INSERT INTO orders (sku, quantity, total_amount)
	VALUES ($1, $2, $3)
	RETURNING id, created_at;`

	err = tx.QueryRowContext(r.Context(), insertOrderSQL, req.SKU, req.Quantity, totalAmount).Scan(&orderID, &orderCreatedAt)
	if err != nil {
		http.Error(w, "Failed to record order", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	resp := OrderResponse{
		OrderID:   orderID,
		SKU:       req.SKU,
		Quantity:  req.Quantity,
		Total:     totalAmount,
		Status:    "CONFIRMED",
		CreatedAt: orderCreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
