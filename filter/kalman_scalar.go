package filter

import c "github.com/vistormu/go-dsa/constraints"

// estimate a scalar signal with a 1d kalman filter
//
// model assumes a random walk with f = 1 and h = 1
//
// this type is not safe for concurrent use
type KalmanScalar[T c.Float] struct {
	q float64
	r float64

	xHat float64
	p    float64
}

// create a 1d kalman filter
//
// processVariance is q
//
// measurementVariance is r
//
// initialErrorCovariance is p0
//
// initialEstimate is x0
func NewKalmanScalar[T c.Float](processVariance, measurementVariance, initialErrorCovariance float64, initialEstimate T) *KalmanScalar[T] {
	return &KalmanScalar[T]{
		q:    processVariance,
		r:    measurementVariance,
		p:    initialErrorCovariance,
		xHat: float64(initialEstimate),
	}
}

// reset the filter state
func (k *KalmanScalar[T]) Reset(initialErrorCovariance float64, initialEstimate T) {
	k.p = initialErrorCovariance
	k.xHat = float64(initialEstimate)
}

// compute the next estimate from a measurement
//
// time: O(1)
func (k *KalmanScalar[T]) Compute(measurement T) T {
	// predict
	xHatPred := k.xHat
	pPred := k.p + k.q

	// update
	den := pPred + k.r
	if den == 0 {
		return T(k.xHat)
	}

	kgain := pPred / den
	k.xHat = xHatPred + kgain*(float64(measurement)-xHatPred)
	k.p = (1 - kgain) * pPred

	return T(k.xHat)
}
