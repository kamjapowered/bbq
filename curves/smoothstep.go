package curves

import c "github.com/kamjapowered/bbq/constraints"

// hermite smoothstep, c1 continuous
//
//	f(u) = 3u^2 - 2u^3
type Smoothstep[T c.Float] struct{}

// create a smoothstep curve
//
// time: O(1)
func NewSmoothstep[T c.Float]() Smoothstep[T] {
	return Smoothstep[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (s Smoothstep[T]) Compute(u T) T {
	return u * u * (3 - 2*u)
}

// hermite smootherstep, c2 continuous
//
//	f(u) = 6u^5 - 15u^4 + 10u^3
type Smootherstep[T c.Float] struct{}

// create a smootherstep curve
//
// time: O(1)
func NewSmootherstep[T c.Float]() Smootherstep[T] {
	return Smootherstep[T]{}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (s Smootherstep[T]) Compute(u T) T {
	return u * u * u * (u*(u*6-15) + 10)
}
