package geometry

import (
	"iter"

	c "github.com/vistormu/go-dsa/constraints"
)

// store a path as ordered points in local space
//
// interpret as 2d on the xy plane
type Path[T c.Number] struct {
	Points []Vector[T]
	Closed bool
}

// create an empty path
func NewPath[T c.Number]() Path[T] {
	return Path[T]{Points: make([]Vector[T], 0)}
}

// add a point
func (p *Path[T]) Add(v Vector[T]) {
	p.Points = append(p.Points, v)
}

// add a point from x and y with z set to zero
func (p *Path[T]) AddXY(x, y T) {
	p.Points = append(p.Points, Vector[T]{X: x, Y: y})
}

// clear all points
func (p *Path[T]) Clear() {
	p.Points = p.Points[:0]
}

// iterate over points in order
func (p *Path[T]) Iter() iter.Seq[Vector[T]] {
	return func(yield func(Vector[T]) bool) {
		for _, v := range p.Points {
			if !yield(v) {
				return
			}
		}
	}
}
