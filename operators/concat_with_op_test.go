// operators/concat_with_op_test.go
package operators

import (
	"context"
	"reflect"
	"testing"

	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestConcatWithOp_Order(t *testing.T) {
	ctx := context.Background()
	a := sources.FromSlice(ctx, []int{1,2})
	b := sources.FromSlice(ctx, []int{10})
	c := sources.FromSlice(ctx, []int{7,8})

	out := stream.PipeOp(a, ConcatWithOp[int](b, c))
	got, err := sinks.ToSlice(out)
	if err != nil { t.Fatalf("unexpected err: %v", err) }

	want := []int{1,2,10,7,8}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}