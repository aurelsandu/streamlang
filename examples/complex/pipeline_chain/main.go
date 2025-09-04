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

	src := sources.FromSlice(ctx, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	map10   := operators.MapOp[int, int](func(x int) int { return x * 10 })                // Op[int,int]
	asStr   := operators.MapOp[int, string](func(x int) string { return strconv.Itoa(x) }) // Op[int,string]
	no30    := operators.FilterOp[string](func(s string) bool { return s != "30" })        // Op[string,string]
	skipLt5 := operators.SkipWhileOp[string](func(s string) bool { return s < "50" })      // Op[string,string]
	take3   := operators.TakeOp[string](3)                               // Op[string,string]  ← IMPORTANT

	stage1   := stream.Compose(map10, asStr)                  // Op[int,string]
	stage2   := stream.Compose3(no30, skipLt5, take3)         // Op[string,string]
	pipeline := stream.Compose(stage1, stage2)                // Op[int,string]

	out := stream.PipeOp(src, pipeline)
	list, err := sinks.ToSlice(out)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("result:", list) // așteptat: ["50","60","70"]
}
/*
Mini-cheatsheet (memorare rapidă)

MapOp[A,B] → Op[A,B]

FilterOp[T], SkipWhileOp[T], TakeOp[T] → Op[T,T]

Compose(a: Op[A,B], b: Op[B,C]) → Op[A,C]

PipeOp(s: Stream[A], op: Op[A,B]) → Stream[B]

Instanțiere generici: XxxOp[Tip](args)
ex: TakeOp, SkipOp
*/