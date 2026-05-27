package filter

import (
	"math"

	c "github.com/vistormu/go-dsa/constraints"
)

// estimate position and velocity from scalar position measurements
//
// use a constant velocity model
//
// state x = [pos, vel]^t
//
// model
//
//	x_k = F*x_{k-1} + w
//	z_k = H*x_k + v
//
// with
//
//	F = [1 dt; 0 1]
//	H = [1 0]
//
// q is the process variance of acceleration noise
// r is the measurement variance of position noise
//
// this type is not safe for concurrent use
type KalmanConstVel[T c.Float] struct {
	q float64
	r float64

	// state estimate
	x0 float64 // pos
	x1 float64 // vel

	// covariance p (2x2 symmetric)
	p00 float64
	p01 float64
	p11 float64
}

// create a constant velocity kalman filter
//
// q is the acceleration noise variance
//
// r is the position measurement noise variance
//
// p0 is the initial position covariance
//
// v0 is the initial velocity covariance
func NewKalmanConstVel[T c.Float](
	q, r float64,
	initialPos T,
	initialVel T,
	p0 float64,
	v0 float64,
) *KalmanConstVel[T] {
	return &KalmanConstVel[T]{
		q: q,
		r: r,

		x0: float64(initialPos),
		x1: float64(initialVel),

		p00: p0,
		p01: 0,
		p11: v0,
	}
}

// reset the filter state
func (k *KalmanConstVel[T]) Reset(initialPos, initialVel T, p0, v0 float64) {
	k.x0 = float64(initialPos)
	k.x1 = float64(initialVel)

	k.p00 = p0
	k.p01 = 0
	k.p11 = v0
}

// compute the next estimate from a position measurement and dt
//
// return the updated position estimate
//
// return zero if dt is not positive
//
// time: O(1)
func (k *KalmanConstVel[T]) Compute(z T, dt T) T {
	dtf := float64(dt)
	if dtf <= 0 {
		return 0
	}

	// ----------
	// predict
	// ----------
	// x = F*x
	x0p := k.x0 + dtf*k.x1
	x1p := k.x1

	// process noise for constant velocity with acceleration noise
	// qd = q * [dt^4/4  dt^3/2
	//          dt^3/2  dt^2]
	dt2 := dtf * dtf
	dt3 := dt2 * dtf
	dt4 := dt2 * dt2

	q00 := k.q * (dt4 / 4)
	q01 := k.q * (dt3 / 2)
	q11 := k.q * dt2

	// P = F*P*F^T + Q
	// with P symmetric (p10 = p01)
	p00p := k.p00 + 2*dtf*k.p01 + dt2*k.p11 + q00
	p01p := k.p01 + dtf*k.p11 + q01
	p11p := k.p11 + q11

	// ----------
	// update
	// ----------
	// innovation y = z - H*x = z - pos
	innov := float64(z) - x0p

	// S = H*P*H^T + R = p00 + r
	s := p00p + k.r
	if s == 0 || math.IsNaN(s) || math.IsInf(s, 0) {
		// keep predicted state if update is ill conditioned
		k.x0 = x0p
		k.x1 = x1p
		k.p00 = p00p
		k.p01 = p01p
		k.p11 = p11p
		return T(k.x0)
	}

	// K = P*H^T / S = [p00; p01] / S
	k0 := p00p / s
	k1 := p01p / s

	// x = x + K*y
	k.x0 = x0p + k0*innov
	k.x1 = x1p + k1*innov

	// P = (I - K*H) P
	// with H = [1 0]
	// p00 = (1-k0)*p00
	// p01 = (1-k0)*p01
	// p11 = p11 - k1*p01
	k.p00 = (1 - k0) * p00p
	k.p01 = (1 - k0) * p01p
	k.p11 = p11p - k1*p01p

	return T(k.x0)
}

// return the current position estimate
//
// time: O(1)
func (k *KalmanConstVel[T]) Pos() T {
	return T(k.x0)
}

// return the current velocity estimate
//
// time: O(1)
func (k *KalmanConstVel[T]) Vel() T {
	return T(k.x1)
}
