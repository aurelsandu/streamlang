// examples/basic/compose_take_slice/main.go
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
	
	map10  := operators.MapOp[int, string](func(x int) string { return strconv.Itoa(x * 10) })
	prefix := operators.MapOp[string, string](func(s string) string { return "val=" + s })
	take3  := operators.TakeOp[string](3) // IMPORTANT: instanțiat și apelat

	// Compose: map -> prefix -> take(3)
	pipeline := stream.Compose3(map10, prefix, take3)

	// Variante de colectare
	out, err := sinks.ToSlice(stream.PipeOp(src, pipeline))
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("ToSlice:", out)

	src2 := sources.FromSlice(ctx, []int{1,2,3,4,5})
	list2, err := sinks.PipeToSlice(src2, pipeline)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("PipeToSlice:", list2)
}