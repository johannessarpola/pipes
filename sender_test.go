package pipes

import (
	"context"
	"testing"
	"time"
)

func TestSendOrDone(t *testing.T) {
	c, cancel := context.WithTimeout(context.TODO(), time.Millisecond)
	defer cancel()
	ch := make(chan int)
	n := time.Now()
	// This should be canceled with context timeout since there is no consumer
	SendOrDone(1, ch, c)

	// It should be over in 1 Millisecond
	if time.Since(n) > time.Millisecond*100 {
		t.Errorf("sending was not canceled")
	}

	var ee []int
	c2, cancel2 := context.WithTimeout(context.TODO(), time.Second)
	defer cancel2()

	// This sending should succedd as context lifetime is 1 second
	go func() {
		defer close(ch)
		SendOrDone(1, ch, c2)
	}()

	for e := range ch {
		ee = append(ee, e)
	}

	if len(ee) != 1 {
		t.Errorf("sending was not canceled")
	}
}

func TestSendOrTimeout(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	c1Cntr := 0
	c2Cntr := 0

	// This goroutine sending should not get canceled since there is consumer within the timeout
	go func() {
		defer close(ch1)
		SendOrTimeout(1, ch1, time.Second)
	}()

	// This goroutine sending should get canceled since there is NO consumer within the timeout
	go func() {
		defer close(ch2)
		SendOrTimeout(1, ch2, time.Microsecond)
	}()

	for range ch1 {
		c1Cntr++
	}

	if c1Cntr == 0 {
		t.Errorf("sending was not successful within timeframe for first channel")
	}

	// Wait for Millisecond to make sure the Microsecond timeout triggers for ch2
	time.Sleep(time.Millisecond)

	for range ch2 {
		c2Cntr++
	}

	if c1Cntr == c2Cntr || c2Cntr != 0 {
		t.Errorf("sending was not canceled for second channel")
	}
}
