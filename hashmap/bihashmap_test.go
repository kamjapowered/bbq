package hashmap

import (
	"slices"
	"testing"
)

// the point of the type: every forward binding has the matching reverse
// binding and nothing else lingers on either side.
func checkBiInvariant[K, V comparable](t *testing.T, b *BiHashMap[K, V]) {
	t.Helper()

	if len(b.fwd) != len(b.rev) {
		t.Fatalf("fwd len = %d, rev len = %d", len(b.fwd), len(b.rev))
	}
	if b.Len() != len(b.fwd) {
		t.Fatalf("Len = %d, fwd len = %d", b.Len(), len(b.fwd))
	}
	for k, v := range b.fwd {
		back, ok := b.rev[v]
		if !ok {
			t.Fatalf("value %v has no reverse entry for key %v", v, k)
		}
		if back != k {
			t.Fatalf("rev[%v] = %v, want %v", v, back, k)
		}
	}
	for v, k := range b.rev {
		fwd, ok := b.fwd[k]
		if !ok {
			t.Fatalf("key %v has no forward entry for value %v", k, v)
		}
		if fwd != v {
			t.Fatalf("fwd[%v] = %v, want %v", k, fwd, v)
		}
	}
}

func biOf(t *testing.T, pairs ...[2]string) *BiHashMap[string, string] {
	t.Helper()

	b := NewBiHashMap[string, string]()
	for _, p := range pairs {
		b.Set(p[0], p[1])
	}
	return b
}

func TestBiHashMapSetReportsInsertion(t *testing.T) {
	b := NewBiHashMap[string, string]()

	if !b.Set("k1", "v1") {
		t.Fatal("Set of a new pair returned false")
	}
	if b.Set("k1", "v2") {
		t.Fatal("Set replacing a key returned true")
	}
	if b.Set("k2", "v2") {
		t.Fatal("Set replacing a value returned true")
	}
	if !b.Set("k3", "v3") {
		t.Fatal("Set of a second new pair returned false")
	}
	checkBiInvariant(t, b)
}

// each rebinding has to evict the pairing it displaces.
func TestBiHashMapSetRebindingEvictsOldPairings(t *testing.T) {
	tests := []struct {
		name   string
		start  [][2]string
		set    [2]string
		want   map[string]string
		goneK  []string
		goneV  []string
		wantLn int
	}{
		{
			name:   "existing key gets a new value",
			start:  [][2]string{{"k1", "v1"}, {"k2", "v2"}},
			set:    [2]string{"k1", "v3"},
			want:   map[string]string{"k1": "v3", "k2": "v2"},
			goneV:  []string{"v1"},
			wantLn: 2,
		},
		{
			name:   "existing value gets a new key",
			start:  [][2]string{{"k1", "v1"}, {"k2", "v2"}},
			set:    [2]string{"k3", "v1"},
			want:   map[string]string{"k3": "v1", "k2": "v2"},
			goneK:  []string{"k1"},
			wantLn: 2,
		},
		{
			name:   "both sides already bound elsewhere",
			start:  [][2]string{{"k1", "v1"}, {"k2", "v2"}},
			set:    [2]string{"k1", "v2"},
			want:   map[string]string{"k1": "v2"},
			goneK:  []string{"k2"},
			goneV:  []string{"v1"},
			wantLn: 1,
		},
		{
			name:   "setting the very same pair again",
			start:  [][2]string{{"k1", "v1"}, {"k2", "v2"}},
			set:    [2]string{"k1", "v1"},
			want:   map[string]string{"k1": "v1", "k2": "v2"},
			wantLn: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := biOf(t, tt.start...)
			b.Set(tt.set[0], tt.set[1])

			checkBiInvariant(t, b)

			if got := b.Len(); got != tt.wantLn {
				t.Fatalf("Len = %d, want %d", got, tt.wantLn)
			}
			for k, v := range tt.want {
				if got, ok := b.GetByKey(k); !ok || got != v {
					t.Fatalf("GetByKey(%s) = %s, %v, want %s, true", k, got, ok, v)
				}
				if got, ok := b.GetByValue(v); !ok || got != k {
					t.Fatalf("GetByValue(%s) = %s, %v, want %s, true", v, got, ok, k)
				}
			}
			for _, k := range tt.goneK {
				if b.HasKey(k) {
					t.Fatalf("key %s survived the rebinding", k)
				}
			}
			for _, v := range tt.goneV {
				if b.HasValue(v) {
					t.Fatalf("value %s survived the rebinding", v)
				}
			}
		})
	}
}

func TestBiHashMapGetAndHasOnAbsentEntries(t *testing.T) {
	b := biOf(t, [2]string{"k1", "v1"})

	if v, ok := b.GetByKey("nope"); ok || v != "" {
		t.Fatalf("GetByKey of an absent key = %q, %v, want \"\", false", v, ok)
	}
	if k, ok := b.GetByValue("nope"); ok || k != "" {
		t.Fatalf("GetByValue of an absent value = %q, %v, want \"\", false", k, ok)
	}
	if b.HasKey("nope") || b.HasValue("nope") {
		t.Fatal("Has reported an absent entry")
	}

	// the two axes are not interchangeable
	if b.HasKey("v1") {
		t.Fatal("HasKey matched a value")
	}
	if b.HasValue("k1") {
		t.Fatal("HasValue matched a key")
	}
	if !b.HasKey("k1") || !b.HasValue("v1") {
		t.Fatal("Has missed a present entry")
	}
}

func TestBiHashMapRemoveByKeyAndByValue(t *testing.T) {
	b := biOf(t, [2]string{"k1", "v1"}, [2]string{"k2", "v2"}, [2]string{"k3", "v3"})

	if !b.RemoveByKey("k1") {
		t.Fatal("RemoveByKey of a present key returned false")
	}
	checkBiInvariant(t, b)

	if b.HasKey("k1") || b.HasValue("v1") {
		t.Fatal("RemoveByKey left an orphan behind")
	}
	if b.RemoveByKey("k1") {
		t.Fatal("RemoveByKey of an already removed key returned true")
	}
	if b.RemoveByValue("v1") {
		t.Fatal("RemoveByValue of an already removed value returned true")
	}

	if !b.RemoveByValue("v2") {
		t.Fatal("RemoveByValue of a present value returned false")
	}
	checkBiInvariant(t, b)

	if b.HasKey("k2") || b.HasValue("v2") {
		t.Fatal("RemoveByValue left an orphan behind")
	}

	// removing by the wrong axis must not match
	if b.RemoveByKey("v3") {
		t.Fatal("RemoveByKey matched a value")
	}
	if b.RemoveByValue("k3") {
		t.Fatal("RemoveByValue matched a key")
	}

	if got := b.Len(); got != 1 {
		t.Fatalf("Len = %d, want 1", got)
	}
}

func TestBiHashMapLenEmptyAndClear(t *testing.T) {
	b := NewBiHashMap[string, string]()

	if !b.Empty() || b.Len() != 0 {
		t.Fatalf("fresh map: Empty = %v, Len = %d", b.Empty(), b.Len())
	}

	b.Set("k1", "v1")
	b.Set("k2", "v2")

	if b.Empty() || b.Len() != 2 {
		t.Fatalf("after Set: Empty = %v, Len = %d", b.Empty(), b.Len())
	}

	b.Clear()
	checkBiInvariant(t, b)

	if !b.Empty() || b.Len() != 0 {
		t.Fatalf("after Clear: Empty = %v, Len = %d", b.Empty(), b.Len())
	}
	if b.HasKey("k1") || b.HasValue("v1") {
		t.Fatal("Clear left an entry behind")
	}

	// still usable afterwards
	b.Set("k3", "v3")
	checkBiInvariant(t, b)
	if v, ok := b.GetByKey("k3"); !ok || v != "v3" {
		t.Fatalf("GetByKey(k3) = %s, %v, want v3, true", v, ok)
	}
}

func TestBiHashMapKeysValuesAndIter(t *testing.T) {
	b := biOf(t, [2]string{"k1", "v1"}, [2]string{"k2", "v2"}, [2]string{"k3", "v3"})

	keys := b.Keys()
	values := b.Values()
	slices.Sort(keys)
	slices.Sort(values)

	if !slices.Equal(keys, []string{"k1", "k2", "k3"}) {
		t.Fatalf("keys = %v, want [k1 k2 k3]", keys)
	}
	if !slices.Equal(values, []string{"v1", "v2", "v3"}) {
		t.Fatalf("values = %v, want [v1 v2 v3]", values)
	}

	seen := make(map[string]string)
	for k, v := range b.Iter() {
		if _, dup := seen[k]; dup {
			t.Fatalf("Iter yielded key %s twice", k)
		}
		seen[k] = v
	}
	if len(seen) != 3 {
		t.Fatalf("Iter yielded %d pairs, want 3", len(seen))
	}
	for k, v := range seen {
		if want, _ := b.GetByKey(k); want != v {
			t.Fatalf("Iter yielded %s -> %s, want %s", k, v, want)
		}
	}
}

func TestBiHashMapIterEarlyBreakDoesNotMutate(t *testing.T) {
	b := biOf(t, [2]string{"k1", "v1"}, [2]string{"k2", "v2"}, [2]string{"k3", "v3"})

	for range b.Iter() {
		break
	}

	checkBiInvariant(t, b)
	if got := b.Len(); got != 3 {
		t.Fatalf("Len = %d after broken iteration, want 3", got)
	}
}

func TestBiHashMapDrainEmptiesAsItGoes(t *testing.T) {
	b := biOf(t, [2]string{"k1", "v1"}, [2]string{"k2", "v2"}, [2]string{"k3", "v3"})

	drained := make(map[string]string)
	for k, v := range b.Drain() {
		drained[k] = v
		checkBiInvariant(t, b)
		if b.HasKey(k) || b.HasValue(v) {
			t.Fatalf("Drain yielded %s -> %s but the pair is still stored", k, v)
		}
	}

	if len(drained) != 3 {
		t.Fatalf("drained %d pairs, want 3", len(drained))
	}
	if !b.Empty() {
		t.Fatalf("Len = %d after a full drain, want 0", b.Len())
	}
}

func TestBiHashMapDrainEarlyBreakKeepsTheRest(t *testing.T) {
	b := biOf(t, [2]string{"k1", "v1"}, [2]string{"k2", "v2"}, [2]string{"k3", "v3"})

	var takenK, takenV string
	for k, v := range b.Drain() {
		takenK, takenV = k, v
		break
	}

	checkBiInvariant(t, b)

	if got := b.Len(); got != 2 {
		t.Fatalf("Len = %d after draining one pair, want 2", got)
	}
	if b.HasKey(takenK) || b.HasValue(takenV) {
		t.Fatalf("drained pair %s -> %s is still stored", takenK, takenV)
	}
}
