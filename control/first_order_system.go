package control

import c "github.com/vistormu/go-dsa/constraints"

// simulate a first order siso lti system in state space form
//
// continuous model
//
//	x' = a*x + b*u
//	y  = c*x + d*u
//
// discretise with forward euler on each compute call
//
// this type is not safe for concurrent use
type FirstOrder[T c.Float] struct {
	a, b, c, d T
	x          T
}

// create a first order system from state space coefficients
//
// time: O(1)
func NewFirstOrderSS[T c.Float](a, b, c, d T) *FirstOrder[T] {
	return &FirstOrder[T]{a: a, b: b, c: c, d: d}
}

// create a first order system from transfer function constants
//
// transfer function
//
//	g(s) = k / (tau*s + 1)
//
// time: O(1)
func NewFirstOrder[T c.Float](k, tau T) *FirstOrder[T] {
	if tau == 0 {
		// degenerate, treat as pure gain in output equation
		return &FirstOrder[T]{a: 0, b: 0, c: 0, d: k}
	}

	return &FirstOrder[T]{
		a: -1 / tau,
		b: k / tau,
		c: 1,
		d: 0,
	}
}

// reset internal state
//
// time: O(1)
func (s *FirstOrder[T]) Reset() {
	s.x = 0
}

// compute output for input u and timestep dt
//
// return 0 if dt is not positive
//
// time: O(1)
func (s *FirstOrder[T]) Compute(u, dt T) T {
	if dt <= 0 {
		return 0
	}

	// forward euler: x_{k+1} = x_k + dt*(a*x_k + b*u_k)
	s.x += dt * (s.a*s.x + s.b*u)

	return s.c*s.x + s.d*u
}
