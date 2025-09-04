package operators

import "github.com/aurelsandu/StreamLang/stream"

// Filter păstrează doar elementele pentru care predicatul e true; erorile trec mai departe.
func Filter[T any](s stream.Stream[T], pred func(T) bool) stream.Stream[T] {
	out := make(chan stream.Item[T])
	go func() {
		defer close(out)
		for it := range s.Ch {
			if it.Err != nil {
				out <- it
				continue
			}
			if pred(it.Value) {
				out <- it
			}
		}
	}()
	return stream.Stream[T]{Ctx: s.Ctx, Ch: out}
}