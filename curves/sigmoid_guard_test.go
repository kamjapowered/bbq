package curves

import (
	"testing"

	bmath "kamjapowered.com/bbq/math"
)

// a non positive k would divide by a zero span and yield NaN, so it is raised
func TestNewSigmoidGuardsNonPositiveK(t *testing.T) {
	for _, k := range []float64{0, -1, -10} {
		s := NewSigmoid(k, 0.5)
		for _, u := range []float64{0, 0.25, 0.5, 0.75, 1} {
			got := s.Compute(u)
			if got != got {
				t.Fatalf("NewSigmoid(%v).Compute(%v) = NaN", k, u)
			}
		}
		if got := s.Compute(0); !bmath.AlmostEqual(got, 0.0, eps) {
			t.Errorf("NewSigmoid(%v).Compute(0) = %v, want 0", k, got)
		}
		if got := s.Compute(1); !bmath.AlmostEqual(got, 1.0, eps) {
			t.Errorf("NewSigmoid(%v).Compute(1) = %v, want 1", k, got)
		}
	}
}
