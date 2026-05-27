package queue

import "iter"

// store elements ordered by priority using a binary heap
//
// less(a, b) must return true if a has higher priority than b
//
// this type is not safe for concurrent use
type PriorityQueue[T any] struct {
	data []T
	less func(a, b T) bool
}

// create an empty priority queue
//
// time: O(1)
func NewPriorityQueue[T any](less func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		less: less,
	}
}

// add an element to the queue
//
// time: O(log n)
func (q *PriorityQueue[T]) Push(v T) {
	q.data = append(q.data, v)
	q.up(len(q.data) - 1)
}

// remove and return the highest priority element
//
// return false if the queue is empty
//
// time: O(log n)
func (q *PriorityQueue[T]) Pop() (T, bool) {
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}

	top := q.data[0]
	last := q.data[len(q.data)-1]

	var zero T
	q.data[len(q.data)-1] = zero
	q.data = q.data[:len(q.data)-1]

	if len(q.data) > 0 {
		q.data[0] = last
		q.down(0)
	}

	return top, true
}

// return the highest priority element without removing it
//
// return false if the queue is empty
//
// time: O(1)
func (q *PriorityQueue[T]) Peek() (T, bool) {
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}

	return q.data[0], true
}

// report whether no elements are stored
//
// time: O(1)
func (q *PriorityQueue[T]) Empty() bool {
	return len(q.data) == 0
}

// return the number of stored elements
//
// time: O(1)
func (q *PriorityQueue[T]) Len() int {
	return len(q.data)
}

// remove all elements but keep allocated storage
//
// time: O(n)
func (q *PriorityQueue[T]) Clear() {
	for i := range q.data {
		var zero T
		q.data[i] = zero
	}
	q.data = q.data[:0]
}

// iterate over stored elements in heap order
//
// this does not guarantee sorted order
//
// iteration stops if yield returns false
//
// time: O(n)
func (q *PriorityQueue[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range q.data {
			if !yield(v) {
				return
			}
		}
	}
}

// iterate over elements in priority order while removing them
//
// iteration stops if yield returns false
//
// time: O(n log n)
func (q *PriorityQueue[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for len(q.data) > 0 {
			v, _ := q.Pop()
			if !yield(v) {
				return
			}
		}
	}
}

func (q *PriorityQueue[T]) up(i int) {
	for {
		p := (i - 1) / 2
		if i == 0 || !q.less(q.data[i], q.data[p]) {
			return
		}

		q.data[i], q.data[p] = q.data[p], q.data[i]
		i = p
	}
}

func (q *PriorityQueue[T]) down(i int) {
	n := len(q.data)

	for {
		l := 2*i + 1
		if l >= n {
			return
		}

		best := l
		r := l + 1
		if r < n && q.less(q.data[r], q.data[l]) {
			best = r
		}

		if !q.less(q.data[best], q.data[i]) {
			return
		}

		q.data[i], q.data[best] = q.data[best], q.data[i]
		i = best
	}
}
