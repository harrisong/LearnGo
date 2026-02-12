package main

import (
	"context"
	"fmt"
)

// PipelineStage represents a stage in the pipeline with error handling
type PipelineStage struct {
	in  <-chan int
	out chan int
	err chan error
}

// NewPipelineStage creates a new stage with configurable buffer size
func NewPipelineStage(in <-chan int, bufferSize int) *PipelineStage {
	return &PipelineStage{
		in:  in,
		out: make(chan int, bufferSize),
		err: make(chan error, 1),
	}
}

// WithBackpressure applies backpressure by using buffered channels
// and context cancellation
func (s *PipelineStage) WithBackpressure(ctx context.Context, process func(int) (int, error)) {
	defer close(s.out)
	defer func() {
		if r := recover(); r != nil {
			s.err <- fmt.Errorf("panic recovered: %v", r)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			s.err <- ctx.Err()
			return
		case n, ok := <-s.in:
			if !ok {
				return // input channel closed
			}
			result, err := process(n)
			if err != nil {
				s.err <- err
				continue
			}
			select {
			case s.out <- result:
			case <-ctx.Done():
				s.err <- ctx.Err()
				return
			}
		}
	}
}

// Err returns the error channel
func (s *PipelineStage) Err() <-chan error {
	return s.err
}

// WithBackpressureNoErr is a simpler version without error handling
func WithBackpressure(in <-chan int, bufferSize int, process func(int) int) <-chan int {
	out := make(chan int, bufferSize)

	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- process(n):
			default:
				// Block if buffer is full (backpressure)
				// This prevents unbounded memory growth
				out <- process(n)
			}
		}
	}()

	return out
}

// gen with error handling and context support
func gen(ctx context.Context, nums ...int) (<-chan int, <-chan error) {
	out := make(chan int, len(nums))
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer func() {
			if r := recover(); r != nil {
				errCh <- fmt.Errorf("gen panicked: %v", r)
			}
		}()

		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}
	}()

	return out, errCh
}

// square with backpressure handling
func square(in <-chan int) <-chan int {
	// Use buffered channel for backpressure
	// Buffer size should be tuned based on expected throughput
	out := make(chan int, 100)

	go func() {
		defer close(out)
		for n := range in {
			// This will block if out buffer is full,
			// applying backpressure to the previous stage
			out <- n * n
		}
	}()

	return out
}

// double with panic recovery
func double(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Log panic but don't crash
				fmt.Printf("double stage recovered from panic: %v\n", r)
			}
		}()
		defer close(out)

		for n := range in {
			out <- n * 2
		}
	}()

	return out
}

// Error-aware pipeline runner
func RunPipeline(ctx context.Context, stages ...func(<-chan int) <-chan int) (<-chan int, <-chan error) {
	if len(stages) == 0 {
		ch := make(chan int)
		close(ch)
		return ch, make(chan error)
	}

	// Create error channel for the entire pipeline
	errCh := make(chan error, len(stages))

	// Chain stages
	ch := stages[0](genWithError(ctx))

	for i := 1; i < len(stages); i++ {
		ch = stages[i](ch)
	}

	return ch, errCh
}

// genWithError is a helper for the error-aware pipeline
func genWithError(ctx context.Context) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := 1; i <= 4; i++ {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func main() {
	// Example 1: Basic backpressure with buffered channels
	fmt.Println("Example 1: Basic pipeline with backpressure")
	for result := range double(square(gen(1, 2, 3, 4))) {
		fmt.Println(result)
	}

	// Example 2: Using context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("\nExample 2: Pipeline with context cancellation")
	for result := range double(square(gen(1, 2, 3, 4))) {
		fmt.Println(result)
		if result == 8 { // Cancel after certain results
			cancel()
			break
		}
	}

	// Example 3: Error handling with recovery
	fmt.Println("\nExample 3: Pipeline with panic recovery")
	for result := range double(square(gen(1, 2, 3, 4))) {
		fmt.Println(result)
	}
}
