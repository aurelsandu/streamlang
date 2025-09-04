//package stream
package stream_test


import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/aurelsandu/StreamLang/operators"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestCompose_MapMap(t *testing.T) {
	ctx := context.Background()

	src := sources.FromSlice(ctx, []int{1, 2, 3})

	op1 := operators.MapOp[int, string](func(x int) string { return strconv.Itoa(x * 10) })
	op2 := operators.MapOp[string, string](func(s string) string { return "v=" + s })

	pipeline := stream.Compose(op1, op2)
	out := stream.PipeOp(src, pipeline)

	var got []string
	for it := range out.Ch {
		if it.Err != nil {
			t.Fatalf("unexpected error: %v", it.Err)
		}
		got = append(got, it.Value)
	}
	want := []string{"v=10", "v=20", "v=30"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}