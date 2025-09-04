package operators

import (
	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

func MapOp[A, B any](fn func(A) B) stream.Op[A, B] {
	return func(in stream.Stream[A]) stream.Stream[B] {
		return rt.Unary[A, B](in, func(it stream.Item[A], em *rt.Emitter[B]) bool {
			if it.Err != nil {
				_ = em.SendErr(it.Err); return false
			}
			return em.SendVal(fn(it.Value))
		})
	}
}