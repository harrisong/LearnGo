package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	count int
	ch chan int
	mutex sync.Mutex
}

func (c * Counter) Run() {
	for v := range(c.ch) {
		c.mutex.Lock()
		c.count += v
		c.mutex.Unlock()
	}
}

func (c *Counter) Increment() {
	c.ch <-1
}

func (c *Counter) Value() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.count
}

func main() {
	c := &Counter{
		count: 0,
		ch: make(chan int),
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