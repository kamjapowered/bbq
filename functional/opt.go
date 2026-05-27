package functional

type Option[T any] = func(*T)

func Apply[T any](v *T, opts ...Option[T]) {
	if v == nil {
		return
	}

	for _, opt := range opts {
		if opt != nil {
			opt(v)
		}
	}
}
