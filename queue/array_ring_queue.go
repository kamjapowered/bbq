package queue

import "iter"

// store elements in fifo order using a fixed size ring buffer
//
// this type is not safe for concurrent use
type RingQueue[T any] struct {
	data []T
	head int
	tail int
	size int
}

// create an empty ring backed queue with the given capacity
//
// if capacity is less than or equal to zero, the queue is created empty
// and will reject all enqueue operations
//
// time: O(n) due to allocation
func NewRingQueue[T any](capacity int) *RingQueue[T] {
	if capacity <= 0 {
		return &RingQueue[T]{}
	}

	return &RingQueue[T]{
		data: make([]T, capacity),
	}
}

// add an element to the end of the queue
//
// return false if the queue is full
//
// time: O(1)
func (q *RingQueue[T]) Enqueue(v T) bool {
	if q.size == len(q.data) {
		return false
	}

	q.data[q.tail] = v
	q.tail++
	if q.tail == len(q.data) {
		q.tail = 0
	}
	q.size++

	return true
}

// remove and return the oldest element
//
// return false if the queue is empty
//
// time: O(1)
func (q *RingQueue[T]) Dequeue() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}

	v := q.data[q.head]

	var zero T
	q.data[q.head] = zero

	q.head++
	if q.head == len(q.data) {
		q.head = 0
	}
	q.size--

	return v, true
}

// return the oldest element without removing it
//
// return false if the queue is empty
//
// time: O(1)
func (q *RingQueue[T]) Peek() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}

	return q.data[q.head], true
}

// report whether no elements are stored
//
// time: O(1)
func (q *RingQueue[T]) Empty() bool {
	return q.size == 0
}

// report whether no more elements can be added
//
// time: O(1)
func (q *RingQueue[T]) Full() bool {
	return q.size == len(q.data)
}

// return the number of stored elements
//
// time: O(1)
func (q *RingQueue[T]) Len() int {
	return q.size
}

// return the maximum number of elements that can be stored
//
// time: O(1)
func (q *RingQueue[T]) Capacity() int {
	return len(q.data)
}

// remove all elements while preserving capacity
//
// time: O(n)
func (q *RingQueue[T]) Clear() {
	if q.size == 0 {
		q.head = 0
		q.tail = 0
		return
	}
	if len(q.data) == 0 {
		q.head = 0
		q.tail = 0
		q.size = 0
		return
	}

	for i := range q.size {
		idx := q.head + i
		if idx >= len(q.data) {
			idx -= len(q.data)
		}
		var zero T
		q.data[idx] = zero
	}

	q.head = 0
	q.tail = 0
	q.size = 0
}

// return a copy of the stored elements in fifo order
//
// time: O(n)
func (q *RingQueue[T]) ToSlice() []T {
	out := make([]T, q.size)
	if q.size == 0 {
		return out
	}
	if len(q.data) == 0 {
		return out
	}

	for i := range q.size {
		idx := q.head + i
		if idx >= len(q.data) {
			idx -= len(q.data)
		}
		out[i] = q.data[idx]
	}

	return out
}

// iterate over stored elements in fifo order
//
// iteration stops if yield returns false
//
// time: O(n)
func (q *RingQueue[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if q.size == 0 || len(q.data) == 0 {
			return
		}

		for i := range q.size {
			idx := q.head + i
			if idx >= len(q.data) {
				idx -= len(q.data)
			}
			if !yield(q.data[idx]) {
				return
			}
		}
	}
}

// iterate over elements in fifo order while removing them from the queue
//
// iteration stops if yield returns false
//
// time: O(n)
func (q *RingQueue[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for q.size > 0 {
			v := q.data[q.head]

			var zero T
			q.data[q.head] = zero

			q.head++
			if q.head == len(q.data) {
				q.head = 0
			}
			q.size--

			if !yield(v) {
				return
			}
		}
	}
}
