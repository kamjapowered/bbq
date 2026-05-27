package queue

import "iter"

// store up to a fixed number of elements in fifo order
//
// this type is not safe for concurrent use
type BoundedQueue[T any] struct {
	data     []T
	head     int
	tail     int
	size     int
	capacity int
}

// create a new fifo queue with a fixed capacity
//
// if capacity is less than or equal to zero, create an empty queue that rejects enqueues
//
// time: O(n) due to allocation
func NewBoundedQueue[T any](capacity int) *BoundedQueue[T] {
	if capacity <= 0 {
		return &BoundedQueue[T]{}
	}

	return &BoundedQueue[T]{
		data:     make([]T, capacity),
		capacity: capacity,
	}
}

// report whether no more elements can be added
//
// time: O(1)
func (q *BoundedQueue[T]) Full() bool {
	return q.size == q.capacity
}

// report whether no elements are stored
//
// time: O(1)
func (q *BoundedQueue[T]) Empty() bool {
	return q.size == 0
}

// return the number of stored elements
//
// time: O(1)
func (q *BoundedQueue[T]) Len() int {
	return q.size
}

// return the maximum number of elements that can be stored
//
// time: O(1)
func (q *BoundedQueue[T]) Capacity() int {
	return q.capacity
}

// add an element to the end of the queue
//
// return false if the queue is full
//
// time: O(1)
func (q *BoundedQueue[T]) Enqueue(v T) bool {
	if q.size == q.capacity {
		return false
	}

	if q.capacity == 0 {
		return false
	}

	q.data[q.tail] = v
	q.tail++
	if q.tail == q.capacity {
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
func (q *BoundedQueue[T]) Dequeue() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}

	v := q.data[q.head]

	var zero T
	q.data[q.head] = zero

	q.head++
	if q.head == q.capacity {
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
func (q *BoundedQueue[T]) Peek() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}

	return q.data[q.head], true
}

// remove all elements while preserving capacity
//
// time: O(n)
func (q *BoundedQueue[T]) Clear() {
	if q.size == 0 || q.capacity == 0 {
		q.head = 0
		q.tail = 0
		q.size = 0
		return
	}

	for i := range q.size {
		idx := q.head + i
		if idx >= q.capacity {
			idx -= q.capacity
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
func (q *BoundedQueue[T]) ToSlice() []T {
	out := make([]T, q.size)
	if q.size == 0 || q.capacity == 0 {
		return out
	}

	for i := range q.size {
		idx := q.head + i
		if idx >= q.capacity {
			idx -= q.capacity
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
func (q *BoundedQueue[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if q.size == 0 || q.capacity == 0 {
			return
		}

		for i := range q.size {
			idx := q.head + i
			if idx >= q.capacity {
				idx -= q.capacity
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
func (q *BoundedQueue[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for q.size > 0 {
			v := q.data[q.head]

			var zero T
			q.data[q.head] = zero

			q.head++
			if q.head == q.capacity {
				q.head = 0
			}
			q.size--

			if !yield(v) {
				return
			}
		}
	}
}
