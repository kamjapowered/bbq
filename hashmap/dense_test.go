package hashmap

import (
	"slices"
	"testing"
)

// the dense slice, the value slice and the index must always agree: same
// length, and every key parked at the row the index claims for it.
func checkDenseInvariant[K comparable, V any](t *testing.T, m *DenseMap[K, V]) {
	t.Helper()

	if len(m.dense) != len(m.values) {
		t.Fatalf("dense len = %d, values len = %d", len(m.dense), len(m.values))
	}
	if len(m.location) != len(m.dense) {
		t.Fatalf("location len = %d, dense len = %d", len(m.location), len(m.dense))
	}
	for i, k := range m.dense {
		row, ok := m.location[k]
		if !ok {
			t.Fatalf("key %v at row %d missing from location", k, i)
		}
		if row != i {
			t.Fatalf("key %v sits at row %d but location says %d", k, i, row)
		}
	}
}

func denseOf(t *testing.T, pairs ...string) *DenseMap[string, int] {
	t.Helper()

	m := NewDenseMap[string, int]()
	for i, k := range pairs {
		m.Set(k, i)
	}
	return m
}

func denseKeys[V any](m *DenseMap[string, V]) []string {
	out := make([]string, 0, m.Len())
	for k := range m.IterKeys() {
		out = append(out, k)
	}
	return out
}

func TestDenseMapSetOverwritesExistingKey(t *testing.T) {
	m := NewDenseMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("a", 99)

	if got := m.Len(); got != 2 {
		t.Fatalf("Len = %d, want 2", got)
	}
	if v, ok := m.Get("a"); !ok || v != 99 {
		t.Fatalf("Get(a) = %d, %v, want 99, true", v, ok)
	}

	// overwriting must not move the key within the dense order
	if got := denseKeys(m); !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("keys = %v, want [a b]", got)
	}
	checkDenseInvariant(t, m)
}

func TestDenseMapGetMissingKeyReturnsZero(t *testing.T) {
	m := denseOf(t, "a")

	v, ok := m.Get("nope")
	if ok {
		t.Fatal("Get of an absent key reported true")
	}
	if v != 0 {
		t.Fatalf("Get of an absent key = %d, want 0", v)
	}

	if p, ok := m.GetRef("nope"); ok || p != nil {
		t.Fatalf("GetRef of an absent key = %v, %v, want nil, false", p, ok)
	}
}

func TestDenseMapGetRefMutationIsVisibleToGet(t *testing.T) {
	m := denseOf(t, "a", "b", "c")

	p, ok := m.GetRef("b")
	if !ok {
		t.Fatal("GetRef(b) returned false")
	}

	*p = 42

	if v, ok := m.Get("b"); !ok || v != 42 {
		t.Fatalf("Get(b) after mutation = %d, %v, want 42, true", v, ok)
	}
	if v := slices.Collect(m.IterValues()); !slices.Equal(v, []int{0, 42, 2}) {
		t.Fatalf("values = %v, want [0 42 2]", v)
	}

	// the other rows must be untouched
	if v, _ := m.Get("a"); v != 0 {
		t.Fatalf("Get(a) = %d, want 0", v)
	}
	if v, _ := m.Get("c"); v != 2 {
		t.Fatalf("Get(c) = %d, want 2", v)
	}
}

func TestDenseMapRemoveReportsPresence(t *testing.T) {
	m := denseOf(t, "a", "b")

	if !m.Remove("a") {
		t.Fatal("Remove of a present key returned false")
	}
	if m.Remove("a") {
		t.Fatal("Remove of an already removed key returned true")
	}
	if m.Remove("never") {
		t.Fatal("Remove of an absent key returned true")
	}
	if got := m.Len(); got != 1 {
		t.Fatalf("Len = %d, want 1", got)
	}
	checkDenseInvariant(t, m)
}

// the swap remove path: pulling any row out must leave every surviving key
// reachable, with its own value, and the index still pointing at the right row.
func TestDenseMapRemoveSwapKeepsSurvivors(t *testing.T) {
	all := []string{"a", "b", "c", "d", "e"}

	tests := []struct {
		name   string
		remove []string
	}{
		{"first", []string{"a"}},
		{"middle", []string{"c"}},
		{"last", []string{"e"}},
		{"first then last", []string{"a", "e"}},
		{"last then first", []string{"e", "a"}},
		{"middle twice", []string{"c", "b"}},
		{"all but one", []string{"a", "b", "c", "d"}},
		{"everything", []string{"a", "b", "c", "d", "e"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := denseOf(t, all...)

			for _, k := range tt.remove {
				if !m.Remove(k) {
					t.Fatalf("Remove(%s) returned false", k)
				}
				checkDenseInvariant(t, m)
			}

			want := make([]string, 0, len(all))
			for _, k := range all {
				if !slices.Contains(tt.remove, k) {
					want = append(want, k)
				}
			}

			if got := m.Len(); got != len(want) {
				t.Fatalf("Len = %d, want %d", got, len(want))
			}
			if m.Empty() != (len(want) == 0) {
				t.Fatalf("Empty = %v, want %v", m.Empty(), len(want) == 0)
			}

			for _, k := range want {
				if !m.Has(k) {
					t.Fatalf("Has(%s) = false, want true", k)
				}
				// the value stored for a key is its index in all
				if v, ok := m.Get(k); !ok || v != slices.Index(all, k) {
					t.Fatalf("Get(%s) = %d, %v, want %d, true", k, v, ok, slices.Index(all, k))
				}
			}
			for _, k := range tt.remove {
				if m.Has(k) {
					t.Fatalf("Has(%s) = true after removal", k)
				}
				if _, ok := m.Get(k); ok {
					t.Fatalf("Get(%s) returned true after removal", k)
				}
			}

			got := denseKeys(m)
			slices.Sort(got)
			if !slices.Equal(got, want) {
				t.Fatalf("surviving keys = %v, want %v", got, want)
			}
		})
	}
}

func TestDenseMapSetAfterRemoveReusesTheHole(t *testing.T) {
	m := denseOf(t, "a", "b", "c")
	m.Remove("b")
	m.Set("d", 3)

	checkDenseInvariant(t, m)

	if got := m.Len(); got != 3 {
		t.Fatalf("Len = %d, want 3", got)
	}
	for k, want := range map[string]int{"a": 0, "c": 2, "d": 3} {
		if v, ok := m.Get(k); !ok || v != want {
			t.Fatalf("Get(%s) = %d, %v, want %d, true", k, v, ok, want)
		}
	}
}

func TestDenseMapHasLenEmpty(t *testing.T) {
	m := NewDenseMap[string, int]()

	if !m.Empty() || m.Len() != 0 {
		t.Fatalf("fresh map: Empty = %v, Len = %d", m.Empty(), m.Len())
	}
	if m.Has("a") {
		t.Fatal("Has on a fresh map returned true")
	}

	m.Set("a", 1)

	if m.Empty() || m.Len() != 1 {
		t.Fatalf("after Set: Empty = %v, Len = %d", m.Empty(), m.Len())
	}
	if !m.Has("a") {
		t.Fatal("Has(a) = false after Set")
	}
}

func TestDenseMapClearEmptiesAndStaysUsable(t *testing.T) {
	m := denseOf(t, "a", "b", "c")
	m.Clear()

	checkDenseInvariant(t, m)

	if !m.Empty() || m.Len() != 0 {
		t.Fatalf("after Clear: Empty = %v, Len = %d", m.Empty(), m.Len())
	}
	if m.Has("a") {
		t.Fatal("Has(a) = true after Clear")
	}
	if got := len(slices.Collect(m.IterKeys())); got != 0 {
		t.Fatalf("iterated %d keys after Clear, want 0", got)
	}

	// clearing twice is harmless, and the map still accepts new entries
	m.Clear()
	m.Set("z", 9)

	checkDenseInvariant(t, m)
	if v, ok := m.Get("z"); !ok || v != 9 {
		t.Fatalf("Get(z) = %d, %v, want 9, true", v, ok)
	}
	if got := m.Len(); got != 1 {
		t.Fatalf("Len = %d, want 1", got)
	}
}

// the doc comment promises dense order, so all three iterators and both
// snapshots must walk the same rows in the same order.
func TestDenseMapIteratorsAgreeOnDenseOrder(t *testing.T) {
	m := denseOf(t, "a", "b", "c", "d")
	m.Remove("b")
	m.Set("e", 4)
	m.Remove("a")

	keys := slices.Collect(m.IterKeys())
	values := slices.Collect(m.IterValues())

	if !slices.Equal(keys, m.Keys()) {
		t.Fatalf("IterKeys = %v, Keys = %v", keys, m.Keys())
	}
	if !slices.Equal(values, m.Values()) {
		t.Fatalf("IterValues = %v, Values = %v", values, m.Values())
	}
	if len(keys) != m.Len() || len(values) != m.Len() {
		t.Fatalf("iterated %d keys and %d values, Len = %d", len(keys), len(values), m.Len())
	}

	i := 0
	for k, v := range m.Iter() {
		if k != keys[i] || v != values[i] {
			t.Fatalf("Iter step %d = %s, %d, want %s, %d", i, k, v, keys[i], values[i])
		}
		i++
	}
	if i != m.Len() {
		t.Fatalf("Iter yielded %d pairs, want %d", i, m.Len())
	}

	// without removals dense order is insertion order
	fresh := denseOf(t, "x", "y", "z")
	if got := fresh.Keys(); !slices.Equal(got, []string{"x", "y", "z"}) {
		t.Fatalf("keys = %v, want [x y z]", got)
	}
}

func TestDenseMapKeysAndValuesAreSnapshots(t *testing.T) {
	m := denseOf(t, "a", "b")

	keys := m.Keys()
	values := m.Values()

	keys[0] = "mutated"
	values[0] = -1

	if !m.Has("a") {
		t.Fatal("mutating the Keys snapshot reached the map")
	}
	if v, _ := m.Get("a"); v != 0 {
		t.Fatalf("mutating the Values snapshot reached the map: Get(a) = %d", v)
	}
}

func TestDenseMapIterEarlyBreakDoesNotMutate(t *testing.T) {
	m := denseOf(t, "a", "b", "c")

	for range m.Iter() {
		break
	}
	for range m.IterKeys() {
		break
	}
	for range m.IterValues() {
		break
	}

	checkDenseInvariant(t, m)

	if got := m.Len(); got != 3 {
		t.Fatalf("Len = %d after broken iteration, want 3", got)
	}
	if got := m.Keys(); !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Fatalf("keys = %v after broken iteration, want [a b c]", got)
	}
}

func TestDenseMapZeroValueIsUsable(t *testing.T) {
	var m DenseMap[string, int]

	if !m.Empty() || m.Len() != 0 || m.Has("a") {
		t.Fatal("zero value did not start empty")
	}

	m.Set("a", 1)

	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Fatalf("Get(a) = %d, %v, want 1, true", v, ok)
	}
	checkDenseInvariant(t, &m)
}
