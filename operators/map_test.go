package operators

import (
	"context"
	"reflect"
	"testing"

	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestMap_Basic(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3})

	out := Map(src, func(x int) int { return x * 10 })

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

func TestMap_PropagatesError(t *testing.T) {
	ctx := context.Background()

	// sursă cu o eroare în mijloc
	ch := make(chan stream.Item[int], 3)
	ch <- stream.Item[int]{Value: 1}
	ch <- stream.Item[int]{Err: assertErr{}}
	ch <- stream.Item[int]{Value: 3}
	close(ch)

	src := stream.Stream[int]{Ctx: ctx, Ch: ch}
	out := Map(src, func(x int) int { return x * 2 })

	seenErr := false
	for it := range out.Ch {
		if it.Err != nil {
			seenErr = true
			break
		}
	}
	if !seenErr {
		t.Fatalf("expected an error item to pass through")
	}
}

type assertErr struct{}
func (assertErr) Error() string { return "boom" }
