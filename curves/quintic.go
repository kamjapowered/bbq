package curves

import c "github.com/vistormu/go-dsa/constraints"

// quintic ease in
//
//	f(u) = u^5
type QuinticIn[T c.Float] struct{}

// create a quintic ease in curve
//
// time: O(1)
func NewQuinticIn[T c.Float]() QuinticIn[T] {
	return QuinticIn[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (q QuinticIn[T]) Compute(u T) T {
	return u * u * u * u * u
}

// quintic ease out
//
//	f(u) = 1 - (1 - u)^5
type QuinticOut[T c.Float] struct{}

// create a quintic ease out curve
//
// time: O(1)
func NewQuinticOut[T c.Float]() QuinticOut[T] {
	return QuinticOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (q QuinticOut[T]) Compute(u T) T {
	v := 1 - u
	return 1 - v*v*v*v*v
}

// quintic ease in out, symmetric around u = 0.5
type QuinticInOut[T c.Float] struct{}

// create a quintic ease in out curve
//
// time: O(1)
func NewQuinticInOut[T c.Float]() QuinticInOut[T] {
	return QuinticInOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (q QuinticInOut[T]) Compute(u T) T {
	if u < 0.5 {
		return 16 * u * u * u * u * u
	}
	v := 1 - u
	return 1 - 16*v*v*v*v*v
}
