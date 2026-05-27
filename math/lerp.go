package math

import c "github.com/vistormu/go-dsa/constraints"

// linearly interpolate from a to b using t
//
// time: O(1)
func Lerp[T c.Float](a, b, t T) T {
	return a + t*(b-a)
}

// return t such that lerp(a, b, t) equals v
//
// return 0 if a equals b
//
// time: O(1)
func InvLerp[T c.Float](a, b, v T) T {
	if a == b {
		return 0
	}
	return (v - a) / (b - a)
}
