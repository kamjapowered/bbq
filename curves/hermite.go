package curves

import (
	c "github.com/vistormu/go-dsa/constraints"
	"github.com/vistormu/go-dsa/geometry"
)

// cubic hermite segment between p0 and p1 with tangents m0 and m1
//
//	p(u) = h00 p0 + h10 m0 + h01 p1 + h11 m1
//
// where the basis polynomials are
//
//	h00 =  2u^3 - 3u^2 + 1
//	h10 =   u^3 - 2u^2 + u
//	h01 = -2u^3 + 3u^2
//	h11 =   u^3 -  u^2
type Hermite[T c.Float] struct {
	p0, p1 geometry.Vector[T]
	m0, m1 geometry.Vector[T]
}

// create a hermite segment from endpoints p0 and p1 and tangents m0 and m1
//
// time: O(1)
func NewHermite[T c.Float](p0, p1, m0, m1 geometry.Vector[T]) Hermite[T] {
	return Hermite[T]{p0: p0, p1: p1, m0: m0, m1: m1}
}

// compute the point at parameter u in [0, 1]
//
// time: O(1)
func (h Hermite[T]) Compute(u T) geometry.Vector[T] {
	u2 := u * u
	u3 := u2 * u

	h00 := 2*u3 - 3*u2 + 1
	h10 := u3 - 2*u2 + u
	h01 := -2*u3 + 3*u2
	h11 := u3 - u2

	return h.p0.Scale(h00).
		Add(h.m0.Scale(h10)).
		Add(h.p1.Scale(h01)).
		Add(h.m1.Scale(h11))
}
