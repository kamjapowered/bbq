package control

import c "kamjapowered.com/bbq/constraints"

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
	y          T
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
	s.y = 0
}

// compute output for input u and timestep dt
//
// hold and return the last output if dt is not positive
//
// time: O(1)
func (s *FirstOrder[T]) Compute(u, dt T) T {
	// a non positive step advances nothing, so hold the last output rather than
	// injecting a zero into the signal
	if dt <= 0 {
		return s.y
	}

	// forward euler: x_{k+1} = x_k + dt*(a*x_k + b*u_k)
	s.x += dt * (s.a*s.x + s.b*u)
	s.y = s.c*s.x + s.d*u

	return s.y
}
