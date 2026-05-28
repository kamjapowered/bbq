package math

import c "github.com/kamjapowered/bbq/constraints"

// return the sign of v as -1, 0, or 1
//
// time: O(1)
//
//microwave:export
func Sign[T c.Signed](v T) T {
	if v > 0 {
		return 1
	}
	if v < 0 {
		return -1
	}
	return 0
}

// return the sign of v as -1, 0, or 1
//
// time: O(1)
//
//microwave:export
func SignFloat[T c.Float](v T) T {
	if v > 0 {
		return 1
	}
	if v < 0 {
		return -1
	}
	return 0
}
