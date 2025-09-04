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

	// Op1: int -> string (înmulțit cu 10 și convertit)
	op1 := operators.MapOp[int, string](func(x int) string { return strconv.Itoa(x * 10) })

	// Op2: string -> string (prefix)
	op2 := operators.MapOp[string, string](func(s string) string { return "val=" + s })

	// Compose: op1 apoi op2 => Op[int, string]
	pipeline := stream.Compose(op1, op2)

	// Aplicare (PipeOp) pe stream-ul sursă
	out := stream.PipeOp(src, pipeline)

	_ = sinks.ForEach(out, func(v string) { fmt.Println(v) })
}
