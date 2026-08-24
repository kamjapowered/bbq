package set

import (
	"slices"
	"testing"
)

func sortedKeys(ks []string) []string {
	out := slices.Clone(ks)
	slices.Sort(out)
	return out
}

func checkChecklist(t *testing.T, label string, c *Checklist[string], required, checked []string) {
	t.Helper()

	if got := c.Required(); got != len(required) {
		t.Fatalf("%s: Required = %d, want %d", label, got, len(required))
	}
	if got := c.Completed(); got != len(checked) {
		t.Fatalf("%s: Completed = %d, want %d", label, got, len(checked))
	}
	if got := sortedKeys(c.ToSlice()); !slices.Equal(got, sortedKeys(required)) {
		t.Fatalf("%s: ToSlice = %v, want %v", label, got, sortedKeys(required))
	}
	if got := sortedKeys(c.Checked()); !slices.Equal(got, sortedKeys(checked)) {
		t.Fatalf("%s: Checked = %v, want %v", label, got, sortedKeys(checked))
	}

	want := []string{}
	for _, k := range required {
		if !slices.Contains(checked, k) {
			want = append(want, k)
		}
	}
	if got := sortedKeys(c.Unchecked()); !slices.Equal(got, sortedKeys(want)) {
		t.Fatalf("%s: Unchecked = %v, want %v", label, got, sortedKeys(want))
	}
}

func TestChecklistNewStartsUnchecked(t *testing.T) {
	c := NewChecklist("a", "b", "c")

	checkChecklist(t, "a fresh checklist", c, []string{"a", "b", "c"}, nil)
	if c.IsAllChecked() {
		t.Fatal("IsAllChecked = true on a fresh checklist")
	}
	if !c.IsAllUnchecked() {
		t.Fatal("IsAllUnchecked = false on a fresh checklist")
	}
}

func TestChecklistDuplicateRequiredKeysCollapse(t *testing.T) {
	c := NewChecklist("a", "a", "b")

	checkChecklist(t, "duplicated required keys", c, []string{"a", "b"}, nil)
}

func TestChecklistCheckIsIdempotentAndVariadic(t *testing.T) {
	c := NewChecklist("a", "b", "c")

	c.Check("a")
	c.Check("a")
	checkChecklist(t, "after checking a twice", c, []string{"a", "b", "c"}, []string{"a"})

	c.Check("b", "c")
	checkChecklist(t, "after checking b and c", c, []string{"a", "b", "c"}, []string{"a", "b", "c"})
	if !c.IsAllChecked() {
		t.Fatal("IsAllChecked = false once every required key is checked")
	}
	if c.IsAllUnchecked() {
		t.Fatal("IsAllUnchecked = true once every required key is checked")
	}

	// no keys is a no op
	c.Check()
	checkChecklist(t, "after Check with no keys", c, []string{"a", "b", "c"}, []string{"a", "b", "c"})
}

func TestChecklistCheckIgnoresKeysThatAreNotRequired(t *testing.T) {
	c := NewChecklist("a", "b")

	c.Check("z")
	checkChecklist(t, "after checking an unrequired key", c, []string{"a", "b"}, nil)
	if c.IsChecked("z") {
		t.Fatal(`IsChecked("z") = true for a key that is not required`)
	}
	if c.IsAllChecked() {
		t.Fatal("IsAllChecked = true after checking only an unrequired key")
	}

	// a mixed call still records the required half
	c.Check("a", "z")
	checkChecklist(t, "after a mixed Check", c, []string{"a", "b"}, []string{"a"})
}

func TestChecklistUncheck(t *testing.T) {
	c := NewChecklist("a", "b", "c")
	c.Check("a", "b", "c")

	c.Uncheck("b")
	checkChecklist(t, "after unchecking b", c, []string{"a", "b", "c"}, []string{"a", "c"})
	if c.IsAllChecked() {
		t.Fatal("IsAllChecked = true with one key unchecked")
	}

	// unchecking a key that was never checked is a no op
	c.Uncheck("b")
	c.Uncheck("z")
	c.Uncheck()
	checkChecklist(t, "after redundant Uncheck calls", c, []string{"a", "b", "c"}, []string{"a", "c"})

	c.Uncheck("a", "c")
	checkChecklist(t, "after unchecking the rest", c, []string{"a", "b", "c"}, nil)
	if !c.IsAllUnchecked() {
		t.Fatal("IsAllUnchecked = false with nothing checked")
	}
}

func TestChecklistClear(t *testing.T) {
	c := NewChecklist("a", "b")
	c.Check("a", "b")

	c.Clear()
	checkChecklist(t, "after Clear", c, []string{"a", "b"}, nil)
	if !c.IsAllUnchecked() {
		t.Fatal("IsAllUnchecked = false after Clear")
	}

	// the required keys survive, so the checklist is still usable
	c.Check("a")
	checkChecklist(t, "after Clear then Check", c, []string{"a", "b"}, []string{"a"})
}

func TestChecklistPredicates(t *testing.T) {
	// a is checked, b and c are required but unchecked, z is not required
	newList := func() *Checklist[string] {
		c := NewChecklist("a", "b", "c")
		c.Check("a")
		return c
	}

	tests := []struct {
		name string
		fn   func(*Checklist[string], ...string) bool
		keys []string
		want bool
	}{
		{"IsChecked of a checked key", (*Checklist[string]).IsChecked, []string{"a"}, true},
		{"IsChecked of an unchecked key", (*Checklist[string]).IsChecked, []string{"b"}, false},
		{"IsChecked of a mixed pair", (*Checklist[string]).IsChecked, []string{"a", "b"}, false},
		{"IsChecked of an unrequired key", (*Checklist[string]).IsChecked, []string{"z"}, false},

		{"IsUnchecked of an unchecked key", (*Checklist[string]).IsUnchecked, []string{"b"}, true},
		{"IsUnchecked of two unchecked keys", (*Checklist[string]).IsUnchecked, []string{"b", "c"}, true},
		{"IsUnchecked of a checked key", (*Checklist[string]).IsUnchecked, []string{"a"}, false},
		{"IsUnchecked of a mixed pair", (*Checklist[string]).IsUnchecked, []string{"a", "b"}, false},
		{"IsUnchecked of an unrequired key", (*Checklist[string]).IsUnchecked, []string{"z"}, true},

		{"IsAnyChecked of a checked key", (*Checklist[string]).IsAnyChecked, []string{"a"}, true},
		{"IsAnyChecked of a mixed pair", (*Checklist[string]).IsAnyChecked, []string{"a", "b"}, true},
		{"IsAnyChecked of unchecked keys", (*Checklist[string]).IsAnyChecked, []string{"b", "c"}, false},

		{"IsAnyUnchecked of an unchecked key", (*Checklist[string]).IsAnyUnchecked, []string{"b"}, true},
		{"IsAnyUnchecked of a mixed pair", (*Checklist[string]).IsAnyUnchecked, []string{"a", "b"}, true},
		{"IsAnyUnchecked of a checked key", (*Checklist[string]).IsAnyUnchecked, []string{"a"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newList()
			if got := tt.fn(c, tt.keys...); got != tt.want {
				t.Fatalf("%v = %v, want %v", tt.keys, got, tt.want)
			}
			// a predicate must not mutate
			checkChecklist(t, "after the predicate", c, []string{"a", "b", "c"}, []string{"a"})
		})
	}
}

// the documented answers when no keys are given at all
func TestChecklistPredicatesWithNoKeys(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*Checklist[string], ...string) bool
		want bool
	}{
		{"IsChecked", (*Checklist[string]).IsChecked, true},
		{"IsUnchecked", (*Checklist[string]).IsUnchecked, true},
		{"IsAnyChecked", (*Checklist[string]).IsAnyChecked, false},
		{"IsAnyUnchecked", (*Checklist[string]).IsAnyUnchecked, false},
	}

	// the answer is the same whatever the checklist holds
	lists := map[string]*Checklist[string]{
		"empty checklist":    NewChecklist[string](),
		"nothing checked":    NewChecklist("a", "b"),
		"partially checked":  NewChecklist("a", "b"),
		"completely checked": NewChecklist("a", "b"),
	}
	lists["partially checked"].Check("a")
	lists["completely checked"].Check("a", "b")

	for listName, c := range lists {
		for _, tt := range tests {
			t.Run(listName+"/"+tt.name, func(t *testing.T) {
				if got := tt.fn(c); got != tt.want {
					t.Fatalf("%s() on a %s = %v, want %v", tt.name, listName, got, tt.want)
				}
			})
		}
	}
}

func TestChecklistWithNoRequiredKeys(t *testing.T) {
	c := NewChecklist[string]()

	checkChecklist(t, "a checklist with no required keys", c, nil, nil)
	if !c.IsAllChecked() {
		t.Fatal("IsAllChecked = false on a checklist with no required keys")
	}
	if !c.IsAllUnchecked() {
		t.Fatal("IsAllUnchecked = false on a checklist with no required keys")
	}

	// nothing can be checked, since nothing is required
	c.Check("a")
	checkChecklist(t, "after checking an unrequired key", c, nil, nil)
	if !c.IsAllChecked() {
		t.Fatal("IsAllChecked = false after a no op Check")
	}
}

func TestChecklistWithIntKeys(t *testing.T) {
	c := NewChecklist(1, 2, 3)

	c.Check(2, 4)
	if got := c.Completed(); got != 1 {
		t.Fatalf("Completed = %d, want 1", got)
	}
	if !c.IsChecked(2) {
		t.Fatal("IsChecked(2) = false, want true")
	}

	missing := c.Unchecked()
	slices.Sort(missing)
	if want := []int{1, 3}; !slices.Equal(missing, want) {
		t.Fatalf("Unchecked = %v, want %v", missing, want)
	}
}
