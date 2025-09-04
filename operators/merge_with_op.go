package operators

import (
	"sync"

	"github.com/aurelsandu/StreamLang/operators/internal/rt"
	"github.com/aurelsandu/StreamLang/stream"
)

type MergeOpts struct {
	// Dacă este true, la primul Item.Err!=nil emis de oricare sursă
	// se semnalează stop pentru toate reader-ele și operatorul se închide.
	StopOnError bool
}

// MergeWithOp: merge(in, others...) → Op[T,T]
func MergeWithOp[T any](others ...stream.Stream[T]) stream.Op[T, T] {
	return MergeWithOpWithOpts[T](MergeOpts{}, others...)
}

// MergeWithOpWithOpts: variantă cu opțiuni (ex. StopOnError).
func MergeWithOpWithOpts[T any](opts MergeOpts, others ...stream.Stream[T]) stream.Op[T, T] {
	return func(in stream.Stream[T]) stream.Stream[T] {
		em := rt.NewEmitter[T](in.Ctx)

		all := append([]stream.Stream[T]{in}, others...)

		stop := make(chan struct{})
		var once sync.Once
		stopAll := func() { once.Do(func() { close(stop) }) }

		var wg sync.WaitGroup
		wg.Add(len(all))

		startReader := func(s stream.Stream[T]) {
			go func() {
				defer wg.Done()
				for it := range s.Ch {
					select {
					case <-in.Ctx.Done():
						return
					case <-stop:
						return
					default:
					}

					// emit item
					if it.Err != nil {
						// propagă eroarea
						if !em.SendErr(it.Err) {
							return
						}
						if opts.StopOnError {
							stopAll()
							return
						}
						continue
					}
					if !em.SendVal(it.Value) {
						return
					}
				}
			}()
		}
		for _, s := range all {
			startReader(s)
		}

		go func() {
			wg.Wait()
			em.Close()
		}()

		return stream.Stream[T]{Ctx: in.Ctx, Ch: em.Out}
	}
}