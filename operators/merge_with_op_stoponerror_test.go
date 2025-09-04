package operators

import (
	"context"
	"testing"

	"github.com/aurelsandu/StreamLang/sinks"
	"github.com/aurelsandu/StreamLang/stream"
)

func TestMergeWithOp_StopOnError(t *testing.T) {
	ctx := context.Background()

	// src1: 1, err, 99 (99 n-ar trebui să mai apară dacă StopOnError)
	ch1 := make(chan stream.Item[int], 3)
	ch1 <- stream.Item[int]{Value: 1}
	ch1 <- stream.Item[int]{Err: errBoom{}}
	ch1 <- stream.Item[int]{Value: 99}
	close(ch1)
	src1 := stream.Stream[int]{Ctx: ctx, Ch: ch1}

	// src2: 10, 20, 30 (ar putea continua, dar ne oprim la eroare)
	ch2 := make(chan stream.Item[int], 3)
	ch2 <- stream.Item[int]{Value: 10}
	ch2 <- stream.Item[int]{Value: 20}
	ch2 <- stream.Item[int]{Value: 30}
	close(ch2)
	src2 := stream.Stream[int]{Ctx: ctx, Ch: ch2}

	out := stream.PipeOp(src1, MergeWithOpWithOpts[int](MergeOpts{StopOnError: true}, src2))
	// folosim ToSlice: va întoarce prima eroare întâlnită
	values, err := sinks.ToSlice(out)
	if err == nil {
		t.Fatalf("expected error, got nil (values=%v)", values)
	}
	// verificăm că 99 NU a ajuns (ne-am oprit la eroare)
	for _, v := range values {
		if v == 99 {
			t.Fatalf("unexpected value after error: 99 (values=%v)", values)
		}
	}
}