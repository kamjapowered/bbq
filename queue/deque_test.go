package queue

import (
	"slices"
	"testing"
)

func TestDequeEmpty(t *testing.T) {
	d := NewDeque[int]()

	if !d.Empty() {
		t.Fatal("Empty = false on a fresh deque")
	}
	if got := d.Len(); got != 0 {
		t.Fatalf("Len = %d, want 0", got)
	}
	if v, ok := d.PeekFront(); ok || v != 0 {
		t.Fatalf("PeekFront = (%d, %v), want (0, false)", v, ok)
	}
	if v, ok := d.PeekBack(); ok || v != 0 {
		t.Fatalf("PeekBack = (%d, %v), want (0, false)", v, ok)
	}
	if v, ok := d.PopFront(); ok || v != 0 {
		t.Fatalf("PopFront = (%d, %v), want (0, false)", v, ok)
	}
	if v, ok := d.PopBack(); ok || v != 0 {
		t.Fatalf("PopBack = (%d, %v), want (0, false)", v, ok)
	}
	equal(t, "Iter", slices.Collect(d.Iter()), nil)
	equal(t, "Drain", slices.Collect(d.Drain()), nil)

	d.Clear()
	if !d.Empty() || d.Len() != 0 {
		t.Fatal("Clear broke an already empty deque")
	}
}

func TestDequeSingleElement(t *testing.T) {
	tests := map[string]func(*Deque[int]){
		"PushFront": func(d *Deque[int]) { d.PushFront(7) },
		"PushBack":  func(d *Deque[int]) { d.PushBack(7) },
	}

	for name, push := range tests {
		t.Run(name, func(t *testing.T) {
			d := NewDeque[int]()
			push(d)

			if d.Empty() || d.Len() != 1 {
				t.Fatalf("Len = %d, want 1", d.Len())
			}
			if v, ok := d.PeekFront(); !ok || v != 7 {
				t.Fatalf("PeekFront = (%d, %v), want (7, true)", v, ok)
			}
			if v, ok := d.PeekBack(); !ok || v != 7 {
				t.Fatalf("PeekBack = (%d, %v), want (7, true)", v, ok)
			}
			equal(t, "Iter", slices.Collect(d.Iter()), []int{7})
		})
	}
}

func TestDequePushBackIsFifoFromFront(t *testing.T) {
	d := NewDeque[int]()
	want := []int{1, 2, 3, 4, 5}
	for _, v := range want {
		d.PushBack(v)
	}

	equal(t, "Iter", slices.Collect(d.Iter()), want)
	if v, _ := d.PeekFront(); v != 1 {
		t.Fatalf("PeekFront = %d, want 1", v)
	}
	if v, _ := d.PeekBack(); v != 5 {
		t.Fatalf("PeekBack = %d, want 5", v)
	}

	var got []int
	for !d.Empty() {
		v, _ := d.PopFront()
		got = append(got, v)
	}
	equal(t, "PopFront order", got, want)
}

func TestDequePushFrontReversesIterOrder(t *testing.T) {
	d := NewDeque[int]()
	for _, v := range []int{1, 2, 3, 4, 5} {
		d.PushFront(v)
	}

	want := []int{5, 4, 3, 2, 1}
	equal(t, "Iter", slices.Collect(d.Iter()), want)

	var got []int
	for !d.Empty() {
		v, _ := d.PopFront()
		got = append(got, v)
	}
	equal(t, "PopFront order", got, want)
}

func TestDequePushBackPopBackIsLifo(t *testing.T) {
	d := NewDeque[int]()
	in := []int{1, 2, 3, 4, 5}
	for _, v := range in {
		d.PushBack(v)
	}

	want := []int{5, 4, 3, 2, 1}
	var got []int
	for !d.Empty() {
		v, _ := d.PopBack()
		got = append(got, v)
	}
	equal(t, "PopBack order", got, want)
}

func TestDequeBothEnds(t *testing.T) {
	d := NewDeque[int]()

	d.PushBack(3)
	d.PushFront(2)
	d.PushBack(4)
	d.PushFront(1)
	d.PushBack(5)

	equal(t, "Iter", slices.Collect(d.Iter()), []int{1, 2, 3, 4, 5})

	if v, _ := d.PopFront(); v != 1 {
		t.Fatalf("PopFront = %d, want 1", v)
	}
	if v, _ := d.PopBack(); v != 5 {
		t.Fatalf("PopBack = %d, want 5", v)
	}
	equal(t, "Iter", slices.Collect(d.Iter()), []int{2, 3, 4})
}

func TestDequeWraparoundAndGrowth(t *testing.T) {
	// alternate the two ends so the ring wraps repeatedly while it grows,
	// checked against a slice model
	d := NewDeque[int]()
	var model []int

	for i := range 40 {
		switch i % 4 {
		case 0, 1:
			d.PushBack(i)
			model = append(model, i)
		case 2:
			d.PushFront(i)
			model = append([]int{i}, model...)
		case 3:
			if len(model) > 0 {
				v, ok := d.PopFront()
				if !ok || v != model[0] {
					t.Fatalf("PopFront = (%d, %v), want (%d, true)", v, ok, model[0])
				}
				model = model[1:]
			}
		}

		if got := d.Len(); got != len(model) {
			t.Fatalf("step %d: Len = %d, want %d", i, got, len(model))
		}
		equal(t, "Iter", slices.Collect(d.Iter()), model)
	}

	// unwind from the back, still against the model
	for len(model) > 0 {
		v, ok := d.PopBack()
		last := model[len(model)-1]
		if !ok || v != last {
			t.Fatalf("PopBack = (%d, %v), want (%d, true)", v, ok, last)
		}
		model = model[:len(model)-1]
		equal(t, "Iter", slices.Collect(d.Iter()), model)
	}

	if !d.Empty() {
		t.Fatal("Empty = false after popping everything")
	}
}

func TestDequePushFrontAfterEmptying(t *testing.T) {
	// emptying leaves the head somewhere in the middle of the ring, and the
	// next push must still land as the only element
	d := NewDeque[int]()
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)
	for range 3 {
		d.PopFront()
	}

	d.PushFront(9)
	d.PushBack(10)
	equal(t, "Iter", slices.Collect(d.Iter()), []int{9, 10})
	if v, _ := d.PeekFront(); v != 9 {
		t.Fatalf("PeekFront = %d, want 9", v)
	}
	if v, _ := d.PeekBack(); v != 10 {
		t.Fatalf("PeekBack = %d, want 10", v)
	}
}

func TestDequeIterIsNonMutating(t *testing.T) {
	d := NewDeque[int]()
	d.PushFront(2)
	d.PushFront(1)
	d.PushBack(3)
	want := []int{1, 2, 3}

	equal(t, "Iter pass 1", slices.Collect(d.Iter()), want)
	equal(t, "Iter pass 2", slices.Collect(d.Iter()), want)
	if got := d.Len(); got != 3 {
		t.Fatalf("Len = %d after Iter, want 3", got)
	}
}

func TestDequeIterEarlyBreak(t *testing.T) {
	d := NewDeque[int]()
	for _, v := range []int{1, 2, 3, 4} {
		d.PushBack(v)
	}

	var got []int
	for v := range d.Iter() {
		got = append(got, v)
		if len(got) == 2 {
			break
		}
	}
	equal(t, "Iter prefix", got, []int{1, 2})
	if d.Len() != 4 {
		t.Fatalf("Len = %d after a broken Iter, want 4", d.Len())
	}
}

func TestDequeDrain(t *testing.T) {
	d := NewDeque[int]()
	d.PushFront(2)
	d.PushFront(1)
	d.PushBack(3)

	equal(t, "Drain", slices.Collect(d.Drain()), []int{1, 2, 3})
	if !d.Empty() || d.Len() != 0 {
		t.Fatalf("Len = %d after Drain, want 0", d.Len())
	}

	// still usable afterwards
	d.PushBack(4)
	equal(t, "Iter", slices.Collect(d.Iter()), []int{4})
}

func TestDequeDrainEarlyBreak(t *testing.T) {
	d := NewDeque[int]()
	for _, v := range []int{1, 2, 3, 4, 5} {
		d.PushBack(v)
	}

	var got []int
	for v := range d.Drain() {
		got = append(got, v)
		if len(got) == 2 {
			break
		}
	}
	equal(t, "Drain prefix", got, []int{1, 2})

	if d.Len() != 3 {
		t.Fatalf("Len = %d after a broken Drain, want 3", d.Len())
	}
	equal(t, "Iter", slices.Collect(d.Iter()), []int{3, 4, 5})
}

func TestDequeClear(t *testing.T) {
	d := NewDeque[int]()
	for _, v := range []int{1, 2, 3, 4} {
		d.PushBack(v)
	}
	// wrap the head before clearing
	d.PopFront()
	d.PopFront()
	d.PushBack(5)

	d.Clear()

	if !d.Empty() || d.Len() != 0 {
		t.Fatalf("Len = %d after Clear, want 0", d.Len())
	}
	if _, ok := d.PeekFront(); ok {
		t.Fatal("PeekFront succeeded after Clear")
	}
	if _, ok := d.PopBack(); ok {
		t.Fatal("PopBack succeeded after Clear")
	}
	equal(t, "Iter", slices.Collect(d.Iter()), nil)

	// clear keeps the allocated storage, so it starts from index zero again
	if cap := len(d.data); cap == 0 {
		t.Fatal("Clear released the allocated storage")
	}
	d.PushFront(9)
	equal(t, "Iter", slices.Collect(d.Iter()), []int{9})
}
