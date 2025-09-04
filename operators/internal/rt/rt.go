package rt

import (
	"context"

	"github.com/aurelsandu/StreamLang/stream"
)

type Emitter[R any] struct {
	Ctx context.Context
	Out chan stream.Item[R]
}

func NewEmitter[R any](ctx context.Context) Emitter[R] {
	return Emitter[R]{Ctx: ctx, Out: make(chan stream.Item[R])}
}
func (e Emitter[R]) Close() { close(e.Out) }

func (e Emitter[R]) SendVal(v R) bool {
	select {
	case <-e.Ctx.Done():
		return false
	case e.Out <- stream.Item[R]{Value: v}:
		return true
	}
}
func (e Emitter[R]) SendErr(err error) bool {
	select {
	case <-e.Ctx.Done():
		return false
	case e.Out <- stream.Item[R]{Err: err}:
		return true
	}
}

// Unary rulează un operator unar generic: citește din in, apelează body(it, em).
// Dacă body returnează false, se oprește.
func Unary[T, R any](in stream.Stream[T], body func(it stream.Item[T], em *Emitter[R]) bool) stream.Stream[R] {
	em := NewEmitter[R](in.Ctx)
	go func() {
		defer em.Close()
		for it := range in.Ch {
			if !body(it, &em) {
				return
			}
		}
	}()
	return stream.Stream[R]{Ctx: in.Ctx, Ch: em.Out}
}