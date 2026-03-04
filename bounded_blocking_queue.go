package main

type BoundedQueue struct {
	ch chan int
}

func NewBoundedQueue(capacity int) *BoundedQueue {
	return &BoundedQueue{
		ch: make(chan int, capacity),
	}
}

func (b *BoundedQueue) Enqueue(item int) {
	b.ch <- item
}

func (b *BoundedQueue) Dequeue() int {
	return <-b.ch
}

func main() {
	queue := NewBoundedQueue(2)
	queue.Enqueue(1)
	queue.Enqueue(2)

	go func() {
		queue.Enqueue(3) // This will block until an item is dequeued
		queue.Enqueue(4) // This will block until an item is dequeued
	}()

	println(queue.Dequeue()) // Output: 1
	println(queue.Dequeue()) // Output: 2
	println(queue.Dequeue()) // Output: 3
	println(queue.Dequeue()) // Output: 4
}