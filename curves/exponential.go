package curves

import (
	"math"

	c "github.com/vistormu/go-dsa/constraints"
)

// exponential ease in
//
//	f(u) = 2^(10 * (u - 1)), f(0) = 0
type ExpoIn[T c.Float] struct{}

// create an exponential ease in curve
//
// time: O(1)
func NewExpoIn[T c.Float]() ExpoIn[T] {
	return ExpoIn[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (e ExpoIn[T]) Compute(u T) T {
	if u == 0 {
		return 0
	}
	return T(math.Pow(2, 10*(float64(u)-1)))
}

// exponential ease out
//
//	f(u) = 1 - 2^(-10 * u), f(1) = 1
type ExpoOut[T c.Float] struct{}

// create an exponential ease out curve
//
// time: O(1)
func NewExpoOut[T c.Float]() ExpoOut[T] {
	return ExpoOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (e ExpoOut[T]) Compute(u T) T {
	if u == 1 {
		return 1
	}
	return T(1 - math.Pow(2, -10*float64(u)))
}

// exponential ease in out
type ExpoInOut[T c.Float] struct{}

// create an exponential ease in out curve
//
// time: O(1)
func NewExpoInOut[T c.Float]() ExpoInOut[T] {
	return ExpoInOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (e ExpoInOut[T]) Compute(u T) T {
	if u == 0 {
		return 0
	}
	if u == 1 {
		return 1
	}
	if u < 0.5 {
		return T(0.5 * math.Pow(2, 20*float64(u)-10))
	}
	return T(1 - 0.5*math.Pow(2, -20*float64(u)+10))
}
