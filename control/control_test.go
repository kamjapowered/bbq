package control

import (
	"math"
	"testing"

	bmath "kamjapowered.com/bbq/math"
)

// ----------
// helpers
// ----------

// a stateful block driven by (input, dt)
type system interface {
	Compute(u, dt float64) float64
	Reset()
}

// a stateless reference generator driven by time
type generator interface {
	Compute(t float64) float64
}

func near(t *testing.T, name string, got, want, eps float64) {
	t.Helper()
	if !bmath.AlmostEqual(got, want, eps) {
		t.Fatalf("%s = %v, want %v (eps %v)", name, got, want, eps)
	}
}

func drive(s system, in []float64, dt float64) []float64 {
	out := make([]float64, len(in))
	for i, u := range in {
		out[i] = s.Compute(u, dt)
	}
	return out
}

// drive, reset, drive the same signal again, outputs must match exactly
func checkReset(t *testing.T, name string, s system, in []float64, dt float64) {
	t.Helper()

	first := drive(s, in, dt)
	s.Reset()
	second := drive(s, in, dt)

	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("%s: sample %d after reset = %v, want %v", name, i, second[i], first[i])
		}
	}
}

// a varied excitation signal for reset checks
func excitation(n int) []float64 {
	in := make([]float64, n)
	for i := range in {
		in[i] = math.Sin(0.3*float64(i)) + 0.5*math.Cos(0.11*float64(i)) + 0.25
	}
	return in
}

// ----------
// first order
// ----------

func TestFirstOrderStepResponse(t *testing.T) {
	const (
		k   = 2.0
		tau = 1.0
		dt  = 1e-3
	)

	s := NewFirstOrder(k, tau)

	// closed form of the forward euler recursion after n steps
	for _, n := range []int{1, 10, 100, 1000} {
		s.Reset()
		var y float64
		for range n {
			y = s.Compute(1, dt)
		}
		want := k * (1 - math.Pow(1-dt/tau, float64(n)))
		near(t, "first order euler response", y, want, 1e-9)
	}

	// steady state is the dc gain k
	s.Reset()
	var y float64
	for range 20000 {
		y = s.Compute(1, dt)
	}
	near(t, "first order steady state", y, k, 1e-6)
}

func TestFirstOrderZeroTau(t *testing.T) {
	// degenerate tau collapses to a pure gain with no state
	s := NewFirstOrder(3.0, 0.0)

	for _, u := range []float64{1, -2, 0.5} {
		near(t, "first order pure gain", s.Compute(u, 0.1), 3*u, 0)
	}
}

func TestFirstOrderSS(t *testing.T) {
	// direct feedthrough only
	d := NewFirstOrderSS(0.0, 0.0, 0.0, 4.0)
	near(t, "feedthrough", d.Compute(2, 0.1), 8, 0)

	// x' = -x + u, y = x, steady state is u
	s := NewFirstOrderSS(-1.0, 1.0, 1.0, 0.0)
	var y float64
	for range 200000 {
		y = s.Compute(1.5, 1e-4)
	}
	near(t, "state space steady state", y, 1.5, 1e-6)
}

func TestFirstOrderNonPositiveDt(t *testing.T) {
	s := NewFirstOrder(2.0, 1.0)
	s.Compute(1, 0.1)
	before := s.Compute(1, 0.1)

	// a rejected step holds the last output instead of emitting a zero
	for _, dt := range []float64{0, -1} {
		if got := s.Compute(1, dt); got != before {
			t.Fatalf("Compute with dt=%v = %v, want the held %v", dt, got, before)
		}
	}

	// state must be untouched by the rejected calls
	after := s.Compute(1, 0.1)
	want := before + 0.1*(-before+2*1)
	near(t, "state after rejected dt", after, want, 1e-12)
}

func TestFirstOrderReset(t *testing.T) {
	checkReset(t, "FirstOrder", NewFirstOrder(2.0, 0.5), excitation(64), 0.01)
	checkReset(t, "FirstOrderSS", NewFirstOrderSS(-3.0, 2.0, 1.5, 0.25), excitation(64), 0.01)
}

// ----------
// second order
// ----------

func TestSecondOrderStepResponse(t *testing.T) {
	const dt = 1e-4

	// critically damped, no overshoot, converges to k
	s := NewSecondOrder(1.0, 1.0, 1.0)
	var y, peak float64
	for range 300000 {
		y = s.Compute(1, dt)
		peak = max(peak, y)
	}
	near(t, "second order steady state", y, 1.0, 1e-3)
	if peak > 1.0+1e-6 {
		t.Fatalf("critically damped peak = %v, want no overshoot", peak)
	}

	// underdamped overshoot matches exp(-pi*zeta/sqrt(1-zeta^2))
	const zeta = 0.3
	u := NewSecondOrder(1.0, 1.0, zeta)
	peak = 0
	for range 300000 {
		y = u.Compute(1, dt)
		peak = max(peak, y)
	}
	near(t, "underdamped steady state", y, 1.0, 1e-3)
	wantPeak := 1 + math.Exp(-math.Pi*zeta/math.Sqrt(1-zeta*zeta))
	near(t, "underdamped overshoot", peak, wantPeak, 1e-3)
}

func TestSecondOrderGain(t *testing.T) {
	// dc gain is k regardless of wn
	s := NewSecondOrder(2.5, 4.0, 0.9)
	var y float64
	for range 200000 {
		y = s.Compute(1, 1e-4)
	}
	near(t, "second order dc gain", y, 2.5, 1e-3)
}

func TestSecondOrderSS(t *testing.T) {
	// direct feedthrough only
	d := NewSecondOrderSS(0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 3.0)
	near(t, "feedthrough", d.Compute(-2, 0.1), -6, 0)

	// two decoupled integrators observed through c
	s := NewSecondOrderSS(-1.0, 0.0, 0.0, -2.0, 1.0, 4.0, 1.0, 0.5, 0.0)
	var y float64
	for range 200000 {
		y = s.Compute(1, 1e-4)
	}
	// x1 -> 1, x2 -> 2, y = 1*1 + 0.5*2
	near(t, "second order state space steady state", y, 2.0, 1e-6)
}

func TestSecondOrderNonPositiveDt(t *testing.T) {
	s := NewSecondOrder(1.0, 1.0, 1.0)
	s.Compute(1, 0.1)

	for _, dt := range []float64{0, -0.5} {
		if got := s.Compute(1, dt); got != 0 {
			t.Fatalf("Compute with dt=%v = %v, want 0", dt, got)
		}
	}
}

func TestSecondOrderReset(t *testing.T) {
	checkReset(t, "SecondOrder", NewSecondOrder(1.5, 3.0, 0.4), excitation(64), 0.01)
	checkReset(t, "SecondOrderSS", NewSecondOrderSS(0.0, 1.0, -2.0, -0.5, 0.0, 1.0, 1.0, 0.0, 0.0), excitation(64), 0.01)
}

// ----------
// pid
// ----------

func TestPidProportional(t *testing.T) {
	p := NewPid(2.0, 0.0, 0.0, 1.0)
	near(t, "p term", p.Compute(3, 0.1), 6, 0)
	near(t, "p term", p.Compute(-1.5, 0.1), -3, 0)
}

func TestPidIntegral(t *testing.T) {
	p := NewPid(0.0, 1.0, 0.0, 1.0)

	// integral accumulates err*dt
	near(t, "i step 1", p.Compute(1, 0.5), 0.5, 1e-12)
	near(t, "i step 2", p.Compute(1, 0.5), 1.0, 1e-12)
	near(t, "i step 3", p.Compute(-2, 0.5), 0.0, 1e-12)
}

func TestPidDerivative(t *testing.T) {
	// alpha = 1 uses the raw derivative
	p := NewPid(0.0, 0.0, 1.0, 1.0)
	near(t, "d first sample", p.Compute(1, 1), 1, 1e-12)
	near(t, "d constant slope", p.Compute(2, 1), 1, 1e-12)
	near(t, "d constant err", p.Compute(2, 1), 0, 1e-12)

	// alpha = 0 freezes the derivative at its previous value, which is zero
	f := NewPid(0.0, 0.0, 1.0, 0.0)
	for i := range 10 {
		if got := f.Compute(float64(i), 0.1); got != 0 {
			t.Fatalf("frozen derivative = %v, want 0", got)
		}
	}

	// alpha in between is a first order blend of the raw derivative
	const alpha = 0.25
	b := NewPid(0.0, 0.0, 1.0, alpha)
	var want float64
	prev := 0.0
	for i := 1; i <= 5; i++ {
		err := float64(i) * float64(i)
		raw := (err - prev) / 0.5
		want = alpha*raw + (1-alpha)*want
		near(t, "blended derivative", b.Compute(err, 0.5), want, 1e-12)
		prev = err
	}
}

func TestPidAntiWindup(t *testing.T) {
	const (
		ki   = 1.0
		iMax = 2.0
		err  = 10.0
		dt   = 0.1
	)

	p := NewPid(0.0, ki, 0.0, 1.0)
	p.AntiWindup(-iMax, iMax)

	// sustained saturation must not grow the integral past the clamp
	for range 1000 {
		if got := p.Compute(err, dt); got > iMax+1e-12 {
			t.Fatalf("output under saturation = %v, want <= %v", got, iMax)
		}
	}
	near(t, "clamped integral", p.Compute(err, dt), iMax, 1e-12)

	// the clamp must not lag: a sign flip unwinds immediately
	near(t, "unwind step", p.Compute(-err, dt), iMax-err*dt, 1e-12)

	// and it clamps symmetrically on the negative side
	for range 1000 {
		p.Compute(-err, dt)
	}
	near(t, "clamped negative integral", p.Compute(-err, dt), -iMax, 1e-12)
}

func TestPidAntiWindupDisabled(t *testing.T) {
	// both limits zero disables the clamp entirely
	p := NewPid(0.0, 1.0, 0.0, 1.0)
	p.AntiWindup(0, 0)

	for range 100 {
		p.Compute(1, 0.1)
	}
	near(t, "unclamped integral", p.Compute(1, 0.1), 10.1, 1e-9)
}

func TestPidAntiWindupAsymmetric(t *testing.T) {
	p := NewPid(0.0, 1.0, 0.0, 1.0)
	p.AntiWindup(-0.5, 3.0)

	for range 100 {
		p.Compute(1, 0.1)
	}
	near(t, "upper clamp", p.Compute(1, 0.1), 3.0, 1e-12)

	for range 100 {
		p.Compute(-1, 0.1)
	}
	near(t, "lower clamp", p.Compute(-1, 0.1), -0.5, 1e-12)
}

func TestPidNonPositiveDt(t *testing.T) {
	p := NewPid(1.0, 1.0, 1.0, 1.0)
	twin := NewPid(1.0, 1.0, 1.0, 1.0)

	held := p.Compute(1, 0.1)
	twin.Compute(1, 0.1)

	// a rejected step holds the last output instead of emitting a zero
	for _, dt := range []float64{0, -0.2} {
		if got := p.Compute(1, dt); got != held {
			t.Fatalf("Compute with dt=%v = %v, want the held %v", dt, got, held)
		}
	}

	// the rejected calls must leave integral, derivative and prevErr alone,
	// so the next valid call must track a twin that never saw them
	for range 5 {
		got := p.Compute(1, 0.1)
		want := twin.Compute(1, 0.1)
		if got != want {
			t.Fatalf("Compute after rejected dt = %v, want %v", got, want)
		}
	}
}

func TestPidReset(t *testing.T) {
	p := NewPid(1.5, 0.8, 0.3, 0.6)
	checkReset(t, "Pid", p, excitation(64), 0.02)

	// reset must clear the clamped integral too
	c := NewPid(0.0, 1.0, 0.0, 1.0)
	c.AntiWindup(-1, 1)
	for range 100 {
		c.Compute(10, 0.1)
	}
	c.Reset()
	near(t, "integral after reset", c.Compute(1, 0.1), 0.1, 1e-12)
}

// ----------
// ramp
// ----------

func TestRamp(t *testing.T) {
	tests := []struct {
		name  string
		slope float64
		delay float64
		t     float64
		want  float64
	}{
		{"before delay", 2, 1, 0.5, 0},
		{"at delay", 2, 1, 1, 0},
		{"after delay", 2, 1, 3, 4},
		{"no delay", 0.5, 0, 4, 2},
		{"negative slope", -1, 0, 3, -3},
		{"negative time", 1, 0, -2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRamp(tt.slope, tt.delay)
			near(t, "ramp", r.Compute(tt.t), tt.want, 1e-12)
		})
	}
}

func TestRampInteger(t *testing.T) {
	// the ramp is generic over numbers, not just floats
	r := NewRamp(3, 2)
	for _, tt := range []struct{ in, want int }{{0, 0}, {1, 0}, {2, 0}, {5, 9}} {
		if got := r.Compute(tt.in); got != tt.want {
			t.Fatalf("Compute(%d) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

// ----------
// step
// ----------

func TestStep(t *testing.T) {
	tests := []struct {
		name  string
		amp   float64
		delay float64
		t     float64
		want  float64
	}{
		{"before delay", 5, 2, 1.99, 0},
		{"at delay", 5, 2, 2, 5},
		{"after delay", 5, 2, 100, 5},
		{"no delay", -3, 0, 0, -3},
		{"negative time", 1, 0, -0.001, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStep(tt.amp, tt.delay)
			near(t, "step", s.Compute(tt.t), tt.want, 0)
		})
	}
}

func TestStepInteger(t *testing.T) {
	s := NewStep(7, 3)
	for _, tt := range []struct{ in, want int }{{2, 0}, {3, 7}, {9, 7}} {
		if got := s.Compute(tt.in); got != tt.want {
			t.Fatalf("Compute(%d) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

// ----------
// sine
// ----------

func TestSine(t *testing.T) {
	const (
		amp    = 2.0
		freq   = 4.0
		offset = 1.0
		period = 1 / freq
	)

	s := NewSine(amp, freq, 0.0, offset)

	// quarter period landmarks
	near(t, "sine at 0", s.Compute(0), offset, 1e-12)
	near(t, "sine peak", s.Compute(period/4), offset+amp, 1e-12)
	near(t, "sine zero crossing", s.Compute(period/2), offset, 1e-12)
	near(t, "sine trough", s.Compute(3*period/4), offset-amp, 1e-12)
}

func TestSinePhase(t *testing.T) {
	// a phase of pi/2 turns the sine into a cosine
	s := NewSine(1.0, 1.0, math.Pi/2, 0.0)
	near(t, "phase shifted", s.Compute(0), 1, 1e-12)

	// a phase of phi is the same as a time shift of phi/(2*pi*f)
	const freq = 3.0
	const phi = 0.7
	shifted := NewSine(1.0, freq, phi, 0.0)
	plain := NewSine(1.0, freq, 0.0, 0.0)
	for _, tv := range []float64{0, 0.05, 0.13, 0.29} {
		near(t, "phase as time shift", shifted.Compute(tv), plain.Compute(tv+phi/(2*math.Pi*freq)), 1e-12)
	}
}

// ----------
// square
// ----------

func TestSquare(t *testing.T) {
	const (
		amp    = 3.0
		freq   = 2.0
		offset = 0.5
		period = 1 / freq
	)

	s := NewSquare(amp, freq, 0.0, offset)

	// exactly at the start the underlying sine is zero, so only the offset shows
	near(t, "square at 0", s.Compute(0), offset, 0)

	// high for the first half of the period, low for the second
	near(t, "square high", s.Compute(period/4), offset+amp, 0)
	near(t, "square low", s.Compute(3*period/4), offset-amp, 0)

	// the output only ever takes the two levels
	for i := range 1000 {
		got := s.Compute(float64(i) * 0.0013)
		if got != offset+amp && got != offset-amp && got != offset {
			t.Fatalf("square produced %v, want one of %v, %v, %v", got, offset+amp, offset-amp, offset)
		}
	}
}

func TestSquarePhase(t *testing.T) {
	// half a period of phase inverts the wave
	s := NewSquare(1.0, 1.0, math.Pi, 0.0)
	near(t, "inverted square", s.Compute(0.25), -1, 0)
	near(t, "inverted square", s.Compute(0.75), 1, 0)
}

// ----------
// triangular
// ----------

func TestTriangular(t *testing.T) {
	const (
		amp    = 2.0
		freq   = 1.0
		offset = -0.5
		period = 1 / freq
	)

	tr := NewTriangular(amp, freq, 0.0, offset)

	near(t, "triangular at 0", tr.Compute(0), offset, 1e-12)
	near(t, "triangular peak", tr.Compute(period/4), offset+amp, 1e-12)
	near(t, "triangular zero crossing", tr.Compute(period/2), offset, 1e-12)
	near(t, "triangular trough", tr.Compute(3*period/4), offset-amp, 1e-12)

	// the rising edge is linear with slope 4*amp/period
	slope := 4 * amp / period
	for _, tv := range []float64{0.05, 0.1, 0.2} {
		near(t, "triangular linear ramp", tr.Compute(tv), offset+slope*tv, 1e-12)
	}

	// never leaves the band
	for i := range 1000 {
		got := tr.Compute(float64(i) * 0.0017)
		if got > offset+amp+1e-12 || got < offset-amp-1e-12 {
			t.Fatalf("triangular produced %v, outside [%v, %v]", got, offset-amp, offset+amp)
		}
	}
}

// ----------
// periodicity
// ----------

func TestGeneratorPeriod(t *testing.T) {
	const freq = 2.5
	const period = 1 / freq

	gens := []struct {
		name string
		g    generator
	}{
		{"Sine", NewSine(1.5, freq, 0.4, 0.25)},
		{"Square", NewSquare(1.5, freq, 0.4, 0.25)},
		{"Triangular", NewTriangular(1.5, freq, 0.4, 0.25)},
	}

	for _, tt := range gens {
		t.Run(tt.name, func(t *testing.T) {
			for i := range 50 {
				tv := float64(i) * 0.007
				near(t, "period", tt.g.Compute(tv+period), tt.g.Compute(tv), 1e-9)
			}
		})
	}
}

// ----------
// float32 instantiation
// ----------

func TestFloat32Instantiation(t *testing.T) {
	s := NewFirstOrder[float32](1, 1)
	if got := s.Compute(1, 0.1); got != float32(0.1) {
		t.Fatalf("float32 first order = %v, want 0.1", got)
	}

	sn := NewSine[float32](1, 1, 0, 0)
	if got := sn.Compute(0.25); !bmath.AlmostEqual(got, 1, 1e-6) {
		t.Fatalf("float32 sine = %v, want 1", got)
	}

	p := NewPid[float32](1, 0, 0, 1)
	if got := p.Compute(2, 0.1); got != 2 {
		t.Fatalf("float32 pid = %v, want 2", got)
	}
}
