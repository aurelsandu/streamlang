// sources/merge.go
package sources

import (
	"context"
	"sync"

	"github.com/aurelsandu/StreamLang/stream"
)

func Merge[T any](ctx context.Context, streams ...stream.Stream[T]) stream.Stream[T] {
	out := make(chan stream.Item[T])
	var wg sync.WaitGroup
	wg.Add(len(streams))

	for _, s := range streams {
		s := s // capture
		go func() {
			defer wg.Done()
			for it := range s.Ch {
				select {
				case <-ctx.Done():
					return
				case out <- it:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return stream.Stream[T]{Ctx: ctx, Ch: out}
}