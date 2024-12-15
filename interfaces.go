package pipes

type ValueOrError[T any] interface {
	Value() T
	Err() error
}
