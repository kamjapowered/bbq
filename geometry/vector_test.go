package geometry

import (
	"testing"

	bmath "kamjapowered.com/bbq/math"
)

const eps = 1e-9

// report a mismatch between two float vectors within eps
func checkVec(t *testing.T, what string, got, want Vector[float64]) {
	t.Helper()
	if !bmath.AlmostEqual(got.X, want.X, eps) ||
		!bmath.AlmostEqual(got.Y, want.Y, eps) ||
		!bmath.AlmostEqual(got.Z, want.Z, eps) {
		t.Fatalf("%s = %v, want %v", what, got, want)
	}
}

func TestNewVector(t *testing.T) {
	tests := []struct {
		name string
		got  Vector[int]
		want Vector[int]
	}{
		{"empty", NewVector[int](), Vector[int]{}},
		{"x", NewVector(1), Vector[int]{X: 1}},
		{"xy", NewVector(1, 2), Vector[int]{X: 1, Y: 2}},
		{"xyz", NewVector(1, 2, 3), Vector[int]{X: 1, Y: 2, Z: 3}},
		// anything past the third component is dropped
		{"extra", NewVector(1, 2, 3, 4, 5), Vector[int]{X: 1, Y: 2, Z: 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("NewVector = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestVectorAddSubScale(t *testing.T) {
	a := NewVector(1, 2, 3)
	b := NewVector(4, -5, 6)

	if got, want := a.Add(b), NewVector(5, -3, 9); got != want {
		t.Fatalf("Add = %v, want %v", got, want)
	}

	if got, want := a.Sub(b), NewVector(-3, 7, -3); got != want {
		t.Fatalf("Sub = %v, want %v", got, want)
	}

	if got, want := a.Scale(3), NewVector(3, 6, 9); got != want {
		t.Fatalf("Scale = %v, want %v", got, want)
	}

	if got, want := a.Scale(0), (Vector[int]{}); got != want {
		t.Fatalf("Scale(0) = %v, want %v", got, want)
	}

	// the operands are values, so none of the above may have touched them
	if a != NewVector(1, 2, 3) || b != NewVector(4, -5, 6) {
		t.Fatalf("operands mutated: a = %v, b = %v", a, b)
	}

	// a-b is the inverse of a+b
	if got := a.Add(b).Sub(b); got != a {
		t.Fatalf("Add then Sub = %v, want %v", got, a)
	}
}

func TestVectorDot(t *testing.T) {
	tests := []struct {
		name string
		a, b Vector[int]
		want int
	}{
		{"known", NewVector(1, 2, 3), NewVector(4, -5, 6), 12},
		{"zero", NewVector(1, 2, 3), Vector[int]{}, 0},
		{"orthogonal", NewVector(1, 0, 0), NewVector(0, 1, 0), 0},
		{"antiparallel", NewVector(2, 0, 0), NewVector(-3, 0, 0), -6},
		{"self is lensq", NewVector(3, 4, 0), NewVector(3, 4, 0), 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Dot(tt.b); got != tt.want {
				t.Fatalf("Dot = %d, want %d", got, tt.want)
			}
			// dot is commutative
			if got := tt.b.Dot(tt.a); got != tt.want {
				t.Fatalf("reversed Dot = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestVectorCross(t *testing.T) {
	tests := []struct {
		name string
		a, b Vector[int]
		want Vector[int]
	}{
		{"known", NewVector(1, 2, 3), NewVector(4, 5, 6), NewVector(-3, 6, -3)},
		{"x cross y", NewVector(1, 0, 0), NewVector(0, 1, 0), NewVector(0, 0, 1)},
		{"y cross z", NewVector(0, 1, 0), NewVector(0, 0, 1), NewVector(1, 0, 0)},
		{"z cross x", NewVector(0, 0, 1), NewVector(1, 0, 0), NewVector(0, 1, 0)},
		{"parallel", NewVector(2, 4, 6), NewVector(1, 2, 3), Vector[int]{}},
		{"self", NewVector(1, 2, 3), NewVector(1, 2, 3), Vector[int]{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Cross(tt.b)
			if got != tt.want {
				t.Fatalf("Cross = %v, want %v", got, tt.want)
			}

			// cross is anticommutative
			if rev, want := tt.b.Cross(tt.a), tt.want.Scale(-1); rev != want {
				t.Fatalf("reversed Cross = %v, want %v", rev, want)
			}

			// the result is orthogonal to both operands
			if d := got.Dot(tt.a); d != 0 {
				t.Fatalf("Cross dot a = %d, want 0", d)
			}
			if d := got.Dot(tt.b); d != 0 {
				t.Fatalf("Cross dot b = %d, want 0", d)
			}
		})
	}
}

func TestVectorLen(t *testing.T) {
	tests := []struct {
		name  string
		v     Vector[int]
		lenSq int
		len   float64
	}{
		{"zero", Vector[int]{}, 0, 0},
		{"unit x", NewVector(1, 0, 0), 1, 1},
		{"3-4-5", NewVector(3, 4, 0), 25, 5},
		{"1-2-2", NewVector(1, 2, 2), 9, 3},
		{"negative", NewVector(-3, -4, 0), 25, 5},
		{"2-3-6", NewVector(2, 3, 6), 49, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.LenSq(); got != tt.lenSq {
				t.Fatalf("LenSq = %d, want %d", got, tt.lenSq)
			}
			if got := tt.v.Len(); !bmath.AlmostEqual(got, tt.len, eps) {
				t.Fatalf("Len = %f, want %f", got, tt.len)
			}
			// len squared must agree with lensq
			if got := tt.v.Len() * tt.v.Len(); !bmath.AlmostEqual(got, float64(tt.lenSq), eps) {
				t.Fatalf("Len squared = %f, want %d", got, tt.lenSq)
			}
		})
	}
}

func TestVectorLenFloat(t *testing.T) {
	v := NewVector(0.3, 0.4, 0.0)

	if got := v.LenSq(); !bmath.AlmostEqual(got, 0.25, eps) {
		t.Fatalf("LenSq = %f, want %f", got, 0.25)
	}
	if got := v.Len(); !bmath.AlmostEqual(got, 0.5, eps) {
		t.Fatalf("Len = %f, want %f", got, 0.5)
	}
}

func TestVectorNorm(t *testing.T) {
	tests := []struct {
		name string
		v    Vector[float64]
		want Vector[float64]
	}{
		{"unit x", NewVector(1.0, 0.0, 0.0), NewVector(1.0, 0.0, 0.0)},
		{"3-4-0", NewVector(3.0, 4.0, 0.0), NewVector(0.6, 0.8, 0.0)},
		{"1-2-2", NewVector(1.0, 2.0, 2.0), NewVector(1.0/3, 2.0/3, 2.0/3)},
		{"negative", NewVector(-3.0, -4.0, 0.0), NewVector(-0.6, -0.8, 0.0)},
		{"already unit", NewVector(0.0, 0.0, -1.0), NewVector(0.0, 0.0, -1.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.Norm()
			checkVec(t, "Norm", got, tt.want)

			// a normalized vector has unit length
			if l := got.Len(); !bmath.AlmostEqual(l, 1.0, eps) {
				t.Fatalf("Norm length = %f, want 1", l)
			}
		})
	}
}

func TestVectorNormFromInt(t *testing.T) {
	// norm widens to float64 whatever the source component type is
	got := NewVector(3, 4, 0).Norm()
	checkVec(t, "Norm", got, NewVector(0.6, 0.8, 0.0))
}

func TestVectorNormZero(t *testing.T) {
	// a zero vector has no direction, and the documented answer is the zero
	// vector rather than a NaN one
	for _, v := range []Vector[float64]{{}, NewVector(0.0, 0.0, 0.0)} {
		got := v.Norm()
		if got != (Vector[float64]{}) {
			t.Fatalf("Norm of zero vector = %v, want %v", got, Vector[float64]{})
		}
		if l := got.Len(); l != 0 {
			t.Fatalf("Norm of zero vector has length %f, want 0", l)
		}
	}

	if got := (Vector[int]{}).Norm(); got != (Vector[float64]{}) {
		t.Fatalf("Norm of zero int vector = %v, want %v", got, Vector[float64]{})
	}
}
