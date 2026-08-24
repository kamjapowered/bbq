package queue

import "iter"

// store elements in fifo order using a singly linked list
//
// this type is not safe for concurrent use
type LinkedQueue[T any] struct {
	head   *node[T]
	tail   *node[T]
	length int
}

type node[T any] struct {
	data T
	next *node[T]
}

// create an empty queue
//
// time: O(1)
func NewLinkedQueue[T any]() *LinkedQueue[T] {
	return &LinkedQueue[T]{}
}

// add an element to the end of the queue
//
// time: O(1)
func (q *LinkedQueue[T]) Push(v T) {
	n := &node[T]{data: v}

	if q.tail != nil {
		q.tail.next = n
	} else {
		q.head = n
	}

	q.tail = n
	q.length++
}

// remove and return the oldest element
//
// return false if the queue is empty
//
// time: O(1)
func (q *LinkedQueue[T]) Pop() (T, bool) {
	if q.length == 0 {
		var zero T
		return zero, false
	}

	n := q.head
	v := n.data

	q.head = n.next
	if q.head == nil {
		q.tail = nil
	}

	q.length--
	return v, true
}

// return the oldest element without removing it
//
// return false if the queue is empty
//
// time: O(1)
func (q *LinkedQueue[T]) Peek() (T, bool) {
	if q.length == 0 {
		var zero T
		return zero, false
	}

	return q.head.data, true
}

// report whether no elements are stored
//
// time: O(1)
func (q *LinkedQueue[T]) Empty() bool {
	return q.length == 0
}

// return the number of stored elements
//
// time: O(1)
func (q *LinkedQueue[T]) Len() int {
	return q.length
}

// remove all elements
//
// time: O(1)
func (q *LinkedQueue[T]) Clear() {
	q.head = nil
	q.tail = nil
	q.length = 0
}

// return a copy of the stored elements in fifo order
//
// time: O(n)
func (q *LinkedQueue[T]) ToSlice() []T {
	out := make([]T, q.length)

	curr := q.head
	for i := range q.length {
		out[i] = curr.data
		curr = curr.next
	}

	return out
}

// iterate over stored elements in fifo order
//
// iteration stops if yield returns false
//
// time: O(n)
func (q *LinkedQueue[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for curr := q.head; curr != nil; curr = curr.next {
			if !yield(curr.data) {
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
func (q *LinkedQueue[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for q.head != nil {
			n := q.head
			v := n.data

			q.head = n.next
			q.length--
			if q.head == nil {
				q.tail = nil
			}

			n.next = nil

			if !yield(v) {
				return
			}
		}
	}
}
