package queue

import "iter"

// store elements in fifo order using a circular singly linked list
//
// this type is not safe for concurrent use
type QueueLinkedListArray[T any] struct {
	head   *node[T]
	tail   *node[T]
	length int
}

// create an empty queue
//
// time: O(1)
func NewQueueLinkedListArray[T any]() *QueueLinkedListArray[T] {
	return &QueueLinkedListArray[T]{}
}

// add an element to the end of the queue
//
// time: O(1)
func (q *QueueLinkedListArray[T]) Enqueue(v T) {
	n := &node[T]{data: v}

	if q.head == nil {
		q.head = n
		q.tail = n
		n.next = n
		q.length = 1
		return
	}

	n.next = q.head
	q.tail.next = n
	q.tail = n
	q.length++
}

// remove and return the oldest element
//
// return false if the queue is empty
//
// time: O(1)
func (q *QueueLinkedListArray[T]) Dequeue() (T, bool) {
	if q.length == 0 {
		var zero T
		return zero, false
	}

	v := q.head.data

	if q.head == q.tail {
		q.head = nil
		q.tail = nil
		q.length = 0
		return v, true
	}

	q.head = q.head.next
	q.tail.next = q.head
	q.length--

	return v, true
}

// return the oldest element without removing it
//
// return false if the queue is empty
//
// time: O(1)
func (q *QueueLinkedListArray[T]) Peek() (T, bool) {
	if q.length == 0 {
		var zero T
		return zero, false
	}

	return q.head.data, true
}

// report whether no elements are stored
//
// time: O(1)
func (q *QueueLinkedListArray[T]) Empty() bool {
	return q.length == 0
}

// return the number of stored elements
//
// time: O(1)
func (q *QueueLinkedListArray[T]) Len() int {
	return q.length
}

// remove all elements
//
// time: O(1)
func (q *QueueLinkedListArray[T]) Clear() {
	q.head = nil
	q.tail = nil
	q.length = 0
}

// return a copy of the stored elements in fifo order
//
// time: O(n)
func (q *QueueLinkedListArray[T]) ToSlice() []T {
	out := make([]T, q.length)
	if q.length == 0 {
		return out
	}

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
func (q *QueueLinkedListArray[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for curr, i := q.head, 0; i < q.length; i++ {
			if !yield(curr.data) {
				return
			}
			curr = curr.next
		}
	}
}

// iterate over elements in fifo order while removing them from the queue
//
// iteration stops if yield returns false
//
// time: O(n)
func (q *QueueLinkedListArray[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for q.length > 0 {
			v := q.head.data

			if q.head == q.tail {
				q.head = nil
				q.tail = nil
				q.length = 0
				if !yield(v) {
					return
				}
				return
			}

			q.head = q.head.next
			q.tail.next = q.head
			q.length--

			if !yield(v) {
				return
			}
		}
	}
}
