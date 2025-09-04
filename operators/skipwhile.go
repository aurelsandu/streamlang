package operators

import "github.com/aurelsandu/StreamLang/stream"

// SkipWhile sare peste elementele de la început cât timp predicatul este true.
// Când găsește primul element pentru care predicatul este false, de acolo înainte
// emite toate elementele (inclusiv pe acela). Erorile sunt propagate și NU schimbă starea.
func SkipWhile[T any](s stream.Stream[T], pred func(T) bool) stream.Stream[T] {
	out := make(chan stream.Item[T])
	go func() {
		defer close(out)
		released := false
		for it := range s.Ch {
			if it.Err != nil {
				out <- it
				continue
			}
			if !released {
				if pred(it.Value) {
					// încă sărim peste
					continue
				}
				// primul care NU îndeplinește condiția: de aici dăm drumul la tot
				released = true
			}
			out <- it
		}
	}()
	return stream.Stream[T]{Ctx: s.Ctx, Ch: out}
}