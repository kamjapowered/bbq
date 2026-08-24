package set

import (
	"slices"
	"testing"
)

// build a bitset from a literal, so a table can name its operands inline
func bitsOf(is ...int) *BitSet {
	b := NewBitSet()
	for _, i := range is {
		b.Add(i)
	}
	return b
}

// pin contents, Len and Empty in one go. Iter is documented as increasing, so
// the yielded order is part of the contract
func checkBits(t *testing.T, label string, b *BitSet, want ...int) {
	t.Helper()

	want = slices.Sorted(slices.Values(want))
	if got := slices.Collect(b.Iter()); !slices.Equal(got, want) {
		t.Fatalf("%s = %v, want %v", label, got, want)
	}
	if got := b.Len(); got != len(want) {
		t.Fatalf("%s: Len = %d, want %d", label, got, len(want))
	}
	if got, w := b.Empty(), len(want) == 0; got != w {
		t.Fatalf("%s: Empty = %v, want %v", label, got, w)
	}
	for _, i := range want {
		if !b.Has(i) {
			t.Fatalf("%s: Has(%d) = false, want true", label, i)
		}
	}
}

func TestBitSetNewIsEmpty(t *testing.T) {
	b := NewBitSet()

	checkBits(t, "new bitset", b)
	if got := b.MaxIndex(); got != -1 {
		t.Fatalf("MaxIndex = %d on an unallocated bitset, want -1", got)
	}
	if b.Has(0) {
		t.Fatal("Has(0) = true on an unallocated bitset")
	}
	if b.Remove(0) {
		t.Fatal("Remove(0) = true on an unallocated bitset")
	}
}

func TestBitSetAddReportsWhetherItInserted(t *testing.T) {
	b := NewBitSet()

	if !b.Add(5) {
		t.Fatal("first Add(5) = false, want true")
	}
	if b.Add(5) {
		t.Fatal("second Add(5) = true, want false")
	}
	checkBits(t, "after adding 5 twice", b, 5)
}

func TestBitSetIndexZero(t *testing.T) {
	b := NewBitSet()

	if !b.Add(0) {
		t.Fatal("Add(0) = false, want true")
	}
	checkBits(t, "after Add(0)", b, 0)
	if got := b.MaxIndex(); got != 63 {
		t.Fatalf("MaxIndex = %d after Add(0), want 63", got)
	}

	if !b.Remove(0) {
		t.Fatal("Remove(0) = false, want true")
	}
	checkBits(t, "after Remove(0)", b)
}

func TestBitSetNegativeIndicesAreRejected(t *testing.T) {
	b := bitsOf(1)

	for _, i := range []int{-1, -64, -65} {
		if b.Add(i) {
			t.Fatalf("Add(%d) = true, want false", i)
		}
		if b.Remove(i) {
			t.Fatalf("Remove(%d) = true, want false", i)
		}
		if b.Has(i) {
			t.Fatalf("Has(%d) = true, want false", i)
		}
	}

	checkBits(t, "after the negative indices", b, 1)
}

func TestBitSetWordBoundaries(t *testing.T) {
	b := bitsOf(63, 64, 65)

	checkBits(t, "across a word boundary", b, 63, 64, 65)
	if got := b.MaxIndex(); got != 127 {
		t.Fatalf("MaxIndex = %d, want 127", got)
	}

	b.Remove(64)
	checkBits(t, "after removing the first bit of word 1", b, 63, 65)
}

func TestBitSetLargeSparseIndex(t *testing.T) {
	b := bitsOf(0, 1000)

	checkBits(t, "a sparse bitset", b, 0, 1000)

	// storage doubles, so it holds 16 words after reaching index 1000
	if got := b.MaxIndex(); got != 1023 {
		t.Fatalf("MaxIndex = %d, want 1023", got)
	}
	if b.Has(999) || b.Has(1001) {
		t.Fatal("a neighbour of the sparse index is set")
	}
	// an index inside the allocated storage but never added
	if b.Has(1023) {
		t.Fatal("Has(1023) = true, want false")
	}
	// and one past it
	if b.Has(4096) {
		t.Fatal("Has(4096) = true, want false")
	}

	b.Remove(1000)
	checkBits(t, "after removing the sparse index", b, 0)
	if got := b.MaxIndex(); got != 1023 {
		t.Fatalf("MaxIndex = %d after Remove, want the storage kept at 1023", got)
	}
}

func TestBitSetRemoveAbsent(t *testing.T) {
	b := bitsOf(1, 2)

	if b.Remove(3) {
		t.Fatal("Remove(3) = true for an absent index inside storage")
	}
	if b.Remove(9999) {
		t.Fatal("Remove(9999) = true for an index past storage")
	}
	if !b.Remove(1) {
		t.Fatal("Remove(1) = false, want true")
	}
	if b.Remove(1) {
		t.Fatal("second Remove(1) = true, want false")
	}
	checkBits(t, "after the removals", b, 2)
}

func TestBitSetClearKeepsStorage(t *testing.T) {
	b := bitsOf(0, 5, 200)
	before := b.MaxIndex()

	b.Clear()
	checkBits(t, "after Clear", b)
	if got := b.MaxIndex(); got != before {
		t.Fatalf("MaxIndex = %d after Clear, want the storage kept at %d", got, before)
	}

	b.Add(5)
	checkBits(t, "after Clear then Add(5)", b, 5)
}

func TestBitSetCloneIsIndependent(t *testing.T) {
	b := bitsOf(1, 70, 200)
	c := b.Clone()

	checkBits(t, "the clone", c, 1, 70, 200)

	c.Add(3)
	c.Remove(70)
	checkBits(t, "the original after mutating the clone", b, 1, 70, 200)
	checkBits(t, "the clone after its own mutation", c, 1, 3, 200)

	b.Add(500)
	b.Remove(1)
	checkBits(t, "the clone after mutating the original", c, 1, 3, 200)
	checkBits(t, "the original after its own mutation", b, 70, 200, 500)
}

func TestBitSetCloneOfAnEmptySet(t *testing.T) {
	b := NewBitSet()
	c := b.Clone()

	checkBits(t, "the clone of an empty bitset", c)
	if got := c.MaxIndex(); got != -1 {
		t.Fatalf("MaxIndex = %d, want -1", got)
	}

	c.Add(0)
	checkBits(t, "the original after the clone grew", b)
}

func TestBitSetIterIsIncreasingAndNonMutating(t *testing.T) {
	b := bitsOf(200, 0, 64, 63, 1)

	got := slices.Collect(b.Iter())
	if want := []int{0, 1, 63, 64, 200}; !slices.Equal(got, want) {
		t.Fatalf("Iter = %v, want %v", got, want)
	}

	// a second pass yields the same, so Iter does not consume
	if again := slices.Collect(b.Iter()); !slices.Equal(again, got) {
		t.Fatalf("second Iter = %v, want %v", again, got)
	}
	checkBits(t, "after Iter", b, 0, 1, 63, 64, 200)
}

func TestBitSetIterEarlyBreakDoesNotMutate(t *testing.T) {
	b := bitsOf(1, 2, 3)

	var seen []int
	for i := range b.Iter() {
		seen = append(seen, i)
		if len(seen) == 2 {
			break
		}
	}

	if want := []int{1, 2}; !slices.Equal(seen, want) {
		t.Fatalf("early break yielded %v, want %v", seen, want)
	}
	checkBits(t, "after an early break out of Iter", b, 1, 2, 3)
}

func TestBitSetDrainEmptiesTheSet(t *testing.T) {
	b := bitsOf(0, 63, 64, 300)

	got := slices.Collect(b.Drain())
	if want := []int{0, 63, 64, 300}; !slices.Equal(got, want) {
		t.Fatalf("Drain = %v, want %v", got, want)
	}
	checkBits(t, "after a full Drain", b)

	if got := slices.Collect(b.Drain()); len(got) != 0 {
		t.Fatalf("Drain of an empty bitset = %v, want nothing", got)
	}
}

func TestBitSetDrainEarlyBreakKeepsTheRest(t *testing.T) {
	b := bitsOf(1, 2, 100)

	var seen []int
	for i := range b.Drain() {
		seen = append(seen, i)
		break
	}

	if want := []int{1}; !slices.Equal(seen, want) {
		t.Fatalf("early break yielded %v, want %v", seen, want)
	}
	// breaking rejects the index, so it stays alongside the rest
	checkBits(t, "after a break out of Drain", b, 1, 2, 100)
}

// the set algebra, over the empty, disjoint, identical and overlapping cases
func TestBitSetAlgebra(t *testing.T) {
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
			a:            []int{0, 2},
			b:            []int{1, 3},
			union:        []int{0, 1, 2, 3},
			differenceAB: []int{0, 2},
			differenceBA: []int{1, 3},
		},
		{
			name:         "identical",
			a:            []int{1, 64, 200},
			b:            []int{1, 64, 200},
			union:        []int{1, 64, 200},
			intersection: []int{1, 64, 200},
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
			name:         "operands of different word counts",
			a:            []int{1, 300},
			b:            []int{1, 5},
			union:        []int{1, 5, 300},
			intersection: []int{1},
			differenceAB: []int{300},
			differenceBA: []int{5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, b := bitsOf(tt.a...), bitsOf(tt.b...)

			checkBits(t, "a|b", a.Union(b), tt.union...)
			checkBits(t, "b|a", b.Union(a), tt.union...)
			checkBits(t, "a&b", a.Intersection(b), tt.intersection...)
			checkBits(t, "b&a", b.Intersection(a), tt.intersection...)
			checkBits(t, "a-b", a.Difference(b), tt.differenceAB...)
			checkBits(t, "b-a", b.Difference(a), tt.differenceBA...)

			// no operation may touch its operands
			checkBits(t, "a after the algebra", a, tt.a...)
			checkBits(t, "b after the algebra", b, tt.b...)
		})
	}
}

func TestBitSetAlgebraReturnsAnIndependentSet(t *testing.T) {
	a, b := bitsOf(1, 2), bitsOf(2, 3)

	u := a.Union(b)
	u.Add(99)
	u.Remove(1)

	checkBits(t, "a after mutating the union", a, 1, 2)
	checkBits(t, "b after mutating the union", b, 2, 3)
}

func TestBitSetDifferenceIsNotCommutative(t *testing.T) {
	a, b := bitsOf(1, 2, 3), bitsOf(3, 4)

	ab := slices.Collect(a.Difference(b).Iter())
	ba := slices.Collect(b.Difference(a).Iter())

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

func TestBitSetAlgebraAgainstItself(t *testing.T) {
	a := bitsOf(1, 64, 200)

	checkBits(t, "a|a", a.Union(a), 1, 64, 200)
	checkBits(t, "a&a", a.Intersection(a), 1, 64, 200)
	checkBits(t, "a-a", a.Difference(a))
	checkBits(t, "a after the algebra", a, 1, 64, 200)
}

func TestBitSetAllBitsOfAWord(t *testing.T) {
	b := NewBitSet()
	for i := range 64 {
		b.Add(i)
	}

	if got := b.Len(); got != 64 {
		t.Fatalf("Len = %d, want 64", got)
	}
	if got := b.MaxIndex(); got != 63 {
		t.Fatalf("MaxIndex = %d, want 63", got)
	}

	got := slices.Collect(b.Iter())
	if len(got) != 64 {
		t.Fatalf("Iter yielded %d indices, want 64", len(got))
	}
	for i, v := range got {
		if v != i {
			t.Fatalf("Iter yielded %d at position %d", v, i)
		}
	}

	// the top bit of a word is the sign bit of the underlying uint64
	if !b.Has(63) {
		t.Fatal("Has(63) = false, want true")
	}
	if !b.Remove(63) {
		t.Fatal("Remove(63) = false, want true")
	}
	if b.Has(63) {
		t.Fatal("Has(63) = true after Remove")
	}
	if got := b.Len(); got != 63 {
		t.Fatalf("Len = %d after removing the top bit, want 63", got)
	}
}
