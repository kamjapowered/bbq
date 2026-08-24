package queue

import (
	"fmt"
	"iter"
	"slices"
	"testing"
)

// the shape shared by the two fixed capacity implementations
type bounded[T any] interface {
	Push(T) bool
	Pop() (T, bool)
	Peek() (T, bool)
	Empty() bool
	Full() bool
	Len() int
	Capacity() int
	Clear()
	Iter() iter.Seq[T]
	Drain() iter.Seq[T]
	ToSlice() []T
}

var boundedFactories = map[string]func(capacity int) bounded[int]{
	"RingQueue": func(c int) bounded[int] { return NewRingQueue[int](c) },
}

func TestBoundedCapacity(t *testing.T) {
	for name, newQueue := range boundedFactories {
		t.Run(name, func(t *testing.T) {
			q := newQueue(3)

			if got := q.Capacity(); got != 3 {
				t.Fatalf("Capacity = %d, want 3", got)
			}
			if q.Full() {
				t.Fatal("Full = true on an empty queue with room")
			}

			for i := range 3 {
				if !q.Push(i) {
					t.Fatalf("Push(%d) rejected below capacity", i)
				}
				if got := q.Len(); got != i+1 {
					t.Fatalf("Len = %d, want %d", got, i+1)
				}
			}

			if !q.Full() {
				t.Fatal("Full = false at capacity")
			}
			if q.Push(99) {
				t.Fatal("Push succeeded on a full queue")
			}
			if got := q.Len(); got != 3 {
				t.Fatalf("Len = %d after a rejected Push, want 3", got)
			}
			equal(t, "ToSlice", q.ToSlice(), []int{0, 1, 2})

			// making room lets the next enqueue through
			if _, ok := q.Pop(); !ok {
				t.Fatal("Pop failed on a full queue")
			}
			if q.Full() {
				t.Fatal("Full = true after making room")
			}
			if !q.Push(99) {
				t.Fatal("Push rejected after making room")
			}
			equal(t, "ToSlice", q.ToSlice(), []int{1, 2, 99})
		})
	}
}

func TestBoundedZeroCapacity(t *testing.T) {
	for name, newQueue := range boundedFactories {
		for _, capacity := range []int{0, -1, -7} {
			t.Run(fmt.Sprintf("%s/%d", name, capacity), func(t *testing.T) {
				q := newQueue(capacity)

				if got := q.Capacity(); got != 0 {
					t.Fatalf("Capacity = %d for requested %d, want 0", got, capacity)
				}
				if q.Push(1) {
					t.Fatal("Push succeeded on a zero capacity queue")
				}
				if !q.Empty() || q.Len() != 0 {
					t.Fatalf("Len = %d, want 0", q.Len())
				}
				if !q.Full() {
					t.Fatal("Full = false on a queue that can never accept an element")
				}
				if _, ok := q.Pop(); ok {
					t.Fatal("Pop succeeded on a zero capacity queue")
				}
				if _, ok := q.Peek(); ok {
					t.Fatal("Peek succeeded on a zero capacity queue")
				}
				equal(t, "Iter", slices.Collect(q.Iter()), nil)
				equal(t, "Drain", slices.Collect(q.Drain()), nil)
				equal(t, "ToSlice", q.ToSlice(), []int{})

				// none of the above may panic on the nil backing array
				q.Clear()
				if q.Len() != 0 {
					t.Fatalf("Len = %d after Clear, want 0", q.Len())
				}
			})
		}
	}
}

func TestBoundedCapacityOne(t *testing.T) {
	for name, newQueue := range boundedFactories {
		t.Run(name, func(t *testing.T) {
			q := newQueue(1)

			for i := range 5 {
				if !q.Push(i) {
					t.Fatalf("Push(%d) rejected on an empty queue of capacity 1", i)
				}
				if q.Push(-1) {
					t.Fatal("Push succeeded on a full queue of capacity 1")
				}
				if v, ok := q.Peek(); !ok || v != i {
					t.Fatalf("Peek = (%d, %v), want (%d, true)", v, ok, i)
				}
				if v, ok := q.Pop(); !ok || v != i {
					t.Fatalf("Pop = (%d, %v), want (%d, true)", v, ok, i)
				}
			}
		})
	}
}

func TestBoundedWraparound(t *testing.T) {
	for name, newQueue := range boundedFactories {
		t.Run(name, func(t *testing.T) {
			q := newQueue(4)

			// fill, then dequeue half so the head sits mid buffer
			for i := range 4 {
				q.Push(i)
			}
			q.Pop()
			q.Pop()

			// these two wrap past the end of the backing array
			if !q.Push(4) || !q.Push(5) {
				t.Fatal("Push rejected while wrapping with free room")
			}

			want := []int{2, 3, 4, 5}
			equal(t, "ToSlice", q.ToSlice(), want)
			equal(t, "Iter", slices.Collect(q.Iter()), want)
			if !q.Full() {
				t.Fatal("Full = false after wrapping back to capacity")
			}
			equal(t, "Drain", slices.Collect(q.Drain()), want)
			if !q.Empty() {
				t.Fatal("Empty = false after Drain")
			}
		})
	}
}

func TestBoundedManyWraps(t *testing.T) {
	for name, newQueue := range boundedFactories {
		t.Run(name, func(t *testing.T) {
			q := newQueue(3)
			var model []int

			// run the head and tail around the buffer many times over
			for i := range 50 {
				if len(model) < 3 {
					if !q.Push(i) {
						t.Fatalf("Push(%d) rejected below capacity", i)
					}
					model = append(model, i)
				}

				if i%2 == 1 && len(model) > 0 {
					v, ok := q.Pop()
					if !ok || v != model[0] {
						t.Fatalf("Pop = (%d, %v), want (%d, true)", v, ok, model[0])
					}
					model = model[1:]
				}

				equal(t, "ToSlice", q.ToSlice(), model)
				if got := q.Len(); got != len(model) {
					t.Fatalf("Len = %d, want %d", got, len(model))
				}
				if got, want := q.Capacity(), 3; got != want {
					t.Fatalf("Capacity = %d, want %d", got, want)
				}
			}
		})
	}
}

func TestBoundedClearPreservesCapacity(t *testing.T) {
	for name, newQueue := range boundedFactories {
		t.Run(name, func(t *testing.T) {
			q := newQueue(4)
			for i := range 4 {
				q.Push(i)
			}
			// leave the head wrapped before clearing
			q.Pop()
			q.Pop()
			q.Push(9)

			q.Clear()

			if got := q.Capacity(); got != 4 {
				t.Fatalf("Capacity = %d after Clear, want 4", got)
			}
			if !q.Empty() || q.Len() != 0 {
				t.Fatalf("Len = %d after Clear, want 0", q.Len())
			}
			if q.Full() {
				t.Fatal("Full = true after Clear")
			}

			// the whole capacity is available again, from the start
			for i := range 4 {
				if !q.Push(i) {
					t.Fatalf("Push(%d) rejected after Clear", i)
				}
			}
			equal(t, "ToSlice", q.ToSlice(), []int{0, 1, 2, 3})
		})
	}
}

func TestBoundedDrainEarlyBreakKeepsCapacity(t *testing.T) {
	for name, newQueue := range boundedFactories {
		t.Run(name, func(t *testing.T) {
			q := newQueue(4)
			for i := range 4 {
				q.Push(i)
			}

			var got []int
			for v := range q.Drain() {
				got = append(got, v)
				if len(got) == 2 {
					break
				}
			}
			equal(t, "Drain prefix", got, []int{0, 1})

			if q.Full() {
				t.Fatal("Full = true after draining two of four")
			}
			if !q.Push(4) || !q.Push(5) {
				t.Fatal("Push rejected after a partial Drain freed room")
			}
			equal(t, "ToSlice", q.ToSlice(), []int{2, 3, 4, 5})
		})
	}
}
