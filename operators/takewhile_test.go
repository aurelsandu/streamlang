package operators

import (
	"context"
	"reflect"
	"testing"

	"github.com/aurelsandu/StreamLang/sources"
)

func TestTakeWhile_Basic(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})

	out := TakeWhile(src, func(x int) bool { return x < 4 })

	var got []int
	for it := range out.Ch {
		if it.Err != nil { t.Fatalf("unexpected error: %v", it.Err) }
		got = append(got, it.Value)
	}

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestTakeWhile_StopsEarly(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{5, 6, 7})

	out := TakeWhile(src, func(x int) bool { return x < 4 }) // niciun element

	got := []int{}
	for it := range out.Ch {
		if it.Err != nil { t.Fatalf("unexpected error: %v", it.Err) }
		got = append(got, it.Value)
	}

	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}