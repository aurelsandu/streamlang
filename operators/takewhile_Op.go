package operators

import (
	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

// Emite cât timp predicatul e TRUE; la primul FALSE oprește.
func TakeWhileOp[T any](pred func(T) bool) stream.Op[T, T] {
	return func(in stream.Stream[T]) stream.Stream[T] {
		return rt.Unary[T, T](in, func(it stream.Item[T], em *rt.Emitter[T]) bool {
			if it.Err != nil {
				_ = em.SendErr(it.Err)
				return false // stop
			}
			if !pred(it.Value) {
				return false // stop fără să emită it
			}
			return em.SendVal(it.Value) // continuă dacă reușește să trimită
		})
	}
}