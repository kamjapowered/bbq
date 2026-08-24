package queue

import "testing"

func TestDequeToSlice(t *testing.T) {
	d := NewDeque[int]()
	if got := d.ToSlice(); len(got) != 0 {
		t.Fatalf("empty deque ToSlice = %v, want empty", got)
	}

	// push both ends so the ring wraps, then read front to back
	d.PushBack(3)
	d.PushBack(4)
	d.PushFront(2)
	d.PushFront(1)

	got := d.ToSlice()
	want := []int{1, 2, 3, 4}
	if len(got) != len(want) {
		t.Fatalf("ToSlice = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ToSlice = %v, want %v", got, want)
		}
	}

	// the copy is independent and the deque is unchanged
	got[0] = 99
	if v, _ := d.PeekFront(); v != 1 {
		t.Fatalf("ToSlice aliased the deque, front = %v", v)
	}
	if d.Len() != 4 {
		t.Fatalf("Len = %d, want 4", d.Len())
	}
}

func TestPriorityQueueToSlice(t *testing.T) {
	q := NewPriorityQueue(func(a, b int) bool { return a < b })
	if got := q.ToSlice(); len(got) != 0 {
		t.Fatalf("empty ToSlice = %v, want empty", got)
	}

	for _, v := range []int{5, 1, 4, 2, 3} {
		q.Push(v)
	}

	got := q.ToSlice()
	if len(got) != 5 {
		t.Fatalf("ToSlice = %v, want 5 elements", got)
	}
	// heap order, not priority order, but the root is the minimum
	if got[0] != 1 {
		t.Fatalf("ToSlice[0] = %d, want the heap root 1", got[0])
	}

	// the copy is independent and the queue still drains by priority
	got[0] = 99
	var drained []int
	for v := range q.Drain() {
		drained = append(drained, v)
	}
	for i, want := range []int{1, 2, 3, 4, 5} {
		if drained[i] != want {
			t.Fatalf("Drain = %v, want ascending", drained)
		}
	}
}
