package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure wg.Done() is called when the worker is finished

	for orderID := range jobs {
		fmt.Printf("Worker %d started processing order %d\n", id, orderID)
		time.Sleep(500 * time.Millisecond) // Simulate processing time
		fmt.Printf("Worker %d finished processing order %d\n", id, orderID)
	}
}

func main() {
	// Create a job queue channel with a capacity of 10
	jobs := make(chan int, 10)
	var wg sync.WaitGroup

	// Start 2 concurrent worker Goroutines (Worker 1 and WOrker 2)
	//We increment the WaitGroup counter for each worker before launching them
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	//Enqueue 5 order IDs (101-105)
	for orderID := 101; orderID <= 105; orderID++ {
		jobs <- orderID
		fmt.Printf("Enqueued order %d\n", orderID)
	}

	close(jobs) // Close the jobs channel to signal no more jobs will be sent

	wg.Wait() // Wait for all workers to finish processing
	fmt.Println("All orders processed sucessfully.")
}
