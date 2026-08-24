package set

import (
	"slices"
	"testing"
)

// build a set from a literal, so a table can name its operands inline
func hashOf(vs ...int) *HashSet[int] {
	h := NewHashSet[int]()
	for _, v := range vs {
		h.Add(v)
	}
	return h
}

// the contents as a sorted slice, since iteration order is unspecified
func hashElems(h *HashSet[int]) []int {
	out := h.ToSlice()
	slices.Sort(out)
	return out
}

// pin contents, Len and Empty in one go
func checkHash(t *testing.T, label string, h *HashSet[int], want ...int) {
	t.Helper()

	// the caller may spread a table field, so never sort in place
	want = slices.Sorted(slices.Values(want))
	if got := hashElems(h); !slices.Equal(got, want) {
		t.Fatalf("%s = %v, want %v", label, got, want)
	}
	if got := h.Len(); got != len(want) {
		t.Fatalf("%s: Len = %d, want %d", label, got, len(want))
	}
	if got, w := h.Empty(), len(want) == 0; got != w {
		t.Fatalf("%s: Empty = %v, want %v", label, got, w)
	}
	for _, v := range want {
		if !h.Has(v) {
			t.Fatalf("%s: Has(%d) = false, want true", label, v)
		}
	}
}

func TestHashSetNewIsEmpty(t *testing.T) {
	h := NewHashSet[int]()

	checkHash(t, "new set", h)
	if h.Has(0) {
		t.Fatal("Has(0) = true on an empty set")
	}
	if got := h.ToSlice(); got == nil {
		t.Fatal("ToSlice = nil on an empty set, want an empty slice")
	}
}

func TestHashSetAddIsIdempotent(t *testing.T) {
	h := NewHashSet[int]()

	for range 3 {
		h.Add(7)
	}
	checkHash(t, "after three Add(7)", h, 7)

	h.Add(8)
	checkHash(t, "after Add(8)", h, 7, 8)
}

func TestHashSetRemovePresentAndAbsent(t *testing.T) {
	h := hashOf(1, 2, 3)

	h.Remove(2)
	checkHash(t, "after Remove(2)", h, 1, 3)

	// removing an absent element is a no op
	h.Remove(2)
	h.Remove(99)
	checkHash(t, "after removing absent elements", h, 1, 3)

	// and from an empty set
	h.Remove(1)
	h.Remove(3)
	checkHash(t, "after removing everything", h)
	h.Remove(1)
	checkHash(t, "after removing from an empty set", h)
}

func TestHashSetClearKeepsTheSetUsable(t *testing.T) {
	h := hashOf(1, 2, 3)

	h.Clear()
	checkHash(t, "after Clear", h)

	h.Add(4)
	checkHash(t, "after Clear then Add(4)", h, 4)

	// clearing an already empty set is a no op
	h.Clear()
	h.Clear()
	checkHash(t, "after clearing twice", h)
}

func TestHashSetToSliceAgreesWithIter(t *testing.T) {
	h := hashOf(3, 1, 2)

	fromSlice := h.ToSlice()
	slices.Sort(fromSlice)

	fromIter := slices.Sorted(h.Iter())

	if !slices.Equal(fromSlice, fromIter) {
		t.Fatalf("ToSlice = %v, Iter = %v", fromSlice, fromIter)
	}
	if want := []int{1, 2, 3}; !slices.Equal(fromIter, want) {
		t.Fatalf("contents = %v, want %v", fromIter, want)
	}

	// the returned slice is a copy, so mutating it must not touch the set
	fromSlice[0] = 99
	checkHash(t, "after mutating the ToSlice result", h, 1, 2, 3)
}

func TestHashSetIterVisitsEachElementOnce(t *testing.T) {
	h := hashOf(1, 2, 3, 4)

	seen := make(map[int]int)
	for v := range h.Iter() {
		seen[v]++
	}

	if len(seen) != 4 {
		t.Fatalf("Iter yielded %d distinct elements, want 4", len(seen))
	}
	for v, n := range seen {
		if n != 1 {
			t.Fatalf("Iter yielded %d %d times, want once", v, n)
		}
	}

	// iterating twice yields the same thing, so Iter does not consume
	if got := slices.Sorted(h.Iter()); !slices.Equal(got, []int{1, 2, 3, 4}) {
		t.Fatalf("second Iter = %v, want [1 2 3 4]", got)
	}
	checkHash(t, "after Iter", h, 1, 2, 3, 4)
}

func TestHashSetIterEarlyBreakDoesNotMutate(t *testing.T) {
	h := hashOf(1, 2, 3)

	n := 0
	for range h.Iter() {
		n++
		break
	}

	if n != 1 {
		t.Fatalf("early break visited %d elements, want 1", n)
	}
	checkHash(t, "after an early break out of Iter", h, 1, 2, 3)
}

func TestHashSetDrainEmptiesTheSet(t *testing.T) {
	h := hashOf(1, 2, 3)

	got := slices.Sorted(h.Drain())
	if want := []int{1, 2, 3}; !slices.Equal(got, want) {
		t.Fatalf("Drain yielded %v, want %v", got, want)
	}
	checkHash(t, "after a full Drain", h)

	// draining an empty set yields nothing
	if got := slices.Sorted(h.Drain()); len(got) != 0 {
		t.Fatalf("Drain of an empty set yielded %v, want nothing", got)
	}
}

func TestHashSetDrainEarlyBreakKeepsTheRest(t *testing.T) {
	h := hashOf(1, 2, 3)

	var taken int
	for v := range h.Drain() {
		taken = v
		break
	}

	// breaking rejects the element, so it stays and nothing was consumed
	if !h.Has(taken) {
		t.Fatalf("Drain yielded %d and removed it despite the break", taken)
	}
	if got := h.Len(); got != 3 {
		t.Fatalf("Len after a break out of Drain = %d, want 3", got)
	}

	rest := slices.Sorted(h.Iter())
	if want := []int{1, 2, 3}; !slices.Equal(rest, want) {
		t.Fatalf("remaining = %v, want %v", rest, want)
	}
}

// the set algebra, over the empty, disjoint, identical and overlapping cases
func TestHashSetAlgebra(t *testing.T) {
	tests := []struct {
		name         string
		a, b         []int
		union        []int
		intersection []int
		differenceAB []int
		differenceBA []int
	}{
		{
			name: "both empty",
		},
		{
			name:         "empty against non empty",
			b:            []int{1, 2},
			union:        []int{1, 2},
			differenceBA: []int{1, 2},
		},
		{
			name:         "non empty against empty",
			a:            []int{1, 2},
			union:        []int{1, 2},
			differenceAB: []int{1, 2},
		},
		{
			name:         "disjoint",
			a:            []int{1, 2},
			b:            []int{3, 4},
			union:        []int{1, 2, 3, 4},
			differenceAB: []int{1, 2},
			differenceBA: []int{3, 4},
		},
		{
			name:         "identical",
			a:            []int{1, 2, 3},
			b:            []int{1, 2, 3},
			union:        []int{1, 2, 3},
			intersection: []int{1, 2, 3},
		},
		{
			name:         "overlapping",
			a:            []int{1, 2, 3},
			b:            []int{2, 3, 4},
			union:        []int{1, 2, 3, 4},
			intersection: []int{2, 3},
			differenceAB: []int{1},
			differenceBA: []int{4},
		},
		{
			name:         "subset",
			a:            []int{1, 2, 3},
			b:            []int{2},
			union:        []int{1, 2, 3},
			intersection: []int{2},
			differenceAB: []int{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, b := hashOf(tt.a...), hashOf(tt.b...)

			checkHash(t, "a|b", a.Union(b), tt.union...)
			checkHash(t, "b|a", b.Union(a), tt.union...)
			checkHash(t, "a&b", a.Intersection(b), tt.intersection...)
			checkHash(t, "b&a", b.Intersection(a), tt.intersection...)
			checkHash(t, "a-b", a.Difference(b), tt.differenceAB...)
			checkHash(t, "b-a", b.Difference(a), tt.differenceBA...)

			// no operation may touch its operands
			checkHash(t, "a after the algebra", a, tt.a...)
			checkHash(t, "b after the algebra", b, tt.b...)
		})
	}
}

func TestHashSetAlgebraReturnsAnIndependentSet(t *testing.T) {
	a, b := hashOf(1, 2), hashOf(2, 3)

	u := a.Union(b)
	u.Add(99)
	u.Remove(1)

	checkHash(t, "a after mutating the union", a, 1, 2)
	checkHash(t, "b after mutating the union", b, 2, 3)
}

func TestHashSetDifferenceIsNotCommutative(t *testing.T) {
	a, b := hashOf(1, 2, 3), hashOf(3, 4)

	ab := hashElems(a.Difference(b))
	ba := hashElems(b.Difference(a))

	if slices.Equal(ab, ba) {
		t.Fatalf("a-b and b-a both = %v, want them to differ", ab)
	}
	if want := []int{1, 2}; !slices.Equal(ab, want) {
		t.Fatalf("a-b = %v, want %v", ab, want)
	}
	if want := []int{4}; !slices.Equal(ba, want) {
		t.Fatalf("b-a = %v, want %v", ba, want)
	}
}

func TestHashSetAlgebraAgainstItself(t *testing.T) {
	a := hashOf(1, 2, 3)

	checkHash(t, "a|a", a.Union(a), 1, 2, 3)
	checkHash(t, "a&a", a.Intersection(a), 1, 2, 3)
	checkHash(t, "a-a", a.Difference(a))
	checkHash(t, "a after the algebra", a, 1, 2, 3)
}

func TestHashSetWithStringElements(t *testing.T) {
	h := NewHashSet[string]()

	h.Add("a")
	h.Add("a")
	h.Add("")

	if got := h.Len(); got != 2 {
		t.Fatalf("Len = %d, want 2", got)
	}
	if !h.Has("") {
		t.Fatal(`Has("") = false, want true`)
	}
	if h.Has("b") {
		t.Fatal(`Has("b") = true, want false`)
	}

	got := slices.Sorted(h.Iter())
	if want := []string{"", "a"}; !slices.Equal(got, want) {
		t.Fatalf("contents = %q, want %q", got, want)
	}
}
