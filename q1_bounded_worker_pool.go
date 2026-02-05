package main

import (
	"fmt"
	"sync"
	"context"
	"time"
	"math/rand"
)

func worker(ctx context.Context, in <-chan int, out chan <-int, wg *sync.WaitGroup, worker int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <- in:
			if !ok {
				return
			}
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

			select {
			case <- ctx.Done():
				return
			case out <- task:
			}
		}
	}
}

func ProcessTasks(tasks []int, numWorkers int, timeout time.Duration) []int {
	var ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var wg sync.WaitGroup
	var in = make(chan int, len(tasks))
	var out = make(chan int, len(tasks))

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, in, out, &wg, i)
	}

	for _, val := range tasks {
		in <-val
	}
	close(in)

	wg.Wait()
	close(out)

	var results []int
	for r := range out {
		results = append(results, r)
	}

	return results
}

func main() {
	rand.Seed(time.Now().UnixNano())
	tasks := []int{1, 2, 3, 4, 5} 
	results := ProcessTasks(tasks, 3, 2*time.Second) 
	fmt.Println(results)
}