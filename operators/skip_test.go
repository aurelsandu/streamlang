package operators

import (
	"context"
	"reflect"
	"testing"

	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

// errBoom comun pentru teste; dacă îl ai deja într-un alt fișier test util, păstrează o singură definiție.
//type errBoom struct{}
//func (errBoom) Error() string { return "boom" }

func TestSkip_Basic(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})

	// SkipOp: sare peste primele 2
	skip2 := SkipOp[int](2)
	out := stream.PipeOp(src, skip2)

	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	want := []int{3, 4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestSkip_PropagatesError(t *testing.T) {
	ctx := context.Background()
	ch := make(chan stream.Item[int], 3)
	ch <- stream.Item[int]{Value: 1}
	ch <- stream.Item[int]{Err: errBoom{}} // eroare în stream
	ch <- stream.Item[int]{Value: 3}
	close(ch)

	src := stream.Stream[int]{Ctx: ctx, Ch: ch}
	out := stream.PipeOp(src, SkipOp[int](2)) // sare peste primul (1), apoi vede eroarea

	_, err := sinks.ToSlice(out)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestSkip_Composable_WithMap(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})

	map10 := MapOp[int, int](func(x int) int { return x * 10 }) // 10,20,30,40,50
	skip2 := SkipOp[int](2)                                     //                 → 30,40,50
	pipe := stream.Compose(map10, skip2)
	out := stream.PipeOp(src, pipe)

	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	want := []int{30, 40, 50}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}