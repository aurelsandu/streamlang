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

	// A și B au dimensiuni diferite; Zip se oprește la lungimea minimă.
	a := sources.FromSlice(ctx, []int{1, 2, 3})
	b := sources.FromSlice(ctx, []string{"A", "B", "C", "D"})

	// 1) Zip clasic: emite perechi (First, Second)
	zipAB := operators.ZipOp[int, string](b) // Op[int, Pair[int,string]]
	outPairs := stream.PipeOp(a, zipAB)

	pairs, err := sinks.ToSlice(outPairs)
	if err != nil {
		panic(err)
	}
	fmt.Println("Zip -> Pair[int,string]:")
	for _, p := range pairs {
		fmt.Printf("  (%d, %s)\n", p.First, p.Second)
	}

	// 2) Compozitie: Zip apoi mapare la string pentru afișare
	//    Op[int, Pair] ∘ Op[Pair, string] => Op[int, string]
	asString := operators.MapOp[operators.Pair[int, string], string](
		func(p operators.Pair[int, string]) string {
			return fmt.Sprintf("%d|%s", p.First, p.Second)
		},
	)
	pipe := stream.Compose(zipAB, asString) // Op[int, string]
	outStrings := stream.PipeOp(sources.FromSlice(ctx, []int{10, 20, 30, 40}), pipe)

	strs, err := sinks.ToSlice(outStrings)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nZip → Map(fmt):")
	fmt.Println(" ", strs)

	// 3) Demonstrație de oprire: dacă unul se închide mai devreme, zip se oprește.
	shortA := sources.FromSlice(ctx, []int{7})
	outShort := stream.PipeOp(shortA, zipAB)
	shortPairs, _ := sinks.ToSlice(outShort)
	fmt.Println("\nZip cu A mai scurt (A=[7], B=[A,B,C,D]) → doar prima pereche:")
	for _, p := range shortPairs {
		fmt.Printf("  (%d, %s)\n", p.First, p.Second)
	}
}