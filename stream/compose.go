package stream

// Op este o transformare: primește un Stream[A] și produce un Stream[B].
type Op[A any, B any] func(Stream[A]) Stream[B]

// PipeOp aplică un Op pe un Stream.
func PipeOp[A any, B any](s Stream[A], op Op[A, B]) Stream[B] {
	return op(s)
}


// Compose leagă două transformări: mai întâi a, apoi b.
func Compose[A any, B any, C any](a Op[A, B], b Op[B, C]) Op[A, C] {
	return func(s Stream[A]) Stream[C] {
		return b(a(s))
	}
}



// Compose3 – util când vrei să eviți paranteze (A->B->C->D).
func Compose3[A any, B any, C any, D any](o1 Op[A, B], o2 Op[B, C], o3 Op[C, D]) Op[A, D] {
	return func(s Stream[A]) Stream[D] {
		return o3(o2(o1(s)))
	}
}

// Compose4 compune patru operatori consecutivi.
func Compose4[A, B, C, D, E any](
	op1 Op[A, B],
	op2 Op[B, C],
	op3 Op[C, D],
	op4 Op[D, E],
) Op[A, E] {
	return func(s Stream[A]) Stream[E] {
		return op4(op3(op2(op1(s))))
	}
}

// Identity – transformare no-op; utilă la construcții dinamice.
func Identity[T any]() Op[T, T] {
	return func(s Stream[T]) Stream[T] { return s }
}
