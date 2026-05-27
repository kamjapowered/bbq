package set

// track completion of a fixed set of required keys
//
// this type is not safe for concurrent use
type Checklist[K comparable] struct {
	required *HashSet[K]
	checked  *HashSet[K]
}

// create a checklist with the given required keys
//
// time: O(n)
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
// time: O(1)
func (c *Checklist[K]) Check(k K) {
	if c.required.Contains(k) {
		c.checked.Add(k)
	}
}

// mark a key as not completed
//
// time: O(1)
func (c *Checklist[K]) Uncheck(k K) {
	c.checked.Remove(k)
}

// remove all completed marks
//
// time: O(n)
func (c *Checklist[K]) Clear() {
	c.checked.Clear()
}

// report whether a required key is completed
//
// time: O(1)
func (c *Checklist[K]) IsChecked(k K) bool {
	return c.checked.Contains(k)
}

// report whether all required keys are completed
//
// time: O(1)
func (c *Checklist[K]) All() bool {
	return c.checked.Len() == c.required.Len()
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

// return all missing required keys
//
// time: O(n)
func (c *Checklist[K]) Missing() []K {
	return c.required.Difference(c.checked).ToSlice()
}
