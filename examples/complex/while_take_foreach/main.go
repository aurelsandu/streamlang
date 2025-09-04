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

	// sursă: 2,4,6,1,2,3,7,8
	src := sources.FromSlice(ctx, []int{2,4,6,1,2,3,7,8})

	// ia cât timp x e par, apoi STOP (TakeWhile)
	takeWhileEven := operators.TakeWhileOp[int](func(x int) bool { return x%2 == 0 })
	// apoi ia doar primele 2 din ce a rămas după TakeWhile (de fapt tot even-run inițial)
	take2 := operators.TakeOp[int](3)

	pipeline := stream.Compose(takeWhileEven, take2)

	out := stream.PipeOp(src, pipeline)

	// ForEach ca sink
	if err := sinks.ForEach(out, func(v int) {
		fmt.Println("value:", v)
	}); err != nil {
		fmt.Println("err:", err)
	}
	// așteptat: value: 2 ; value: 4
}
