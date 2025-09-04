package operators

import (
	"context"
	"testing"
	"time"

	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/stream"
)

// newTestCtx returnează un context cu timeout generos pentru teste.
func newTestCtx(t testing.TB, d time.Duration) (context.Context, context.CancelFunc) {
	t.Helper()
	if d <= 0 {
		d = 2 * time.Second
	}
	return context.WithTimeout(context.Background(), d)
}

// fromVals creează rapid un Stream[T] dintr-un slice de valori (fără erori).
func fromVals[T any](ctx context.Context, vals ...T) stream.Stream[T] {
	ch := make(chan stream.Item[T], len(vals))
	for _, v := range vals {
		ch <- stream.Item[T]{Value: v}
	}
	close(ch)
	return stream.Stream[T]{Ctx: ctx, Ch: ch}
}

// fromItems creează un Stream[T] dintr-o listă de Item[T] (poți mixa valori + erori).
func fromItems[T any](ctx context.Context, items ...stream.Item[T]) stream.Stream[T] {
	ch := make(chan stream.Item[T], len(items))
	for _, it := range items {
		ch <- it
	}
	close(ch)
	return stream.Stream[T]{Ctx: ctx, Ch: ch}
}

// makeTimedStream emite valorile din vals cu o pauză (pauses[i]) înaintea FIECĂREI valori.
// Dacă pauses are mai puține elemente, restul sunt considerate 0.
func makeTimedStream[T any](ctx context.Context, vals []T, pauses []time.Duration) stream.Stream[T] {
	out := make(chan stream.Item[T])
	go func() {
		defer close(out)
		for i, v := range vals {
			var p time.Duration
			if i < len(pauses) {
				p = pauses[i]
			}
			if p > 0 {
				select {
				case <-ctx.Done():
					return
				case <-time.After(p):
				}
			}
			select {
			case <-ctx.Done():
				return
			case out <- stream.Item[T]{Value: v}:
			}
		}
	}()
	return stream.Stream[T]{Ctx: ctx, Ch: out}
}

// toSliceMust colectează toate valorile (ignoră erorile) și face t.Fatal dacă sink-ul întoarce eroare.
func toSliceMust[T any](t testing.TB, s stream.Stream[T]) []T {
	t.Helper()
	vals, err := sinks.ToSlice(s)
	if err != nil {
		t.Fatalf("unexpected error from sink: %v", err)
	}
	return vals
}

// readAll colectează separat valorile și erorile (fără să oprească la prima eroare).
func readAll[T any](s stream.Stream[T]) (values []T, errs []error) {
	for it := range s.Ch {
		if it.Err != nil {
			errs = append(errs, it.Err)
			continue
		}
		values = append(values, it.Value)
	}
	return
}

// mustEqual helper simplu pentru comparație de slices (fără reflect.DeepEqual, doar pentru teste mici).
func mustEqual[T comparable](t testing.TB, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%d want=%d | got=%v want=%v", len(got), len(want), got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("mismatch at %d: got=%v want=%v | gotSlice=%v wantSlice=%v", i, got[i], want[i], got, want)
		}
	}
}
