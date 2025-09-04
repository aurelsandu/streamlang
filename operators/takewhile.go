package operators

import "github.com/aurelsandu/StreamLang/stream"

// TakeWhile emite elemente cât timp predicatul este true.
// La primul element pentru care predicatul este false, închide ieșirea.
// Erorile sunt propagate și NU schimbă starea (nu "opresc" fluxul).
func TakeWhile[T any](s stream.Stream[T], pred func(T) bool) stream.Stream[T] {
	out := make(chan stream.Item[T])
	go func() {
		defer close(out)
		for it := range s.Ch {
			if it.Err != nil {
				out <- it
				continue
			}
			if !pred(it.Value) {
				return
			}
			out <- it
		}
	}()
	return stream.Stream[T]{Ctx: s.Ctx, Ch: out}
}