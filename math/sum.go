package math

import c "github.com/kamjapowered/bbq/constraints"

// return the sum of values
//
// time: O(n)
//
//microwave:export
func Sum[T c.Number](values []T) T {
	var sum T
	for _, v := range values {
		sum += v
	}
	return sum
}
