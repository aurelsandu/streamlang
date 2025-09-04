package operators

import (
	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

func TakeOp[T any](n int) stream.Op[T, T] {
	if n < 0 { n = 0 }
	return func(in stream.Stream[T]) stream.Stream[T] {
		count := 0
		return rt.Unary[T, T](in, func(it stream.Item[T], em *rt.Emitter[T]) bool {
			if it.Err != nil {
				_ = em.SendErr(it.Err); return false
			}
			if count >= n { return false }
			if ok := em.SendVal(it.Value); !ok { return false }
			count++
			return count < n
		})
	}
}


/*package operators

import "github.com/aurelsandu/StreamLang/stream"

// TakeOp ridică Take într-un Op[T,T] compozabil cu stream.Compose.
func TakeOp[T any](n int) stream.Op[T, T] {
	return func(s stream.Stream[T]) stream.Stream[T] {
		return Take(s, n)
	}
}*/