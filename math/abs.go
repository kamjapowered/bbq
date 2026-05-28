package math

import c "github.com/kamjapowered/bbq/constraints"

// return the absolute value
//
// time: O(1)
//
//microwave:export
func Abs[T c.Number](v T) T {
	if v < 0 {
		return -v
	}
	return v
}
