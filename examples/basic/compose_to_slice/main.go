// examples/basic/compose_to_slice/main.go
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
	map10 := operators.MapOp[int, string](func(x int) string { return strconv.Itoa(x * 10) })
	prefix := operators.MapOp[string, string](func(s string) string { return "val=" + s })

	pipe := stream.Compose(map10, prefix)
	out := stream.PipeOp(src, pipe)

	list, err := sinks.ToSlice(out)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("slice:", list)
}