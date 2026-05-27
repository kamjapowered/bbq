package set

import "iter"

// store unique elements using a hash table
//
// this type is not safe for concurrent use
type HashSet[T comparable] struct {
	data map[T]struct{}
}

// create an empty set
//
// time: O(1)
func NewHashSet[T comparable]() *HashSet[T] {
	return &HashSet[T]{data: make(map[T]struct{})}
}

// add an element to the set
//
// time: O(1) amortised
func (h *HashSet[T]) Add(v T) {
	h.data[v] = struct{}{}
}

// remove an element from the set
//
// time: O(1)
func (h *HashSet[T]) Remove(v T) {
	delete(h.data, v)
}

// report whether an element is stored
//
// time: O(1)
func (h *HashSet[T]) Contains(v T) bool {
	_, ok := h.data[v]
	return ok
}

// remove all elements while preserving capacity
//
// time: O(n)
func (h *HashSet[T]) Clear() {
	for k := range h.data {
		delete(h.data, k)
	}
}

// report whether no elements are stored
//
// time: O(1)
func (h *HashSet[T]) Empty() bool {
	return len(h.data) == 0
}

// return the number of stored elements
//
// time: O(1)
func (h *HashSet[T]) Len() int {
	return len(h.data)
}

// return a copy of the stored elements
//
// time: O(n)
func (h *HashSet[T]) ToSlice() []T {
	out := make([]T, 0, len(h.data))
	for k := range h.data {
		out = append(out, k)
	}
	return out
}

// iterate over stored elements
//
// iteration stops if yield returns false
//
// order is unspecified
//
// time: O(n)
func (h *HashSet[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range h.data {
			if !yield(k) {
				return
			}
		}
	}
}

// iterate over elements while removing them from the set
//
// iteration stops if yield returns false
//
// time: O(n)
func (h *HashSet[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range h.data {
			delete(h.data, k)
			if !yield(k) {
				return
			}
		}
	}
}

// return a new set containing elements present in either set
//
// time: O(n + m)
func (h *HashSet[T]) Union(other *HashSet[T]) *HashSet[T] {
	out := NewHashSet[T]()
	for k := range h.data {
		out.data[k] = struct{}{}
	}
	for k := range other.data {
		out.data[k] = struct{}{}
	}
	return out
}

// return a new set containing elements present in both sets
//
// time: O(n)
func (h *HashSet[T]) Intersection(other *HashSet[T]) *HashSet[T] {
	out := NewHashSet[T]()
	for k := range h.data {
		if _, ok := other.data[k]; ok {
			out.data[k] = struct{}{}
		}
	}
	return out
}

// return a new set containing elements present in this set but not the other
//
// time: O(n)
func (h *HashSet[T]) Difference(other *HashSet[T]) *HashSet[T] {
	out := NewHashSet[T]()
	for k := range h.data {
		if _, ok := other.data[k]; !ok {
			out.data[k] = struct{}{}
		}
	}
	return out
}
