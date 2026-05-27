package stack

import "iter"

// store unique elements in lifo order
//
// pushing an existing element moves it to the top
//
// this type is not safe for concurrent use
type UniqueStack[T comparable] struct {
	head   *dnode[T]
	tail   *dnode[T]
	length int
	index  map[T]*dnode[T]
}

type dnode[T comparable] struct {
	data T
	next *dnode[T]
	prev *dnode[T]
}

// create an empty unique stack
//
// time: O(1)
func NewUniqueStack[T comparable]() *UniqueStack[T] {
	return &UniqueStack[T]{index: make(map[T]*dnode[T])}
}

// report whether no elements are stored
//
// time: O(1)
func (s *UniqueStack[T]) Empty() bool { return s.length == 0 }

// return the number of stored elements
//
// time: O(1)
func (s *UniqueStack[T]) Len() int { return s.length }

// return the top element without removing it
//
// return false if the stack is empty
//
// time: O(1)
func (s *UniqueStack[T]) Peek() (T, bool) {
	if s.length == 0 {
		var zero T
		return zero, false
	}
	return s.tail.data, true
}

// remove all elements
//
// time: O(n)
func (s *UniqueStack[T]) Clear() {
	for n := s.head; n != nil; {
		next := n.next
		n.next = nil
		n.prev = nil
		n = next
	}

	s.head = nil
	s.tail = nil
	s.length = 0

	for k := range s.index {
		delete(s.index, k)
	}
}

// add v to the top
//
// if v already exists, move it to the top
//
// time: O(1) amortised
func (s *UniqueStack[T]) Push(v T) {
	if n, ok := s.index[v]; ok {
		s.unlink(n)
	}

	n := &dnode[T]{data: v, prev: s.tail}
	if s.tail != nil {
		s.tail.next = n
	} else {
		s.head = n
	}
	s.tail = n

	s.index[v] = n
	s.length++
}

// remove and return the top element
//
// return false if the stack is empty
//
// time: O(1)
func (s *UniqueStack[T]) Pop() (T, bool) {
	if s.length == 0 {
		var zero T
		return zero, false
	}

	n := s.tail
	s.unlink(n)
	delete(s.index, n.data)

	v := n.data
	n.next = nil
	n.prev = nil

	return v, true
}

// report whether v is stored
//
// time: O(1)
func (s *UniqueStack[T]) Contains(v T) bool {
	_, ok := s.index[v]
	return ok
}

// remove v if present
//
// return false if v is not present
//
// time: O(1)
func (s *UniqueStack[T]) Remove(v T) bool {
	n, ok := s.index[v]
	if !ok {
		return false
	}

	s.unlink(n)
	delete(s.index, v)

	n.next = nil
	n.prev = nil

	return true
}

// return a copy of the stored elements from bottom to top
//
// time: O(n)
func (s *UniqueStack[T]) ToSlice() []T {
	out := make([]T, 0, s.length)
	for n := s.head; n != nil; n = n.next {
		out = append(out, n.data)
	}
	return out
}

// iterate over elements from bottom to top
//
// iteration stops if yield returns false
//
// time: O(n)
func (s *UniqueStack[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for n := s.head; n != nil; n = n.next {
			if !yield(n.data) {
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
func (s *UniqueStack[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for s.length > 0 {
			v, _ := s.Pop()
			if !yield(v) {
				return
			}
		}
	}
}

func (s *UniqueStack[T]) unlink(n *dnode[T]) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		s.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		s.tail = n.prev
	}

	n.next = nil
	n.prev = nil

	s.length--
}
