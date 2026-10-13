package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

// CreateOrderRequest represents the strict incoming payload schema
type CreateOrderRequest struct {
	Item     string  `json:"item"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// Custom validation method
func (req *CreateOrderRequest) Validate() error {
	if strings.TrimSpace(req.Item) == "" {
		return errors.New("field 'item' cannot be blank")
	}
	if req.Quantity <= 0 {
		return errors.New("field 'quantity' must be greater than zero")
	}
	if req.Price <= 0 {
		return errors.New("field 'price' must be greater than zero")
	}
	return nil
}

func handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	// 1. Guard against payload bombs (limit request body to 1MB)
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var order CreateOrderRequest

	// 2. Decode JSON stream directly into struct
	decoder := json.NewDecoder(r.Body)

	// CRITICAL SECURITY RULE: Reject any keys not defined in our struct
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&order)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request body: ` + escapeJSONString(err.Error()) + `"}`))
		return
	}

	// 3. Execute business validation rules
	if err := order.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // 422 Unprocessable Entity
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	// 4. Success Response
	slog.Info("order successfully created", "item", order.Item, "qty", order.Quantity)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	w.Write([]byte(`{"status":"success","message":"order placed successfully"}`))
}

// Simple helper to sanitize error strings for JSON output
func escapeJSONString(s string) string {
	s = strings.ReplaceAll(s, `"`, `\'`)
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

func main() {
	mux := http.NewServeMux()

	// Go 1.22+ method-based routing (`METHOD /path`)
	mux.HandleFunc("POST /api/orders", handleCreateOrder)

	slog.Info("Validation server running on port :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("server crashed", slog.String("err", err.Error()))
	}
}
