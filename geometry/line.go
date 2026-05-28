package geometry

import c "github.com/kamjapowered/bbq/constraints"

// store a line as a point plus direction vector
//
//microwave:export
type Line[T c.Number] struct {
	Point     Vector[T]
	Direction Vector[T]
}

// create a line from two points
//
//microwave:export
func NewLine[T c.Number](p0, p1 Vector[T]) Line[T] {
	return Line[T]{
		Point:     p0,
		Direction: p1.Sub(p0),
	}
}

// return a point at parameter t
func (l Line[T]) At(t T) Vector[T] {
	return l.Point.Add(l.Direction.Scale(t))
}
