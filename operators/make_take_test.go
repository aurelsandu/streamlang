package operators

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestMapIntToString_ThenTake(t *testing.T) {
	ctx := context.Background()

	// Sursă: 1,2,3,4
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4})

	// Op-uri: int->string (x*10 și apoi strconv) + Take primele 2 pe string
	mapIntToStr := MapOp[int, string](func(x int) string { return strconv.Itoa(x * 10) }) // Op[int,string]
	take2Str    := TakeOp[string](2)                                                       // Op[string,string]

	pipe := stream.Compose(mapIntToStr, take2Str) // Op[int,string]
	out  := stream.PipeOp(src, pipe)

	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"10", "20"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}