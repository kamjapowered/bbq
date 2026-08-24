package geometry

import (
	"iter"

	c "kamjapowered.com/bbq/constraints"
)

// store a polygon as points in local space
//
// interpret as 2d on the xy plane
//
//microwave:export
type Polygon[T c.Number] struct {
	Points []Vector[T]
}

// create an empty polygon
//
//microwave:export
func NewPolygon[T c.Number]() Polygon[T] {
	return Polygon[T]{Points: make([]Vector[T], 0)}
}

// add a point
func (p *Polygon[T]) Append(v Vector[T]) {
	p.Points = append(p.Points, v)
}

// clear all points
func (p *Polygon[T]) Clear() {
	p.Points = p.Points[:0]
}

// append a point built from x and y
func (p *Polygon[T]) AppendXY(x, y T) {
	p.Points = append(p.Points, Vector[T]{X: x, Y: y})
}

// iterate the points in insertion order
func (p *Polygon[T]) Iter() iter.Seq[Vector[T]] {
	return func(yield func(Vector[T]) bool) {
		for _, v := range p.Points {
			if !yield(v) {
				return
			}
		}
	}
}
