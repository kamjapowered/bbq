package hashmap_test

import (
	"testing"

	"kamjapowered.com/bbq/hashmap"
	"kamjapowered.com/bbq/set"
)

// go 1.27 mangles a generic method's symbol wrongly when the type argument's
// type name ends with the method's own name, and emits a relocation to a name
// nothing defines. golang/go#80976. TypeMap's generic methods are therefore
// named so that no bbq container collides with them: a method named Set would
// break every caller storing a value that holds a HashSet.
//
// the failure is a link error, not a compile error or an assertion, so this
// test guards by being linked at all.
type holder struct {
	set  *set.HashSet[int]
	list *set.HashSet[string]
}

func TestPutAcceptsValuesHoldingGenericContainers(t *testing.T) {
	tm := hashmap.NewTypeMap()

	if !tm.Put(holder{set: set.NewHashSet[int]()}) {
		t.Fatal("Put returned false")
	}

	h, ok := tm.Get[holder]()
	if !ok {
		t.Fatal("Get returned false")
	}
	if h.set == nil {
		t.Fatal("stored value did not survive")
	}
}
