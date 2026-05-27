package math

import c "github.com/vistormu/go-dsa/constraints"

// return the absolute value
//
// time: O(1)
func Abs[T c.Number](v T) T {
	if v < 0 {
		return -v
	}
	return v
}
