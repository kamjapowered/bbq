package math

import c "github.com/vistormu/go-dsa/constraints"

// return the sum of values
//
// time: O(n)
func Sum[T c.Number](values []T) T {
	var sum T
	for _, v := range values {
		sum += v
	}
	return sum
}
