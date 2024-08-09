package pipes

import (
	"context"
	"time"
)

// ThrottleChannel  throttles the input channel with interval to the output channel
func ThrottleChannel[T any](in <-chan T, interval time.Duration, ctx context.Context) chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		rl := time.Tick(interval)
		for {
			select {
			case <-ctx.Done():
				return
			case <-rl:
				v, ok := <-in
				if !ok {
					return
				}
				out <- v
			}
		}
	}()
	return out
}
