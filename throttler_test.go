package pipes

import (
	"context"
	"testing"
	"time"
)

func testThrottleChannel(t *testing.T) {
	l := 10
	ch := make(chan int)

	go func(ch chan int) {
		defer close(ch)
		for i := 0; i < l; i++ {
			ch <- i
		}
	}(ch)

	interval := time.Microsecond * 100
	rl := ThrottleChannel(context.TODO(), ch, interval)

	cont := true
	cntr := 0
	for cont {
		select {
		case <-time.After(time.Second):
			t.Errorf("ratelimit took to long")
			return
		case _, ok := <-rl:
			if !ok {
				cont = false
			} else {
				cntr++
			}
		}
	}
	if cntr != l {
		t.Errorf("expected %d elements, got %d", l, cntr)
	}
}
