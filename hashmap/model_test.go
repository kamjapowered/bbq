package hashmap

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// a deterministic script of operations run against both the container and a
// plain go map, compared after every step. the key space is kept small so
// collisions, rebindings and removals of absent entries all happen often.
const (
	modelSteps = 4000
	modelKeys  = 24
)

func newModelRand() *rand.Rand {
	return rand.New(rand.NewPCG(0x5eed, 0xf00d))
}

func TestDenseMapMatchesAPlainMap(t *testing.T) {
	m := NewDenseMap[int, int]()
	model := make(map[int]int)

	rng := newModelRand()

	for step := range modelSteps {
		k := rng.IntN(modelKeys)

		switch rng.IntN(10) {
		case 0, 1, 2, 3, 4: // set
			v := rng.IntN(1000)
			m.Set(k, v)
			model[k] = v

		case 5, 6: // remove
			_, want := model[k]
			if got := m.Remove(k); got != want {
				t.Fatalf("step %d: Remove(%d) = %v, want %v", step, k, got, want)
			}
			delete(model, k)

		case 7: // mutate through a reference
			p, ok := m.GetRef(k)
			if _, want := model[k]; ok != want {
				t.Fatalf("step %d: GetRef(%d) ok = %v, want %v", step, k, ok, want)
			}
			if ok {
				*p += 1
				model[k]++
			}

		case 8: // read back one key
			got, ok := m.Get(k)
			want, wantOK := model[k]
			if ok != wantOK || (ok && got != want) {
				t.Fatalf("step %d: Get(%d) = %d, %v, want %d, %v", step, k, got, ok, want, wantOK)
			}

		case 9: // clear, rarely
			if rng.IntN(20) == 0 {
				m.Clear()
				clear(model)
			}
		}

		if got := m.Len(); got != len(model) {
			t.Fatalf("step %d: Len = %d, want %d", step, got, len(model))
		}
		if got := m.Empty(); got != (len(model) == 0) {
			t.Fatalf("step %d: Empty = %v, want %v", step, got, len(model) == 0)
		}
		checkDenseInvariant(t, m)
	}

	// the full contents must agree at the end, iterators included
	for k, want := range model {
		if got, ok := m.Get(k); !ok || got != want {
			t.Fatalf("final: Get(%d) = %d, %v, want %d, true", k, got, ok, want)
		}
		if !m.Has(k) {
			t.Fatalf("final: Has(%d) = false", k)
		}
	}

	seen := make(map[int]int)
	for k, v := range m.Iter() {
		if _, dup := seen[k]; dup {
			t.Fatalf("final: Iter yielded key %d twice", k)
		}
		seen[k] = v
	}
	if len(seen) != len(model) {
		t.Fatalf("final: Iter yielded %d pairs, want %d", len(seen), len(model))
	}
	for k, want := range model {
		if seen[k] != want {
			t.Fatalf("final: Iter yielded %d -> %d, want %d", k, seen[k], want)
		}
	}

	keys := slices.Collect(m.IterKeys())
	if !slices.Equal(keys, m.Keys()) {
		t.Fatalf("final: IterKeys = %v, Keys = %v", keys, m.Keys())
	}
	values := slices.Collect(m.IterValues())
	if !slices.Equal(values, m.Values()) {
		t.Fatalf("final: IterValues = %v, Values = %v", values, m.Values())
	}
	for i, k := range keys {
		if values[i] != model[k] {
			t.Fatalf("final: row %d holds %d -> %d, want %d", i, k, values[i], model[k])
		}
	}
}

func TestBiHashMapMatchesAPlainMapPair(t *testing.T) {
	b := NewBiHashMap[int, int]()
	fwd := make(map[int]int)
	rev := make(map[int]int)

	rng := newModelRand()

	for step := range modelSteps {
		k := rng.IntN(modelKeys)
		v := rng.IntN(modelKeys)

		switch rng.IntN(10) {
		case 0, 1, 2, 3, 4: // set, rebinding both sides as needed
			_, hadK := fwd[k]
			_, hadV := rev[v]

			got := b.Set(k, v)
			if want := !(hadK || hadV); got != want {
				t.Fatalf("step %d: Set(%d, %d) = %v, want %v", step, k, v, got, want)
			}

			if oldV, ok := fwd[k]; ok {
				delete(rev, oldV)
			}
			if oldK, ok := rev[v]; ok {
				delete(fwd, oldK)
			}
			fwd[k] = v
			rev[v] = k

		case 5, 6: // remove by key
			_, want := fwd[k]
			if got := b.RemoveByKey(k); got != want {
				t.Fatalf("step %d: RemoveByKey(%d) = %v, want %v", step, k, got, want)
			}
			if oldV, ok := fwd[k]; ok {
				delete(fwd, k)
				delete(rev, oldV)
			}

		case 7, 8: // remove by value
			_, want := rev[v]
			if got := b.RemoveByValue(v); got != want {
				t.Fatalf("step %d: RemoveByValue(%d) = %v, want %v", step, v, got, want)
			}
			if oldK, ok := rev[v]; ok {
				delete(rev, v)
				delete(fwd, oldK)
			}

		case 9: // clear, rarely
			if rng.IntN(20) == 0 {
				b.Clear()
				clear(fwd)
				clear(rev)
			}
		}

		if got := b.Len(); got != len(fwd) {
			t.Fatalf("step %d: Len = %d, want %d", step, got, len(fwd))
		}
		if got := b.Empty(); got != (len(fwd) == 0) {
			t.Fatalf("step %d: Empty = %v, want %v", step, got, len(fwd) == 0)
		}

		gotV, gotOK := b.GetByKey(k)
		wantV, wantOK := fwd[k]
		if gotOK != wantOK || (gotOK && gotV != wantV) {
			t.Fatalf("step %d: GetByKey(%d) = %d, %v, want %d, %v", step, k, gotV, gotOK, wantV, wantOK)
		}

		gotK, gotOK := b.GetByValue(v)
		wantK, wantOK := rev[v]
		if gotOK != wantOK || (gotOK && gotK != wantK) {
			t.Fatalf("step %d: GetByValue(%d) = %d, %v, want %d, %v", step, v, gotK, gotOK, wantK, wantOK)
		}

		checkBiInvariant(t, b)
	}

	for k, v := range fwd {
		if got, ok := b.GetByKey(k); !ok || got != v {
			t.Fatalf("final: GetByKey(%d) = %d, %v, want %d, true", k, got, ok, v)
		}
		if !b.HasKey(k) {
			t.Fatalf("final: HasKey(%d) = false", k)
		}
	}
	for v, k := range rev {
		if got, ok := b.GetByValue(v); !ok || got != k {
			t.Fatalf("final: GetByValue(%d) = %d, %v, want %d, true", v, got, ok, k)
		}
		if !b.HasValue(v) {
			t.Fatalf("final: HasValue(%d) = false", v)
		}
	}

	keys := b.Keys()
	values := b.Values()
	if len(keys) != len(fwd) || len(values) != len(rev) {
		t.Fatalf("final: %d keys and %d values, want %d and %d", len(keys), len(values), len(fwd), len(rev))
	}

	drained := 0
	for k, v := range b.Drain() {
		if fwd[k] != v {
			t.Fatalf("final: Drain yielded %d -> %d, want %d", k, v, fwd[k])
		}
		drained++
	}
	if drained != len(fwd) {
		t.Fatalf("final: drained %d pairs, want %d", drained, len(fwd))
	}
	if !b.Empty() {
		t.Fatalf("final: Len = %d after a full drain, want 0", b.Len())
	}
}
