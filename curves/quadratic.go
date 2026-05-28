package curves

import c "github.com/kamjapowered/bbq/constraints"

// quadratic ease in
//
//	f(u) = u^2
type QuadIn[T c.Float] struct{}

// create a quadratic ease in curve
//
// time: O(1)
func NewQuadIn[T c.Float]() QuadIn[T] {
	return QuadIn[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (q QuadIn[T]) Compute(u T) T {
	return u * u
}

// quadratic ease out
//
//	f(u) = 1 - (1 - u)^2
type QuadOut[T c.Float] struct{}

// create a quadratic ease out curve
//
// time: O(1)
func NewQuadOut[T c.Float]() QuadOut[T] {
	return QuadOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (q QuadOut[T]) Compute(u T) T {
	return 1 - (1-u)*(1-u)
}

// quadratic ease in out, symmetric around u = 0.5
type QuadInOut[T c.Float] struct{}

// create a quadratic ease in out curve
//
// time: O(1)
func NewQuadInOut[T c.Float]() QuadInOut[T] {
	return QuadInOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (q QuadInOut[T]) Compute(u T) T {
	if u < 0.5 {
		return 2 * u * u
	}
	v := 1 - u
	return 1 - 2*v*v
}
