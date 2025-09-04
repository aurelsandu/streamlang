package e2e

import (
	"context"
	"reflect"
	"testing"

	"github.com/aurelsandu/StreamLang/operators"
	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestCompose_Map_Map_Take_ToSlice(t *testing.T) {
	ctx := context.Background()

	// Pipeline: int -> string (x*10) -> string ("val=" + s) -> take(3)
	map10  := operators.MapOp[int, string](func(x int) string { return itoa(x * 10) })
	prefix := operators.MapOp[string, string](func(s string) string { return "val=" + s })
	take3  := operators.TakeOp[string](3) // ← instanțierea genericii + apel

	pipeline := stream.Compose3(map10, prefix, take3)

	// 1) Prima execuție
	src1 := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})
	out1, err := sinks.ToSlice(stream.PipeOp(src1, pipeline))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"val=10", "val=20", "val=30"}
	if !reflect.DeepEqual(out1, want) {
		t.Fatalf("unexpected result.\n got: %v\nwant: %v", out1, want)
	}

	// 2) Re-executăm pe o sursă recreată (stream-urile sunt single-use)
	src2 := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})
	out2, err := sinks.ToSlice(stream.PipeOp(src2, pipeline))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(out2, want) {
		t.Fatalf("unexpected result on rerun.\n got: %v\nwant: %v", out2, want)
	}
}

// itoa: helper local ca să evităm strconv în test
func itoa(x int) string {
	if x == 0 {
		return "0"
	}
	sign := ""
	if x < 0 {
		sign = "-"
		x = -x
	}
	buf := [20]byte{}
	i := len(buf)
	for x > 0 {
		i--
		buf[i] = byte('0' + x%10)
		x /= 10
	}
	return sign + string(buf[i:])
}
