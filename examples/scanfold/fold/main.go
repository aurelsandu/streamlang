package main

import (
	"context"
	"fmt"

	"github.com/aurelsandu/StreamLang/operators"
	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func main() {
	ctx := context.Background()

	fold := operators.FoldOp[int, int](0, func(acc, x int) int { return acc + x })

	// sursă dintr-un slice
	src := sources.FromSlice(ctx, []int{1, 2, 3})

	// aplică operatorul
	out := stream.PipeOp(src, fold)

	// colectează rezultatele
	vals, err := sinks.ToSlice(out) // [6]
	if err != nil {
		panic(err)
	}
	fmt.Println(vals)
}