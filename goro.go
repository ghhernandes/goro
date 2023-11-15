package goro

import (
	"context"
	"sync"
)

// Fan-out receives a channel and multiplexes it into N channels. N goroutines
// will be started to read from the input and write to the output channels.
func FanOut[T any](ctx context.Context, input <-chan T, n int) []<-chan T {
	outputs := make([]<-chan T, n)

	for i := 0; i < n; i++ {
		out := make(chan T)
		go func() {
			defer close(out)
			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-input:
					if !ok {
						return
					}
					out <- value
				}
			}
		}()
		outputs[i] = out
	}

	return outputs
}

// FanIn receives multiple channels and returns a single channel. N+1 goroutines
// will be started. Where N is equal to the length of input channels and another
// channel to wait goroutines to finish.
func FanIn[T any](ctx context.Context, inputs ...<-chan T) <-chan T {
	out := make(chan T)
	wg := sync.WaitGroup{}

	for i := 0; i < len(inputs); i++ {
		wg.Add(1)
		go func(ctx context.Context, input <-chan T, output chan<- T) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-input:
					if !ok {
						return
					}
					output <- value
				}
			}
		}(ctx, inputs[i], out)
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}

// Map start a new goroutine will be created to receive input data, apply the
// map function and send to the output channel.
func Map[T1, T2 any](ctx context.Context, input <-chan T1, mapFn func(T1) T2) <-chan T2 {
	out := make(chan T2)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-input:
				if !ok {
					return
				}
				out <- mapFn(value)
			}
		}
	}()

	return out
}

// Filter start a new goroutine and returns a channel that only returns data
// thats the filter function returns `true`.
func Filter[T any](ctx context.Context, input <-chan T, filter func(T) bool) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-input:
				if !ok {
					return
				}
				if filter(value) {
					out <- value
				}
			}
		}
	}()

	return out
}
