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

func TestMap_Filter_Take_WithCompose3(t *testing.T) {
	ctx := context.Background()

	// sursă: 1..6
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5, 6})

	// Op-uri:
	// 1) int -> string (x * 10)
	mapIntToStr := MapOp[int, string](func(x int) string { return strconv.Itoa(x * 10) }) // Op[int,string]
	// 2) păstrăm doar valorile care nu sunt "30" și "50"
	filterStr := FilterOp[string](func(s string) bool { return s != "30" && s != "50" })   // Op[string,string]
	// 3) luăm primele 3 stringuri rămase
	take3 := TakeOp[string](3)                                                             // Op[string,string]

	// compunere: Map -> Filter -> Take
	pipe := stream.Compose3(mapIntToStr, filterStr, take3) // Op[int,string]

	out := stream.PipeOp(src, pipe)

	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// după map: ["10","20","30","40","50","60"]
	// după filter: ["10","20","40","60"]
	// după take3: ["10","20","40"]
	want := []string{"10", "20", "40"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}