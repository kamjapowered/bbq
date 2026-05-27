package control

import c "github.com/vistormu/go-dsa/constraints"

// compute a pid control output from the current error signal
//
// this type is not safe for concurrent use
type Pid[T c.Float] struct {
	kp, ki, kd T
	alpha      T

	prevErr    T
	integral   T
	derivative T

	// integral clamp range for anti windup
	iMin T
	iMax T
}

// create a pid controller
//
// alpha controls derivative smoothing in [0, 1]
//
// alpha = 1 uses the raw derivative
//
// alpha = 0 freezes the derivative at its previous value
func NewPid[T c.Float](kp, ki, kd, alpha T) *Pid[T] {
	return &Pid[T]{kp: kp, ki: ki, kd: kd, alpha: alpha}
}

// reset internal state
func (p *Pid[T]) Reset() {
	p.prevErr = 0
	p.integral = 0
	p.derivative = 0
}

// set integral clamp limits for anti windup
//
// if both min and max are zero, anti windup is disabled
func (p *Pid[T]) AntiWindup(min, max T) {
	p.iMin = min
	p.iMax = max
}

// compute output given an error value and dt
//
// return 0 if dt is not positive
//
// time: O(1)
func (p *Pid[T]) Compute(err, dt T) T {
	if dt <= 0 {
		return 0
	}

	p.integral += err * dt

	if !(p.iMin == 0 && p.iMax == 0) {
		p.integral = min(p.iMax, max(p.iMin, p.integral))
	}

	rawD := (err - p.prevErr) / dt
	p.derivative = p.alpha*rawD + (1-p.alpha)*p.derivative

	out := p.kp*err + p.ki*p.integral + p.kd*p.derivative

	p.prevErr = err

	return out
}
