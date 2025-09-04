package sinks

import "github.com/aurelsandu/StreamLang/stream"

// ToSlice colectează toate valorile într-un slice și returnează prima eroare întâlnită (dacă există).
func ToSlice[T any](s stream.Stream[T]) ([]T, error) {
	var out []T
	for {
		select {
		case <-s.Ctx.Done():
			return out, s.Ctx.Err()
		case item, ok := <-s.Ch:
			if !ok {
				return out, nil
			}
			if item.Err != nil {
				return out, item.Err
			}
			out = append(out, item.Value)
		}
	}
}
