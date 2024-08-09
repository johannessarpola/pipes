package pipes

import "fmt"

// Result is a struct to wrap either a element or a error
type Result[T any] struct {
	Val T
	Err error
}

// NewErrResult builds a result with error
func NewErrResult[T any](err error) Result[T] {
	var noop T // Zero value
	return Result[T]{
		Val: noop,
		Err: err,
	}
}

// NewResult builds a result with successful value
func NewResult[T any](val T) Result[T] {
	return Result[T]{
		Val: val,
		Err: nil,
	}
}

func (r Result[T]) WithError(err error) Result[T] {
	r.Err = err
	return r
}

func (r Result[T]) WithValue(val T) Result[T] {
	r.Val = val
	return r
}

// Error
func (r Result[T]) Error() string {
	if r.Err != nil {
		return r.Err.Error()
	}
	return ""
}

// String
func (r Result[T]) String() string {
	if r.Err != nil {
		return fmt.Sprintf("%v", r.Err)
	}
	return fmt.Sprintf("%v", r.Val)
}
