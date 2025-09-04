package stream

import "context"

// Item: valoare + eroare propagată prin stream.
type Item[T any] struct {
	Value T
	Err   error
}

// Stream: (Ctx, Ch) tipizat.
type Stream[T any] struct {
	Ctx context.Context
	Ch  <-chan Item[T]
}
