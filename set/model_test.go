package set

import (
	"fmt"
	"maps"
	"math/rand/v2"
	"slices"
	"testing"
)

// a deterministic script of operations run against both the container and a
// plain go map, compared after every step. the key space is kept small so
// repeated adds, removals of absent elements and overlapping operands all
// happen often.
const (
	modelSteps = 4000
	modelKeys  = 24
)

func newModelRand() *rand.Rand {
	return rand.New(rand.NewPCG(0x5eed, 0xf00d))
}

func modelSorted(m map[int]struct{}) []int {
	return slices.Sorted(maps.Keys(m))
}

func TestHashSetMatchesAPlainMap(t *testing.T) {
	h := NewHashSet[int]()
	model := make(map[int]struct{})

	rng := newModelRand()

	for step := range modelSteps {
		k := rng.IntN(modelKeys)

		switch rng.IntN(10) {
		case 0, 1, 2, 3, 4: // add
			h.Add(k)
			model[k] = struct{}{}

		case 5, 6: // remove
			h.Remove(k)
			delete(model, k)

		case 7: // read one element back
			_, want := model[k]
			if got := h.Has(k); got != want {
				t.Fatalf("step %d: Has(%d) = %v, want %v", step, k, got, want)
			}

		case 8: // peek the head of a drain, then reject it with a break
			for v := range h.Drain() {
				if _, ok := model[v]; !ok {
					t.Fatalf("step %d: Drain yielded %d, which is not in the model", step, v)
				}
				// the break rejects v, so the model keeps it
				break
			}

		case 9: // clear, rarely
			if rng.IntN(20) == 0 {
				h.Clear()
				clear(model)
			}
		}

		if got := h.Len(); got != len(model) {
			t.Fatalf("step %d: Len = %d, want %d", step, got, len(model))
		}
		if got, want := h.Empty(), len(model) == 0; got != want {
			t.Fatalf("step %d: Empty = %v, want %v", step, got, want)
		}

		want := modelSorted(model)
		if got := slices.Sorted(h.Iter()); !slices.Equal(got, want) {
			t.Fatalf("step %d: Iter = %v, want %v", step, got, want)
		}
		if got := hashElems(h); !slices.Equal(got, want) {
			t.Fatalf("step %d: ToSlice = %v, want %v", step, got, want)
		}
	}
}

func TestHashSetAlgebraMatchesPlainMaps(t *testing.T) {
	rng := newModelRand()

	for round := range 200 {
		a, b := NewHashSet[int](), NewHashSet[int]()
		ma, mb := make(map[int]struct{}), make(map[int]struct{})

		for range rng.IntN(modelKeys) {
			k := rng.IntN(modelKeys)
			a.Add(k)
			ma[k] = struct{}{}
		}
		for range rng.IntN(modelKeys) {
			k := rng.IntN(modelKeys)
			b.Add(k)
			mb[k] = struct{}{}
		}

		union := make(map[int]struct{})
		maps.Copy(union, ma)
		maps.Copy(union, mb)

		intersection := make(map[int]struct{})
		difference := make(map[int]struct{})
		for k := range ma {
			if _, ok := mb[k]; ok {
				intersection[k] = struct{}{}
			} else {
				difference[k] = struct{}{}
			}
		}

		if got, want := hashElems(a.Union(b)), modelSorted(union); !slices.Equal(got, want) {
			t.Fatalf("round %d: Union = %v, want %v", round, got, want)
		}
		if got, want := hashElems(a.Intersection(b)), modelSorted(intersection); !slices.Equal(got, want) {
			t.Fatalf("round %d: Intersection = %v, want %v", round, got, want)
		}
		if got, want := hashElems(a.Difference(b)), modelSorted(difference); !slices.Equal(got, want) {
			t.Fatalf("round %d: Difference = %v, want %v", round, got, want)
		}

		// the operands survive untouched
		if got, want := hashElems(a), modelSorted(ma); !slices.Equal(got, want) {
			t.Fatalf("round %d: a = %v after the algebra, want %v", round, got, want)
		}
		if got, want := hashElems(b), modelSorted(mb); !slices.Equal(got, want) {
			t.Fatalf("round %d: b = %v after the algebra, want %v", round, got, want)
		}
	}
}

// indices are drawn mostly from a dense low range, with the occasional sparse
// one, so word growth and multi word operands are exercised
func modelIndex(rng *rand.Rand) int {
	switch rng.IntN(8) {
	case 0:
		return rng.IntN(1000)
	default:
		return rng.IntN(130)
	}
}

func TestBitSetMatchesAPlainMap(t *testing.T) {
	b := NewBitSet()
	model := make(map[int]struct{})

	rng := newModelRand()

	for step := range modelSteps {
		i := modelIndex(rng)

		switch rng.IntN(10) {
		case 0, 1, 2, 3, 4: // add, reporting whether it inserted
			_, had := model[i]
			if got := b.Add(i); got != !had {
				t.Fatalf("step %d: Add(%d) = %v, want %v", step, i, got, !had)
			}
			model[i] = struct{}{}

		case 5, 6: // remove, reporting whether it removed
			_, had := model[i]
			if got := b.Remove(i); got != had {
				t.Fatalf("step %d: Remove(%d) = %v, want %v", step, i, got, had)
			}
			delete(model, i)

		case 7: // read one index back
			_, want := model[i]
			if got := b.Has(i); got != want {
				t.Fatalf("step %d: Has(%d) = %v, want %v", step, i, got, want)
			}

		case 8: // peek the head of a drain, then reject it with a break
			for v := range b.Drain() {
				if want := modelSorted(model); len(want) == 0 || v != want[0] {
					t.Fatalf("step %d: Drain yielded %d first, want %v", step, v, want)
				}
				// the break rejects v, so the model keeps it
				break
			}

		case 9: // clear, rarely
			if rng.IntN(20) == 0 {
				b.Clear()
				clear(model)
			}
		}

		if got := b.Len(); got != len(model) {
			t.Fatalf("step %d: Len = %d, want %d", step, got, len(model))
		}
		if got, want := b.Empty(), len(model) == 0; got != want {
			t.Fatalf("step %d: Empty = %v, want %v", step, got, want)
		}

		want := modelSorted(model)
		if got := slices.Collect(b.Iter()); !slices.Equal(got, want) {
			t.Fatalf("step %d: Iter = %v, want %v", step, got, want)
		}
		if got := slices.Collect(b.Clone().Iter()); !slices.Equal(got, want) {
			t.Fatalf("step %d: Clone Iter = %v, want %v", step, got, want)
		}
		// storage never shrinks below the largest index ever stored
		if len(want) > 0 && b.MaxIndex() < want[len(want)-1] {
			t.Fatalf("step %d: MaxIndex = %d, want at least %d", step, b.MaxIndex(), want[len(want)-1])
		}
	}
}

func TestBitSetAlgebraMatchesPlainMaps(t *testing.T) {
	rng := newModelRand()

	for round := range 200 {
		a, b := NewBitSet(), NewBitSet()
		ma, mb := make(map[int]struct{}), make(map[int]struct{})

		for range rng.IntN(modelKeys) {
			i := modelIndex(rng)
			a.Add(i)
			ma[i] = struct{}{}
		}
		for range rng.IntN(modelKeys) {
			i := modelIndex(rng)
			b.Add(i)
			mb[i] = struct{}{}
		}

		union := make(map[int]struct{})
		maps.Copy(union, ma)
		maps.Copy(union, mb)

		intersection := make(map[int]struct{})
		difference := make(map[int]struct{})
		for i := range ma {
			if _, ok := mb[i]; ok {
				intersection[i] = struct{}{}
			} else {
				difference[i] = struct{}{}
			}
		}

		label := fmt.Sprintf("round %d, a = %v, b = %v", round, modelSorted(ma), modelSorted(mb))

		checkBits(t, label+": union", a.Union(b), modelSorted(union)...)
		checkBits(t, label+": intersection", a.Intersection(b), modelSorted(intersection)...)
		checkBits(t, label+": difference", a.Difference(b), modelSorted(difference)...)

		// the operands survive untouched
		checkBits(t, label+": a after the algebra", a, modelSorted(ma)...)
		checkBits(t, label+": b after the algebra", b, modelSorted(mb)...)
	}
}
