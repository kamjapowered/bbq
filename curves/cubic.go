package curves

import c "github.com/kamjapowered/bbq/constraints"

// cubic ease in
//
//	f(u) = u^3
type CubicIn[T c.Float] struct{}

// create a cubic ease in curve
//
// time: O(1)
func NewCubicIn[T c.Float]() CubicIn[T] {
	return CubicIn[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (c CubicIn[T]) Compute(u T) T {
	return u * u * u
}

// cubic ease out
//
//	f(u) = 1 - (1 - u)^3
type CubicOut[T c.Float] struct{}

// create a cubic ease out curve
//
// time: O(1)
func NewCubicOut[T c.Float]() CubicOut[T] {
	return CubicOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (c CubicOut[T]) Compute(u T) T {
	v := 1 - u
	return 1 - v*v*v
}

// cubic ease in out, symmetric around u = 0.5
type CubicInOut[T c.Float] struct{}

// create a cubic ease in out curve
//
// time: O(1)
func NewCubicInOut[T c.Float]() CubicInOut[T] {
	return CubicInOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (c CubicInOut[T]) Compute(u T) T {
	if u < 0.5 {
		return 4 * u * u * u
	}
	v := 1 - u
	return 1 - 4*v*v*v
}
