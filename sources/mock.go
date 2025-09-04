package sources

import (
	"context"
	"math/rand"
	"time"

	"github.com/aurelsandu/StreamLang/stream"
)

// MockQuoteSource: generează random walk pentru simbolurile date.
func MockQuoteSource(ctx context.Context, symbols []string, every time.Duration) stream.Stream[Quote] {
	if every <= 0 {
		every = 500 * time.Millisecond
	}
	// seed de bază
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	type state struct{ sym string; px float64 }
	st := make([]state, 0, len(symbols))
	for _, s := range symbols {
		if s == "" {
			continue
		}
		st = append(st, state{sym: s, px: 100 + r.Float64()*50})
	}

	ch := make(chan stream.Item[Quote])
	go func() {
		defer close(ch)
		t := time.NewTicker(every)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				now := time.Now()
				for i := range st {
					// mică variație
					delta := (r.Float64() - 0.5) * 0.8
					st[i].px += delta
					if st[i].px < 1 {
						st[i].px = 1
					}
					q := Quote{Symbol: st[i].sym, Price: st[i].px, Time: now}
					select {
					case <-ctx.Done():
						return
					case ch <- stream.Item[Quote]{Value: q}:
					}
				}
			}
		}
	}()
	return stream.Stream[Quote]{Ctx: ctx, Ch: ch}
}