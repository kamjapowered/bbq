package geometry

import c "github.com/vistormu/go-dsa/constraints"

// store a line as a point plus direction vector
type Line[T c.Number] struct {
	Point     Vector[T]
	Direction Vector[T]
}

// create a line from two points
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
