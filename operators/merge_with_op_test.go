package operators

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestMergeWithOp_FanIn(t *testing.T) {
	ctx := context.Background()

	src1 := sources.FromSlice(ctx, []int{1,2,3})
	src2 := sources.FromSlice(ctx, []int{10,20})
	src3 := sources.FromSlice(ctx, []int{7})

	pipe := MergeWithOp[int](src2, src3) // Op[int,int]
	out  := stream.PipeOp(src1, pipe)

	got, err := sinks.ToSlice(out)
	if err != nil { t.Fatalf("unexpected err: %v", err) }

	// ordinea e nedeterministă; comparăm setul sortat
	sort.Ints(got)
	want := []int{1,2,3,7,10,20}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
