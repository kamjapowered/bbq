package filter

import (
	"slices"

	c "github.com/vistormu/go-dsa/constraints"
)

// compute a sliding window median
//
// this type is not safe for concurrent use
type Median[T c.Number] struct {
	window int
	values []T
}

// create a median filter with a fixed window size
//
// if windowSize is less than or equal to zero, it creates an empty filter that returns zero
func NewMedian[T c.Number](windowSize int) *Median[T] {
	if windowSize <= 0 {
		return &Median[T]{}
	}

	return &Median[T]{
		window: windowSize,
		values: make([]T, 0, windowSize),
	}
}

// reset the stored samples
func (m *Median[T]) Reset() {
	m.values = m.values[:0]
}

// compute the median of the current window after inserting value
//
// time: O(w log w) where w is the window size
func (m *Median[T]) Compute(value T) T {
	if m.window <= 0 {
		return 0
	}

	m.values = append(m.values, value)
	if len(m.values) > m.window {
		copy(m.values, m.values[1:])
		m.values = m.values[:m.window]
	}

	tmp := make([]T, len(m.values))
	copy(tmp, m.values)
	slices.Sort(tmp)

	n := len(tmp)
	if n%2 == 1 {
		return tmp[n/2]
	}

	return (tmp[n/2-1] + tmp[n/2]) / 2
}
