package operators

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestCombineLatest2Op_Basic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// B emite imediat 10, apoi după 80ms 20
	b := makeTimedStream(ctx, []int{10, 20}, []time.Duration{0, 80 * time.Millisecond})

	// A emite 1,2,3 la 10ms distanță (după ce B a emis deja 10)
	a := makeTimedStream(ctx, []int{1, 2, 3}, []time.Duration{10 * time.Millisecond, 10 * time.Millisecond, 10 * time.Millisecond})

	pipe := CombineLatest2Op[int, int, string](b, func(x, y int) string {
		return fmt.Sprintf("%d|%d", x, y)
	})

	out := stream.PipeOp(a, pipe)
	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	// A: 1,2,3 ; B: 10,...20 mai târziu
	// așteptăm: 1|10, 2|10, 3|10, 3|20 (în această ordine)
	want := []string{"1|10", "2|10", "3|10", "3|20"}
	if len(got) != len(want) {
		t.Fatalf("got %d items (%v), want %d", len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("at %d: got %v, want %v | full=%v", i, got[i], want[i], got)
		}
	}
}