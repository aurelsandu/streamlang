package main

import (
	"context"
	"fmt"

	"github.com/aurelsandu/StreamLang/operators"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

func main() {
	ctx := context.Background()

	// Op-uri compozabile
	map10 := operators.MapOp[int, int](func(x int) int { return x * 10 }) // Op[int,int]
	skip2 := operators.SkipOp[int](2)                                     // Op[int,int]  ← instanțiere + apel

	// compunere: map -> skip
	pipe := stream.Compose(map10, skip2)

	// sursa și rularea pipeline-ului
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})
	out := stream.PipeOp(src, pipe)

	for it := range out.Ch {
		if it.Err != nil {
			fmt.Println("err:", it.Err)
			continue
		}
		fmt.Println("val=", it.Value)
	}
}