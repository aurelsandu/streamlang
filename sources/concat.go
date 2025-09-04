package sources

import (
	"context"

	"github.com/aurelsandu/StreamLang/stream"
)

// Concat emite toate elementele din primul stream, apoi din al doilea, etc.
// Se oprește dacă se anulează ctx.
func Concat[T any](ctx context.Context, streams ...stream.Stream[T]) stream.Stream[T] {
	out := make(chan stream.Item[T])
	go func() {
		defer close(out)
		for _, s := range streams {
			for it := range s.Ch {
				select {
				case <-ctx.Done():
					return
				case out <- it:
				}
			}
		}
	}()
	return stream.Stream[T]{Ctx: ctx, Ch: out}
}