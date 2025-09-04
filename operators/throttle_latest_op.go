package operators

import (
	"time"

	"github.com/aurelsandu/StreamLang/operators/internal/clock"
	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

// ThrottleLatestOp: emite cel mult o valoare la fiecare d, folosind ULTIMUL element primit.
func ThrottleLatestOp[T any](d time.Duration) stream.Op[T, T] {
	return ThrottleLatestOpWithClock[T](d, clock.NewReal())
}

func ThrottleLatestOpWithClock[T any](d time.Duration, clk clock.Clock) stream.Op[T, T] {
	return func(in stream.Stream[T]) stream.Stream[T] {
		em := rt.NewEmitter[T](in.Ctx)

		go func() {
			defer em.Close()

			tk := clk.NewTicker(d)
			defer tk.Stop()

			var (
				latest stream.Item[T]
				has    bool
				open   = true
			)

			for open || has {
				select {
				case <-in.Ctx.Done():
					return

				case it, ok := <-in.Ch:
					if !ok {
						open = false
						continue
					}
					latest = it
					has = true

				case <-tk.C():
					if !has {
						continue
					}
					if latest.Err != nil {
						if !em.SendErr(latest.Err) {
							return
						}
						// dacă preferi stop la eroare, poți return aici
						has = false
						continue
					}
					if !em.SendVal(latest.Value) {
						return
					}
					has = false // „consumat” pentru fereastra curentă
				}
			}
		}()

		return stream.Stream[T]{Ctx: in.Ctx, Ch: em.Out}
	}
}
