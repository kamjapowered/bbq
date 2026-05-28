package functional

//microwave:export
type Option[T any] = func(*T)

//microwave:export
func ApplyOptions[T any](v *T, opts ...Option[T]) {
	if v == nil {
		return
	}

	for _, opt := range opts {
		if opt != nil {
			opt(v)
		}
	}
}
