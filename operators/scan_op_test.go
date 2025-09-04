package operators

import (
	"context"
	"reflect"
	"testing"

	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestScanOp_RunningSum(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4})

	scan := ScanOp[int, int](0, func(acc, x int) int { return acc + x })
	out := stream.PipeOp(src, scan)

	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	want := []int{1, 3, 6, 10}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}