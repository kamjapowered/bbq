package curves

import (
	"math"

	c "github.com/vistormu/go-dsa/constraints"
	"github.com/vistormu/go-dsa/geometry"
)

// uniform catmull rom spline through n points, n >= 4
//
// the curve interpolates p1 through p_{n-2}, with p0 and p_{n-1} acting as tangent hints
//
// the parameter u in [0, 1] spans the n - 3 internal segments
type CatmullRom[T c.Float] struct {
	points []geometry.Vector[T]
}

// create a catmull rom spline from n control points
//
// points must contain at least 4 entries, otherwise the spline degenerates to a single point
//
// time: O(n)
func NewCatmullRom[T c.Float](points ...geometry.Vector[T]) CatmullRom[T] {
	cp := make([]geometry.Vector[T], len(points))
	copy(cp, points)
	return CatmullRom[T]{points: cp}
}

// compute the point at parameter u in [0, 1]
//
// return zero vector if there are fewer than 4 control points
//
// time: O(1)
func (s CatmullRom[T]) Compute(u T) geometry.Vector[T] {
	n := len(s.points)
	if n < 4 {
		return geometry.Vector[T]{}
	}

	segments := n - 3
	scaled := float64(u) * float64(segments)
	idx := int(math.Floor(scaled))
	if idx >= segments {
		idx = segments - 1
	}
	if idx < 0 {
		idx = 0
	}
	t := T(scaled - float64(idx))

	p0 := s.points[idx]
	p1 := s.points[idx+1]
	p2 := s.points[idx+2]
	p3 := s.points[idx+3]

	t2 := t * t
	t3 := t2 * t

	// uniform catmull rom with tension 0.5
	//
	// p(t) = 0.5 * (2 p1 + (-p0 + p2) t + (2 p0 - 5 p1 + 4 p2 - p3) t^2 + (-p0 + 3 p1 - 3 p2 + p3) t^3)
	a := p1.Scale(2)
	b := p2.Sub(p0).Scale(t)
	cterm := p0.Scale(2).Sub(p1.Scale(5)).Add(p2.Scale(4)).Sub(p3).Scale(t2)
	d := p1.Scale(3).Sub(p0).Sub(p2.Scale(3)).Add(p3).Scale(t3)

	return a.Add(b).Add(cterm).Add(d).Scale(0.5)
}
