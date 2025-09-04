package operators

import (
	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

//
// 1) ConcatWithOp: concatenează stream-ul de intrare cu alte stream-uri date, secvențial.
//    Exemplu: out = in ⧺ others[0] ⧺ others[1] ⧺ ...
//

func ConcatWithOp[T any](others ...stream.Stream[T]) stream.Op[T, T] {
	return func(in stream.Stream[T]) stream.Stream[T] {
		em := rt.NewEmitter[T](in.Ctx)

		emitAll := func(s stream.Stream[T]) bool {
			for it := range s.Ch {
				// Respectă contextul
				select {
				case <-in.Ctx.Done():
					return false
				default:
				}
				// Propagă erorile ca item-uri de eroare și oprește concat-ul
				if it.Err != nil {
					return em.SendErr(it.Err)
				}
				if !em.SendVal(it.Value) {
					return false
				}
			}
			return true
		}

		go func() {
			defer em.Close()

			// 1) emite tot din input
			if !emitAll(in) {
				return
			}
			// 2) emite, pe rând, din fiecare „other”
			for _, s := range others {
				if !emitAll(s) {
					return
				}
			}
		}()

		return stream.Stream[T]{Ctx: in.Ctx, Ch: em.Out}
	}
}

//
// 2) ConcatStreamsOp: pentru Stream[Stream[T]] → emită inner-urile pe rând (nu paralel).
//    Diferența față de FlattenOp: aici aștepți să se termine complet inner curent, apoi treci la următorul.
//

func ConcatStreamsOp[T any]() stream.Op[stream.Stream[T], T] {
	return func(outer stream.Stream[stream.Stream[T]]) stream.Stream[T] {
		em := rt.NewEmitter[T](outer.Ctx)

		emitAll := func(s stream.Stream[T]) bool {
			for it := range s.Ch {
				select {
				case <-outer.Ctx.Done():
					return false
				default:
				}
				if it.Err != nil {
					return em.SendErr(it.Err)
				}
				if !em.SendVal(it.Value) {
					return false
				}
			}
			return true
		}

		go func() {
			defer em.Close()
			for innerOrErr := range outer.Ch {
				select {
				case <-outer.Ctx.Done():
					return
				default:
				}
				// Eroare la nivel de outer → propagă și oprește
				if innerOrErr.Err != nil {
					_ = em.SendErr(innerOrErr.Err)
					return
				}
				// Procesează complet inner-ul curent
				if !emitAll(innerOrErr.Value) {
					return
				}
			}
		}()

		return stream.Stream[T]{Ctx: outer.Ctx, Ch: em.Out}
	}
}
