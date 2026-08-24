package queue

import (
	"slices"
	"testing"
)

func minInt(a, b int) bool { return a < b }
func maxInt(a, b int) bool { return a > b }

func TestPriorityQueueEmpty(t *testing.T) {
	q := NewPriorityQueue(minInt)

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
	equal(t, "Iter", slices.Collect(q.Iter()), nil)
	equal(t, "Drain", slices.Collect(q.Drain()), nil)

	q.Clear()
	if !q.Empty() {
		t.Fatal("Clear broke an already empty queue")
	}
}

func TestPriorityQueueSingleElement(t *testing.T) {
	q := NewPriorityQueue(minInt)
	q.Push(7)

	if q.Empty() || q.Len() != 1 {
		t.Fatalf("Len = %d, want 1", q.Len())
	}
	if v, ok := q.Peek(); !ok || v != 7 {
		t.Fatalf("Peek = (%d, %v), want (7, true)", v, ok)
	}
	if v, ok := q.Pop(); !ok || v != 7 {
		t.Fatalf("Pop = (%d, %v), want (7, true)", v, ok)
	}
	if !q.Empty() {
		t.Fatal("Empty = false after popping the only element")
	}
	if _, ok := q.Pop(); ok {
		t.Fatal("Pop succeeded on an emptied queue")
	}
}

func TestPriorityQueueOrdering(t *testing.T) {
	input := []int{5, 3, 9, 1, 7, 2, 8, 6, 4, 0}

	tests := []struct {
		name string
		less func(a, b int) bool
		want []int
	}{
		{"min", minInt, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"max", maxInt, []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewPriorityQueue(tt.less)
			for _, v := range input {
				q.Push(v)
			}
			if got := q.Len(); got != len(input) {
				t.Fatalf("Len = %d, want %d", got, len(input))
			}
			if v, _ := q.Peek(); v != tt.want[0] {
				t.Fatalf("Peek = %d, want %d", v, tt.want[0])
			}

			var got []int
			for !q.Empty() {
				v, ok := q.Pop()
				if !ok {
					t.Fatal("Pop failed on a non empty queue")
				}
				got = append(got, v)
			}
			equal(t, "Pop order", got, tt.want)
		})
	}
}

func TestPriorityQueuePeekTracksTheBest(t *testing.T) {
	q := NewPriorityQueue(minInt)

	best := 0
	for i, v := range []int{5, 8, 3, 9, 1, 4} {
		q.Push(v)
		if i == 0 || v < best {
			best = v
		}
		if got, ok := q.Peek(); !ok || got != best {
			t.Fatalf("Peek = (%d, %v) after pushing %d, want (%d, true)", got, ok, v, best)
		}
	}
}

func TestPriorityQueueDuplicates(t *testing.T) {
	q := NewPriorityQueue(minInt)
	input := []int{3, 1, 3, 1, 2, 2, 3, 1}
	for _, v := range input {
		q.Push(v)
	}

	want := slices.Sorted(slices.Values(input))
	equal(t, "Drain", slices.Collect(q.Drain()), want)
	if !q.Empty() {
		t.Fatal("Empty = false after Drain")
	}
}

func TestPriorityQueueDrainIsSorted(t *testing.T) {
	q := NewPriorityQueue(minInt)
	input := []int{42, 7, 13, 99, 1, 56, 23, 8}
	for _, v := range input {
		q.Push(v)
	}

	got := slices.Collect(q.Drain())
	equal(t, "Drain", got, slices.Sorted(slices.Values(input)))
	if !q.Empty() || q.Len() != 0 {
		t.Fatalf("Len = %d after Drain, want 0", q.Len())
	}

	// a drained queue is still usable
	q.Push(1)
	if v, ok := q.Peek(); !ok || v != 1 {
		t.Fatalf("Peek = (%d, %v) after refilling, want (1, true)", v, ok)
	}
}

func TestPriorityQueueDrainEarlyBreak(t *testing.T) {
	q := NewPriorityQueue(minInt)
	input := []int{5, 3, 9, 1, 7}
	for _, v := range input {
		q.Push(v)
	}

	var got []int
	for v := range q.Drain() {
		got = append(got, v)
		if len(got) == 2 {
			break
		}
	}
	equal(t, "Drain prefix", got, []int{1, 3})

	if n := q.Len(); n != 3 {
		t.Fatalf("Len = %d after a broken Drain, want 3", n)
	}
	// the heap invariant survives the break
	equal(t, "Drain rest", slices.Collect(q.Drain()), []int{5, 7, 9})
}

func TestPriorityQueueIterIsHeapOrder(t *testing.T) {
	q := NewPriorityQueue(minInt)
	input := []int{5, 3, 9, 1, 7, 2}
	for _, v := range input {
		q.Push(v)
	}

	// iter makes no ordering promise beyond the root, but it must yield
	// every element exactly once and leave the queue alone
	got := slices.Collect(q.Iter())
	if len(got) != len(input) {
		t.Fatalf("len(Iter) = %d, want %d", len(got), len(input))
	}
	if got[0] != 1 {
		t.Fatalf("Iter[0] = %d, want the root 1", got[0])
	}
	equal(t, "Iter as a multiset", slices.Sorted(slices.Values(got)), slices.Sorted(slices.Values(input)))

	if q.Len() != len(input) {
		t.Fatalf("Len = %d after Iter, want %d", q.Len(), len(input))
	}
	equal(t, "Iter pass 2", slices.Collect(q.Iter()), got)
}

func TestPriorityQueueIterEarlyBreak(t *testing.T) {
	q := NewPriorityQueue(minInt)
	for _, v := range []int{5, 3, 9, 1} {
		q.Push(v)
	}

	n := 0
	for range q.Iter() {
		n++
		break
	}
	if n != 1 {
		t.Fatalf("Iter yielded %d elements after a break, want 1", n)
	}
	if q.Len() != 4 {
		t.Fatalf("Len = %d after a broken Iter, want 4", q.Len())
	}
}

func TestPriorityQueueClear(t *testing.T) {
	q := NewPriorityQueue(minInt)
	for _, v := range []int{5, 3, 9, 1} {
		q.Push(v)
	}
	before := cap(q.data)

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
	if got := cap(q.data); got != before {
		t.Fatalf("cap after Clear = %d, want the allocated %d kept", got, before)
	}

	// still usable, and the old elements are gone
	q.Push(4)
	q.Push(2)
	equal(t, "Drain", slices.Collect(q.Drain()), []int{2, 4})
}

func TestPriorityQueueInterleaved(t *testing.T) {
	// push and pop in a fixed script, checked against a sorted slice model
	q := NewPriorityQueue(minInt)
	var model []int

	script := []int{17, 4, 23, 8, 15, 42, 1, 9, 30, 6, 11, 2, 38, 20}
	for i, v := range script {
		q.Push(v)
		model = append(model, v)
		slices.Sort(model)

		if i%3 == 2 {
			got, ok := q.Pop()
			if !ok || got != model[0] {
				t.Fatalf("step %d: Pop = (%d, %v), want (%d, true)", i, got, ok, model[0])
			}
			model = model[1:]
		}

		if n := q.Len(); n != len(model) {
			t.Fatalf("step %d: Len = %d, want %d", i, n, len(model))
		}
		if len(model) > 0 {
			if got, _ := q.Peek(); got != model[0] {
				t.Fatalf("step %d: Peek = %d, want %d", i, got, model[0])
			}
		}
	}

	equal(t, "Drain", slices.Collect(q.Drain()), model)
}

func TestPriorityQueueStableUnderStructs(t *testing.T) {
	type task struct {
		name     string
		priority int
	}

	q := NewPriorityQueue(func(a, b task) bool { return a.priority < b.priority })
	for _, tk := range []task{
		{"deploy", 3},
		{"compile", 1},
		{"test", 2},
		{"lint", 0},
	} {
		q.Push(tk)
	}

	var got []string
	for tk := range q.Drain() {
		got = append(got, tk.name)
	}

	want := []string{"lint", "compile", "test", "deploy"}
	if !slices.Equal(got, want) {
		t.Fatalf("Drain = %v, want %v", got, want)
	}
}

func TestPriorityQueueSortedAndReverseInput(t *testing.T) {
	// ascending and descending inputs are the worst cases for sift up
	tests := map[string][]int{
		"ascending":  {0, 1, 2, 3, 4, 5, 6, 7},
		"descending": {7, 6, 5, 4, 3, 2, 1, 0},
		"equal":      {4, 4, 4, 4, 4},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			q := NewPriorityQueue(minInt)
			for _, v := range input {
				q.Push(v)
			}
			equal(t, "Drain", slices.Collect(q.Drain()), slices.Sorted(slices.Values(input)))
		})
	}
}
