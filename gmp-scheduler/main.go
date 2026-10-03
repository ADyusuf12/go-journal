package main

import "fmt"

type Account struct {
	ID      int
	Balance float64
}

// Case 1: Value return — Stays on Stack
func createAccountValue(id int, balance float64) Account {
	acc := Account{ID: id, Balance: balance}
	return acc
}

// Case 2: Pointer return — Escapes to Heap!
func createAccountPointer(id int, balance float64) *Account {
	acc := Account{ID: id, Balance: balance}
	return &acc // Returning address causes memory to outlive stack frame
}

func main() {
	// Stack allocation
	acc1 := createAccountValue(1, 150.0)

	// Heap allocation
	acc2 := createAccountPointer(2, 300.0)

	// Interfaces force heap escaping
	fmt.Println("Account 1 Balance:", acc1.Balance)
	fmt.Println("Account 2 Pointer:", acc2)
}
