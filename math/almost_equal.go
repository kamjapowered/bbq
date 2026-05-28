package math

import c "github.com/kamjapowered/bbq/constraints"

// report whether a and b are within eps
//
// time: O(1)
//
//microwave:export
func AlmostEqual[T c.Float](a, b, eps T) bool {
	if eps < 0 {
		eps = -eps
	}

	if a > b {
		return a-b <= eps
	}
	return b-a <= eps
}
