package math

import c "github.com/vistormu/go-dsa/constraints"

// return the minimum and maximum value in values
//
// return 0, 0 if values is empty
//
// time: O(n)
func MinMax[T c.Number](values []T) (T, T) {
	if len(values) == 0 {
		return 0, 0
	}

	lo, hi := values[0], values[0]
	for _, v := range values[1:] {
		lo = min(lo, v)
		hi = max(hi, v)
	}
	return lo, hi
}
