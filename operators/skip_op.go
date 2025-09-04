package operators

import (
	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

func SkipOp[T any](n int) stream.Op[T, T] {
	if n < 0 { n = 0 }
	return func(in stream.Stream[T]) stream.Stream[T] {
		skipped := 0
		return rt.Unary[T, T](in, func(it stream.Item[T], em *rt.Emitter[T]) bool {
			if it.Err != nil {
				_ = em.SendErr(it.Err); return false
			}
			if skipped < n {
				skipped++
				return true // continuă fără emit
			}
			return em.SendVal(it.Value)
		})
	}
}
/*
package operators

import "github.com/aurelsandu/StreamLang/stream"

// SkipOp ridică Skip într-un Op[T,T] compozabil cu stream.Compose.
func SkipOp[T any](n int) stream.Op[T, T] {
	return func(s stream.Stream[T]) stream.Stream[T] {
		return Skip(s, n)
	}
}
*/