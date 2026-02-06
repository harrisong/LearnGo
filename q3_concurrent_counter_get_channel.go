package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	count int
	ch chan int
	get chan chan int
}

func (c * Counter) Run() {
	for {
		select {
		case v := <-c.ch:
			c.count += v
		case resp := <-c.get:
			// Received read channel (inbox) from Value()
			resp <-c.count
		}
	}
}

func (c *Counter) Increment() {
	c.ch <-1
}

func (c *Counter) Value() int {
	resp := make(chan int)
	c.get <- resp
	return <-resp
}

func main() {
	c := &Counter{
		count: 0,
		ch: make(chan int),
		get: make(chan chan int),
	}

	go c.Run()

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}

	wg.Wait()
	close(c.ch)
	fmt.Println("Final count: ", c.Value())
}