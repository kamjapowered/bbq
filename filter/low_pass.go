package filter

import c "github.com/vistormu/go-dsa/constraints"

// smooth a signal with a first order low pass filter
//
// update uses y += alpha * (x - y)
//
// this type is not safe for concurrent use
type LowPass[T c.Float] struct {
	alpha T
	y     T
	init  bool
}

// create a low pass filter with a fixed alpha in [0, 1]
//
// alpha = 1 follows input with no smoothing
//
// alpha = 0 holds the previous output
func NewLowPass[T c.Float](alpha T) *LowPass[T] {
	return &LowPass[T]{alpha: alpha}
}

// create a low pass filter from a time constant tau and timestep dt
//
// if tau is not positive, alpha becomes 1
func NewLowPassTau[T c.Float](tau, dt T) *LowPass[T] {
	if tau <= 0 || dt <= 0 {
		return &LowPass[T]{alpha: 1}
	}

	alpha := dt / (tau + dt)
	return &LowPass[T]{alpha: alpha}
}

// reset internal state
func (f *LowPass[T]) Reset() {
	f.y = 0
	f.init = false
}

// compute the filtered value
//
// time: O(1)
func (f *LowPass[T]) Compute(x T) T {
	if !f.init {
		f.y = x
		f.init = true
		return f.y
	}

	f.y += f.alpha * (x - f.y)
	return f.y
}
