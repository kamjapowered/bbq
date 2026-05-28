package curves

import c "github.com/kamjapowered/bbq/constraints"

// quartic ease in
//
//	f(u) = u^4
type QuarticIn[T c.Float] struct{}

// create a quartic ease in curve
//
// time: O(1)
func NewQuarticIn[T c.Float]() QuarticIn[T] {
	return QuarticIn[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (q QuarticIn[T]) Compute(u T) T {
	return u * u * u * u
}

// quartic ease out
//
//	f(u) = 1 - (1 - u)^4
type QuarticOut[T c.Float] struct{}

// create a quartic ease out curve
//
// time: O(1)
func NewQuarticOut[T c.Float]() QuarticOut[T] {
	return QuarticOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (q QuarticOut[T]) Compute(u T) T {
	v := 1 - u
	return 1 - v*v*v*v
}

// quartic ease in out, symmetric around u = 0.5
type QuarticInOut[T c.Float] struct{}

// create a quartic ease in out curve
//
// time: O(1)
func NewQuarticInOut[T c.Float]() QuarticInOut[T] {
	return QuarticInOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (q QuarticInOut[T]) Compute(u T) T {
	if u < 0.5 {
		return 8 * u * u * u * u
	}
	v := 1 - u
	return 1 - 8*v*v*v*v
}
