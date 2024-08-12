package pipes

import (
	"context"
	"time"
)

// SendOrDone tries to send a element into chan or cancels if context is done
func SendOrDone[T any](ctx context.Context, ele T, dst chan T) {
	select {
	case <-ctx.Done():
		// ctx canceled before sent successfully
	case dst <- ele:
		// sent successfully within the time out context
	}
	return
}

// SendOrTimeout tries to send a element into chan or cancels after timeout
func SendOrTimeout[T any](ele T, dst chan T, duration time.Duration) {
	select {
	case <-time.After(duration):
		// timeout before sent
	case dst <- ele:
		// sent successfully within the time out context
	}
	return
}
