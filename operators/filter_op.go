package operators

import (
	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

func FilterOp[T any](pred func(T) bool) stream.Op[T, T] {
	return func(in stream.Stream[T]) stream.Stream[T] {
		return rt.Unary[T, T](in, func(it stream.Item[T], em *rt.Emitter[T]) bool {
			if it.Err != nil {
				_ = em.SendErr(it.Err); return false
			}
			if pred(it.Value) {
				return em.SendVal(it.Value)
			}
			return true // skip & continue
		})
	}
}
/*
package operators

import "github.com/aurelsandu/StreamLang/stream"

// FilterOp ridică predicatul într-un Op[T,T] compozabil cu stream.Compose.
func FilterOp[T any](pred func(T) bool) stream.Op[T, T] {
	return func(s stream.Stream[T]) stream.Stream[T] {
		return Filter(s, pred)
	}
}
*/