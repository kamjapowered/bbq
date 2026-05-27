package hashmap

import (
	"iter"
	"reflect"
	"sync"
)

// store values keyed by base reflect type
//
// store values as pointers to their base type to allow in place mutation
//
// allow at most one value per base type
type TypeMap struct {
	mu   sync.RWMutex
	data map[reflect.Type]any
}

// create an empty typemap
func NewTypeMap() *TypeMap {
	return &TypeMap{
		data: make(map[reflect.Type]any),
	}
}

// return the number of entries
func (tm *TypeMap) Len() int {
	tm.mu.RLock()
	n := len(tm.data)
	tm.mu.RUnlock()
	return n
}

// report whether no entries are stored
func (tm *TypeMap) Empty() bool {
	return tm.Len() == 0
}

// return a snapshot of all base type keys
func (tm *TypeMap) Keys() []reflect.Type {
	tm.mu.RLock()
	keys := make([]reflect.Type, 0, len(tm.data))
	for k := range tm.data {
		keys = append(keys, k)
	}
	tm.mu.RUnlock()
	return keys
}

// return a snapshot of all stored values
//
// each value is the stored pointer to its base type
func (tm *TypeMap) Values() []any {
	tm.mu.RLock()
	values := make([]any, 0, len(tm.data))
	for _, v := range tm.data {
		values = append(values, v)
	}
	tm.mu.RUnlock()
	return values
}

// iterate over a snapshot of the current contents
//
// do not hold locks while yielding
//
// each yielded value is the stored pointer to its base type
func (tm *TypeMap) Iter() iter.Seq2[reflect.Type, any] {
	return func(yield func(reflect.Type, any) bool) {
		type pair struct {
			k reflect.Type
			v any
		}

		tm.mu.RLock()
		pairs := make([]pair, 0, len(tm.data))
		for k, v := range tm.data {
			pairs = append(pairs, pair{k: k, v: v})
		}
		tm.mu.RUnlock()

		for _, p := range pairs {
			if !yield(p.k, p.v) {
				return
			}
		}
	}
}

// remove all entries while preserving capacity
func (tm *TypeMap) Clear() {
	tm.mu.Lock()
	for k := range tm.data {
		delete(tm.data, k)
	}
	tm.mu.Unlock()
}

// =======
// helpers
// =======

func baseType(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}
	if t.Kind() == reflect.Pointer {
		return t.Elem()
	}
	return t
}

// =================
// generic functions
// =================

// add stores v under the base type of T
//
// require T to be a non pointer type
//
// store the value as *T
func Add[T any](tm *TypeMap, v T) bool {
	rt := reflect.TypeFor[T]()
	if rt.Kind() == reflect.Pointer {
		return false
	}

	p := new(T)
	*p = v

	tm.mu.Lock()
	tm.data[rt] = p
	tm.mu.Unlock()

	return true
}

// get retrieves the stored pointer for the base type T
//
// require T to be a non pointer type
//
// return false if the entry does not exist or the cast fails
func Get[T any](tm *TypeMap) (*T, bool) {
	rt := reflect.TypeFor[T]()
	if rt.Kind() == reflect.Pointer {
		return nil, false
	}

	tm.mu.RLock()
	v, ok := tm.data[rt]
	tm.mu.RUnlock()
	if !ok {
		return nil, false
	}

	out, ok := v.(*T)
	return out, ok
}

// remove deletes the entry for the base type T
//
// require T to be a non pointer type
func Remove[T any](tm *TypeMap) bool {
	rt := reflect.TypeFor[T]()
	if rt.Kind() == reflect.Pointer {
		return false
	}

	tm.mu.Lock()
	_, ok := tm.data[rt]
	if ok {
		delete(tm.data, rt)
	}
	tm.mu.Unlock()
	return ok
}

// checks the existance of T in the TypeMap
//
// require T to be a non pointer type
func Has[T any](tm *TypeMap) bool {
	rt := reflect.TypeFor[T]()
	if rt.Kind() == reflect.Pointer {
		return false
	}

	tm.mu.RLock()
	_, ok := tm.data[rt]
	tm.mu.RUnlock()

	return ok
}

// ====================
// reflection functions
// ====================

// addbytype stores v under the base type of t
//
// accept v as either base value or pointer to base value
func AddByType(tm *TypeMap, v any, t reflect.Type) bool {
	if tm == nil || v == nil || t == nil {
		return false
	}

	key := baseType(t)
	if key == nil {
		return false
	}

	rv := reflect.ValueOf(v)
	rt := rv.Type()

	ptrType := reflect.PointerTo(key)

	tm.mu.Lock()
	defer tm.mu.Unlock()

	// pointer input
	if rt == ptrType || rt.AssignableTo(ptrType) {
		if rt != ptrType && rv.CanConvert(ptrType) {
			rv = rv.Convert(ptrType)
		}
		tm.data[key] = rv.Interface()
		return true
	}

	// value input
	if rt == key || rt.AssignableTo(key) || rv.CanConvert(key) {
		if rt != key {
			if !rv.CanConvert(key) {
				return false
			}
			rv = rv.Convert(key)
		}
		ptr := reflect.New(key)
		ptr.Elem().Set(rv)
		tm.data[key] = ptr.Interface()
		return true
	}

	return false
}

// getbytype retrieves the stored pointer for the base type of t
func GetByType(tm *TypeMap, t reflect.Type) (any, bool) {
	if tm == nil || t == nil {
		return nil, false
	}

	key := baseType(t)

	tm.mu.RLock()
	raw, ok := tm.data[key]
	tm.mu.RUnlock()

	return raw, ok
}

// deletebytype removes the entry keyed by the base type of t
func DeleteByType(tm *TypeMap, t reflect.Type) bool {
	if tm == nil || t == nil {
		return false
	}

	key := baseType(t)

	tm.mu.Lock()
	_, ok := tm.data[key]
	if ok {
		delete(tm.data, key)
	}
	tm.mu.Unlock()

	return ok
}
