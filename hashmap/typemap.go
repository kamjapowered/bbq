package hashmap

import (
	"iter"
	"reflect"
)

// store values keyed by base reflect type
//
// store values as pointers to their base type to allow in place mutation
//
// allow at most one value per base type
//
//microwave:export
type TypeMap struct {
	data map[reflect.Type]any
}

// create an empty typemap
//
//microwave:export
func NewTypeMap() *TypeMap {
	return &TypeMap{
		data: make(map[reflect.Type]any),
	}
}

// add stores v under the base type of T
//
// require T to be a non pointer type
//
// store the value as *T
//
// named Put rather than Set on purpose. go 1.27 mangles a generic method's
// symbol wrongly when the type argument's type name ends with the method's own
// name, and emits a relocation nothing defines, so a method named Set breaks
// every caller storing a value that holds a HashSet. golang/go#80976, guarded
// by TestPutAcceptsValuesHoldingGenericContainers
func (tm *TypeMap) Put[T any](v T) bool {
	rt := reflect.TypeFor[T]()
	if rt.Kind() == reflect.Pointer {
		return false
	}

	p := new(T)
	*p = v

	if tm.data == nil {
		tm.data = make(map[reflect.Type]any)
	}

	tm.data[rt] = p

	return true
}

// get retrieves the stored pointer for the base type T
//
// require T to be a non pointer type
//
// return false if the entry does not exist or the cast fails
func (tm *TypeMap) Get[T any]() (*T, bool) {
	rt := reflect.TypeFor[T]()
	if rt.Kind() == reflect.Pointer {
		return nil, false
	}

	v, ok := tm.data[rt]
	if !ok {
		return nil, false
	}

	out, ok := v.(*T)
	return out, ok
}

// remove deletes the entry for the base type T
//
// require T to be a non pointer type
func (tm *TypeMap) Remove[T any]() bool {
	rt := reflect.TypeFor[T]()
	if rt.Kind() == reflect.Pointer {
		return false
	}

	_, ok := tm.data[rt]
	if ok {
		delete(tm.data, rt)
	}

	return ok
}

// checks the existance of T in the TypeMap
//
// require T to be a non pointer type
func (tm *TypeMap) Has[T any]() bool {
	rt := reflect.TypeFor[T]()
	if rt.Kind() == reflect.Pointer {
		return false
	}

	_, ok := tm.data[rt]

	return ok
}

// return the number of entries
func (tm *TypeMap) Len() int {
	return len(tm.data)
}

// report whether no entries are stored
func (tm *TypeMap) Empty() bool {
	return tm.Len() == 0
}

// return a snapshot of all base type keys
func (tm *TypeMap) Keys() []reflect.Type {
	keys := make([]reflect.Type, 0, len(tm.data))
	for k := range tm.data {
		keys = append(keys, k)
	}
	return keys
}

// return a snapshot of all stored values
//
// each value is the stored pointer to its base type
func (tm *TypeMap) Values() []any {
	values := make([]any, 0, len(tm.data))
	for _, v := range tm.data {
		values = append(values, v)
	}
	return values
}

// iterate over a snapshot of the current contents
//
// each yielded value is the stored pointer to its base type
func (tm *TypeMap) Iter() iter.Seq2[reflect.Type, any] {
	return func(yield func(reflect.Type, any) bool) {
		type pair struct {
			t reflect.Type
			v any
		}

		pairs := make([]pair, 0, len(tm.data))
		for t, v := range tm.data {
			pairs = append(pairs, pair{t: t, v: v})
		}

		for _, p := range pairs {
			if !yield(p.t, p.v) {
				return
			}
		}
	}
}

// remove all entries while preserving capacity
func (tm *TypeMap) Clear() {
	clear(tm.data)
}
