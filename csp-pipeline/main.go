package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Order struct {
	ID     int
	Amount float64
	Status string
}

func main() {
	// Create a context that can be canceled mid-flight
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	// 1. Stage 1: Order Generator (Produces order IDs on a channel)
	ordersChan := generateOrders(ctx, 10)

	// 2. Stage 2: Fan-Out to a Worker Pool of 3 concurrent processors
	workerCount := 3
	resultChans := make([]<-chan Order, workerCount)

	for i := 0; i < workerCount; i++ {
		resultChans[i] = processOrdersWorker(ctx, i+1, ordersChan)
	}

	// 3. Stage 3: Fan-In (Merge multiple worker output channels into 1 stream)
	finalStream := fanIn(ctx, resultChans...)

	// 4. Consume the processed orders stream
	fmt.Println("--- Starting Order Pipeline Processing ---")
	for order := range finalStream {
		fmt.Printf(" [SUCCESS] Order #%d processed successfully. Status: %s\n", order.ID, order.Status)
	}

	fmt.Println("Pipeline execution finished.")
}

// Stage 1: Generator Pattern
func generateOrders(ctx context.Context, count int) <-chan Order {
	out := make(chan Order)

	go func() {
		defer close(out) // Always close output channel when generator finishes!
		for i := 1; i <= count; i++ {
			order := Order{ID: i, Amount: float64(i * 25)}

			select {
			case <-ctx.Done():
				fmt.Println("[Generator] Cancellation received! Stopping order generation.")
				return
			case out <- order:
			}
		}
	}()

	return out
}

// Stage 2: Worker Pool Stage
func processOrdersWorker(ctx context.Context, workerID int, in <-chan Order) <-chan Order {
	out := make(chan Order)

	go func() {
		defer close(out)
		for order := range in {
			// Simulate work (e.g., checking inventory or calling a payment gateway)
			time.Sleep(100 * time.Millisecond)

			order.Status = fmt.Sprintf("Processed by Worker #%d", workerID)

			select {
			case <-ctx.Done():
				fmt.Printf("[Worker #%d] Cancellation received! Halting active task.\n", workerID)
				return
			case out <- order:
			}
		}
	}()

	return out
}

// Stage 3: Fan-In Aggregator Pattern
func fanIn(ctx context.Context, channels ...<-chan Order) <-chan Order {
	out := make(chan Order)
	var wg sync.WaitGroup

	// Multiplex each inbound worker channel onto the single output channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan Order) {
			defer wg.Done()
			for order := range c {
				select {
				case <-ctx.Done():
					return
				case out <- order:
				}
			}
		}(ch)
	}

	// Wait for all worker channels to drain and close, then close output channel
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
