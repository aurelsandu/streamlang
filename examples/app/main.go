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

	// sursă: un slice simplu
	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5})

	// 1) Map: multiplică fiecare element cu 10
	map10 := operators.MapOp[int, int](func(x int) int { return x * 10 })

	// 2) Filter: păstrează doar elementele > 20
	filter := operators.FilterOp[int](func(x int) bool { return x > 20 })

	// 3) Scan: sumă rulantă (running sum)
	scan := operators.ScanOp[int, int](0, func(acc, x int) int { return acc + x })

	// 4) Fold: suma totală (toate elementele)
	fold := operators.FoldOp[int, int](0, func(acc, x int) int { return acc + x })

	// === Pipeline A: Map + Filter + Scan ===
	pipeA := stream.Compose3(map10, filter, scan)
	outA := stream.PipeOp(src, pipeA)
	valsA, err := sinks.ToSlice(outA)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pipeline A (map→filter→scan):", valsA)
	// pentru [1,2,3,4,5]:
	// map10 => [10,20,30,40,50]
	// filter >20 => [30,40,50]
	// scan => [30,70,120]

	// === Pipeline B: Map + Fold ===
	pipeB := stream.Compose(map10, fold)
	outB := stream.PipeOp(src, pipeB)
	valsB, err := sinks.ToSlice(outB)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pipeline B (map→fold):", valsB)
	// map10 => [10,20,30,40,50]
	// fold sum => [150]
}