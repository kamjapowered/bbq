package control

import (
	"math"

	c "github.com/vistormu/go-dsa/constraints"
)

// generate a square reference signal
type Square[T c.Float] struct {
	amp    float64
	freq   float64
	phi    float64
	offset float64
}

// create a square reference
//
// amp sets the amplitude
//
// freq sets the frequency in hz
//
// phi sets the phase offset in radians
//
// offset adds a constant bias
func NewSquare[T c.Float](amp, freq, phi, offset T) Square[T] {
	return Square[T]{
		amp:    float64(amp),
		freq:   float64(freq),
		phi:    float64(phi),
		offset: float64(offset),
	}
}

// compute the reference value at time t
//
// time: O(1)
func (s Square[T]) Compute(t T) T {
	x := math.Sin(2*math.Pi*s.freq*float64(t) + s.phi)
	if x == 0 {
		return T(s.offset)
	}
	return T(s.amp*math.Copysign(1, x) + s.offset)
}
