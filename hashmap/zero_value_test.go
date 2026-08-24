package hashmap

import "testing"

// DenseMap already allowed it, so BiHashMap and TypeMap do too
func TestZeroValuesAreUsable(t *testing.T) {
	t.Run("BiHashMap", func(t *testing.T) {
		var b BiHashMap[string, int]
		b.Set("a", 1)

		if v, ok := b.GetByKey("a"); !ok || v != 1 {
			t.Fatalf("GetByKey = %v, %v, want 1, true", v, ok)
		}
		if k, ok := b.GetByValue(1); !ok || k != "a" {
			t.Fatalf("GetByValue = %v, %v, want a, true", k, ok)
		}
		if b.Len() != 1 {
			t.Fatalf("Len = %d, want 1", b.Len())
		}
	})

	t.Run("TypeMap", func(t *testing.T) {
		var tm TypeMap
		if !tm.Put(42) {
			t.Fatal("Put returned false")
		}

		v, ok := tm.Get[int]()
		if !ok || *v != 42 {
			t.Fatalf("Get = %v, %v, want 42, true", v, ok)
		}
		if tm.Len() != 1 {
			t.Fatalf("Len = %d, want 1", tm.Len())
		}
	})
}
