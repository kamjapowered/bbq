package filter

import c "github.com/vistormu/go-dsa/constraints"

// compute a sliding window mean
//
// this type is not safe for concurrent use
type Mean[T c.Float] struct {
	window int
	values []T
	sum    T
}

// create a mean filter with a fixed window size
//
// if windowSize is less than or equal to zero, it creates an empty filter that returns zero
func NewMean[T c.Float](windowSize int) *Mean[T] {
	if windowSize <= 0 {
		return &Mean[T]{}
	}

	return &Mean[T]{
		window: windowSize,
		values: make([]T, 0, windowSize),
	}
}

// reset the stored samples
func (m *Mean[T]) Reset() {
	m.values = m.values[:0]
	m.sum = 0
}

// compute the mean of the current window after inserting value
//
// time: O(1)
func (m *Mean[T]) Compute(value T) T {
	if m.window <= 0 {
		return 0
	}

	if len(m.values) == m.window {
		oldest := m.values[0]
		copy(m.values, m.values[1:])
		m.values[m.window-1] = value
		m.sum = m.sum - oldest + value
		return m.sum / T(m.window)
	}

	m.values = append(m.values, value)
	m.sum += value
	return m.sum / T(len(m.values))
}
