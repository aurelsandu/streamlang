package operators

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/stream"
)

// makeTimedStream emite valorile din vals cu pauza corespunzătoare din pauses înaintea FIECĂREI valori.
// Ultima valoare este urmată de închiderea canalului.

func TestDebounce_Trailing_EmitsLast(t *testing.T) {
	// Ideea: trimitem câteva valori rapid, apoi o perioadă de liniște.
	// Debounce trebuie să emită DOAR ultima valoare, după liniște.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Emitem:   (delay)  0ms   40ms   40ms     (apoi tăcere)
	// Valori:            1     2      3
	src := makeTimedStream(ctx, []int{1, 2, 3}, []time.Duration{0, 40 * time.Millisecond, 40 * time.Millisecond})

	// Debounce = 100ms  (după ultimul 3, așteaptă 100ms și emite 3)
	deb := Debounce(src, 100*time.Millisecond)

	got, err := sinks.ToSlice(deb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []int{3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDebounceOp_Composable_WithMap(t *testing.T) {
	// Map (x*10) -> Debounce(80ms). Trimitem 1..5 la 20ms; rezultatul trebuie să fie doar 50.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	src := makeTimedStream(ctx, []int{1, 2, 3, 4, 5}, []time.Duration{0, 20 * time.Millisecond, 20 * time.Millisecond, 20 * time.Millisecond, 20 * time.Millisecond})

	// Map *10
	map10 := MapOp[int, int](func(x int) int { return x * 10 })
	// Debounce op pe int
	debOp := DebounceOp[int](80 * time.Millisecond)

	pipe := stream.Compose(map10, debOp)
	out := stream.PipeOp(src, pipe)

	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []int{50}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}