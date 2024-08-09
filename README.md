# pipes

library with general purpose tools to work with golang `chan`nels

## details

`pipes.go`
- `Pour` consumes channel into a sink function
- `Map` transforms elements in channel
- `Materialize` copies a pointer stream into concrete types (dubious utility)
- `Collect` consumes a channel into a slice
- `Filter` filters elements that fail the predicate (returns false)
- `FanIn` collects multiple channels into single channel concurrently
- `ParMap` same as map but runs a goroutine for each input channel
- `RoundRobinFanOut` uses round robin to fan a channel out into multiple
- `FanOut` transforms a channel into two channels (copies)
- `FilterError` filters Results which have errors and calls a function for each error

`throttler.go`
- `ThrottleChannel` ratelimits the input channel with a duration, useful to not dos when for example fetching data from web resource

`result.go`
- mainly utility to have a single channel to combine values & errors, e.g. if a fetch fails the result can be only zero value in `Val` -field and error in the `Err` -field

`ring.go`
- mainly utility to have a ring which gives a next element with `Next()` and loops back to beginning at the end

`sender.go`
- `SendOrDone` tries to send element into a channel, cancels if context is done
- `SendOrTimeout`  tries to send element into a channel, cancels if timeout exceeded
