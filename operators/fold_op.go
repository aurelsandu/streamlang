package operators

import (
	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

// FoldOp: reduce fluxul la O SINGURĂ valoare finală.
//  - init: seed-ul
//  - f: func(acc, x) -> acc'
//
// Emită un singur Item (valoarea finală) când input-ul se închide.
// Dacă apare o eroare pe parcurs, o propagă și se oprește.
func FoldOp[T, R any](init R, f func(R, T) R) stream.Op[T, R] {
	return func(in stream.Stream[T]) stream.Stream[R] {
		// Implementăm manual, pentru că fold-ul emite o singură valoare la final,
		// nu la fiecare pas (ca Scan).
		em := rt.NewEmitter[R](in.Ctx)
		go func() {
			defer em.Close()
			acc := init
			for it := range in.Ch {
				if it.Err != nil {
					_ = em.SendErr(it.Err)
					return
				}
				acc = f(acc, it.Value)
			}
			_ = em.SendVal(acc)
		}()
		return stream.Stream[R]{Ctx: in.Ctx, Ch: em.Out}
	}
}