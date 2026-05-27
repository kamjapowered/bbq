package control

import c "github.com/vistormu/go-dsa/constraints"

// simulate a second order siso lti system in state space form
//
// continuous model
//
//	x' = a*x + b*u, where x is 2x1
//	y  = c*x + d*u, where c is 1x2
//
// discretise with forward euler on each compute call
//
// this type is not safe for concurrent use
type SecondOrder[T c.Float] struct {
	a11, a12 T
	a21, a22 T

	b1, b2 T

	c1, c2 T
	d      T

	x1, x2 T
}

// create a second order system from state space coefficients
//
// time: O(1)
func NewSecondOrderSS[T c.Float](
	a11, a12, a21, a22 T,
	b1, b2 T,
	c1, c2 T,
	d T,
) *SecondOrder[T] {
	return &SecondOrder[T]{
		a11: a11, a12: a12,
		a21: a21, a22: a22,
		b1: b1, b2: b2,
		c1: c1, c2: c2,
		d: d,
	}
}

// create a second order system from standard transfer function constants
//
// transfer function
//
//	g(s) = k*wn^2 / (s^2 + 2*zeta*wn*s + wn^2)
//
// time: O(1)
func NewSecondOrder[T c.Float](k, wn, zeta T) *SecondOrder[T] {
	// canonical controllable form with y = x1
	// x1' = x2
	// x2' = -wn^2*x1 - 2*zeta*wn*x2 + k*wn^2*u
	return &SecondOrder[T]{
		a11: 0, a12: 1,
		a21: -(wn * wn), a22: -(2 * zeta * wn),

		b1: 0, b2: k * wn * wn,

		c1: 1, c2: 0,
		d: 0,
	}
}

// reset internal state
//
// time: O(1)
func (s *SecondOrder[T]) Reset() {
	s.x1 = 0
	s.x2 = 0
}

// compute output for input u and timestep dt
//
// return 0 if dt is not positive
//
// time: O(1)
func (s *SecondOrder[T]) Compute(u, dt T) T {
	if dt <= 0 {
		return 0
	}

	// x' = a*x + b*u
	dx1 := s.a11*s.x1 + s.a12*s.x2 + s.b1*u
	dx2 := s.a21*s.x1 + s.a22*s.x2 + s.b2*u

	// forward euler
	s.x1 += dt * dx1
	s.x2 += dt * dx2

	return s.c1*s.x1 + s.c2*s.x2 + s.d*u
}
