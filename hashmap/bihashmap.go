package hashmap

import "iter"

// store a one to one mapping between keys and values
//
// putting an existing key replaces its value
// putting an existing value replaces its key
//
// this type is not safe for concurrent use
//
//microwave:export
type BiHashMap[K, V comparable] struct {
	fwd map[K]V
	rev map[V]K
}

// create an empty bimap
//
// time: O(1)
//
//microwave:export
func NewBiHashMap[K, V comparable]() *BiHashMap[K, V] {
	return &BiHashMap[K, V]{
		fwd: make(map[K]V),
		rev: make(map[V]K),
	}
}

// set k -> v and v -> k
//
// return true if it inserted a new pair, false if it replaced something
//
// time: O(1) amortised
func (b *BiHashMap[K, V]) Set(k K, v V) bool {
	replaced := false

	if oldV, ok := b.fwd[k]; ok {
		replaced = true
		delete(b.rev, oldV)
	}

	if oldK, ok := b.rev[v]; ok {
		replaced = true
		delete(b.fwd, oldK)
	}

	if b.fwd == nil {
		b.fwd = make(map[K]V)
		b.rev = make(map[V]K)
	}

	b.fwd[k] = v
	b.rev[v] = k

	return !replaced
}

// get the value associated with k
//
// return false if k is not present
//
// time: O(1)
func (b *BiHashMap[K, V]) GetByKey(k K) (V, bool) {
	v, ok := b.fwd[k]
	return v, ok
}

// get the key associated with v
//
// return false if v is not present
//
// time: O(1)
func (b *BiHashMap[K, V]) GetByValue(v V) (K, bool) {
	k, ok := b.rev[v]
	return k, ok
}

// remove the pair associated with k
//
// return false if k is not present
//
// time: O(1)
func (b *BiHashMap[K, V]) RemoveByKey(k K) bool {
	v, ok := b.fwd[k]
	if !ok {
		return false
	}

	delete(b.fwd, k)
	delete(b.rev, v)
	return true
}

// remove the pair associated with v
//
// return false if v is not present
//
// time: O(1)
func (b *BiHashMap[K, V]) RemoveByValue(v V) bool {
	k, ok := b.rev[v]
	if !ok {
		return false
	}

	delete(b.rev, v)
	delete(b.fwd, k)
	return true
}

// report whether k is present
//
// time: O(1)
func (b *BiHashMap[K, V]) HasKey(k K) bool {
	_, ok := b.fwd[k]
	return ok
}

// report whether v is present
//
// time: O(1)
func (b *BiHashMap[K, V]) HasValue(v V) bool {
	_, ok := b.rev[v]
	return ok
}

// report whether no pairs are stored
//
// time: O(1)
func (b *BiHashMap[K, V]) Empty() bool {
	return len(b.fwd) == 0
}

// return the number of stored pairs
//
// time: O(1)
func (b *BiHashMap[K, V]) Len() int {
	return len(b.fwd)
}

// remove all pairs while preserving capacity
//
// time: O(n)
func (b *BiHashMap[K, V]) Clear() {
	for k, v := range b.fwd {
		delete(b.fwd, k)
		delete(b.rev, v)
	}
}

// return a copy of keys in unspecified order
//
// time: O(n)
func (b *BiHashMap[K, V]) Keys() []K {
	out := make([]K, 0, len(b.fwd))
	for k := range b.fwd {
		out = append(out, k)
	}
	return out
}

// return a copy of values in unspecified order
//
// time: O(n)
func (b *BiHashMap[K, V]) Values() []V {
	out := make([]V, 0, len(b.rev))
	for v := range b.rev {
		out = append(out, v)
	}
	return out
}

// iterate over stored key value pairs
//
// iteration stops if yield returns false
//
// time: O(n)
func (b *BiHashMap[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range b.fwd {
			if !yield(k, v) {
				return
			}
		}
	}
}

// iterate over pairs while removing them
//
// iteration stops if yield returns false
//
// time: O(n)
func (b *BiHashMap[K, V]) Drain() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range b.fwd {
			delete(b.fwd, k)
			delete(b.rev, v)
			if !yield(k, v) {
				return
			}
		}
	}
}
