package goro

import (
	"context"
	"sync"
)

// FanOut receives a single channel and multiplexes to N channels.
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

// FanIn receives multiples channels and returns one that will receive all data.
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

// Map receives a input channel, apply the mapFn function and write to the
// output channel.
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

// Filter receives a input channel, apply the filter function and write matches
// to the output channel.
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
