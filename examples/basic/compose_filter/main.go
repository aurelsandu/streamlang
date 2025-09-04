package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aurelsandu/StreamLang/operators"
	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func main() {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})

	// int -> string (x*10)
	map10 := operators.MapOp[int, string](func(x int) string { return strconv.Itoa(x * 10) })
	// păstrează doar valorile diferite de "30"
	neq30 := operators.FilterOp[string](func(s string) bool { return s != "30" })

	// compunem: mai întâi map10, apoi filter
	pipeline := stream.Compose(map10, neq30)

	// aplicăm pe stream
	out := stream.PipeOp(src, pipeline)

	_ = sinks.ForEach(out, func(v string) { fmt.Println("=>", v) })
}