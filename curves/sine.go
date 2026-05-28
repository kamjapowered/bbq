package curves

import (
	"math"

	c "github.com/kamjapowered/bbq/constraints"
)

// sinusoidal ease in
//
//	f(u) = 1 - cos(u * pi / 2)
type SineIn[T c.Float] struct{}

// create a sinusoidal ease in curve
//
// time: O(1)
func NewSineIn[T c.Float]() SineIn[T] {
	return SineIn[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (s SineIn[T]) Compute(u T) T {
	return T(1 - math.Cos(float64(u)*math.Pi/2))
}

// sinusoidal ease out
//
//	f(u) = sin(u * pi / 2)
type SineOut[T c.Float] struct{}

// create a sinusoidal ease out curve
//
// time: O(1)
func NewSineOut[T c.Float]() SineOut[T] {
	return SineOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (s SineOut[T]) Compute(u T) T {
	return T(math.Sin(float64(u) * math.Pi / 2))
}

// sinusoidal ease in out
//
//	f(u) = (1 - cos(u * pi)) / 2
type SineInOut[T c.Float] struct{}

// create a sinusoidal ease in out curve
//
// time: O(1)
func NewSineInOut[T c.Float]() SineInOut[T] {
	return SineInOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (s SineInOut[T]) Compute(u T) T {
	return T(0.5 * (1 - math.Cos(float64(u)*math.Pi)))
}
