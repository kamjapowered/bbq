package set

import "iter"

// store membership of non negative integers using a bitset
//
// this type is not safe for concurrent use
type BitSet struct {
	words []uint64
	size  int
}

// create an empty bitset
//
// time: O(1)
func NewBitSet() *BitSet {
	return &BitSet{}
}

// add i to the set
//
// return false if i is negative
//
// time: O(1) amortised
func (b *BitSet) Add(i int) bool {
	if i < 0 {
		return false
	}

	w, m := bitIndex(i)
	if w >= len(b.words) {
		b.growToWord(w)
	}

	before := b.words[w]
	after := before | m
	if before != after {
		b.words[w] = after
		b.size++
		return true
	}

	return false
}

// remove i from the set
//
// return false if i is negative or not present
//
// time: O(1)
func (b *BitSet) Remove(i int) bool {
	if i < 0 {
		return false
	}

	w, m := bitIndex(i)
	if w >= len(b.words) {
		return false
	}

	before := b.words[w]
	if before&m == 0 {
		return false
	}

	b.words[w] = before &^ m
	b.size--
	return true
}

// report whether i is present
//
// return false if i is negative
//
// time: O(1)
func (b *BitSet) Has(i int) bool {
	if i < 0 {
		return false
	}

	w, m := bitIndex(i)
	if w >= len(b.words) {
		return false
	}

	return b.words[w]&m != 0
}

// report whether no elements are stored
//
// time: O(1)
func (b *BitSet) Empty() bool {
	return b.size == 0
}

// return the number of stored elements
//
// time: O(1)
func (b *BitSet) Len() int {
	return b.size
}

// remove all elements but keep allocated storage
//
// time: O(n)
func (b *BitSet) Clear() {
	for i := range b.words {
		b.words[i] = 0
	}
	b.size = 0
}

// return the maximum representable index without growing
//
// return -1 if the set has no allocated storage
//
// time: O(1)
func (b *BitSet) MaxIndex() int {
	if len(b.words) == 0 {
		return -1
	}
	return len(b.words)*64 - 1
}

// iterate over stored indices in increasing order
//
// iteration stops if yield returns false
//
// time: O(nwords + k) where k is number of set bits
func (b *BitSet) Iter() iter.Seq[int] {
	return func(yield func(int) bool) {
		for wi, word := range b.words {
			for word != 0 {
				lsb := word & -word
				bit := trailingZeros64(word)
				idx := wi*64 + bit
				if !yield(idx) {
					return
				}
				word ^= lsb
			}
		}
	}
}

// iterate over stored indices in increasing order while removing them
//
// iteration stops if yield returns false
//
// time: O(nwords + k) where k is number of set bits
func (b *BitSet) Drain() iter.Seq[int] {
	return func(yield func(int) bool) {
		for wi := range b.words {
			word := b.words[wi]
			for word != 0 {
				lsb := word & -word
				bit := trailingZeros64(word)
				idx := wi*64 + bit

				b.words[wi] &^= lsb
				b.size--

				if !yield(idx) {
					return
				}

				word = b.words[wi]
			}
		}
	}
}

// return a new set containing all elements present in either set
//
// time: O(n)
func (b *BitSet) Union(other *BitSet) *BitSet {
	out := NewBitSet()
	n := max(len(b.words), len(other.words))

	out.words = make([]uint64, n)

	for i := range n {
		var a uint64
		var c uint64

		if i < len(b.words) {
			a = b.words[i]
		}
		if i < len(other.words) {
			c = other.words[i]
		}

		out.words[i] = a | c
		out.size += popcount64(out.words[i])
	}

	return out
}

// return a new set containing all elements present in both sets
//
// time: O(n)
func (b *BitSet) Intersection(other *BitSet) *BitSet {
	out := NewBitSet()
	n := min(len(b.words), len(other.words))

	out.words = make([]uint64, n)

	for i := range n {
		out.words[i] = b.words[i] & other.words[i]
		out.size += popcount64(out.words[i])
	}

	return out
}

// return a new set containing all elements present in this set but not the other
//
// time: O(n)
func (b *BitSet) Difference(other *BitSet) *BitSet {
	out := NewBitSet()
	n := len(b.words)

	out.words = make([]uint64, n)

	for i := range n {
		a := b.words[i]
		if i < len(other.words) {
			a &^= other.words[i]
		}
		out.words[i] = a
		out.size += popcount64(a)
	}

	return out
}

// copy the set into a new bitset
//
// time: O(n)
func (b *BitSet) Clone() *BitSet {
	out := NewBitSet()
	out.words = make([]uint64, len(b.words))
	copy(out.words, b.words)
	out.size = b.size
	return out
}

func (b *BitSet) growToWord(w int) {
	if w < len(b.words) {
		return
	}

	newLen := len(b.words)
	if newLen == 0 {
		newLen = 1
	}
	for newLen <= w {
		newLen *= 2
	}

	words := make([]uint64, newLen)
	copy(words, b.words)
	b.words = words
}

func bitIndex(i int) (word int, mask uint64) {
	word = i >> 6
	mask = uint64(1) << uint(i&63)
	return word, mask
}

func trailingZeros64(x uint64) int {
	// x must be non zero
	n := 0
	for x&1 == 0 {
		x >>= 1
		n++
	}
	return n
}

func popcount64(x uint64) int {
	c := 0
	for x != 0 {
		x &= x - 1
		c++
	}
	return c
}
