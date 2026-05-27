package stack

import "iter"

// store elements in lifo order using a singly linked list
//
// this type is not safe for concurrent use
type LinkedStack[T any] struct {
	head   *node[T]
	length int
}

type node[T any] struct {
	data T
	next *node[T]
}

// create an empty linked stack
//
// time: O(1)
func NewLinkedStack[T any]() *LinkedStack[T] {
	return &LinkedStack[T]{}
}

// add an element to the top of the stack
//
// time: O(1)
func (s *LinkedStack[T]) Push(v T) {
	n := &node[T]{data: v, next: s.head}
	s.head = n
	s.length++
}

// remove and return the top element
//
// return false if the stack is empty
//
// time: O(1)
func (s *LinkedStack[T]) Pop() (T, bool) {
	if s.length == 0 {
		var zero T
		return zero, false
	}

	n := s.head
	v := n.data

	s.head = n.next
	n.next = nil
	s.length--

	return v, true
}

// return the top element without removing it
//
// return false if the stack is empty
//
// time: O(1)
func (s *LinkedStack[T]) Peek() (T, bool) {
	if s.length == 0 {
		var zero T
		return zero, false
	}

	return s.head.data, true
}

// report whether no elements are stored
//
// time: O(1)
func (s *LinkedStack[T]) Empty() bool {
	return s.length == 0
}

// return the number of stored elements
//
// time: O(1)
func (s *LinkedStack[T]) Len() int {
	return s.length
}

// remove all elements
//
// time: O(1)
func (s *LinkedStack[T]) Clear() {
	s.head = nil
	s.length = 0
}

// return a copy of the stored elements from top to bottom
//
// time: O(n)
func (s *LinkedStack[T]) ToSlice() []T {
	out := make([]T, 0, s.length)
	for curr := s.head; curr != nil; curr = curr.next {
		out = append(out, curr.data)
	}
	return out
}

// iterate over elements from top to bottom
//
// iteration stops if yield returns false
//
// time: O(n)
func (s *LinkedStack[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for curr := s.head; curr != nil; curr = curr.next {
			if !yield(curr.data) {
				return
			}
		}
	}
}

// iterate over elements from top to bottom while removing them
//
// iteration stops if yield returns false
//
// time: O(n)
func (s *LinkedStack[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for s.length > 0 {
			v, _ := s.Pop()
			if !yield(v) {
				return
			}
		}
	}
}
