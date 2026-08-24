package filter

import (
	"math"
	"testing"

	bmath "kamjapowered.com/bbq/math"
)

// ----------
// helpers
// ----------

// a stateful filter driven by a single sample
type filter interface {
	Compute(x float64) float64
	Reset()
}

func near(t *testing.T, name string, got, want, eps float64) {
	t.Helper()
	if !bmath.AlmostEqual(got, want, eps) {
		t.Fatalf("%s = %v, want %v (eps %v)", name, got, want, eps)
	}
}

// a varied excitation signal for reset checks
func excitation(n int) []float64 {
	in := make([]float64, n)
	for i := range in {
		in[i] = math.Sin(0.37*float64(i)) + 0.4*math.Cos(0.13*float64(i)) - 0.2
	}
	return in
}

// drive, reset, drive the same signal again, outputs must match exactly
func checkReset(t *testing.T, name string, f filter, in []float64) {
	t.Helper()

	first := make([]float64, len(in))
	for i, x := range in {
		first[i] = f.Compute(x)
	}

	f.Reset()

	for i, x := range in {
		if got := f.Compute(x); got != first[i] {
			t.Fatalf("%s: sample %d after reset = %v, want %v", name, i, got, first[i])
		}
	}
}

// ----------
// low pass
// ----------

func TestLowPassFirstSample(t *testing.T) {
	// the first sample seeds the state instead of being smoothed from zero
	f := NewLowPass(0.1)
	near(t, "first sample", f.Compute(7), 7, 0)
	near(t, "second sample", f.Compute(7), 7, 0)
}

func TestLowPassRecursion(t *testing.T) {
	const alpha = 0.25
	f := NewLowPass(alpha)

	in := []float64{2, 5, -1, 3, 3, 0}
	want := in[0]
	for i, x := range in {
		if i > 0 {
			want += alpha * (x - want)
		}
		near(t, "low pass recursion", f.Compute(x), want, 1e-12)
	}
}

func TestLowPassStepConvergence(t *testing.T) {
	f := NewLowPass(0.05)
	f.Compute(0)

	var y float64
	for range 1000 {
		y = f.Compute(10)
	}
	near(t, "low pass steady state", y, 10, 1e-9)
}

func TestLowPassAlphaExtremes(t *testing.T) {
	// alpha = 1 follows the input exactly
	one := NewLowPass(1.0)
	for _, x := range []float64{3, -2, 8} {
		near(t, "alpha one", one.Compute(x), x, 0)
	}

	// alpha = 0 holds the first sample forever
	zero := NewLowPass(0.0)
	near(t, "alpha zero seed", zero.Compute(4), 4, 0)
	for _, x := range []float64{100, -100, 0} {
		near(t, "alpha zero hold", zero.Compute(x), 4, 0)
	}
}

func TestNewLowPassTau(t *testing.T) {
	// alpha = dt / (tau + dt)
	const tau, dt = 0.5, 0.1
	f := NewLowPassTau(tau, dt)
	f.Compute(0)
	near(t, "tau alpha", f.Compute(1), dt/(tau+dt), 1e-12)

	// degenerate tau or dt gives alpha = 1, so the filter passes through
	for _, c := range [][2]float64{{0, 0.1}, {-1, 0.1}, {0.5, 0}, {0.5, -1}} {
		d := NewLowPassTau(c[0], c[1])
		d.Compute(0)
		near(t, "degenerate tau passthrough", d.Compute(9), 9, 0)
	}
}

func TestLowPassReset(t *testing.T) {
	checkReset(t, "LowPass", NewLowPass(0.3), excitation(64))

	// reset must clear the seeded flag, not just the value
	f := NewLowPass(0.1)
	f.Compute(100)
	f.Reset()
	near(t, "reseed after reset", f.Compute(5), 5, 0)
}

// ----------
// mean
// ----------

func TestMeanBeforeWindowIsFull(t *testing.T) {
	// averages over the samples seen so far
	m := NewMean[float64](4)
	near(t, "mean of 1", m.Compute(1), 1, 1e-12)
	near(t, "mean of 2", m.Compute(3), 2, 1e-12)
	near(t, "mean of 3", m.Compute(5), 3, 1e-12)
	near(t, "mean of 4", m.Compute(7), 4, 1e-12)
}

func TestMeanSlidingWindow(t *testing.T) {
	const w = 3
	m := NewMean[float64](w)

	in := []float64{1, 2, 3, 4, 5, 6, 7, 0, -3}
	for i, x := range in {
		got := m.Compute(x)

		lo := max(0, i+1-w)
		var sum float64
		for _, v := range in[lo : i+1] {
			sum += v
		}
		near(t, "sliding mean", got, sum/float64(i+1-lo), 1e-9)
	}
}

func TestMeanWindowOne(t *testing.T) {
	m := NewMean[float64](1)
	for _, x := range []float64{4, -2, 9} {
		near(t, "window of one", m.Compute(x), x, 1e-12)
	}
}

func TestMeanEmptyWindow(t *testing.T) {
	for _, w := range []int{0, -1, -10} {
		m := NewMean[float64](w)
		if got := m.Compute(42); got != 0 {
			t.Fatalf("Compute on window %d = %v, want 0", w, got)
		}
	}
}

func TestMeanReset(t *testing.T) {
	checkReset(t, "Mean", NewMean[float64](5), excitation(64))

	m := NewMean[float64](3)
	for _, x := range []float64{100, 200, 300} {
		m.Compute(x)
	}
	m.Reset()
	near(t, "mean after reset", m.Compute(1), 1, 1e-12)
}

// ----------
// median
// ----------

func TestMedianOddWindow(t *testing.T) {
	m := NewMedian[float64](3)

	// before the window is full the median is taken over what is there
	near(t, "one sample", m.Compute(5), 5, 0)
	near(t, "two samples", m.Compute(1), 3, 0)
	near(t, "three samples", m.Compute(9), 5, 0)

	// the window slides
	near(t, "slid window", m.Compute(2), 2, 0)
	near(t, "slid window", m.Compute(3), 3, 0)
}

func TestMedianEvenWindow(t *testing.T) {
	m := NewMedian[float64](4)

	near(t, "one sample", m.Compute(1), 1, 0)
	near(t, "two samples", m.Compute(2), 1.5, 0)
	near(t, "three samples", m.Compute(3), 2, 0)
	near(t, "four samples", m.Compute(4), 2.5, 0)

	// window slides to 2,3,4,10
	near(t, "slid window", m.Compute(10), 3.5, 0)
}

func TestMedianRejectsSpikes(t *testing.T) {
	m := NewMedian[float64](3)
	in := []float64{1, 1, 1, 1000, 1, 1, 1}
	var got float64
	for _, x := range in {
		got = m.Compute(x)
		if got > 1.0001 {
			t.Fatalf("median let a spike through: %v", got)
		}
	}
	near(t, "median after spike", got, 1, 0)
}

func TestMedianInteger(t *testing.T) {
	// integer windows truncate the even case
	m := NewMedian[int](2)
	if got := m.Compute(1); got != 1 {
		t.Fatalf("Compute(1) = %d, want 1", got)
	}
	// (1 + 2) / 2 in integer arithmetic
	if got := m.Compute(2); got != 1 {
		t.Fatalf("Compute(2) = %d, want 1", got)
	}
}

func TestMedianEmptyWindow(t *testing.T) {
	for _, w := range []int{0, -1} {
		m := NewMedian[float64](w)
		if got := m.Compute(42); got != 0 {
			t.Fatalf("Compute on window %d = %v, want 0", w, got)
		}
	}
}

func TestMedianReset(t *testing.T) {
	checkReset(t, "Median", NewMedian[float64](5), excitation(64))
	checkReset(t, "MedianEven", NewMedian[float64](4), excitation(64))
}

// ----------
// rate limiter
// ----------

func TestRateLimiterNeverExceedsRate(t *testing.T) {
	const (
		rate = 2.0
		dt   = 0.01
	)

	r := NewRateLimiter(rate)

	// a hostile input: big jumps, sign flips, holds
	in := make([]float64, 2000)
	for i := range in {
		switch {
		case i%200 < 50:
			in[i] = 100
		case i%200 < 100:
			in[i] = -100
		case i%200 < 150:
			in[i] = 0
		default:
			in[i] = 30 * math.Sin(0.7*float64(i))
		}
	}

	prev := r.Compute(in[0], dt)
	maxDelta := rate * dt
	for _, x := range in[1:] {
		y := r.Compute(x, dt)
		if math.Abs(y-prev) > maxDelta+1e-12 {
			t.Fatalf("rate limiter moved %v in one step, limit is %v", math.Abs(y-prev), maxDelta)
		}
		prev = y
	}
}

func TestRateLimiterFollowsSlowInput(t *testing.T) {
	// an input changing slower than the limit passes through untouched
	r := NewRateLimiter(10.0)
	const dt = 0.1

	r.Compute(0, dt)
	for i := 1; i <= 20; i++ {
		x := float64(i) * 0.5 // 5 per second, under the 10 per second limit
		near(t, "slow input passthrough", r.Compute(x, dt), x, 1e-12)
	}
}

func TestRateLimiterSlewsToTarget(t *testing.T) {
	const rate, dt = 2.0, 0.1
	r := NewRateLimiter(rate)

	r.Compute(0, dt)

	// a step of 10 needs 10 / (rate*dt) = 50 steps
	for i := 1; i <= 50; i++ {
		near(t, "slew", r.Compute(10, dt), float64(i)*rate*dt, 1e-9)
	}
	near(t, "arrived", r.Compute(10, dt), 10, 1e-9)
}

func TestRateLimiterFirstSample(t *testing.T) {
	// the first sample seeds the state without any limiting
	r := NewRateLimiter(0.001)
	near(t, "seed", r.Compute(500, 0.01), 500, 0)
}

func TestRateLimiterZeroAndNegativeRate(t *testing.T) {
	// a negative rate is clamped to zero, which freezes the output
	for _, rate := range []float64{0, -5} {
		r := NewRateLimiter(rate)
		r.Compute(1, 0.1)
		for _, x := range []float64{100, -100, 3} {
			near(t, "frozen output", r.Compute(x, 0.1), 1, 0)
		}
	}
}

func TestRateLimiterNonPositiveDt(t *testing.T) {
	r := NewRateLimiter(1.0)
	held := r.Compute(5, 0.1)

	// a rejected step holds the last output instead of emitting a zero
	for _, dt := range []float64{0, -1} {
		if got := r.Compute(9, dt); got != held {
			t.Fatalf("Compute with dt=%v = %v, want the held %v", dt, got, held)
		}
	}

	// the rejected calls must not move the state
	near(t, "state after rejected dt", r.Compute(9, 0.1), 5.1, 1e-12)
}

func TestRateLimiterReset(t *testing.T) {
	const dt = 0.05
	r := NewRateLimiter(3.0)
	in := excitation(64)

	first := make([]float64, len(in))
	for i, x := range in {
		first[i] = r.Compute(x, dt)
	}

	r.Reset()

	for i, x := range in {
		if got := r.Compute(x, dt); got != first[i] {
			t.Fatalf("RateLimiter: sample %d after reset = %v, want %v", i, got, first[i])
		}
	}
}

// ----------
// dead zone
// ----------

func TestDeadZone(t *testing.T) {
	const w = 1.5
	d := NewDeadZone(w)

	tests := []struct {
		name string
		in   float64
		want float64
	}{
		{"inside positive", 1.0, 0},
		{"inside negative", -1.0, 0},
		{"zero", 0, 0},
		{"on positive boundary", w, 0},
		{"on negative boundary", -w, 0},
		{"just outside positive", w + 0.25, 0.25},
		{"just outside negative", -w - 0.25, -0.25},
		{"far positive", 10, 10 - w},
		{"far negative", -10, -10 + w},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			near(t, "dead zone", d.Compute(tt.in), tt.want, 1e-12)
		})
	}
}

func TestDeadZoneIsContinuous(t *testing.T) {
	// the output is the input shifted by the width, so there is no jump at
	// the boundary
	const w = 2.0
	d := NewDeadZone(w)
	near(t, "just above boundary", d.Compute(w+1e-9), 1e-9, 1e-15)
	near(t, "just below boundary", d.Compute(w-1e-9), 0, 0)
}

func TestDeadZoneWidths(t *testing.T) {
	// zero width is a pass through
	z := NewDeadZone(0.0)
	for _, x := range []float64{-3, 0, 3} {
		near(t, "zero width", z.Compute(x), x, 0)
	}

	// a negative width is clamped to zero, so it also passes through
	n := NewDeadZone(-4.0)
	for _, x := range []float64{-3, 0, 3} {
		near(t, "negative width", n.Compute(x), x, 0)
	}
}

// ----------
// kalman scalar
// ----------

func TestKalmanScalarIsRunningMeanWithoutProcessNoise(t *testing.T) {
	// with q = 0, r = p0 the filter reduces to the running mean of the
	// measurements, so x_n = z * n / (n + 1)
	const z = 5.0
	k := NewKalmanScalar(0.0, 1.0, 1.0, 0.0)

	for n := 1; n <= 200; n++ {
		got := k.Compute(z)
		want := z * float64(n) / float64(n+1)
		near(t, "kalman running mean", got, want, 1e-9)
	}
}

func TestKalmanScalarConvergesToConstant(t *testing.T) {
	const truth = -3.25
	k := NewKalmanScalar(1e-5, 0.1, 1.0, 0.0)

	var got float64
	for range 5000 {
		got = k.Compute(truth)
	}
	near(t, "kalman convergence", got, truth, 1e-6)
}

func TestKalmanScalarTracksAStep(t *testing.T) {
	k := NewKalmanScalar(1e-3, 0.05, 1.0, 0.0)

	// settle on the first level
	var got float64
	for range 2000 {
		got = k.Compute(1)
	}
	near(t, "first level", got, 1, 1e-3)

	// then follow a step to the second level
	for range 2000 {
		got = k.Compute(-2)
	}
	near(t, "second level", got, -2, 1e-3)
}

func TestKalmanScalarInitialEstimate(t *testing.T) {
	// with p0 = 0 and q = 0 the gain is zero, so the estimate never moves
	k := NewKalmanScalar(0.0, 1.0, 0.0, 7.0)
	for range 10 {
		near(t, "frozen estimate", k.Compute(100), 7, 0)
	}
}

func TestKalmanScalarDegenerateDenominator(t *testing.T) {
	// q = r = p0 = 0 makes the innovation covariance zero, the filter must
	// hold its estimate instead of dividing by zero
	k := NewKalmanScalar(0.0, 0.0, 0.0, 2.5)
	for range 10 {
		got := k.Compute(100)
		if got != 2.5 || math.IsNaN(got) {
			t.Fatalf("degenerate kalman = %v, want 2.5", got)
		}
	}
}

func TestKalmanScalarReset(t *testing.T) {
	const p0, x0 = 1.0, 0.5
	k := NewKalmanScalar(1e-4, 0.2, p0, x0)
	in := excitation(64)

	first := make([]float64, len(in))
	for i, x := range in {
		first[i] = k.Compute(x)
	}

	k.Reset(p0, x0)

	for i, x := range in {
		if got := k.Compute(x); got != first[i] {
			t.Fatalf("KalmanScalar: sample %d after reset = %v, want %v", i, got, first[i])
		}
	}
}

// ----------
// kalman constant velocity
// ----------

func TestKalmanConstVelConvergesToConstant(t *testing.T) {
	const (
		truth = 4.0
		dt    = 0.01
	)

	k := NewKalmanConstVel(1e-6, 0.01, 0.0, 0.0, 1.0, 1.0)

	var got float64
	for range 5000 {
		got = k.Compute(truth, dt)
	}

	near(t, "const vel position", got, truth, 1e-3)
	near(t, "const vel position getter", k.Pos(), got, 0)
	near(t, "const vel velocity", k.Vel(), 0, 1e-3)
}

func TestKalmanConstVelTracksARamp(t *testing.T) {
	const (
		slope = 2.0
		dt    = 0.01
	)

	k := NewKalmanConstVel(1e-6, 0.01, 0.0, 0.0, 1.0, 1.0)

	var got float64
	n := 5000
	for i := 1; i <= n; i++ {
		got = k.Compute(slope*float64(i)*dt, dt)
	}

	near(t, "ramp position", got, slope*float64(n)*dt, 1e-2)
	near(t, "ramp velocity", k.Vel(), slope, 1e-2)
}

func TestKalmanConstVelCovarianceStaysSane(t *testing.T) {
	// the covariance must stay positive semi definite and finite
	k := NewKalmanConstVel(1e-4, 0.05, 0.0, 0.0, 1.0, 1.0)

	for i := range 2000 {
		k.Compute(math.Sin(0.05*float64(i)), 0.01)

		if math.IsNaN(k.Pos()) || math.IsInf(k.Pos(), 0) {
			t.Fatalf("position went non finite at step %d", i)
		}
		if math.IsNaN(k.Vel()) || math.IsInf(k.Vel(), 0) {
			t.Fatalf("velocity went non finite at step %d", i)
		}
		if k.p00 < 0 || k.p11 < 0 {
			t.Fatalf("negative variance at step %d: p00=%v p11=%v", i, k.p00, k.p11)
		}
		if k.p00*k.p11-k.p01*k.p01 < -1e-12 {
			t.Fatalf("covariance not positive semi definite at step %d: det=%v", i, k.p00*k.p11-k.p01*k.p01)
		}
	}
}

func TestKalmanConstVelIllConditioned(t *testing.T) {
	// with every variance zero the innovation covariance is zero, so the
	// filter must fall back to dead reckoning instead of dividing by zero
	const dt = 0.1
	k := NewKalmanConstVel(0.0, 0.0, 1.0, 2.0, 0.0, 0.0)

	for i := 1; i <= 10; i++ {
		got := k.Compute(1000, dt)
		want := 1 + 2*dt*float64(i)
		near(t, "dead reckoned position", got, want, 1e-9)
		near(t, "dead reckoned velocity", k.Vel(), 2, 0)
	}
}

func TestKalmanConstVelNonPositiveDt(t *testing.T) {
	k := NewKalmanConstVel(1e-4, 0.05, 1.0, 2.0, 1.0, 1.0)
	k.Compute(1.5, 0.1)

	pos, vel := k.Pos(), k.Vel()

	for _, dt := range []float64{0, -0.1} {
		if got := k.Compute(9, dt); got != 0 {
			t.Fatalf("Compute with dt=%v = %v, want 0", dt, got)
		}
	}

	// the rejected calls must not move the state
	near(t, "pos after rejected dt", k.Pos(), pos, 0)
	near(t, "vel after rejected dt", k.Vel(), vel, 0)
}

func TestKalmanConstVelReset(t *testing.T) {
	const (
		p0, v0 = 1.0, 1.0
		dt     = 0.02
	)

	k := NewKalmanConstVel(1e-4, 0.05, 0.0, 0.0, p0, v0)
	in := excitation(64)

	first := make([]float64, len(in))
	for i, x := range in {
		first[i] = k.Compute(x, dt)
	}

	k.Reset(0, 0, p0, v0)

	for i, x := range in {
		if got := k.Compute(x, dt); got != first[i] {
			t.Fatalf("KalmanConstVel: sample %d after reset = %v, want %v", i, got, first[i])
		}
	}
}

// ----------
// float32 instantiation
// ----------

func TestFloat32Instantiation(t *testing.T) {
	l := NewLowPass[float32](0.5)
	l.Compute(0)
	if got := l.Compute(1); got != 0.5 {
		t.Fatalf("float32 low pass = %v, want 0.5", got)
	}

	d := NewDeadZone[float32](1)
	if got := d.Compute(2); got != 1 {
		t.Fatalf("float32 dead zone = %v, want 1", got)
	}

	k := NewKalmanScalar[float32](0, 1, 1, 0)
	if got := k.Compute(2); got != 1 {
		t.Fatalf("float32 kalman = %v, want 1", got)
	}
}
