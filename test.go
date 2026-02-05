package main

import "fmt"
import "sync"
import "time"

func worker(c chan int) {
    fmt.Println("Working...")
    c <- 1024 // signal completion
}

func main() {
    c := make(chan int, 1)
    c2 := make(chan int, 1)

    c <- 41
    c2 <- 42

    // go worker(c) // leaked goroutine
    // fmt.Println("Worker finished", <-c)

    // ch := make(chan int, 1)
    // go func() {
    //     ch <- 42
    // }()

    // select {
    // case msg := <-c:
    //     fmt.Println(msg)
    // case msg := <-c2:
    //     fmt.Println("c2 finished:", msg)
    // }
    
    fmt.Println(<-c)
    fmt.Println(<-c2)

    var wg sync.WaitGroup

    // wg.Add(1)
    go func() {
        time.Sleep(500 * time.Millisecond)
        fmt.Println("Goroutine working...")
        defer wg.Done()
    }()
    fmt.Println("Continue")

    wg.Wait()
    fmt.Println("WaitGroup done")
}
