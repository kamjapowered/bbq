package queue

import "iter"

// store elements in a double ended fifo structure using a ring buffer
//
// this type is not safe for concurrent use
type Deque[T any] struct {
	data []T
	head int
	size int
}

// create an empty deque
//
// time: O(1)
func NewDeque[T any]() *Deque[T] {
	return &Deque[T]{}
}

// add an element to the front
//
// time: O(1) amortised
func (d *Deque[T]) PushFront(v T) {
	if d.size == len(d.data) {
		d.grow()
	}

	if d.size == 0 {
		d.data[0] = v
		d.head = 0
		d.size = 1
		return
	}

	d.head--
	if d.head < 0 {
		d.head = len(d.data) - 1
	}

	d.data[d.head] = v
	d.size++
}

// add an element to the back
//
// time: O(1) amortised
func (d *Deque[T]) PushBack(v T) {
	if d.size == len(d.data) {
		d.grow()
	}

	idx := (d.head + d.size) % len(d.data)
	d.data[idx] = v
	d.size++
}

// remove and return the front element
//
// return false if empty
//
// time: O(1)
func (d *Deque[T]) PopFront() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}

	v := d.data[d.head]

	var zero T
	d.data[d.head] = zero

	d.head = (d.head + 1) % len(d.data)
	d.size--

	return v, true
}

// remove and return the back element
//
// return false if empty
//
// time: O(1)
func (d *Deque[T]) PopBack() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}

	idx := (d.head + d.size - 1) % len(d.data)
	v := d.data[idx]

	var zero T
	d.data[idx] = zero

	d.size--

	return v, true
}

// return the front element without removing it
//
// return false if empty
//
// time: O(1)
func (d *Deque[T]) PeekFront() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}

	return d.data[d.head], true
}

// return the back element without removing it
//
// return false if empty
//
// time: O(1)
func (d *Deque[T]) PeekBack() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}

	idx := (d.head + d.size - 1) % len(d.data)
	return d.data[idx], true
}

// report whether no elements are stored
//
// time: O(1)
func (d *Deque[T]) Empty() bool {
	return d.size == 0
}

// return the number of stored elements
//
// time: O(1)
func (d *Deque[T]) Len() int {
	return d.size
}

// remove all elements but keep allocated storage
//
// time: O(n)
func (d *Deque[T]) Clear() {
	if d.size == 0 || len(d.data) == 0 {
		d.head = 0
		d.size = 0
		return
	}

	for i := range d.size {
		idx := (d.head + i) % len(d.data)
		var zero T
		d.data[idx] = zero
	}

	d.head = 0
	d.size = 0
}

// iterate over elements from front to back
//
// iteration stops if yield returns false
//
// time: O(n)
func (d *Deque[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range d.size {
			idx := (d.head + i) % len(d.data)
			if !yield(d.data[idx]) {
				return
			}
		}
	}
}

// iterate over elements from front to back while removing them
//
// iteration stops if yield returns false
//
// time: O(n)
func (d *Deque[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for d.size > 0 {
			v := d.data[d.head]

			var zero T
			d.data[d.head] = zero

			d.head = (d.head + 1) % len(d.data)
			d.size--

			if !yield(v) {
				return
			}
		}
	}
}

func (d *Deque[T]) grow() {
	newCap := 1
	if len(d.data) > 0 {
		newCap = len(d.data) * 2
	}

	buf := make([]T, newCap)
	for i := range d.size {
		buf[i] = d.data[(d.head+i)%len(d.data)]
	}

	d.data = buf
	d.head = 0
}
