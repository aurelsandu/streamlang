/*
super—iată un ZipOp curat (binary, 2→1) implementat cu rt pentru ieșire sigură și ctx.Done() respectat. 
„Zipează” elementele în ordinea lor: primul din A cu primul din B, al doilea cu al doilea, etc. Se oprește când oricare dintre stream-uri se închide. 
Erorile sunt propagate imediat.
*/
package operators

import (
	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

// Pair este rezultatul pentru Zip(A,B) – echivalentul unui tuplu (A,B).
type Pair[A, B any] struct {
	First  A
	Second B
}

// ZipOp: aliniază elementele "în perechi": (a1,b1), (a2,b2), ...
// - Se oprește când unul dintre stream-uri se închide (comportamentul clasic de zip).
// - Dacă apare o eroare pe oricare intrare, o propagă și se oprește.
// - Respectă contextul (ctx.Done()) pentru închidere rapidă.
func ZipOp[A, B any](other stream.Stream[B]) stream.Op[A, Pair[A, B]] {
	return func(in stream.Stream[A]) stream.Stream[Pair[A, B]] {
		em := rt.NewEmitter[Pair[A, B]](in.Ctx)

		recvA := func() (it stream.Item[A], ok bool) {
			select {
			case <-in.Ctx.Done():
				return stream.Item[A]{}, false
			case it, ok = <-in.Ch:
				return it, ok
			}
		}
		recvB := func() (it stream.Item[B], ok bool) {
			select {
			case <-in.Ctx.Done():
				return stream.Item[B]{}, false
			case it, ok = <-other.Ch:
				return it, ok
			}
		}

		go func() {
			defer em.Close()
			for {
				ai, ok := recvA()
				if !ok {
					return // A s-a închis sau context anulat
				}
				if ai.Err != nil {
					_ = em.SendErr(ai.Err)
					return
				}

				bi, ok := recvB()
				if !ok {
					return // B s-a închis sau context anulat
				}
				if bi.Err != nil {
					_ = em.SendErr(bi.Err)
					return
				}

				if !em.SendVal(Pair[A, B]{First: ai.Value, Second: bi.Value}) {
					return
				}
			}
		}()

		return stream.Stream[Pair[A, B]]{Ctx: in.Ctx, Ch: em.Out}
	}
}