package control

import (
	"math"

	c "github.com/vistormu/go-dsa/constraints"
)

// generate a sinusoidal reference signal
type Sine[T c.Float] struct {
	amp    float64
	freq   float64
	phi    float64
	offset float64
}

// create a sine reference
//
// amp sets the amplitude
//
// freq sets the frequency in hz
//
// phi sets the phase offset in radians
//
// offset adds a constant bias
func NewSine[T c.Float](amp, freq, phi, offset T) Sine[T] {
	return Sine[T]{
		amp:    float64(amp),
		freq:   float64(freq),
		phi:    float64(phi),
		offset: float64(offset),
	}
}

// compute the reference value at time t
//
// time: O(1)
func (s Sine[T]) Compute(t T) T {
	return T(
		s.amp*math.Sin(2*math.Pi*s.freq*float64(t)+s.phi) +
			s.offset,
	)
}
