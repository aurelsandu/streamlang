package operators

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

// TestDebounceOp_Composable verifică integrarea DebounceOp într-un pipeline Compose
func TestDebounceOp_Composable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// sursa emite rapid 5 valori
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})

	// pipeline = map *10 → debounce 100ms
	pipe := stream.Compose(
		MapOp[int, int](func(x int) int { return x * 10 }),
		DebounceOp[int](100*time.Millisecond),
	)

	out := stream.PipeOp(src, pipe)
	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// doar ultima valoare (5*10=50) ar trebui să rămână
	want := []int{50}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}