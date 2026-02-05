package main

import "fmt"

func main() {
	c := make(chan int, 1)

	c <- 40

	for i := 0; i < 1; i++ {
		fmt.Println(<-c)
	}
}