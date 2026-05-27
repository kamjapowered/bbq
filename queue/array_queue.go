package queue

import (
	"iter"
)

// store elements in fifo order using a ring buffer.
//
// this type is not safe for concurrent use.
type Queue[T any] struct {
	data []T
	head int
	size int
}

// create an empty queue.
//
// time: O(1)
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

// add an element to the end of the queue.
//
// time: O(1) amortised
func (q *Queue[T]) Enqueue(v T) {
	if q.size == len(q.data) {
		q.grow()
	}

	idx := (q.head + q.size) % len(q.data)
	q.data[idx] = v
	q.size++
}

// remove and return the oldest element.
//
// return false if the queue is empty.
//
// time: O(1)
func (q *Queue[T]) Dequeue() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}

	v := q.data[q.head]

	var zero T
	q.data[q.head] = zero

	q.head = (q.head + 1) % len(q.data)
	q.size--

	return v, true
}

// return the oldest element without removing it.
//
// return false if the queue is empty.
//
// time: O(1)
func (q *Queue[T]) Peek() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}

	return q.data[q.head], true
}

// report whether no elements are stored.
//
// time: O(1)
func (q *Queue[T]) Empty() bool {
	return q.size == 0
}

// return the number of stored elements.
//
// time: O(1)
func (q *Queue[T]) Len() int {
	return q.size
}

// remove all elements but keep allocated storage.
//
// time: O(n)
func (q *Queue[T]) Clear() {
	for i := 0; i < q.size; i++ {
		idx := (q.head + i) % len(q.data)
		var zero T
		q.data[idx] = zero
	}

	q.head = 0
	q.size = 0
}

// return a copy of the stored elements in fifo order.
//
// time: O(n)
func (q *Queue[T]) ToSlice() []T {
	out := make([]T, 0, q.size)
	for v := range q.Iter() {
		out = append(out, v)
	}
	return out
}

func (q *Queue[T]) grow() {
	newCap := 1
	if len(q.data) > 0 {
		newCap = len(q.data) * 2
	}

	buf := make([]T, newCap)
	for i := 0; i < q.size; i++ {
		buf[i] = q.data[(q.head+i)%len(q.data)]
	}

	q.data = buf
	q.head = 0
}

// iterate over stored elements in fifo order.
//
// iteration stops if yield returns false.
//
// time: O(n)
func (q *Queue[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < q.size; i++ {
			v := q.data[(q.head+i)%len(q.data)]
			if !yield(v) {
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
func (q *Queue[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for q.size > 0 {
			v := q.data[q.head]

			var zero T
			q.data[q.head] = zero

			q.head = (q.head + 1) % len(q.data)
			q.size--

			if !yield(v) {
				return
			}
		}
	}
}
