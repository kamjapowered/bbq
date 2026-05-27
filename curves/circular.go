package curves

import (
	"math"

	c "github.com/vistormu/go-dsa/constraints"
)

// circular ease in
//
//	f(u) = 1 - sqrt(1 - u^2)
type CircIn[T c.Float] struct{}

// create a circular ease in curve
//
// time: O(1)
func NewCircIn[T c.Float]() CircIn[T] {
	return CircIn[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (c CircIn[T]) Compute(u T) T {
	uf := float64(u)
	return T(1 - math.Sqrt(1-uf*uf))
}

// circular ease out
//
//	f(u) = sqrt(1 - (1 - u)^2)
type CircOut[T c.Float] struct{}

// create a circular ease out curve
//
// time: O(1)
func NewCircOut[T c.Float]() CircOut[T] {
	return CircOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (c CircOut[T]) Compute(u T) T {
	v := 1 - float64(u)
	return T(math.Sqrt(1 - v*v))
}

// circular ease in out
type CircInOut[T c.Float] struct{}

// create a circular ease in out curve
//
// time: O(1)
func NewCircInOut[T c.Float]() CircInOut[T] {
	return CircInOut[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (c CircInOut[T]) Compute(u T) T {
	uf := float64(u)
	if uf < 0.5 {
		return T(0.5 * (1 - math.Sqrt(1-4*uf*uf)))
	}
	v := 2*uf - 2
	return T(0.5 * (math.Sqrt(1-v*v) + 1))
}
