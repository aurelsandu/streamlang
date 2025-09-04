package operators

import (
	"context"
	"reflect"
	"testing"
	"strconv"
	
	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/sources"
	"github.com/aurelsandu/StreamLang/stream"
)

// helper local ca să evităm conflicte cu alte fișiere de test
//type errBoom struct{}
//func (errBoom) Error() string { return "boom_zip" }

// mic utilitar pentru a crea un stream dintr-o listă de Item-uri
func streamFromItems[T any](items ...stream.Item[T]) stream.Stream[T] {
	ch := make(chan stream.Item[T], len(items))
	for _, it := range items {
		ch <- it
	}
	close(ch)
	return stream.Stream[T]{Ctx: context.Background(), Ch: ch}
}

func TestZipOp_BasicPairs(t *testing.T) {
	ctx := context.Background()
	a := sources.FromSlice(ctx, []int{1, 2, 3})
	b := sources.FromSlice(ctx, []string{"A", "B", "C", "D"}) // mai lung

	zipAB := ZipOp[int, string](b) // Op[int, Pair[int,string]]
	out := stream.PipeOp(a, zipAB)

	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Pair[int, string]{
		{First: 1, Second: "A"},
		{First: 2, Second: "B"},
		{First: 3, Second: "C"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestZipOp_StopsAtShortest(t *testing.T) {
	ctx := context.Background()
	a := sources.FromSlice(ctx, []int{10, 20, 30, 40})
	b := sources.FromSlice(ctx, []int{1}) // mult mai scurt

	zipAB := ZipOp[int, int](b)
	out := stream.PipeOp(a, zipAB)

	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []Pair[int, int]{
		{First: 10, Second: 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestZipOp_ErrorFromA_FirstItem(t *testing.T) {
	// A pornește cu eroare -> zip propagă eroarea imediat, fără valori emise
	a := streamFromItems[int](stream.Item[int]{Err: errBoom{}})
	b := sources.FromSlice(context.Background(), []int{1, 2, 3})

	zipAB := ZipOp[int, int](b)
	out := stream.PipeOp(a, zipAB)

	got, err := sinks.ToSlice(out)
	if err == nil {
		t.Fatalf("expected error, got nil; got=%v", got)
	}
}

func TestZipOp_ErrorFromB_FirstItem(t *testing.T) {
	// B pornește cu eroare -> zip propagă eroarea imediat, fără valori emise
	a := sources.FromSlice(context.Background(), []int{1, 2, 3})
	b := streamFromItems[int](stream.Item[int]{Err: errBoom{}})

	zipAB := ZipOp[int, int](b)
	out := stream.PipeOp(a, zipAB)

	got, err := sinks.ToSlice(out)
	if err == nil {
		t.Fatalf("expected error, got nil; got=%v", got)
	}
}

func TestZipOp_ComposeWithMap(t *testing.T) {
	ctx := context.Background()
	a := sources.FromSlice(ctx, []int{1, 2, 3})
	b := sources.FromSlice(ctx, []string{"x", "y", "z"})

	zipAB := ZipOp[int, string](b) // Op[int, Pair[int,string]]
	toStr := MapOp[Pair[int, string], string](func(p Pair[int, string]) string {
			return p.Second + "|" + strconv.Itoa(p.First)})
	
	pipe := stream.Compose(zipAB, toStr) // Op[int, string]
	out := stream.PipeOp(a, pipe)

	got, err := sinks.ToSlice(out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"x|1", "y|2", "z|3"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
