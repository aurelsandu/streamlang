package sinks

import "github.com/aurelsandu/StreamLang/stream"

// ForEach consumă stream-ul și aplică fn la fiecare valoare;
// returnează prima eroare întâlnită sau nil la final.
func ForEach[T any](s stream.Stream[T], fn func(T)) error {
	for {
		select {
		case <-s.Ctx.Done():
			return s.Ctx.Err()
		case it, ok := <-s.Ch:
			if !ok {
				return nil
			}
			if it.Err != nil {
				return it.Err
			}
			fn(it.Value)
		}
	}
}
