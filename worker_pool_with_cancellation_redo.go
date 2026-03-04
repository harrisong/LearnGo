package main

import (
	"fmt"
	"sync"
    "context"
)

func worker(wg *sync.WaitGroup, jobs chan int, i int, context context.Context) {
	defer wg.Done()

    for {
        select {
            case <-context.Done():
                fmt.Printf("Context cancelled in Worker %d\n", i)
                return
            case job, ok := <-jobs:
                if !ok {
                    return
                }
                fmt.Printf("Worker %d processing job %d\n", i, job)
        }
    }
}

func main() {
    context, cancel := context.WithCancel(context.Background())
    defer cancel()

	jobs := make(chan int)
	var wg sync.WaitGroup

	var numJobs = 100
	var numWorkers = 3

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, jobs, i, context)
	}

	go func() {
		for j := 0; j < numJobs; j++ {
			jobs <- j
            if j == 10 {
                cancel()
            }
		}
		close(jobs)
	}()
	
	wg.Wait()
}