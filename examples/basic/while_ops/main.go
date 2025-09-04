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

	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})

	// TakeWhile: ia până când x < 4 (=> 1,2,3)
	takeWhile := operators.TakeWhileOp[int](func(x int) bool { return x < 4 })
	out1 := stream.PipeOp(src, takeWhile)
	fmt.Println("TakeWhile(x<4):")
	for it := range out1.Ch {
		if it.Err != nil {
			fmt.Println("err:", it.Err)
			continue
		}
		fmt.Println(" ", it.Value)
	}

	// Re-creăm sursa (stream-urile sunt single-use)
	src2 := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})

	// SkipWhile: sare peste cât timp x < 4, apoi dă drumul la restul (=> 4,5)
	skipWhile := operators.SkipWhileOp[int](func(x int) bool { return x < 4 })
	out2 := stream.PipeOp(src2, skipWhile)
	fmt.Println("SkipWhile(x<4):")
	for it := range out2.Ch {
		if it.Err != nil {
			fmt.Println("err:", it.Err)
			continue
		}
		fmt.Println(" ", it.Value)
	}
}