# naming

the vocabulary every exported name in bbq is built from.

the shared verbs match `xpeto/docs/standards/naming.md`. bbq cannot reference it,
since the dependency runs the other way, so the two are kept in step by hand.

## the container quartet

every keyed or set-like container answers these four, and spells them the same way.

the quartet covers container membership only. `Vector.Add` is arithmetic and
`Logger.AddSink` registers a thing, and neither is bound by it.

### `Add`
- inserts a value into a collection
- old value wins on collision

### `Set`
- inserts a value into a collection
- new value wins on collision

### `Get`
- returns `(V, bool)` or `(*V, bool)` from a container
- non-mutating
- a verb rather than a noun, because a container does not know what it holds.
  where the receiver does know, the method takes the thing's name, as in
  `geometry.Ray.At`
- the `Ref` suffix marks the pointer form, and is warranted only where both forms
  exist. `DenseMap` carries `Get` and `GetRef`. `TypeMap` returns `*T` from `Get`
  alone, since it stores pointers so a caller can mutate in place

### `Remove`
- deletes a value from a collection

`Add` against `Set` turns on what happens to an existing entry, not on whether the
container is a map or a set.

### the exceptions to `Get`

two, both deliberate.

`Get` is a verb here where the rest of the library prefers the noun. xpeto's
`docs/standards/naming.md` says a method is named for the thing rather than for
the fetching, so a world exposes `Component[T]()` and not `GetComponent[T]()`.
that rule works because the receiver knows what it holds. a container does not,
so the only noun available is "value", and `m.Value(k)` reads as a property that
is simply there rather than as a lookup that can fail. the fallible shape needs
the verb. `geometry.Ray.At` keeps the noun, because a ray does know what lies
along it.

`TypeMap.Get` returns `*T` without the `Ref` suffix. the suffix distinguishes a
pointer form from a value form, and a `TypeMap` has no value form to distinguish
it from, since it stores pointers so a caller can mutate in place. adding `Ref`
there would name a difference that does not exist.

## the shared vocabulary

### `Has`
- checks existence of a value in a container
- returns `bool`
- never `Contains`

### `Len`
- number of elements, matching the builtin `len`
- never `Length`, `Size` or `Count`. the ban covers element counts only, so
  `terminal.Size`, which reports a window's columns and rows, is untouched

### `Empty`
- reports whether `Len` is zero
- returns `bool`, and does not mutate

### `Capacity`
- the bound of a container that has one, where `Len` is the current fill

### `Clear`
- returns to the empty state, keeping allocated capacity

### `Reset`
- returns to the initial state, which for a stateful computation is its starting
  conditions rather than emptiness

### `Iter`
- returns an `iter.Seq` or `iter.Seq2` over the elements
- non-mutating. `IterKeys` and `IterValues` narrow it when a map needs one axis
- a container whose iteration order is meaningful says so in its doc comment

### `Drain`
- returns an `iter.Seq` that empties the container as it goes

### `ToSlice`
- copies the elements into a fresh slice

### `Is*`
- a boolean predicate reporting a derived condition
- non-mutating
- the bare noun is a different method, not a duplicate. `Checklist.IsChecked`
  answers a question about given keys, `Checklist.Checked` returns the keys

### `New`
- constructs and returns a new value

## positional containers

a stack, queue, deque or list decides for itself which end an operation acts on,
so the verb names the operation and the type names the end. `Push` adds and `Pop`
removes at whichever end the type defines, the same end for a stack and opposite
ends for a queue.

- stack: `Push`, `Pop`, `Peek`
- queue: `Push`, `Pop`, `Peek`
- priority queue: `Push`, `Pop`, `Peek`
- deque: `PushFront`, `PushBack`, `PopFront`, `PopBack`, `PeekFront`, `PeekBack`
- list: `Append`, `Prepend`, `Insert`, `Pop`, `Remove`

`Enqueue` and `Dequeue` are not used. they name the container rather than the
operation, and they leave a priority queue with no spelling of its own.

`Len`, `Empty`, `Clear`, `Iter`, `Drain` and `ToSlice` follow the shared
vocabulary above.

## containers with more than one axis

a container addressable two ways suffixes the axis and keeps the verb.

- correct: `GetByKey`, `GetByValue`, `RemoveByKey`, `RemoveByValue`
- wrong: `DeleteByKey`, which invents a second verb for `Remove`

## type names

- a map-like type ends in `Map`, capital M: `DenseMap`, `TypeMap`, `BiHashMap`
- a set-like type ends in `Set`: `HashSet`, `BitSet`
- the package already carries the domain, so a type does not repeat it. inside
  `queue`, the type is `Queue` rather than `ArrayQueue`, unless a package holds
  several implementations that a caller picks between

## the minimum go version

go1.27.0 mangles a generic method's symbol when the type argument's type name
ends with the method's own name, so `TypeMap.Set` fails to link for any caller
storing a value that holds a `HashSet`. golang/go#80976, fixed in go1.27.1
(backport golang/go#80981). the method was named `Put` in v0.2.0 to dodge it and
is `Set` again from v0.3.0, which is why `go.mod` requires go 1.27.1.

`hashmap/typemap_test.go` guards it: on go1.27.0 the test binary fails to link.

## what does not follow this yet

nothing. the names above are what the library uses.
