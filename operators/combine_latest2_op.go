package operators

import (
	"sync"

	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

// CombineLatest2Op: combină stream-ul de intrare (A) cu "other" (B) folosind f(A,B)->R.
// Emite după ce are cel puțin o valoare din A și B; apoi la orice nou A sau B.
// Erorile din A sau B sunt propagate ca Item{Err:...} (nu opresc definitiv decât dacă sink-ul decide).
func CombineLatest2Op[A, B, R any](other stream.Stream[B], f func(A, B) R) stream.Op[A, R] {
	return func(in stream.Stream[A]) stream.Stream[R] {
		em := rt.NewEmitter[R](in.Ctx)

		type tagA struct{ v stream.Item[A] }
		type tagB struct{ v stream.Item[B] }

		events := make(chan any, 16) // mic buffer ca să evităm lockstep
		var wg sync.WaitGroup
		wg.Add(2)

		// reader A
		go func() {
			defer wg.Done()
			for it := range in.Ch {
				select {
				case <-in.Ctx.Done():
					return
				case events <- tagA{v: it}:
				}
			}
		}()

		// reader B
		go func() {
			defer wg.Done()
			for it := range other.Ch {
				select {
				case <-in.Ctx.Done():
					return
				case events <- tagB{v: it}:
				}
			}
		}()

		// closer pt. events
		go func() {
			wg.Wait()
			close(events)
		}()

		// combiner + emitere
		go func() {
			defer em.Close()

			var (
				haveA, haveB bool
				lastA A
				lastB B
			)

			for ev := range events {
				select {
				case <-in.Ctx.Done():
					return
				default:
				}

				switch e := ev.(type) {
				case tagA:
					if e.v.Err != nil {
						if !em.SendErr(e.v.Err) { return }
						continue
					}
					lastA, haveA = e.v.Value, true

				case tagB:
					if e.v.Err != nil {
						if !em.SendErr(e.v.Err) { return }
						continue
					}
					lastB, haveB = e.v.Value, true
				}

				if haveA && haveB {
					if !em.SendVal(f(lastA, lastB)) {
						return
					}
				}
			}
		}()

		return stream.Stream[R]{Ctx: in.Ctx, Ch: em.Out}
	}
}