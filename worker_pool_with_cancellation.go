package main

import (
	"context"
	"fmt"
	"sync"
)

func worker(ctx context.Context, wg *sync.WaitGroup, ch chan int, count chan int, worker int) {
	defer wg.Done()
	for {
		select {
			case job, ok := <-ch:
				if !ok 
					return
				
				fmt.Printf("Worker %d processing job %d\n", worker, job)
				count <- 1
			case <-ctx.Done():
				fmt.Printf("Worker %d received cancellation signal\n", worker)
				return
		}
	}
}

func main() {
	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan int)
	count := make(chan int)

	var numWorkers = 3
	var numJobs = 100

	var wg sync.WaitGroup
	// Start worker pool with context
	for i := 0; i< numWorkers; i++ {
		wg.Add(1)
		go worker(context, &wg, ch, count, i)
	}

	go func() {
		for i := 0; i < numJobs; i++ {
			ch <- i
		}

		close(ch)
	}()

	go func() {
		wg.Wait()
		close(count)
	}()

	total := 0
	for i := range count {
		total += i
	}

	fmt.Printf("Total processed jobs: %d\n", total)
}