package control

import c "github.com/vistormu/go-dsa/constraints"

// produce a ramp reference after an optional delay
type Ramp[T c.Number] struct {
	slope T
	delay T
}

// create a ramp reference
func NewRamp[T c.Number](slope, delay T) Ramp[T] {
	return Ramp[T]{slope: slope, delay: delay}
}

// compute the reference value at time t
//
// time: O(1)
func (r Ramp[T]) Compute(t T) T {
	if t < r.delay {
		return 0
	}
	return r.slope * (t - r.delay)
}
