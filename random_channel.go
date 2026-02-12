package main

import "time"

func main() {
	ch := make(chan int)
	ch2 := make(chan int)

	// go func() {
	// 	ch <- 1
	// }()

	// go func() {
	// 	ch <- 2
	// }()

	select {
	case val := <-ch:
		println("Received from ch:", val)
	case val := <-ch2:
		println("Received from ch2:", val)
	case <-time.After(1 * time.Second):
		println("Timeout occurred")
	}

	close(ch)
	close(ch2)
}