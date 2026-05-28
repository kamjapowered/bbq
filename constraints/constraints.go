package constraints

//microwave:export
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

//microwave:export
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

//microwave:export
type Integer interface {
	Signed | Unsigned
}

//microwave:export
type Float interface {
	~float32 | ~float64
}

//microwave:export
type Number interface {
	Integer | Float
}

//microwave:export
type Ordered interface {
	Integer | ~string | Float
}
