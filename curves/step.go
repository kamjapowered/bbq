package curves

import (
	"math"

	c "github.com/kamjapowered/bbq/constraints"
)

// single threshold step
//
//	f(u) = 0 if u < u0, else 1
type Step[T c.Float] struct {
	u0 T
}

// create a step curve
//
// u0 sets the transition point in [0, 1]
//
// time: O(1)
func NewStep[T c.Float](u0 T) Step[T] {
	return Step[T]{u0: u0}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (s Step[T]) Compute(u T) T {
	if u < s.u0 {
		return 0
	}
	return 1
}

// uniform staircase with n steps from 0 to 1
//
// outputs the discrete levels 0, 1/(n-1), 2/(n-1), ..., 1
type Staircase[T c.Float] struct {
	n T
}

// create a staircase curve
//
// n sets the number of steps, must be at least 2
//
// time: O(1)
func NewStaircase[T c.Float](n int) Staircase[T] {
	if n < 2 {
		n = 2
	}
	return Staircase[T]{n: T(n)}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (s Staircase[T]) Compute(u T) T {
	if u >= 1 {
		return 1
	}
	k := math.Floor(float64(u) * float64(s.n))
	return T(k / float64(s.n-1))
}
