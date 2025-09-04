package operators

import (
	"context"
	"reflect"
	"testing"

	"github.com/aurelsandu/StreamLang/sources"
)

func TestFilter_KeepEven(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5, 6})

	out := Filter(src, func(x int) bool { return x%2 == 0 })

	var got []int
	for it := range out.Ch {
		if it.Err != nil {
			t.Fatalf("unexpected error: %v", it.Err)
		}
		got = append(got, it.Value)
	}
	want := []int{2, 4, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
