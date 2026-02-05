package main

import (
	"time"
	"fmt"
)

func RateLimitedLogger(logs <-chan string, rate int) {
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case l, ok := <-logs:
				if (!ok) {
					return
				}

				fmt.Println(l)
			}
		}
	}
}

func main() {
	logs := make(chan string)
	go func() {
		for i := 0; i < 10; i++ {
			logs <- fmt.Sprintf("Log #%d", i)
		}
		close(logs)
	}()

	RateLimitedLogger(logs, 2) // Should print at most 2 logs per second
}