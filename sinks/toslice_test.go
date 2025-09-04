package sinks

import (
	"context"
	"reflect"
	"testing"

	"github.com/aurelsandu/StreamLang/sources"
)

func TestToSlice_CollectsAll(t *testing.T) {
	ctx := context.Background()
	src := sources.FromSlice(ctx, []int{1, 2, 3})

	got, err := ToSlice(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}


