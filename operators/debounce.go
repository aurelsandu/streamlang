/*
Operatorul debounce este foarte util în procesarea streamurilor de evenimente — întârzie emiterea unui element până când a trecut o perioadă de timp fără ca altul nou să apară. 
E folosit frecvent pentru a evita suprasolicitarea sistemului cu evenimente foarte dese (ex: tastări, clickuri, etc).

Ce va face operatorul debounce în StreamLang:

- Primește un duration (în milisecunde).
- Dacă un nou eveniment vine înainte ca perioada să expire, îl înlocuiește pe cel anterior.
- Dacă perioada expiră fără alt eveniment, atunci evenimentul este trimis mai departe.
*/
package operators

import (
	"time"

	"github.com/aurelsandu/StreamLang/stream"
)

// Debounce (trailing): emite ultima valoare după ce nu s-au mai primit valori timp de d.
// Erorile sunt propagate imediat; la închiderea sursei, dacă există pending, se emite.
func Debounce[T any](s stream.Stream[T], d time.Duration) stream.Stream[T] {
	out := make(chan stream.Item[T])

	go func() {
		defer close(out)

		var (
			timer      *time.Timer
			timerC     <-chan time.Time
			pending    stream.Item[T]
			hasPending bool
		)

		stopTimer := func() {
			if timer == nil {
				return
			}
			if !timer.Stop() {
				// golim dacă e deja pornit și a declanșat
				select { case <-timer.C: default: }
			}
			timer = nil
			timerC = nil
		}

		armTimer := func() {
			stopTimer()
			timer = time.NewTimer(d)
			timerC = timer.C
		}

		flush := func() {
			if hasPending {
				select {
				case <-s.Ctx.Done():
					return
				case out <- pending:
				}
				hasPending = false
			}
			stopTimer()
		}

		for {
			select {
			case <-s.Ctx.Done():
				return

			case it, ok := <-s.Ch:
				if !ok {
					// sursa s-a închis: emitem ce e în pending (dacă există)
					flush()
					return
				}
				if it.Err != nil {
					// dacă vine o eroare, emitem întâi pending (ca să nu pierdem ultima valoare),
					// apoi propagăm eroarea imediat.
					flush()
					select {
					case <-s.Ctx.Done():
						return
					case out <- stream.Item[T]{Err: it.Err}:
					}
					continue
				}

				// memorăm ultima valoare și resetăm timer-ul
				pending, hasPending = it, true
				armTimer()

			case <-timerC:
				flush()
			}
		}
	}()

	return stream.Stream[T]{Ctx: s.Ctx, Ch: out}
}

// DebounceOp: adaptor pentru compoziție (stream.Op[T,T]).
/*
func DebounceOp[T any](d time.Duration) stream.Op[T, T] {
	return func(s stream.Stream[T]) stream.Stream[T] {
		return Debounce(s, d)
	}
}
*/