package operators

import (
	"sync"

	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

// FlattenOp: primește Stream[Stream[T]] și emite elementele din toate inner-urile (fan-in).
// Se închide când outer-ul s-a închis și TOATE inner-urile au terminat.
// Erorile inner-urilor sunt propagate ca Item{Err:...}.
func FlattenOp[T any]() stream.Op[stream.Stream[T], T] {
	return func(outer stream.Stream[stream.Stream[T]]) stream.Stream[T] {
		em := rt.NewEmitter[T](outer.Ctx)

		var wg sync.WaitGroup

		// dispatcher pentru inner streams
		go func() {
			defer func() {
				wg.Wait()
				em.Close()
			}()

			for s := range outer.Ch {
				select {
				case <-outer.Ctx.Done():
					return
				default:
				}

				// Dacă elementul outer e eroare, propagate și continuă
				if s.Err != nil {
					if !em.SendErr(s.Err) {
						return
					}
					continue
				}

				inner := s.Value
				wg.Add(1)
				go func(inS stream.Stream[T]) {
					defer wg.Done()
					for it := range inS.Ch {
						select {
						case <-outer.Ctx.Done():
							return
						default:
						}
						if it.Err != nil {
							if !em.SendErr(it.Err) {
								return
							}
							continue
						}
						if !em.SendVal(it.Value) {
							return
						}
					}
				}(inner)
			}
		}()

		return stream.Stream[T]{Ctx: outer.Ctx, Ch: em.Out}
	}
}