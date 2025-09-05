package stream

func ForEach[T any](s Stream[T], fn func(T)) error {
	for item := range s.Ch {
		if item.Err != nil {
			return item.Err
		}
		fn(item.Value)
	}
	return nil
}