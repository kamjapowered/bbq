package control

import c "github.com/vistormu/go-dsa/constraints"

// generate a step reference signal
type Step[T c.Number] struct {
	amp   T
	delay T
}

// create a step reference
//
// amp sets the output value after the delay
//
// delay sets the activation time
func NewStep[T c.Number](amp, delay T) Step[T] {
	return Step[T]{amp: amp, delay: delay}
}

// compute the reference value at time t
//
// time: O(1)
func (s Step[T]) Compute(t T) T {
	if t < s.delay {
		return 0
	}
	return s.amp
}
