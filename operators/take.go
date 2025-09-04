package operators

import "github.com/aurelsandu/StreamLang/stream"

// Take emite primele n valori (fără erori). Erorile sunt propagate imediat,
// dar NU sunt numărate la limită. După ce a emis n valori, închide out.
func Take[T any](s stream.Stream[T], n int) stream.Stream[T] {
	out := make(chan stream.Item[T])
	if n <= 0 {
		close(out)
		return stream.Stream[T]{Ctx: s.Ctx, Ch: out}
	}
	go func() {
		defer close(out)
		count := 0
		for it := range s.Ch {
			// propagăm erorile imediat
			if it.Err != nil {
				out <- it
				continue
			}
			out <- it
			count++
			if count >= n {
				return
			}
		}
	}()
	return stream.Stream[T]{Ctx: s.Ctx, Ch: out}
}
