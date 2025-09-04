package main

import (
    "context"
	"fmt"
	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
//	"github.com/aurelsandu/StreamLang/stream"
)

func main(){
    ctx := context.Background()
    a := sources.FromSlice(ctx, []int{1, 3, 5})
    b := sources.FromSlice(ctx, []int{2, 4, 6})

    merged := sources.Merge(ctx, a, b) // Stream[int]

	// colectăm totul într-un slice
	vals, err := sinks.ToSlice(merged)
	if err != nil {
		panic(err)
	}

	fmt.Println("merged:", vals)
}