package math

import c "github.com/kamjapowered/bbq/constraints"

// clamp value to the inclusive range [lo, hi]
//
// if lo is greater than hi, it swaps them
//
// time: O(1)
//
//microwave:export
func Clip[T c.Number](value, lo, hi T) T {
	if lo > hi {
		lo, hi = hi, lo
	}

	return min(hi, max(lo, value))
}
