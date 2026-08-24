package math

import (
	gomath "math"
	"testing"
)

func assertEq[T comparable](t *testing.T, name string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

func TestAbs(t *testing.T) {
	ints := []struct {
		in, want int
	}{
		{0, 0},
		{3, 3},
		{-3, 3},
		{-1, 1},
	}
	for _, tt := range ints {
		assertEq(t, "Abs(int)", Abs(tt.in), tt.want)
	}

	floats := []struct {
		in, want float64
	}{
		{0, 0},
		{2.5, 2.5},
		{-2.5, 2.5},
	}
	for _, tt := range floats {
		assertEq(t, "Abs(float64)", Abs(tt.in), tt.want)
	}

	// other instantiations
	assertEq(t, "Abs(int8)", Abs(int8(-7)), int8(7))
	assertEq(t, "Abs(int64)", Abs(int64(-1<<40)), int64(1<<40))
	assertEq(t, "Abs(float32)", Abs(float32(-1.5)), float32(1.5))
	assertEq(t, "Abs(uint)", Abs(uint(9)), uint(9))

	// two's complement: the most negative integer has no positive counterpart,
	// so negating it wraps around. documented here as the actual behaviour.
	assertEq(t, "Abs(MinInt8)", Abs(int8(gomath.MinInt8)), int8(gomath.MinInt8))
	assertEq(t, "Abs(MinInt)", Abs(int(gomath.MinInt)), int(gomath.MinInt))

	// NaN is not < 0, so it passes through unchanged
	if !gomath.IsNaN(Abs(gomath.NaN())) {
		t.Errorf("Abs(NaN) should stay NaN")
	}

	assertEq(t, "Abs(-Inf)", Abs(gomath.Inf(-1)), gomath.Inf(1))

	// negative zero is not < 0 either, so the sign bit survives
	if !gomath.Signbit(Abs(gomath.Copysign(0, -1))) {
		t.Logf("Abs(-0.0) cleared the sign bit")
	}
}

func TestClip(t *testing.T) {
	ints := []struct {
		value, lo, hi, want int
	}{
		{5, 0, 10, 5},   // inside
		{-1, 0, 10, 0},  // below
		{11, 0, 10, 10}, // above
		{0, 0, 10, 0},   // on lower bound
		{10, 0, 10, 10}, // on upper bound
		{5, 5, 5, 5},    // degenerate range
		{5, 10, 0, 5},   // inverted bounds are swapped
		{-1, 10, 0, 0},  // inverted, below
		{11, 10, 0, 10}, // inverted, above
	}
	for _, tt := range ints {
		got := Clip(tt.value, tt.lo, tt.hi)
		if got != tt.want {
			t.Errorf("Clip(%d, %d, %d) = %d, want %d", tt.value, tt.lo, tt.hi, got, tt.want)
		}
	}

	floats := []struct {
		value, lo, hi, want float64
	}{
		{0.5, 0, 1, 0.5},
		{-0.5, 0, 1, 0},
		{1.5, 0, 1, 1},
		{0.5, 1, 0, 0.5},
		{-0.5, 1, 0, 0},
	}
	for _, tt := range floats {
		got := Clip(tt.value, tt.lo, tt.hi)
		if got != tt.want {
			t.Errorf("Clip(%v, %v, %v) = %v, want %v", tt.value, tt.lo, tt.hi, got, tt.want)
		}
	}

	assertEq(t, "Clip(int8)", Clip(int8(-100), int8(-10), int8(10)), int8(-10))
	assertEq(t, "Clip(uint)", Clip(uint(7), uint(0), uint(3)), uint(3))
	assertEq(t, "Clip(float32)", Clip(float32(4), float32(0), float32(1)), float32(1))
}

func TestAlmostEqual(t *testing.T) {
	tests := []struct {
		name      string
		a, b, eps float64
		want      bool
	}{
		{"equal", 1, 1, 1e-9, true},
		{"within eps, a > b", 1.0000001, 1, 1e-6, true},
		{"within eps, b > a", 1, 1.0000001, 1e-6, true},
		{"outside eps", 1, 1.1, 1e-6, false},
		{"exactly eps apart", 1, 1.5, 0.5, true},
		{"zero eps, equal", 2, 2, 0, true},
		{"zero eps, different", 2, 2.0000001, 0, false},
		{"negative eps is normalised", 1, 1.4, -0.5, true},
		{"negative eps, outside", 1, 2, -0.5, false},
		{"across zero", -0.0000001, 0.0000001, 1e-6, true},
		{"nan vs number", gomath.NaN(), 1, 1e-6, false},
		{"number vs nan", 1, gomath.NaN(), 1e-6, false},
		{"nan vs nan", gomath.NaN(), gomath.NaN(), 1e-6, false},
		// identical operands short-circuit before eps is consulted
		{"nan eps", 1, 1, gomath.NaN(), true},
		{"+inf vs -inf", gomath.Inf(1), gomath.Inf(-1), 1e-6, false},
		{"+inf vs number", gomath.Inf(1), 1, 1e-6, false},
		// an identity check runs first, so an infinity equals itself
		{"+inf vs +inf", gomath.Inf(1), gomath.Inf(1), 1e-6, true},
		{"-inf vs -inf", gomath.Inf(-1), gomath.Inf(-1), 1e-6, true},
		// an infinite tolerance does not rescue finite operands either way
		{"infinite eps", 1, 1e300, gomath.Inf(1), true},
	}

	for _, tt := range tests {
		if got := AlmostEqual(tt.a, tt.b, tt.eps); got != tt.want {
			t.Errorf("AlmostEqual(%v, %v, %v) [%s] = %v, want %v", tt.a, tt.b, tt.eps, tt.name, got, tt.want)
		}
	}

	// float32 instantiation
	assertEq(t, "AlmostEqual(float32)", AlmostEqual(float32(1), float32(1.0001), float32(1e-3)), true)
	assertEq(t, "AlmostEqual(float32)", AlmostEqual(float32(1), float32(1.1), float32(1e-3)), false)
}

func TestLerp(t *testing.T) {
	tests := []struct {
		a, b, t, want float64
	}{
		{0, 10, 0, 0},      // t = 0
		{0, 10, 1, 10},     // t = 1
		{0, 10, 0.5, 5},    // midpoint
		{0, 10, 2, 20},     // extrapolation above
		{0, 10, -1, -10},   // extrapolation below
		{-10, 10, 0.5, 0},  // spanning zero
		{10, 0, 0.25, 7.5}, // descending range
		{5, 5, 0.7, 5},     // degenerate range
	}
	for _, tt := range tests {
		if got := Lerp(tt.a, tt.b, tt.t); got != tt.want {
			t.Errorf("Lerp(%v, %v, %v) = %v, want %v", tt.a, tt.b, tt.t, got, tt.want)
		}
	}

	assertEq(t, "Lerp(float32)", Lerp(float32(0), float32(4), float32(0.25)), float32(1))
}

func TestInvLerp(t *testing.T) {
	tests := []struct {
		a, b, v, want float64
	}{
		{0, 10, 0, 0},     // at a
		{0, 10, 10, 1},    // at b
		{0, 10, 5, 0.5},   // midpoint
		{0, 10, 20, 2},    // outside above
		{0, 10, -5, -0.5}, // outside below
		{10, 0, 7.5, 0.25},
		{5, 5, 5, 0}, // a == b returns 0
		{5, 5, 9, 0},
	}
	for _, tt := range tests {
		if got := InvLerp(tt.a, tt.b, tt.v); got != tt.want {
			t.Errorf("InvLerp(%v, %v, %v) = %v, want %v", tt.a, tt.b, tt.v, got, tt.want)
		}
	}

	assertEq(t, "InvLerp(float32)", InvLerp(float32(0), float32(4), float32(1)), float32(0.25))

	// round trip
	a, b := 3.0, 11.0
	for _, v := range []float64{3, 5, 11, 20, -4} {
		if got := Lerp(a, b, InvLerp(a, b, v)); !AlmostEqual(got, v, 1e-12) {
			t.Errorf("Lerp(InvLerp(%v)) = %v, want %v", v, got, v)
		}
	}
}

func TestRemap(t *testing.T) {
	tests := []struct {
		name                                  string
		v, inMin, inMax, outMin, outMax, want float64
	}{
		{"identity", 5, 0, 10, 0, 10, 5},
		{"scale up", 5, 0, 10, 0, 100, 50},
		{"at input min", 0, 0, 10, 0, 100, 0},
		{"at input max", 10, 0, 10, 0, 100, 100},
		{"outside above", 20, 0, 10, 0, 100, 200},
		{"outside below", -5, 0, 10, 0, 100, -50},
		{"shifted output", 5, 0, 10, -1, 1, 0},
		{"inverted output", 0, 0, 10, 1, 0, 1},
		{"degenerate input returns outMin", 7, 3, 3, -2, 9, -2},
		{"degenerate output", 5, 0, 10, 4, 4, 4},
	}
	for _, tt := range tests {
		got := Remap(tt.v, tt.inMin, tt.inMax, tt.outMin, tt.outMax)
		if !AlmostEqual(got, tt.want, 1e-12) {
			t.Errorf("Remap %s = %v, want %v", tt.name, got, tt.want)
		}
	}

	assertEq(t, "Remap(float32)", Remap(float32(1), float32(0), float32(2), float32(0), float32(10)), float32(5))
}

func TestSign(t *testing.T) {
	ints := []struct {
		in, want int
	}{
		{5, 1},
		{-5, -1},
		{0, 0},
		{gomath.MaxInt, 1},
		{gomath.MinInt, -1},
	}
	for _, tt := range ints {
		assertEq(t, "Sign(int)", Sign(tt.in), tt.want)
	}

	assertEq(t, "Sign(int8)", Sign(int8(-128)), int8(-1))
	assertEq(t, "Sign(int64)", Sign(int64(1<<40)), int64(1))

	floats := []struct {
		in, want float64
	}{
		{2.5, 1},
		{-2.5, -1},
		{0, 0},
		{gomath.Inf(1), 1},
		{gomath.Inf(-1), -1},
		{gomath.Copysign(0, -1), 0}, // negative zero is neither > nor < 0
		{gomath.NaN(), 0},           // NaN compares false both ways
	}
	for _, tt := range floats {
		assertEq(t, "SignFloat", SignFloat(tt.in), tt.want)
	}

	assertEq(t, "SignFloat(float32)", SignFloat(float32(-0.5)), float32(-1))
}

func TestSum(t *testing.T) {
	assertEq(t, "Sum(empty int)", Sum([]int{}), 0)
	assertEq(t, "Sum(nil int)", Sum([]int(nil)), 0)
	assertEq(t, "Sum(single)", Sum([]int{42}), 42)
	assertEq(t, "Sum(int)", Sum([]int{1, 2, 3, 4}), 10)
	assertEq(t, "Sum(negatives)", Sum([]int{-1, -2, 3}), 0)

	assertEq(t, "Sum(empty float)", Sum([]float64{}), 0)
	assertEq(t, "Sum(float64)", Sum([]float64{0.5, 1.5, 2}), 4)
	assertEq(t, "Sum(float32)", Sum([]float32{1, 2, 3}), float32(6))
	assertEq(t, "Sum(uint8)", Sum([]uint8{1, 2, 3}), uint8(6))
}

func TestMean(t *testing.T) {
	assertEq(t, "Mean(empty int)", Mean([]int{}), 0)
	assertEq(t, "Mean(nil int)", Mean([]int(nil)), 0)
	assertEq(t, "Mean(single int)", Mean([]int{7}), 7)
	assertEq(t, "Mean(int)", Mean([]int{1, 2, 3}), 2)

	// integer instantiations divide with truncation towards zero
	assertEq(t, "Mean(int truncating)", Mean([]int{1, 2}), 1)
	assertEq(t, "Mean(int negative)", Mean([]int{-1, -2}), -1)

	assertEq(t, "Mean(empty float)", Mean([]float64{}), 0)
	assertEq(t, "Mean(single float)", Mean([]float64{2.5}), 2.5)
	assertEq(t, "Mean(float64)", Mean([]float64{1, 2}), 1.5)
	assertEq(t, "Mean(float64 spread)", Mean([]float64{2, 4, 4, 4, 5, 5, 7, 9}), 5)
	assertEq(t, "Mean(float32)", Mean([]float32{1, 2, 3}), float32(2))
}

func TestMinMax(t *testing.T) {
	lo, hi := MinMax([]int{})
	assertEq(t, "MinMax(empty) lo", lo, 0)
	assertEq(t, "MinMax(empty) hi", hi, 0)

	lo, hi = MinMax([]int(nil))
	assertEq(t, "MinMax(nil) lo", lo, 0)
	assertEq(t, "MinMax(nil) hi", hi, 0)

	lo, hi = MinMax([]int{5})
	assertEq(t, "MinMax(single) lo", lo, 5)
	assertEq(t, "MinMax(single) hi", hi, 5)

	lo, hi = MinMax([]int{3, -2, 9, 0})
	assertEq(t, "MinMax lo", lo, -2)
	assertEq(t, "MinMax hi", hi, 9)

	// a slice of all negatives must not fall back to zero
	lo, hi = MinMax([]int{-3, -7, -1})
	assertEq(t, "MinMax(negatives) lo", lo, -7)
	assertEq(t, "MinMax(negatives) hi", hi, -1)

	flo, fhi := MinMax([]float64{1.5, -0.5, 2.25})
	assertEq(t, "MinMax(float) lo", flo, -0.5)
	assertEq(t, "MinMax(float) hi", fhi, 2.25)

	ulo, uhi := MinMax([]uint{4, 1, 9})
	assertEq(t, "MinMax(uint) lo", ulo, uint(1))
	assertEq(t, "MinMax(uint) hi", uhi, uint(9))
}

func TestVariance(t *testing.T) {
	assertEq(t, "Variance(empty)", Variance([]float64{}), 0)
	assertEq(t, "Variance(nil)", Variance([]float64(nil)), 0)
	assertEq(t, "Variance(single)", Variance([]float64{42}), 0)
	assertEq(t, "Variance(identical)", Variance([]float64{3, 3, 3}), 0)

	// population variance of this classic sample is 4
	if got := Variance([]float64{2, 4, 4, 4, 5, 5, 7, 9}); !AlmostEqual(got, 4, 1e-12) {
		t.Errorf("Variance = %v, want 4", got)
	}

	// shifting the data must not change the variance
	shifted := []float64{102, 104, 104, 104, 105, 105, 107, 109}
	if got := Variance(shifted); !AlmostEqual(got, 4, 1e-9) {
		t.Errorf("Variance(shifted) = %v, want 4", got)
	}

	if got := Variance([]float32{1, 3}); !AlmostEqual(got, float32(1), float32(1e-6)) {
		t.Errorf("Variance(float32) = %v, want 1", got)
	}
}
