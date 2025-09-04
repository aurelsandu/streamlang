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

	// Factory: dă un stream nou la fiecare apel (cold source)
	srcFactory := func() stream.Stream[int] {
		return sources.FromSlice(ctx, []int{1, 2, 3, 4, 5, 6})
	}

	// Pipeline A: filtrează parele → *10 → take(3)
	onlyEven := operators.FilterOp[int](func(x int) bool { return x%2 == 0 }) // Op[int,int]
	mul10    := operators.MapOp[int, int](func(x int) int { return x * 10 })  // Op[int,int]
	take3    := operators.TakeOp[int](3)                                       // Op[int,int]  ← instanțiat

	pipeA := stream.Compose3(onlyEven, mul10, take3)

	// Pipeline B: skip(2) → *100
	skip2  := operators.SkipOp[int](2)                                         // Op[int,int]  ← instanțiat
	mul100 := operators.MapOp[int, int](func(x int) int { return x * 100 })    // Op[int,int]

	pipeB := stream.Compose(skip2, mul100)

	// Rulează A
	outA := stream.PipeOp(srcFactory(), pipeA)
	listA, _ := sinks.ToSlice(outA)
	fmt.Println("A:", listA) // așteptat: [20 40 60]

	// Rulează B (sursă proaspătă)
	outB := stream.PipeOp(srcFactory(), pipeB)
	listB, _ := sinks.ToSlice(outB)
	fmt.Println("B:", listB) // așteptat: [300 400 500 600]
}
