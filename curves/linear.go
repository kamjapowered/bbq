package curves

import (
	c "github.com/vistormu/go-dsa/constraints"
	"github.com/vistormu/go-dsa/geometry"
)

// linear interpolation between two points in space
//
//	p(u) = (1 - u) * p0 + u * p1
type Linear[T c.Float] struct {
	p0, p1 geometry.Vector[T]
}

// create a linear segment from p0 to p1
//
// time: O(1)
func NewLinear[T c.Float](p0, p1 geometry.Vector[T]) Linear[T] {
	return Linear[T]{p0: p0, p1: p1}
}

// compute the point at parameter u in [0, 1]
//
// time: O(1)
func (l Linear[T]) Compute(u T) geometry.Vector[T] {
	return l.p0.Scale(1 - u).Add(l.p1.Scale(u))
}
