package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: context done, exiting\n", id)
			return
		default:
			fmt.Printf("Worker %d: working...\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

	go worker(ctx, 1)

    time.Sleep(1 * time.Second)
    fmt.Println("Main done")
	time.Sleep(1 * time.Second)
}
