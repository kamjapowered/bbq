package math

import c "github.com/vistormu/go-dsa/constraints"

// return the arithmetic mean of values
//
// return 0 if values is empty
//
// time: O(n)
func Mean[T c.Number](values []T) T {
	if len(values) == 0 {
		return 0
	}
	return Sum(values) / T(len(values))
}
