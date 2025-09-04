package operators

import (
	"context"
	"testing"

	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestCompose_Map_Filter_Skip_Take(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// definim operatorii
	map10   := MapOp[int, int](func(x int) int { return x * 10 })
	filter20 := FilterOp[int](func(x int) bool { return x%20 == 0 })
	skip1   := SkipOp[int](1)
	take2   := TakeOp[int](2)

	// compunerea: Map → Filter → Skip → Take
	pipe := stream.Compose4(map10, filter20, skip1, take2)

	// aplicăm pipeline-ul pe sursă
	out := stream.PipeOp(src, pipe)

	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []int{40, 60}
	if len(got) != len(want) {
		t.Fatalf("got len=%d, want len=%d (got=%v)", len(got), len(want), got)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("at %d: got %v, want %v | full got=%v", i, got[i], want[i], got)
		}
	}
}