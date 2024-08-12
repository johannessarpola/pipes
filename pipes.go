// Package pipes contains bunch of FP styled functions for channels
package pipes

import (
	"context"
	"fmt"
	"sync"
)

// Pour consumes a channel, collects them into array and calls the sink func with it, respecting context cancellation
func Pour[T any](ctx context.Context, in <-chan T, sink func([]T) error, initial ...T) error {
	collect, err := Collect(ctx, in, initial...)
	if err != nil {
		return err
	}
	return sink(collect)
}

// Map transforms elements in channel to another type
func Map[T any, O any](ctx context.Context, in <-chan T, fn func(T) (O, error)) chan Result[O] {
	out := make(chan Result[O])
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				out <- Result[O]{Err: ctx.Err()}
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				func() {
					defer func() {
						if r := recover(); r != nil {
							out <- Result[O]{Err: fmt.Errorf("panic in transformation: %v", r)}
						}
					}()
					mappedValue, err := fn(value)
					select {
					case out <- Result[O]{Val: mappedValue, Err: err}:
					case <-ctx.Done():
						return
					}
				}()
			}
		}
	}()
	return out
}

// Materialize copies the pointer values into another channel as a concrete type, respecting context cancellation
func Materialize[T any](ctx context.Context, in <-chan *T) chan T {
	out := make(chan T)
	go func(in <-chan *T) {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				out <- *value
			}
		}
	}(in)
	return out
}

// Collect reads from the input channel and collects elements into a slice, respecting context cancellation
func Collect[T any](ctx context.Context, in <-chan T, initial ...T) ([]T, error) {
	result := initial
	for {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case v, ok := <-in:
			if !ok {
				// Channel closed, return the collected result
				return result, nil
			}
			result = append(result, v)
		}
	}
}

// Filter filters elements from channels that return `false`  from predicate function
func Filter[T any](ctx context.Context, in <-chan T, predicate func(T) bool) chan T {
	out := make(chan T)
	go func(ch <-chan T) {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case res, ok := <-ch:
				if !ok {
					return
				}
				if predicate(res) {
					out <- res
				}
			}
		}
	}(in)

	return out
}

// FanIn collects multiple channels into single output channel
func FanIn[T any](ctx context.Context, chans ...chan T) chan T {
	wg := &sync.WaitGroup{}
	wg.Add(len(chans))

	out := make(chan T)
	for _, ch := range chans {
		go func(ch <-chan T, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-ch:
					if ok {
						out <- v
					} else {
						return
					}
				}
			}
		}(ch, wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ParMap processes each input channel into single output channel, second chanel for errors if such occur with the process
func ParMap[T any, O any](ctx context.Context, process func(T) (O, error), channels ...chan T) (chan O, chan error) {
	out := make(chan O)
	errs := make(chan error, len(channels))
	wg := &sync.WaitGroup{}
	wg.Add(len(channels))

	// Start a goroutine to process each input channel into single output channel result
	for _, ch := range channels {
		go func(ch chan T, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				case v, ok := <-ch:
					if !ok {
						return
					}
					nv, err := process(v)

					// Early quit with error, better safe than sorry
					if err != nil {
						errs <- err
						return
					}
					out <- nv
				}
			}

		}(ch, wg)
	}

	// Wait for the subprocesses to finish and the close channels
	go func() {
		wg.Wait()
		close(out)
		close(errs)
	}()

	return out, errs
}

// RoundRobinFanOut fans a channel out into `outChannelCounts` - number of output channels in round robin manner
func RoundRobinFanOut[T any](
	ctx context.Context,
	in <-chan T,
	outChannelCounts int,
) []chan T {
	ocl := make([]chan T, outChannelCounts)
	for i := 0; i < outChannelCounts; i++ {
		ocl[i] = make(chan T)
	}

	go func() {
		defer func() {
			for _, c := range ocl {
				close(c)
			}
		}()

		r, err := NewRing(ocl)
		if err != nil {
			panic(err)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				r.Next() <- v
			}

		}

	}()

	return ocl
}

// FanOut fans a input channel out into two channels, respecting context cancellation
func FanOut[T any](ctx context.Context, in <-chan T) (chan T, chan T) {
	o1 := make(chan T)
	o2 := make(chan T)
	go func() {
		var arr []T

	consume:
		for {
			select {
			case <-ctx.Done():
				break consume
			case v, ok := <-in:
				if !ok {
					break consume
				}
				arr = append(arr, v)
			}
		}

		go func(arr []T) {
			defer close(o1)

			for _, e := range arr {
				select {
				case o1 <- e:
				case <-ctx.Done():
					return
				}
			}
		}(arr)

		go func(arr []T) {
			defer close(o2)

			for _, e := range arr {
				select {
				case o2 <- e:
				case <-ctx.Done():
					return
				}
			}
		}(arr)

	}()

	return o1, o2
}

// FilterError filters errored Results from channel and calls onError for each, respecting context cancellation
func FilterError[T any](ctx context.Context, resChan <-chan Result[T], onError func(err error)) chan T {
	out := make(chan T)
	go func(onError func(err error)) {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case res, ok := <-resChan:
				if !ok {
					return
				}
				if res.Err != nil {
					onError(res.Err)
				} else {
					out <- res.Val
				}
			}
		}
	}(onError)

	return out
}
