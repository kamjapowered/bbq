package math

import c "github.com/vistormu/go-dsa/constraints"

// return the population variance of values
//
// return 0 if values is empty
//
// time: O(n)
func Variance[T c.Float](values []T) T {
	if len(values) == 0 {
		return 0
	}

	m := Mean(values)

	var acc T
	for _, v := range values {
		d := v - m
		acc += d * d
	}

	return acc / T(len(values))
}
