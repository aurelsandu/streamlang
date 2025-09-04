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

	// Compoziție prin apeluri succesive la stream.Pipe
	step1 := stream.Pipe(src, func(s stream.Stream[int]) stream.Stream[string] {
		return operators.Map(s, func(x int) string { return strconv.Itoa(x * 10) })
	})

	out := stream.Pipe(step1, func(s stream.Stream[string]) stream.Stream[string] {
		return operators.Map(s, func(v string) string { return "val=" + v })
	})

	_ = sinks.ForEach(out, func(v string) {
		fmt.Println(v)
	})
}
