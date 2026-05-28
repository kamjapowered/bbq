package math

import c "github.com/kamjapowered/bbq/constraints"

// return the arithmetic mean of values
//
// return 0 if values is empty
//
// time: O(n)
//
//microwave:export
func Mean[T c.Number](values []T) T {
	if len(values) == 0 {
		return 0
	}
	return Sum(values) / T(len(values))
}
