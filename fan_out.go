package main

import (
    "fmt"
    "sync"
)

func squareWorker(id int, in <-chan int, out chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for n := range in {
        out <- n * n
    }
}

func main() {
    in := make(chan int, 5)
    out := make(chan int, 5)
    var wg sync.WaitGroup

    // Fan-out: 2 workers
    for i := 0; i < 2; i++ {
        wg.Add(1)
        go squareWorker(i, in, out, &wg)
    }

    // Send input
    for i := 1; i <= 5; i++ {
        in <- i
    }
    close(in)

    // Wait for workers to finish and close output
    go func() {
        wg.Wait()
        close(out)
    }()

    // Fan-in: collect results
    for result := range out {
        fmt.Println("Result:", result)
    }
}
