package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// LedgerItem represents our incoming financial ledger entry payload
type LedgerItem struct {
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
	Category  string  `json:"category"`
}

// Validate ensures business rules are met before hitting PostgreSQL
func (item *LedgerItem) Validate() error {
	if strings.TrimSpace(item.AccountID) == "" {
		return errors.New("field 'account_id' cannot be blank")
	}
	if item.Amount <= 0 {
		return errors.New("field 'amount' must be greater than zero")
	}
	if strings.TrimSpace(item.Category) == "" {
		return errors.New("field 'category' cannot be blank")
	}
	return nil
}

// AppContext holds our shared application dependencies (e.g., database pool, logger)
type AppContext struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// HandleCreateLedgerEntry is our clean HTTP handler bound to AppContext
func (app *AppContext) HandleCreateLedgerEntry(w http.ResponseWriter, r *http.Request) {
	// 1. Guard against payload bombs (1MB limit)
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var item LedgerItem
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Reject unexpected JSON keys

	if err := decoder.Decode(&item); err != nil {
		app.logger.Error("failed to decode ledger request", slog.String("err", err.Error()))
		writeJSONError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	// 2. Business validation
	if err := item.Validate(); err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// 3. Database persistence with strict request-scoped timeout context (3 seconds max)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO ledger_entries (account_id, amount, category, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id`

	var insertedID int64
	err := app.db.QueryRow(ctx, query, item.AccountID, item.Amount, item.Category).Scan(&insertedID)
	if err != nil {
		app.logger.Error("database insert failed", slog.String("err", err.Error()))
		writeJSONError(w, http.StatusInternalServerError, "database persistence error")
		return
	}

	// 4. Success Response
	app.logger.Info("ledger entry created successfully", slog.Int64("id", insertedID), slog.String("account", item.AccountID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"id":     insertedID,
	})
}

// Helper for clean JSON error responses
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Database connection string (adjust credentials as needed for your local setup)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/go_journal_dev?sslmode=disable"
	}

	// Initialize pgxpool connection pool
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Error("unable to connect to database", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		logger.Error("database ping failed", slog.String("err", err.Error()))
		os.Exit(1)
	}

	app := &AppContext{
		db:     dbPool,
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/ledger", app.HandleCreateLedgerEntry)

	logger.Info("Starting database-integrated web server on port :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Error("server crashed", slog.String("err", err.Error()))
	}
}
