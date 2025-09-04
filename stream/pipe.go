package stream

// Pipe: variantă simplă (fără Op) – aplică o transformare funcțională.
func Pipe[A any, B any](s Stream[A], op func(Stream[A]) Stream[B]) Stream[B] {
	return op(s)
}