package operators

import (
	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

// Sare cât timp predicatul e TRUE; la primul FALSE emită acel element și restul.
func SkipWhileOp[T any](pred func(T) bool) stream.Op[T, T] {
	return func(in stream.Stream[T]) stream.Stream[T] {
		skipping := true
		return rt.Unary[T, T](in, func(it stream.Item[T], em *rt.Emitter[T]) bool {
			if it.Err != nil {
				_ = em.SendErr(it.Err)
				return false
			}
			if skipping {
				if pred(it.Value) {
					return true // continuă să citească, nu emite
				}
				skipping = false // a găsit primul FALSE: de aici încolo emite
			}
			return em.SendVal(it.Value)
		})
	}
}