package set

import "slices"

// track completion of a fixed set of required keys
//
// this type is not safe for concurrent use
//
//microwave:export
type Checklist[K comparable] struct {
	required *HashSet[K]
	checked  *HashSet[K]
}

// create a checklist with the given required keys
//
// time: O(n)
//
//microwave:export
func NewChecklist[K comparable](required ...K) *Checklist[K] {
	req := NewHashSet[K]()
	for _, k := range required {
		req.Add(k)
	}

	return &Checklist[K]{
		required: req,
		checked:  NewHashSet[K](),
	}
}

// mark a key as completed if it is required
//
// time: O(n)
func (c *Checklist[K]) Check(ks ...K) {
	for _, k := range ks {
		if c.required.Has(k) {
			c.checked.Add(k)
		}
	}
}

// mark a key as not completed
//
// time: O(n)
func (c *Checklist[K]) Uncheck(ks ...K) {
	for _, k := range ks {
		c.checked.Remove(k)
	}
}

// remove all completed marks
//
// time: O(n)
func (c *Checklist[K]) Clear() {
	c.checked.Clear()
}

// report whether all of the given keys are completed
//
// returns true if no keys are given
//
// time: O(n)
func (c *Checklist[K]) IsChecked(ks ...K) bool {
	return !slices.ContainsFunc(ks, func(k K) bool {
		return !c.checked.Has(k)
	})
}

// report whether at least one of the given keys is completed
//
// returns false if no keys are given
//
// time: O(n)
func (c *Checklist[K]) IsAnyChecked(ks ...K) bool {
	return slices.ContainsFunc(ks, c.checked.Has)
}

// report whether all of the given keys are not completed
//
// returns true if no keys are given
//
// time: O(n)
func (c *Checklist[K]) IsUnchecked(ks ...K) bool {
	return !slices.ContainsFunc(ks, c.checked.Has)
}

// report whether at least one of the given keys is not completed
//
// returns false if no keys are given
//
// time: O(n)
func (c *Checklist[K]) IsAnyUnchecked(ks ...K) bool {
	return slices.ContainsFunc(ks, func(k K) bool {
		return !c.checked.Has(k)
	})
}

// report whether all required keys are completed
//
// time: O(1)
func (c *Checklist[K]) IsAllChecked() bool {
	return c.checked.Len() == c.required.Len()
}

// report whether no required keys are completed
//
// time: O(1)
func (c *Checklist[K]) IsAllUnchecked() bool {
	return c.checked.Len() == 0
}

// return the number of required keys
//
// time: O(1)
func (c *Checklist[K]) Required() int {
	return c.required.Len()
}

// return the number of completed keys
//
// time: O(1)
func (c *Checklist[K]) Completed() int {
	return c.checked.Len()
}

// return all checked keys
//
// time: O(n)
func (c *Checklist[K]) Checked() []K {
	return c.checked.ToSlice()
}

// return all missing required keys
//
// time: O(n)
func (c *Checklist[K]) Unchecked() []K {
	return c.required.Difference(c.checked).ToSlice()
}

// return the slice of required checks
//
// time: O(n)
func (c *Checklist[K]) ToSlice() []K {
	return c.required.ToSlice()
}
