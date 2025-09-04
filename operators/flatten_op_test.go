package operators

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestFlattenOp_FanIn(t *testing.T) {
	ctx := context.Background()
	s1 := sources.FromSlice(ctx, []int{1, 3})
	s2 := sources.FromSlice(ctx, []int{2})
	s3 := sources.FromSlice(ctx, []int{7, 8})

	// outer: Stream[Stream[int]]
	outerCh := make(chan stream.Item[stream.Stream[int]], 3)
	outerCh <- stream.Item[stream.Stream[int]]{Value: s1}
	outerCh <- stream.Item[stream.Stream[int]]{Value: s2}
	outerCh <- stream.Item[stream.Stream[int]]{Value: s3}
	close(outerCh)
	outer := stream.Stream[stream.Stream[int]]{Ctx: ctx, Ch: outerCh}

	out := FlattenOp[int]()(outer)
	got, err := sinks.ToSlice(out)
	if err != nil { t.Fatalf("unexpected err: %v", err) }

	sort.Ints(got)
	want := []int{1, 2, 3, 7, 8}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}