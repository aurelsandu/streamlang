package operators

import (
	"context"
	"testing"

	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestFoldOp_Sum(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4})

	sumOp := FoldOp[int, int](0, func(acc, x int) int { return acc + x })
	out := stream.PipeOp(src, sumOp)

	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != 1 || got[0] != 10 {
		t.Fatalf("got=%v, want=[10]", got)
	}
}

func TestFoldOp_PropagatesError(t *testing.T) {
	ctx := context.Background()
	ch := make(chan stream.Item[int], 2)
	ch <- stream.Item[int]{Value: 1}
	ch <- stream.Item[int]{Err: errBoom{}} // errBoom e în testutil_errors_test.go
	close(ch)
	src := stream.Stream[int]{Ctx: ctx, Ch: ch}

	sumOp := FoldOp[int, int](0, func(acc, x int) int { return acc + x })
	out := stream.PipeOp(src, sumOp)

	_, err := sinks.ToSlice(out)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
