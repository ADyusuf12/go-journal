package main

import (
	"fmt"
	"sync"
)

type BankAccount struct {
	mu      sync.Mutex // Protects access to Balance field
	Balance int
}

func (a *BankAccount) Deposit(amount int) {
	// Acquire lock before touching shared memory
	a.mu.Lock()
	defer a.mu.Unlock() // Ensure lock releases even if code panics

	a.Balance += amount
}

func main() {
	account := &BankAccount{Balance: 1000}
	var wg sync.WaitGroup

	totalTransactions := 1000

	for i := 0; i < totalTransactions; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			account.Deposit(10)
		}()
	}

	wg.Wait()

	fmt.Printf("Final Account Balance: $%d\n", account.Balance)
}
