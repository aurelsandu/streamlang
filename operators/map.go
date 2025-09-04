package operators

import "github.com/aurelsandu/StreamLang/stream"

// Map aplică fn pe fiecare valoare; erorile se propagă neschimbate.
func Map[A any, B any](s stream.Stream[A], fn func(A) B) stream.Stream[B] {
	out := make(chan stream.Item[B])

	go func() {
		defer close(out)
		for it := range s.Ch {
			if it.Err != nil {
				out <- stream.Item[B]{Err: it.Err}
				continue
			}
			out <- stream.Item[B]{Value: fn(it.Value)}
		}
	}()

	return stream.Stream[B]{Ctx: s.Ctx, Ch: out}
}