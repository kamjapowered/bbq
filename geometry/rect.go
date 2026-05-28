package geometry

import c "github.com/kamjapowered/bbq/constraints"

// store a local axis aligned rectangle size
//
// interpret as 2d on the xy plane with center at origin
//
//microwave:export
type Rect[T c.Number] struct {
	Width  T
	Height T
}

// create a rectangle size
//
//microwave:export
func NewRect[T c.Number](w, h T) Rect[T] {
	return Rect[T]{Width: w, Height: h}
}

// create a square size
//
//microwave:export
func NewSquare[T c.Number](s T) Rect[T] {
	return Rect[T]{Width: s, Height: s}
}

// resize the rectangle
func (r *Rect[T]) Resize(w, h T) {
	r.Width = w
	r.Height = h
}

// return half width
func (r Rect[T]) HalfWidth() T {
	return r.Width / 2
}

// return half height
func (r Rect[T]) HalfHeight() T {
	return r.Height / 2
}
