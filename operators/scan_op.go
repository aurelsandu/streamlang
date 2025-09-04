package operators

import (
	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

// ScanOp: fold progresiv (emite fiecare acumulare intermediară).
//  - init: seed-ul acumulării
//  - f: func(acc, x) -> acc'
// Propagă erorile și respectă contextul.
func ScanOp[T, R any](init R, f func(R, T) R) stream.Op[T, R] {
	return func(in stream.Stream[T]) stream.Stream[R] {
		acc := init
		return rt.Unary[T, R](in, func(it stream.Item[T], em *rt.Emitter[R]) bool {
			if it.Err != nil {
				_ = em.SendErr(it.Err)
				return false
			}
			acc = f(acc, it.Value)
			return em.SendVal(acc)
		})
	}
}