package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to configure DB driver: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Target product and deduction quantity
	targetSKU := "PRD-101"
	purchaseQty := 5

	fmt.Println("--- Starting Atomic Inventory Transaction ---")

	// 1. Begin Database Transaction
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}

	// Safety Net: Always defer Rollback! If Commit() succeeds first, Rollback is a safe no-op.
	defer tx.Rollback()

	// 2. Lock & Check stock level (FOR UPDATE locks the row against concurrent writes)
	var currentStock int
	checkStockSQL := `SELECT stock FROM products WHERE sku = $1 FOR UPDATE;`

	err = tx.QueryRowContext(ctx, checkStockSQL, targetSKU).Scan(&currentStock)
	if err != nil {
		log.Fatalf("Failed to fetch product stock: %v", err)
	}

	fmt.Printf("Current stock for %s: %d items\n", targetSKU, currentStock)

	if currentStock < purchaseQty {
		log.Fatalf("Insufficient stock! Available: %d, Requested: %d", currentStock, purchaseQty)
	}

	// 3. Deduct stock inside the same transaction session
	updateStockSQL := `UPDATE products SET stock = stock - $1 WHERE sku = $2 RETURNING stock;`
	var newStock int

	err = tx.QueryRowContext(ctx, updateStockSQL, purchaseQty, targetSKU).Scan(&newStock)
	if err != nil {
		log.Fatalf("Failed to update stock: %v", err)
	}

	// 4. Commit Transaction to persist changes to disk permanently
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Printf("Transaction Committed Successfully! Product %s new stock: %d\n", targetSKU, newStock)
}
