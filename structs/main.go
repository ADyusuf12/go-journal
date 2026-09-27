package main

import (
	"fmt"
	"time"
)

// AuditHeader holds standard auditing fields for ERP entities.
type AuditHeader struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Invoice embeds AuditHeader directly.
type Invoice struct {
	AuditHeader  // Embedded field (anonymous field)
	CustomerName string
	TotalCents   int64
}

func main() {
	now := time.Now()

	inv := Invoice{
		AuditHeader: AuditHeader{
			ID:        1001,
			CreatedAt: now,
			UpdatedAt: now,
		},
		CustomerName: "Acme Logistics",
		TotalCents:   450000,
	}

	// Direct access to outer fields
	fmt.Println("Customer:", inv.CustomerName)

	// Field Promotion: You can access AuditHeader fields directly on inv!
	fmt.Println("Invoice ID:", inv.ID)
	fmt.Println("Created At:", inv.CreatedAt)

	// You can also access it via the full explicit path if you prefer
	fmt.Println("Explicit ID:", inv.AuditHeader.ID)
}
