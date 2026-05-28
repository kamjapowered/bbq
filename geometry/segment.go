package geometry

import c "github.com/kamjapowered/bbq/constraints"

// store a segment as start and end points
//
//microwave:export
type Segment[T c.Number] struct {
	Start Vector[T]
	End   Vector[T]
}

// create a segment from start and end
//
//microwave:export
func NewSegment[T c.Number](start, end Vector[T]) Segment[T] {
	return Segment[T]{Start: start, End: end}
}

// return the direction vector end-start
func (s Segment[T]) Direction() Vector[T] {
	return s.End.Sub(s.Start)
}
