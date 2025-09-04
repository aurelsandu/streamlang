package operators

import (
	"context"
	"reflect"
	"testing"
	"fmt"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestTake_LimitsResults(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})

	out := Take(src, 3)

	var got []int
	for it := range out.Ch {
		if it.Err != nil {
			t.Fatalf("unexpected error: %v", it.Err)
		}
		got = append(got, it.Value)
	}

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
func TestTakeOp_Composable(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{10, 20, 30, 40})

	take_op := TakeOp[int](3)
	out := stream.PipeOp(src, take_op )


	var got []int
	for it := range out.Ch {
		if it.Err != nil {
			t.Fatalf("unexpected error: %v", it.Err)
		}
		got = append(got, it.Value)
	}

	want := []int{10, 20, 30}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
func TestTake_ZeroOrNegative(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3})
	out := Take(src, 0)

	var got []int
	for it := range out.Ch {
		if it.Err != nil {
			t.Fatalf("unexpected error: %v", it.Err)
		}
		got = append(got, it.Value)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}


func TestTake_PropagatesError_NotCounted(t *testing.T) {
	ctx := context.Background()

	ch := make(chan stream.Item[int], 5)
	ch <- stream.Item[int]{Value: 1}
	ch <- stream.Item[int]{Err: errBoom{}}
	ch <- stream.Item[int]{Value: 2}
	ch <- stream.Item[int]{Value: 3}
	close(ch)

	src := stream.Stream[int]{Ctx: ctx, Ch: ch}
	out := Take(src, 2) // trebuie să emită 1 și 2 (eroarea nu se numără)

	var vals []int
	seenErr := false
	for it := range out.Ch {
		if it.Err != nil {
			seenErr = true
			continue
		}
		vals = append(vals, it.Value)
	}

	if !seenErr {
		t.Fatalf("expected error to pass through")
	}
	want := []int{1, 2}
	if fmt.Sprint(vals) != fmt.Sprint(want) {
		t.Fatalf("got %v, want %v", vals, want)
	}
}