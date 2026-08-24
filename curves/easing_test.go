package curves

import (
	stdmath "math"
	"testing"

	bmath "kamjapowered.com/bbq/math"
)

const eps = 1e-9

// sample n + 1 evenly spaced values covering [0, 1]
func samples(n int) []float64 {
	out := make([]float64, n+1)
	for i := range out {
		out[i] = float64(i) / float64(n)
	}
	return out
}

// every scalar curve in the package, tagged with the properties its
// implementation actually guarantees
var easings = []struct {
	name string
	fn   func(float64) float64

	normalized bool // f(0) == 0 and f(1) == 1
	strict     bool // strictly increasing on [0, 1], every curve here is at least non decreasing
	symmetric  bool // f(u) + f(1-u) == 1
}{
	{"QuadIn", NewQuadIn[float64]().Compute, true, true, false},
	{"QuadOut", NewQuadOut[float64]().Compute, true, true, false},
	{"QuadInOut", NewQuadInOut[float64]().Compute, true, true, true},

	{"CubicIn", NewCubicIn[float64]().Compute, true, true, false},
	{"CubicOut", NewCubicOut[float64]().Compute, true, true, false},
	{"CubicInOut", NewCubicInOut[float64]().Compute, true, true, true},

	{"QuarticIn", NewQuarticIn[float64]().Compute, true, true, false},
	{"QuarticOut", NewQuarticOut[float64]().Compute, true, true, false},
	{"QuarticInOut", NewQuarticInOut[float64]().Compute, true, true, true},

	{"QuinticIn", NewQuinticIn[float64]().Compute, true, true, false},
	{"QuinticOut", NewQuinticOut[float64]().Compute, true, true, false},
	{"QuinticInOut", NewQuinticInOut[float64]().Compute, true, true, true},

	{"SineIn", NewSineIn[float64]().Compute, true, true, false},
	{"SineOut", NewSineOut[float64]().Compute, true, true, false},
	{"SineInOut", NewSineInOut[float64]().Compute, true, true, true},

	{"ExpoIn", NewExpoIn[float64]().Compute, true, true, false},
	{"ExpoOut", NewExpoOut[float64]().Compute, true, true, false},
	{"ExpoInOut", NewExpoInOut[float64]().Compute, true, true, true},

	{"CircIn", NewCircIn[float64]().Compute, true, true, false},
	{"CircOut", NewCircOut[float64]().Compute, true, true, false},
	{"CircInOut", NewCircInOut[float64]().Compute, true, true, true},

	{"Smoothstep", NewSmoothstep[float64]().Compute, true, true, true},
	{"Smootherstep", NewSmootherstep[float64]().Compute, true, true, true},

	{"Sigmoid(10,0.5)", NewSigmoid[float64](10, 0.5).Compute, true, true, true},
	{"Sigmoid(4,0.3)", NewSigmoid[float64](4, 0.3).Compute, true, true, false},
	{"Sigmoid(1,0.5)", NewSigmoid[float64](1, 0.5).Compute, true, true, true},

	// discrete curves, monotonic but only weakly and never symmetric
	{"Step(0.5)", NewStep[float64](0.5).Compute, true, false, false},
	{"Step(0.25)", NewStep[float64](0.25).Compute, true, false, false},
	{"Staircase(2)", NewStaircase[float64](2).Compute, true, false, false},
	{"Staircase(4)", NewStaircase[float64](4).Compute, true, false, false},
	{"Staircase(10)", NewStaircase[float64](10).Compute, true, false, false},
}

func TestEndpoints(t *testing.T) {
	for _, tt := range easings {
		if !tt.normalized {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(0); !bmath.AlmostEqual(got, 0, eps) {
				t.Errorf("f(0) = %v, want 0", got)
			}
			if got := tt.fn(1); !bmath.AlmostEqual(got, 1, eps) {
				t.Errorf("f(1) = %v, want 1", got)
			}
		})
	}
}

func TestMonotonic(t *testing.T) {
	us := samples(200)
	for _, tt := range easings {
		t.Run(tt.name, func(t *testing.T) {
			prev := tt.fn(us[0])
			for _, u := range us[1:] {
				v := tt.fn(u)
				if v < prev {
					t.Fatalf("f(%v) = %v decreased from %v", u, v, prev)
				}
				if tt.strict && v <= prev {
					t.Fatalf("f(%v) = %v did not increase from %v", u, v, prev)
				}
				prev = v
			}
		})
	}
}

func TestBounds(t *testing.T) {
	for _, tt := range easings {
		t.Run(tt.name, func(t *testing.T) {
			for _, u := range samples(200) {
				v := tt.fn(u)
				if stdmath.IsNaN(v) {
					t.Fatalf("f(%v) is NaN", u)
				}
				if v < -eps || v > 1+eps {
					t.Fatalf("f(%v) = %v, outside [0, 1]", u, v)
				}
			}
		})
	}
}

func TestSymmetry(t *testing.T) {
	for _, tt := range easings {
		if !tt.symmetric {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			for _, u := range samples(200) {
				a, b := tt.fn(u), tt.fn(1-u)
				if !bmath.AlmostEqual(a+b, 1, eps) {
					t.Fatalf("f(%v) + f(%v) = %v, want 1", u, 1-u, a+b)
				}
			}
		})
	}
}

// every in/out pair must satisfy out(u) == 1 - in(1-u)
func TestInOutDuality(t *testing.T) {
	pairs := []struct {
		name    string
		in, out func(float64) float64
	}{
		{"Quad", NewQuadIn[float64]().Compute, NewQuadOut[float64]().Compute},
		{"Cubic", NewCubicIn[float64]().Compute, NewCubicOut[float64]().Compute},
		{"Quartic", NewQuarticIn[float64]().Compute, NewQuarticOut[float64]().Compute},
		{"Quintic", NewQuinticIn[float64]().Compute, NewQuinticOut[float64]().Compute},
		{"Sine", NewSineIn[float64]().Compute, NewSineOut[float64]().Compute},
		{"Expo", NewExpoIn[float64]().Compute, NewExpoOut[float64]().Compute},
		{"Circ", NewCircIn[float64]().Compute, NewCircOut[float64]().Compute},
	}

	for _, tt := range pairs {
		t.Run(tt.name, func(t *testing.T) {
			for _, u := range samples(200) {
				want := 1 - tt.in(1-u)
				if got := tt.out(u); !bmath.AlmostEqual(got, want, eps) {
					t.Fatalf("out(%v) = %v, want %v", u, got, want)
				}
			}
		})
	}
}

// the in out variants must meet in the middle without a jump
//
// the gap across the midpoint must shrink as the offset shrinks, this is plain
// continuity rather than a bounded slope, CircInOut has a vertical tangent at
// u = 0.5 so its gap only closes like sqrt(h)
func TestInOutMidpointContinuity(t *testing.T) {
	fns := []struct {
		name string
		fn   func(float64) float64
	}{
		{"QuadInOut", NewQuadInOut[float64]().Compute},
		{"CubicInOut", NewCubicInOut[float64]().Compute},
		{"QuarticInOut", NewQuarticInOut[float64]().Compute},
		{"QuinticInOut", NewQuinticInOut[float64]().Compute},
		{"SineInOut", NewSineInOut[float64]().Compute},
		{"ExpoInOut", NewExpoInOut[float64]().Compute},
		{"CircInOut", NewCircInOut[float64]().Compute},
	}

	offsets := []float64{1e-3, 1e-5, 1e-7, 1e-9}

	for _, tt := range fns {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(0.5); !bmath.AlmostEqual(got, 0.5, eps) {
				t.Errorf("f(0.5) = %v, want 0.5", got)
			}

			prev := stdmath.Inf(1)
			for _, h := range offsets {
				gap := tt.fn(0.5+h) - tt.fn(0.5-h)
				if gap < 0 {
					t.Fatalf("h = %v, f(0.5+h) - f(0.5-h) = %v, want non negative", h, gap)
				}
				if gap > prev {
					t.Fatalf("h = %v, gap %v grew from %v, not continuous", h, gap, prev)
				}
				prev = gap
			}
			if prev > 1e-4 {
				t.Errorf("gap at h = %v is %v, want it closed", offsets[len(offsets)-1], prev)
			}
		})
	}
}

// spot checks against values worked out by hand
func TestKnownValues(t *testing.T) {
	tests := []struct {
		name string
		got  float64
		want float64
	}{
		{"QuadIn(0.5)", NewQuadIn[float64]().Compute(0.5), 0.25},
		{"QuadOut(0.5)", NewQuadOut[float64]().Compute(0.5), 0.75},
		{"CubicIn(0.5)", NewCubicIn[float64]().Compute(0.5), 0.125},
		{"QuarticIn(0.5)", NewQuarticIn[float64]().Compute(0.5), 0.0625},
		{"QuinticIn(0.5)", NewQuinticIn[float64]().Compute(0.5), 0.03125},
		{"SineOut(0.5)", NewSineOut[float64]().Compute(0.5), stdmath.Sqrt2 / 2},
		{"SineIn(0.5)", NewSineIn[float64]().Compute(0.5), 1 - stdmath.Sqrt2/2},
		{"ExpoOut(0.5)", NewExpoOut[float64]().Compute(0.5), 0.96875},
		{"ExpoIn(0.5)", NewExpoIn[float64]().Compute(0.5), 0.03125},
		{"CircIn(0.5)", NewCircIn[float64]().Compute(0.5), 1 - stdmath.Sqrt(0.75)},
		{"CircOut(0.5)", NewCircOut[float64]().Compute(0.5), stdmath.Sqrt(0.75)},
		{"Smoothstep(0.25)", NewSmoothstep[float64]().Compute(0.25), 0.15625},
		{"Smootherstep(0.25)", NewSmootherstep[float64]().Compute(0.25), 0.103515625},
		{"Step(0.5)@0.4999", NewStep[float64](0.5).Compute(0.4999), 0},
		{"Step(0.5)@0.5", NewStep[float64](0.5).Compute(0.5), 1},
		{"Staircase(4)@0.5", NewStaircase[float64](4).Compute(0.5), 2.0 / 3.0},
		{"Staircase(4)@0.1", NewStaircase[float64](4).Compute(0.1), 0},
		{"Staircase(4)@0.3", NewStaircase[float64](4).Compute(0.3), 1.0 / 3.0},
	}

	for _, tt := range tests {
		if !bmath.AlmostEqual(tt.got, tt.want, eps) {
			t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

// a staircase must only ever emit the levels k / (n - 1)
func TestStaircaseLevels(t *testing.T) {
	for _, n := range []int{2, 3, 5, 8} {
		s := NewStaircase[float64](n)
		seen := map[float64]bool{}
		for _, u := range samples(1000) {
			v := s.Compute(u)
			level := v * float64(n-1)
			if !bmath.AlmostEqual(level, stdmath.Round(level), 1e-9) {
				t.Fatalf("n = %d, f(%v) = %v is not a multiple of 1/%d", n, u, v, n-1)
			}
			seen[stdmath.Round(level)] = true
		}
		if len(seen) != n {
			t.Errorf("n = %d, got %d distinct levels, want %d", n, len(seen), n)
		}
	}
}

// fewer than two steps is clamped to two
func TestStaircaseClampsN(t *testing.T) {
	for _, n := range []int{-1, 0, 1, 2} {
		s := NewStaircase[float64](n)
		if got := s.Compute(0.25); got != 0 {
			t.Errorf("n = %d, f(0.25) = %v, want 0", n, got)
		}
		if got := s.Compute(0.75); got != 1 {
			t.Errorf("n = %d, f(0.75) = %v, want 1", n, got)
		}
	}
}

// the threshold itself belongs to the upper branch
func TestStepThreshold(t *testing.T) {
	for _, u0 := range []float64{0, 0.25, 0.5, 1} {
		s := NewStep[float64](u0)
		if got := s.Compute(u0); got != 1 {
			t.Errorf("u0 = %v, f(u0) = %v, want 1", u0, got)
		}
		if u0 > 0 {
			if got := s.Compute(stdmath.Nextafter(u0, 0)); got != 0 {
				t.Errorf("u0 = %v, f(u0-) = %v, want 0", u0, got)
			}
		}
	}
}

// the normalization must hold for any steepness and midpoint
func TestSigmoidNormalization(t *testing.T) {
	for _, k := range []float64{0.5, 1, 4, 10, 25} {
		for _, x0 := range []float64{0.1, 0.3, 0.5, 0.7, 0.9} {
			s := NewSigmoid[float64](k, x0)
			if got := s.Compute(0); !bmath.AlmostEqual(got, 0, eps) {
				t.Errorf("k = %v, x0 = %v, f(0) = %v, want 0", k, x0, got)
			}
			if got := s.Compute(1); !bmath.AlmostEqual(got, 1, eps) {
				t.Errorf("k = %v, x0 = %v, f(1) = %v, want 1", k, x0, got)
			}
		}
	}
}

// the curves are generic over float32 as well
func TestFloat32Easings(t *testing.T) {
	const eps32 = float32(1e-6)
	fns := []func(float32) float32{
		NewQuadInOut[float32]().Compute,
		NewCubicInOut[float32]().Compute,
		NewQuarticInOut[float32]().Compute,
		NewQuinticInOut[float32]().Compute,
		NewSineInOut[float32]().Compute,
		NewExpoInOut[float32]().Compute,
		NewCircInOut[float32]().Compute,
		NewSmoothstep[float32]().Compute,
		NewSmootherstep[float32]().Compute,
		NewSigmoid[float32](8, 0.5).Compute,
		NewStaircase[float32](4).Compute,
	}

	for i, fn := range fns {
		if got := fn(0); !bmath.AlmostEqual(got, 0, eps32) {
			t.Errorf("curve %d, f(0) = %v, want 0", i, got)
		}
		if got := fn(1); !bmath.AlmostEqual(got, 1, eps32) {
			t.Errorf("curve %d, f(1) = %v, want 1", i, got)
		}
	}
}

func BenchmarkEasings(b *testing.B) {
	for _, tt := range easings {
		b.Run(tt.name, func(b *testing.B) {
			for b.Loop() {
				_ = tt.fn(0.37)
			}
		})
	}
}
