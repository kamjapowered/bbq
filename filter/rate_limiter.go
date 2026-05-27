package filter

import c "github.com/vistormu/go-dsa/constraints"

// limit the rate of change of a signal
//
// this type is not safe for concurrent use
type RateLimiter[T c.Float] struct {
	rate T
	y    T
	init bool
}

// create a rate limiter with a maximum absolute rate per second
//
// rate must be non negative
func NewRateLimiter[T c.Float](rate T) *RateLimiter[T] {
	if rate < 0 {
		rate = 0
	}
	return &RateLimiter[T]{rate: rate}
}

// reset internal state
func (r *RateLimiter[T]) Reset() {
	r.y = 0
	r.init = false
}

// compute the limited value given input x and timestep dt
//
// return zero if dt is not positive
//
// time: O(1)
func (r *RateLimiter[T]) Compute(x, dt T) T {
	if dt <= 0 {
		return 0
	}

	if !r.init {
		r.y = x
		r.init = true
		return r.y
	}

	maxDelta := r.rate * dt
	delta := x - r.y
	delta = min(maxDelta, max(-maxDelta, delta))

	r.y += delta
	return r.y
}
