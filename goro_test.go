package goro_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/ghhernandes/goro"
)

func TestMap(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	toStr := func(in []byte) string {
		return string(in) + " stringfied"
	}

	byteStream := func() []byte {
		return []byte("Hello world")
	}

	for value := range goro.Map(ctx, limit(ctx, stream(ctx, byteStream), 10), toStr) {
		if value != "Hello world stringfied" {
			t.Error(value)
		}
	}
}

func TestFilter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	even := func(in int) bool {
		return in%2 == 0
	}

	odd := func(in int) bool {
		return !even(in)
	}

	intStream := func() int {
		return rand.Int()
	}

	for value := range goro.Filter(ctx, limit(ctx, stream(ctx, intStream), 10), odd) {
		if even(value) {
			t.Fail()
		}
	}
}

func TestFanOutFanIn(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	intStream := func() int {
		return rand.Int()
	}

	inputs := goro.FanOut(ctx, limit(ctx, stream(ctx, intStream), 10), 9)

	if len(inputs) != 9 {
		t.Error("FanOut len")
	}

	var i int
	for range goro.FanIn(ctx, inputs...) {
		i++
	}

	if i != 10 {
		t.Error(i)
	}
}

func stream[T any](ctx context.Context, value func() T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- value()
			}
		}
	}()
	return ch
}

func limit[T any](ctx context.Context, input <-chan T, n int) <-chan T {
	output := make(chan T)
	go func() {
		defer close(output)

		for i := 0; i < n; i++ {
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
	}()
	return output
}
