package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Transaction represents an incoming financial ledger item
type Transaction struct {
	ID     string
	Amount float64
	Status string
}

// EngineMetrics tracks thread-safe operational counters using atomic counters & mutexes
type EngineMetrics struct {
	mu          sync.Mutex
	processed   int64
	failed      int64
	totalAmount float64
}

func (m *EngineMetrics) RecordSuccess(amount float64) {
	atomic.AddInt64(&m.processed, 1)
	m.mu.Lock()
	m.totalAmount += amount
	m.mu.Unlock()
}

func (m *EngineMetrics) RecordFailure() {
	atomic.AddInt64(&m.failed, 1)
}

func main() {
	// 1. Establish context with a strict operational timeout (graceful cancellation window)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	metrics := &EngineMetrics{}

	// 2. Stage 1: Ingestion Generator (Streams 15 raw transactions)
	rawTransactions := generateTransactions(ctx, 15)

	// 3. Stage 2: Fan-Out to Worker Pool (4 concurrent validation workers)
	workerCount := 4
	workerStreams := make([]<-chan Transaction, workerCount)
	for i := 0; i < workerCount; i++ {
		workerStreams[i] = workerProcessor(ctx, i+1, rawTransactions, metrics)
	}

	// 4. Stage 3: Fan-In Aggregator (Multiplexes worker output channels into a single stream)
	finalStream := fanIn(ctx, workerStreams...)

	// 5. Consume and Log Processed Results
	fmt.Println("--- Starting High-Throughput Ledger Ingestion Engine ---")
	for tx := range finalStream {
		fmt.Printf(" [LEDGER SYNC] ID: %s | Amount: $%.2f | Status: %s\n", tx.ID, tx.Amount, tx.Status)
	}

	// 6. Report Final Thread-Safe Metrics
	fmt.Println("\n--- Ingestion Pipeline Finished ---")
	fmt.Printf("Total Successfully Processed: %d\n", metrics.processed)
	fmt.Printf("Total Failed / Aborted:       %d\n", metrics.failed)
	fmt.Printf("Total Ledger Volume Synced:   $%.2f\n", metrics.totalAmount)
}

// Stage 1: Stream Generator
func generateTransactions(ctx context.Context, count int) <-chan Transaction {
	out := make(chan Transaction)

	go func() {
		defer close(out)
		for i := 1; i <= count; i++ {
			txID := fmt.Sprintf("TX-99%02d", i)
			amount := float64(i * 150)

			tx := Transaction{ID: txID, Amount: amount, Status: "PENDING"}

			select {
			case <-ctx.Done():
				fmt.Println("[Generator] Context canceled. Halting ingestion stream.")
				return
			case out <- tx:
			}
		}
	}()

	return out
}

// Stage 2: Concurrent Worker Processor (Fan-Out)
func workerProcessor(ctx context.Context, workerID int, in <-chan Transaction, metrics *EngineMetrics) <-chan Transaction {
	out := make(chan Transaction)

	go func() {
		defer close(out)
		for tx := range in {
			// Simulate cryptographic validation or database lookup latency
			time.Sleep(50 * time.Millisecond)

			// Validate payload bounds
			if tx.Amount > 1500.0 {
				tx.Status = fmt.Sprintf("REJECTED (Exceeds Limit) by Worker #%d", workerID)
				metrics.RecordFailure()
			} else {
				tx.Status = fmt.Sprintf("APPROVED by Worker #%d", workerID)
				metrics.RecordSuccess(tx.Amount)
			}

			select {
			case <-ctx.Done():
				fmt.Printf("[Worker #%d] Context closed. Aborting active task %s.\n", workerID, tx.ID)
				return
			case out <- tx:
			}
		}
	}()

	return out
}

// Stage 3: Multiplexer (Fan-In)
func fanIn(ctx context.Context, channels ...<-chan Transaction) <-chan Transaction {
	out := make(chan Transaction)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan Transaction) {
			defer wg.Done()
			for tx := range c {
				select {
				case <-ctx.Done():
					return
				case out <- tx:
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
