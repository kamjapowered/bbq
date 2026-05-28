package curves

import (
	c "github.com/kamjapowered/bbq/constraints"
	"github.com/kamjapowered/bbq/geometry"
)

// bezier curve of arbitrary degree evaluated with de casteljau
//
// the curve passes through the first and last control points only
type Bezier[T c.Float] struct {
	points []geometry.Vector[T]
}

// create a bezier curve from n control points
//
// the degree of the curve is n - 1
//
// time: O(n)
func NewBezier[T c.Float](points ...geometry.Vector[T]) Bezier[T] {
	cp := make([]geometry.Vector[T], len(points))
	copy(cp, points)
	return Bezier[T]{points: cp}
}

// compute the point at parameter u in [0, 1] using de casteljau
//
// return zero vector if there are no control points
//
// time: O(n^2) where n is the number of control points
func (b Bezier[T]) Compute(u T) geometry.Vector[T] {
	n := len(b.points)
	if n == 0 {
		return geometry.Vector[T]{}
	}

	buf := make([]geometry.Vector[T], n)
	copy(buf, b.points)

	for k := 1; k < n; k++ {
		for i := 0; i < n-k; i++ {
			buf[i] = buf[i].Scale(1 - u).Add(buf[i+1].Scale(u))
		}
	}

	return buf[0]
}
