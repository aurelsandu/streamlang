// sinks/pipe_to_slice.go
package sinks

import "github.com/aurelsandu/StreamLang/stream"

// PipeToSlice aplică un Op pe s, apoi colectează rezultatul în []B.
func PipeToSlice[A any, B any](s stream.Stream[A], op stream.Op[A, B]) ([]B, error) {
	return ToSlice(op(s))
}
