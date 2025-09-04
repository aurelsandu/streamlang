package operators

import (
	"context"
	"reflect"
	"testing"

	"github.com/aurelsandu/StreamLang/sources"
//	"github.com/aurelsandu/StreamLang/stream"
)

func TestSkipWhile_Basic(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})

	out := SkipWhile(src, func(x int) bool { return x < 4 })

	var got []int
	for it := range out.Ch {
		if it.Err != nil { t.Fatalf("unexpected error: %v", it.Err) }
		got = append(got, it.Value)
	}

	want := []int{4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSkipWhile_AllSkipped(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3})

	out := SkipWhile(src, func(x int) bool { return x < 10 }) // toate < 10 → totul sărit

	got := []int{}
	for it := range out.Ch {
		if it.Err != nil { t.Fatalf("unexpected error: %v", it.Err) }
		got = append(got, it.Value)
	}

	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}
