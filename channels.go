package main

import "fmt"

func main() {
	c := make(chan int, 1)
	c <- 40
	close(c)

	val, ok := <-c
	fmt.Println(val, ok)

	val2, ok2 := <-c
	fmt.Println(val2, ok2)
}