package sources

import (
	"context"

	"github.com/aurelsandu/StreamLang/stream"
)

// FromSlice emite elementele din slice în ordinea lor.
func FromSlice[T any](ctx context.Context, data []T) stream.Stream[T] {
	out := make(chan stream.Item[T])
	go func() {
		defer close(out)
		for _, v := range data {
			select {
			case <-ctx.Done():
				return
			case out <- stream.Item[T]{Value: v}:
			}
		}
	}()
	return stream.Stream[T]{Ctx: ctx, Ch: out}
}
