package hashmap

import (
	"reflect"
	"slices"
	"testing"
)

type config struct {
	name  string
	level int
}

type counter int

type tally int

func TestTypeMapPutAndGetAcrossTypes(t *testing.T) {
	tm := NewTypeMap()

	if !tm.Put(7) {
		t.Fatal("Put(int) returned false")
	}
	if !tm.Put("hello") {
		t.Fatal("Put(string) returned false")
	}
	if !tm.Put(config{name: "a", level: 3}) {
		t.Fatal("Put(config) returned false")
	}
	if !tm.Put(counter(1)) {
		t.Fatal("Put(counter) returned false")
	}
	if !tm.Put(tally(2)) {
		t.Fatal("Put(tally) returned false")
	}

	if got := tm.Len(); got != 5 {
		t.Fatalf("Len = %d, want 5", got)
	}

	if p, ok := tm.Get[int](); !ok || *p != 7 {
		t.Fatalf("Get[int] = %v, %v, want 7, true", p, ok)
	}
	if p, ok := tm.Get[string](); !ok || *p != "hello" {
		t.Fatalf("Get[string] = %v, %v, want hello, true", p, ok)
	}
	if p, ok := tm.Get[config](); !ok || p.name != "a" || p.level != 3 {
		t.Fatalf("Get[config] = %v, %v, want {a 3}, true", p, ok)
	}

	// two named types sharing an underlying type are distinct keys
	if p, ok := tm.Get[counter](); !ok || *p != 1 {
		t.Fatalf("Get[counter] = %v, %v, want 1, true", p, ok)
	}
	if p, ok := tm.Get[tally](); !ok || *p != 2 {
		t.Fatalf("Get[tally] = %v, %v, want 2, true", p, ok)
	}
}

func TestTypeMapGetMissingType(t *testing.T) {
	tm := NewTypeMap()
	tm.Put(1)

	p, ok := tm.Get[string]()
	if ok {
		t.Fatal("Get of a type never stored reported true")
	}
	if p != nil {
		t.Fatalf("Get of a type never stored = %v, want nil", p)
	}
}

func TestTypeMapPutOverwritesSameType(t *testing.T) {
	tm := NewTypeMap()

	tm.Put(config{name: "first", level: 1})
	old, _ := tm.Get[config]()

	if !tm.Put(config{name: "second", level: 2}) {
		t.Fatal("Put overwriting a type returned false")
	}
	if got := tm.Len(); got != 1 {
		t.Fatalf("Len = %d after overwriting, want 1", got)
	}

	p, ok := tm.Get[config]()
	if !ok || p.name != "second" || p.level != 2 {
		t.Fatalf("Get[config] = %v, %v, want {second 2}, true", p, ok)
	}
	// the stored pointer is fresh, so the previous one is not aliased
	if p == old {
		t.Fatal("overwriting reused the previous pointer")
	}
	if old.name != "first" {
		t.Fatalf("the previous value was mutated: %v", old)
	}
}

func TestTypeMapGetReturnsAMutablePointer(t *testing.T) {
	tm := NewTypeMap()
	tm.Put(config{name: "a", level: 1})

	p, ok := tm.Get[config]()
	if !ok {
		t.Fatal("Get[config] returned false")
	}
	p.level = 99

	again, ok := tm.Get[config]()
	if !ok || again.level != 99 {
		t.Fatalf("Get[config] after mutation = %v, %v, want level 99, true", again, ok)
	}
	if again != p {
		t.Fatal("Get returned a different pointer on the second call")
	}

	// Put copies its argument, so the caller's value is not aliased
	c := config{name: "b", level: 1}
	tm.Put(c)
	stored, _ := tm.Get[config]()
	stored.level = 5
	if c.level != 1 {
		t.Fatalf("Put aliased the caller's value: %v", c)
	}
}

// the doc requires a non pointer type argument, and every method must refuse
// one instead of storing it.
func TestTypeMapRejectsPointerTypes(t *testing.T) {
	tm := NewTypeMap()

	v := 3
	if tm.Put(&v) {
		t.Fatal("Put of a pointer type returned true")
	}
	if got := tm.Len(); got != 0 {
		t.Fatalf("Len = %d after a rejected Put, want 0", got)
	}

	if p, ok := tm.Get[*int](); ok || p != nil {
		t.Fatalf("Get of a pointer type = %v, %v, want nil, false", p, ok)
	}
	if tm.Has[*int]() {
		t.Fatal("Has of a pointer type returned true")
	}
	if tm.Remove[*int]() {
		t.Fatal("Remove of a pointer type returned true")
	}

	// the base type is unaffected by the rejection
	if !tm.Put(v) {
		t.Fatal("Put(int) returned false")
	}
	if p, ok := tm.Get[int](); !ok || *p != 3 {
		t.Fatalf("Get[int] = %v, %v, want 3, true", p, ok)
	}
	if tm.Has[*int]() {
		t.Fatal("Has of a pointer type matched the base type entry")
	}
}

func TestTypeMapHasAndRemove(t *testing.T) {
	tm := NewTypeMap()
	tm.Put(1)
	tm.Put("a")

	if !tm.Has[int]() || !tm.Has[string]() {
		t.Fatal("Has missed a stored type")
	}
	if tm.Has[config]() {
		t.Fatal("Has reported a type never stored")
	}

	if !tm.Remove[int]() {
		t.Fatal("Remove of a stored type returned false")
	}
	if tm.Has[int]() {
		t.Fatal("Has reported a removed type")
	}
	if _, ok := tm.Get[int](); ok {
		t.Fatal("Get reported a removed type")
	}
	if tm.Remove[int]() {
		t.Fatal("Remove of an already removed type returned true")
	}
	if tm.Remove[config]() {
		t.Fatal("Remove of a type never stored returned true")
	}

	if got := tm.Len(); got != 1 {
		t.Fatalf("Len = %d, want 1", got)
	}
	if !tm.Has[string]() {
		t.Fatal("Remove took out an unrelated type")
	}
}

func TestTypeMapLenEmptyAndClear(t *testing.T) {
	tm := NewTypeMap()

	if !tm.Empty() || tm.Len() != 0 {
		t.Fatalf("fresh map: Empty = %v, Len = %d", tm.Empty(), tm.Len())
	}

	tm.Put(1)
	tm.Put("a")

	if tm.Empty() || tm.Len() != 2 {
		t.Fatalf("after Put: Empty = %v, Len = %d", tm.Empty(), tm.Len())
	}

	tm.Clear()

	if !tm.Empty() || tm.Len() != 0 {
		t.Fatalf("after Clear: Empty = %v, Len = %d", tm.Empty(), tm.Len())
	}
	if tm.Has[int]() || tm.Has[string]() {
		t.Fatal("Clear left an entry behind")
	}
	if got := len(tm.Keys()); got != 0 {
		t.Fatalf("Keys returned %d entries after Clear, want 0", got)
	}

	// still usable afterwards
	tm.Put(config{name: "z"})
	if p, ok := tm.Get[config](); !ok || p.name != "z" {
		t.Fatalf("Get[config] = %v, %v, want {z 0}, true", p, ok)
	}
}

func TestTypeMapKeysValuesAndIter(t *testing.T) {
	tm := NewTypeMap()
	tm.Put(1)
	tm.Put("a")
	tm.Put(config{name: "c"})

	want := []string{"*hashmap.config", "*int", "*string"}

	keys := make([]string, 0, tm.Len())
	for _, k := range tm.Keys() {
		keys = append(keys, reflect.PointerTo(k).String())
	}
	slices.Sort(keys)
	if !slices.Equal(keys, want) {
		t.Fatalf("keys = %v, want %v", keys, want)
	}

	values := make([]string, 0, tm.Len())
	for _, v := range tm.Values() {
		values = append(values, reflect.TypeOf(v).String())
	}
	slices.Sort(values)
	if !slices.Equal(values, want) {
		t.Fatalf("values = %v, want %v", values, want)
	}

	// every yielded value is the stored pointer to its key type
	seen := make(map[reflect.Type]bool)
	for k, v := range tm.Iter() {
		if seen[k] {
			t.Fatalf("Iter yielded %s twice", k)
		}
		seen[k] = true

		if got := reflect.TypeOf(v); got != reflect.PointerTo(k) {
			t.Fatalf("Iter yielded a %s for key %s", got, k)
		}
	}
	if len(seen) != 3 {
		t.Fatalf("Iter yielded %d entries, want 3", len(seen))
	}
}

func TestTypeMapIterEarlyBreakDoesNotMutate(t *testing.T) {
	tm := NewTypeMap()
	tm.Put(1)
	tm.Put("a")
	tm.Put(config{})

	for range tm.Iter() {
		break
	}

	if got := tm.Len(); got != 3 {
		t.Fatalf("Len = %d after broken iteration, want 3", got)
	}
	if !tm.Has[int]() || !tm.Has[string]() || !tm.Has[config]() {
		t.Fatal("broken iteration dropped an entry")
	}
}

func TestTypeMapIterToleratesRemovalWhileIterating(t *testing.T) {
	tm := NewTypeMap()
	tm.Put(1)
	tm.Put("a")

	// Iter walks a snapshot, so removing during the walk is safe
	count := 0
	for range tm.Iter() {
		tm.Remove[int]()
		tm.Remove[string]()
		count++
	}

	if count != 2 {
		t.Fatalf("Iter yielded %d entries, want 2", count)
	}
	if !tm.Empty() {
		t.Fatalf("Len = %d, want 0", tm.Len())
	}
}
