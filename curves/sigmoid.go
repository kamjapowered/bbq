package curves

import (
	"math"

	c "github.com/kamjapowered/bbq/constraints"
)

// logistic sigmoid normalized so that compute(0) = 0 and compute(1) = 1
//
// raw form
//
//	r(u) = 1 / (1 + exp(-k*(u - x0)))
//
// normalized form
//
//	f(u) = (r(u) - r(0)) / (r(1) - r(0))
type Sigmoid[T c.Float] struct {
	k    float64
	x0   float64
	r0   float64
	span float64
}

// create a sigmoid curve
//
// k sets the steepness, larger values produce a sharper s-shape
//
// x0 sets the midpoint, the inflection of the raw sigmoid in [0, 1]
//
// time: O(1)
func NewSigmoid[T c.Float](k, x0 T) Sigmoid[T] {
	kf := float64(k)
	x0f := float64(x0)

	r0 := 1 / (1 + math.Exp(kf*x0f))
	r1 := 1 / (1 + math.Exp(-kf*(1-x0f)))

	return Sigmoid[T]{
		k:    kf,
		x0:   x0f,
		r0:   r0,
		span: r1 - r0,
	}
}

// compute the curve value at progress u in [0, 1]
//
// time: O(1)
func (s Sigmoid[T]) Compute(u T) T {
	r := 1 / (1 + math.Exp(-s.k*(float64(u)-s.x0)))
	return T((r - s.r0) / s.span)
}
