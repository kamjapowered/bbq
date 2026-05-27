package filter

import c "github.com/vistormu/go-dsa/constraints"

// suppress values around zero with a dead zone
type DeadZone[T c.Float] struct {
	width T
}

// create a dead zone with half width w
//
// if w is negative, it is treated as zero
func NewDeadZone[T c.Float](w T) DeadZone[T] {
	if w < 0 {
		w = 0
	}
	return DeadZone[T]{width: w}
}

// compute the dead zoned value
//
// time: O(1)
func (d DeadZone[T]) Compute(x T) T {
	if x > d.width {
		return x - d.width
	}
	if x < -d.width {
		return x + d.width
	}
	return 0
}
