package stack

import "iter"

// store elements in lifo order using a slice as backing storage
//
// this type is not safe for concurrent use
type Stack[T any] struct {
	data []T
}

// create an empty stack
//
// time: O(1)
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// add an element to the top of the stack
//
// time: O(1) amortised
func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

// remove and return the top element
//
// return false if the stack is empty
//
// time: O(1)
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}

	i := len(s.data) - 1
	v := s.data[i]

	var zero T
	s.data[i] = zero
	s.data = s.data[:i]

	return v, true
}

// return the top element without removing it
//
// return false if the stack is empty
//
// time: O(1)
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}

	return s.data[len(s.data)-1], true
}

// report whether no elements are stored
//
// time: O(1)
func (s *Stack[T]) Empty() bool {
	return len(s.data) == 0
}

// return the number of stored elements
//
// time: O(1)
func (s *Stack[T]) Len() int {
	return len(s.data)
}

// remove all elements but keep allocated storage
//
// time: O(n)
func (s *Stack[T]) Clear() {
	for i := range s.data {
		var zero T
		s.data[i] = zero
	}
	s.data = s.data[:0]
}

// return a copy of the stored elements from bottom to top
//
// time: O(n)
func (s *Stack[T]) ToSlice() []T {
	out := make([]T, len(s.data))
	copy(out, s.data)
	return out
}

// iterate over elements from bottom to top
//
// iteration stops if yield returns false
//
// time: O(n)
func (s *Stack[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s.data {
			if !yield(v) {
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
func (s *Stack[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for len(s.data) > 0 {
			v, _ := s.Pop()
			if !yield(v) {
				return
			}
		}
	}
}
