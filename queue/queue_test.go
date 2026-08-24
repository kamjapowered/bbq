package queue

import (
	"iter"
	"slices"
	"testing"
)

// the shape every fifo implementation in this package shares
type fifo[T any] interface {
	Pop() (T, bool)
	Peek() (T, bool)
	Empty() bool
	Len() int
	Clear()
	Iter() iter.Seq[T]
	Drain() iter.Seq[T]
	ToSlice() []T
}

// a fifo paired with an enqueue adapter, since the bounded implementations
// report rejection and the growable ones cannot fail
type impl struct {
	fifo[int]
	enqueue func(int) bool
}

func (q impl) push(t *testing.T, values ...int) {
	t.Helper()
	for _, v := range values {
		if !q.enqueue(v) {
			t.Fatalf("Push(%d) rejected", v)
		}
	}
}

// capacity of the bounded implementations under the shared tests, large
// enough that no shared test ever hits the bound
const roomy = 32

var factories = map[string]func() impl{
	"Queue": func() impl {
		q := NewQueue[int]()
		return impl{q, func(v int) bool { q.Push(v); return true }}
	},
	"RingQueue": func() impl {
		q := NewRingQueue[int](roomy)
		return impl{q, q.Push}
	},
	"LinkedQueue": func() impl {
		q := NewLinkedQueue[int]()
		return impl{q, func(v int) bool { q.Push(v); return true }}
	},
	"LinkedRingQueue": func() impl {
		q := NewLinkedRingQueue[int]()
		return impl{q, func(v int) bool { q.Push(v); return true }}
	},
}

func equal(t *testing.T, what string, got, want []int) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("%s = %v, want %v", what, got, want)
	}
}

func TestEmpty(t *testing.T) {
	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			q := newQueue()

			if !q.Empty() {
				t.Fatal("Empty = false on a fresh queue")
			}
			if got := q.Len(); got != 0 {
				t.Fatalf("Len = %d, want 0", got)
			}
			if v, ok := q.Peek(); ok || v != 0 {
				t.Fatalf("Peek = (%d, %v), want (0, false)", v, ok)
			}
			if v, ok := q.Pop(); ok || v != 0 {
				t.Fatalf("Pop = (%d, %v), want (0, false)", v, ok)
			}
			if got := len(q.ToSlice()); got != 0 {
				t.Fatalf("len(ToSlice) = %d, want 0", got)
			}
			equal(t, "Iter", slices.Collect(q.Iter()), nil)
			equal(t, "Drain", slices.Collect(q.Drain()), nil)

			// clear on an empty queue is a no-op
			q.Clear()
			if !q.Empty() || q.Len() != 0 {
				t.Fatal("Clear broke an already empty queue")
			}
		})
	}
}

func TestSingleElement(t *testing.T) {
	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			q := newQueue()
			q.push(t, 7)

			if q.Empty() {
				t.Fatal("Empty = true after one enqueue")
			}
			if got := q.Len(); got != 1 {
				t.Fatalf("Len = %d, want 1", got)
			}
			if v, ok := q.Peek(); !ok || v != 7 {
				t.Fatalf("Peek = (%d, %v), want (7, true)", v, ok)
			}
			equal(t, "ToSlice", q.ToSlice(), []int{7})

			if v, ok := q.Pop(); !ok || v != 7 {
				t.Fatalf("Pop = (%d, %v), want (7, true)", v, ok)
			}
			if !q.Empty() {
				t.Fatal("Empty = false after draining the only element")
			}
			if _, ok := q.Pop(); ok {
				t.Fatal("Pop succeeded on an emptied queue")
			}
		})
	}
}

func TestFifoOrder(t *testing.T) {
	want := []int{1, 2, 3, 4, 5}

	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			q := newQueue()
			q.push(t, want...)

			if got := q.Len(); got != len(want) {
				t.Fatalf("Len = %d, want %d", got, len(want))
			}
			if v, _ := q.Peek(); v != want[0] {
				t.Fatalf("Peek = %d, want %d", v, want[0])
			}
			equal(t, "ToSlice", q.ToSlice(), want)
			equal(t, "Iter", slices.Collect(q.Iter()), want)

			var got []int
			for !q.Empty() {
				v, ok := q.Pop()
				if !ok {
					t.Fatal("Pop failed on a non empty queue")
				}
				got = append(got, v)
			}
			equal(t, "Pop order", got, want)
		})
	}
}

func TestPeekDoesNotRemove(t *testing.T) {
	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			q := newQueue()
			q.push(t, 1, 2, 3)

			for range 3 {
				if v, ok := q.Peek(); !ok || v != 1 {
					t.Fatalf("Peek = (%d, %v), want (1, true)", v, ok)
				}
			}
			if got := q.Len(); got != 3 {
				t.Fatalf("Len = %d, want 3", got)
			}
		})
	}
}

func TestIterIsNonMutating(t *testing.T) {
	want := []int{1, 2, 3, 4}

	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			q := newQueue()
			q.push(t, want...)

			// two full passes must agree and leave the queue untouched
			equal(t, "Iter pass 1", slices.Collect(q.Iter()), want)
			equal(t, "Iter pass 2", slices.Collect(q.Iter()), want)

			if got := q.Len(); got != len(want) {
				t.Fatalf("Len = %d after Iter, want %d", got, len(want))
			}
			equal(t, "ToSlice", q.ToSlice(), want)
		})
	}
}

func TestIterEarlyBreak(t *testing.T) {
	want := []int{1, 2, 3, 4}

	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			q := newQueue()
			q.push(t, want...)

			var got []int
			for v := range q.Iter() {
				got = append(got, v)
				if len(got) == 2 {
					break
				}
			}
			equal(t, "Iter prefix", got, []int{1, 2})

			if n := q.Len(); n != len(want) {
				t.Fatalf("Len = %d after a broken Iter, want %d", n, len(want))
			}
			equal(t, "ToSlice", q.ToSlice(), want)
		})
	}
}

func TestToSliceIsACopy(t *testing.T) {
	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			q := newQueue()
			q.push(t, 1, 2, 3)

			s := q.ToSlice()
			s[0] = 99

			if v, _ := q.Peek(); v != 1 {
				t.Fatalf("Peek = %d after mutating the ToSlice result, want 1", v)
			}
			equal(t, "ToSlice", q.ToSlice(), []int{1, 2, 3})
		})
	}
}

func TestDrainEmpties(t *testing.T) {
	want := []int{1, 2, 3, 4, 5}

	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			q := newQueue()
			q.push(t, want...)

			equal(t, "Drain", slices.Collect(q.Drain()), want)

			if !q.Empty() || q.Len() != 0 {
				t.Fatalf("Len = %d after Drain, want 0", q.Len())
			}
			if _, ok := q.Peek(); ok {
				t.Fatal("Peek succeeded after Drain")
			}
			equal(t, "Iter", slices.Collect(q.Iter()), nil)

			// a drained queue is still usable
			q.push(t, 9)
			if v, ok := q.Pop(); !ok || v != 9 {
				t.Fatalf("Pop = (%d, %v) after refilling a drained queue, want (9, true)", v, ok)
			}
		})
	}
}

func TestDrainEarlyBreak(t *testing.T) {
	want := []int{1, 2, 3, 4, 5}

	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			q := newQueue()
			q.push(t, want...)

			var got []int
			for v := range q.Drain() {
				got = append(got, v)
				if len(got) == 2 {
					break
				}
			}
			equal(t, "Drain prefix", got, []int{1, 2})

			// the yielded elements are gone, the rest stays in order
			if n := q.Len(); n != 3 {
				t.Fatalf("Len = %d after a broken Drain, want 3", n)
			}
			equal(t, "ToSlice", q.ToSlice(), []int{3, 4, 5})
			equal(t, "Drain rest", slices.Collect(q.Drain()), []int{3, 4, 5})
		})
	}
}

func TestClearResetsAndReuses(t *testing.T) {
	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			q := newQueue()
			q.push(t, 1, 2, 3)
			q.Clear()

			if !q.Empty() || q.Len() != 0 {
				t.Fatalf("Len = %d after Clear, want 0", q.Len())
			}
			if _, ok := q.Peek(); ok {
				t.Fatal("Peek succeeded after Clear")
			}
			if _, ok := q.Pop(); ok {
				t.Fatal("Pop succeeded after Clear")
			}
			equal(t, "Iter", slices.Collect(q.Iter()), nil)
			equal(t, "ToSlice", q.ToSlice(), []int{})

			// the storage is reusable and starts from a clean state
			q.push(t, 4, 5)
			equal(t, "ToSlice", q.ToSlice(), []int{4, 5})
		})
	}
}

func TestClearMidWraparound(t *testing.T) {
	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			q := newQueue()
			q.push(t, 1, 2, 3, 4)
			// move the head off zero so clear has to normalise it
			q.Pop()
			q.Pop()
			q.Clear()

			q.push(t, 7, 8)
			equal(t, "ToSlice", q.ToSlice(), []int{7, 8})
			equal(t, "Drain", slices.Collect(q.Drain()), []int{7, 8})
		})
	}
}

func TestInterleavedEnqueueDequeue(t *testing.T) {
	// a fixed script of pushes and pops, checked against a slice model
	script := []int{3, 1, 2, 4, 1, 5, 3, 2, 6, 1, 1, 1, 4, 2}

	for name, newQueue := range factories {
		t.Run(name, func(t *testing.T) {
			q := newQueue()

			var model []int
			next := 0

			for i, n := range script {
				if i%2 == 0 {
					// enqueue n elements
					for range n {
						q.push(t, next)
						model = append(model, next)
						next++
					}
				} else {
					// dequeue up to n elements
					for range n {
						v, ok := q.Pop()
						if len(model) == 0 {
							if ok {
								t.Fatalf("Pop = (%d, true) on an empty queue", v)
							}
							continue
						}
						if !ok || v != model[0] {
							t.Fatalf("Pop = (%d, %v), want (%d, true)", v, ok, model[0])
						}
						model = model[1:]
					}
				}

				if got := q.Len(); got != len(model) {
					t.Fatalf("step %d: Len = %d, want %d", i, got, len(model))
				}
				equal(t, "ToSlice", q.ToSlice(), model)
			}
		})
	}
}

func TestDequeueDropsReference(t *testing.T) {
	// dequeue and drain must zero the slot they vacate, so an element does
	// not stay reachable through the queue's storage. white box, since no
	// exported call can observe it
	t.Run("Queue", func(t *testing.T) {
		q := NewQueue[*int]()
		v := 1
		q.Push(&v)
		q.Pop()
		if q.data[0] != nil {
			t.Fatal("Pop left the vacated slot populated")
		}
	})

	t.Run("RingQueue", func(t *testing.T) {
		q := NewRingQueue[*int](2)
		v := 1
		q.Push(&v)
		for range q.Drain() {
		}
		if q.data[0] != nil {
			t.Fatal("Drain left the vacated slot populated")
		}
	})

	t.Run("LinkedQueue", func(t *testing.T) {
		q := NewLinkedQueue[int]()
		q.Push(1)
		q.Push(2)
		q.Pop()
		q.Pop()
		if q.head != nil || q.tail != nil {
			t.Fatal("emptied queue kept a node reference")
		}
	})

	t.Run("LinkedRingQueue", func(t *testing.T) {
		q := NewLinkedRingQueue[int]()
		q.Push(1)
		q.Push(2)
		q.Pop()
		q.Pop()
		if q.head != nil || q.tail != nil {
			t.Fatal("emptied queue kept a node reference")
		}
	})
}
