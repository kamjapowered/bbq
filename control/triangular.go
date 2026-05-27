package control

import (
	"math"

	c "github.com/vistormu/go-dsa/constraints"
)

// generate a triangular reference signal
type Triangular[T c.Float] struct {
	amp    float64
	freq   float64
	phi    float64
	offset float64
}

// create a triangular reference
//
// amp sets the amplitude
//
// freq sets the frequency in hz
//
// phi sets the phase offset in radians
//
// offset adds a constant bias
func NewTriangular[T c.Float](amp, freq, phi, offset T) Triangular[T] {
	return Triangular[T]{
		amp:    float64(amp),
		freq:   float64(freq),
		phi:    float64(phi),
		offset: float64(offset),
	}
}

// compute the reference value at time t
//
// time: O(1)
func (tr Triangular[T]) Compute(t T) T {
	return T(tr.amp*(2/math.Pi*math.Asin(math.Sin(2*math.Pi*tr.freq*float64(t)+tr.phi))) + tr.offset)
}
