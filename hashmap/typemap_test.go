package hashmap_test

import (
	"testing"

	"kamjapowered.com/bbq/hashmap"
	"kamjapowered.com/bbq/set"
)

// go1.27.0 mangles a generic method's symbol wrongly when the type argument's
// type name ends with the method's own name, and emits a relocation to a name
// nothing defines. golang/go#80976, fixed in go1.27.1 (backport
// golang/go#80981). TypeMap.Set storing a value that holds a HashSet is
// exactly that collision, so this guards the minimum go version in go.mod.
//
// the failure is a link error, not a compile error or an assertion, so this
// test guards by being linked at all.
type holder struct {
	set  *set.HashSet[int]
	list *set.HashSet[string]
}

func TestSetAcceptsValuesHoldingGenericContainers(t *testing.T) {
	tm := hashmap.NewTypeMap()

	if !tm.Set(holder{set: set.NewHashSet[int]()}) {
		t.Fatal("Set returned false")
	}

	h, ok := tm.Get[holder]()
	if !ok {
		t.Fatal("Get returned false")
	}
	if h.set == nil {
		t.Fatal("stored value did not survive")
	}
}
