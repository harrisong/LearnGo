package main

import (
	"fmt"
	"sync"
)

func worker(wg *sync.WaitGroup, jobs chan int, i int) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", i, job)
	}
}

func main() {
	jobs := make(chan int)
	var wg sync.WaitGroup

	var numJobs = 100
	var numWorkers = 3

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, jobs, i)
	}

	go func() {
		for j := 0; j < numJobs; j++ {
			jobs <- j
		}
		close(jobs)
	}()
	
	wg.Wait()
}