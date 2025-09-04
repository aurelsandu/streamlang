package operators

import (
	"time"

	"github.com/aurelsandu/StreamLang/operators/internal/clock"
	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

// DebounceOp: emite DOAR ultimul element după o pauză de d fără noi elemente.
func DebounceOp[T any](d time.Duration) stream.Op[T, T] {
	return DebounceOpWithClock[T](d, clock.NewReal())
}

// DebounceOpWithClock: variantă testabilă (injecție de Clock).
func DebounceOpWithClock[T any](d time.Duration, clk clock.Clock) stream.Op[T, T] {
	return func(in stream.Stream[T]) stream.Stream[T] {
		em := rt.NewEmitter[T](in.Ctx)

		go func() {
			defer em.Close()

			var (
				pending *stream.Item[T]
				timer   clock.Timer // nil la început
			)

			fire := func() bool {
				if pending == nil {
					return true
				}
				ok := true
				if pending.Err != nil {
					ok = em.SendErr(pending.Err)
				} else {
					ok = em.SendVal(pending.Value)
				}
				pending = nil
				return ok
			}

			for {
				var timerC <-chan time.Time
				if timer != nil {
					timerC = timer.C()
				}

				select {
				case <-in.Ctx.Done():
					return

				case it, ok := <-in.Ch:
					if !ok {
						// sursa s-a închis → dacă avem un pending, îl emitem
						_ = fire()
						return
					}
					// memorăm ultimul it
					pending = &it

					// reset timer
					if timer == nil {
						timer = clk.NewTimer(d)
					} else {
						_ = timer.Stop()
						_ = timer.Reset(d)
					}

				case <-timerC:
					// pauză d -> emitem pending
					if !fire() {
						return
					}
				}
			}
		}()

		return stream.Stream[T]{Ctx: in.Ctx, Ch: em.Out}
	}
}