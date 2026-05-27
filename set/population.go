package set

import (
	"github.com/vistormu/go-dsa/queue"
)

// ======
// member
// ======

// identify an element using an index and a generation.
//
// the high 32 bits store the generation and the low 32 bits store the index.
//
// the tag parameter gives nominal typing so ids from different domains do not
// mix (for example, entity and asset).
type Member[Tag any] uint64

// pack an index and a generation into a member.
//
// time: O(1)
func newMember[Tag any](index, gen uint32) Member[Tag] {
	return Member[Tag](uint64(gen)<<32 | uint64(index))
}

// return the generation stored in the member.
//
// time: O(1)
func (m Member[Tag]) Gen() uint32 {
	return uint32(uint64(m) >> 32)
}

// return the index stored in the member.
//
// time: O(1)
func (m Member[Tag]) Index() uint32 {
	return uint32(uint64(m) & 0xffffffff)
}

// report whether this is the zero member.
//
// the zero member is reserved and is never spawned.
//
// time: O(1)
func (m Member[Tag]) Zero() bool {
	return m == 0
}

// ==========
// population
// ==========

// allocate and validate generational members.
//
// removed members become invalid and may be reused with a higher generation.
//
// member 0 is reserved and is never spawned. it does not live inside a
// population. this makes the zero value useful as an "unset" id.
//
// this type is not safe for concurrent use.
type Population[Tag any] struct {
	free  *queue.Queue[uint32]
	gens  []uint32
	alive uint32
}

// create an empty population.
//
// time: O(1)
func NewPopulation[Tag any]() *Population[Tag] {
	return &Population[Tag]{
		free:  queue.NewQueue[uint32](),
		gens:  make([]uint32, 0),
		alive: 0,
	}
}

// add a new member.
//
// time: O(1) amortised
func (p *Population[Tag]) Add() Member[Tag] {
	var index uint32

	if !p.free.Empty() {
		index, _ = p.free.Dequeue()
	} else {
		index = uint32(len(p.gens))
		p.gens = append(p.gens, 1) // start at 1 so member 0 is never produced
	}

	gen := p.gens[index]
	p.alive++

	return newMember[Tag](index, gen)
}

// remove a member.
//
// return false if the member is not alive.
//
// time: O(1)
func (p *Population[Tag]) Remove(m Member[Tag]) bool {
	index := m.Index()

	if index >= uint32(len(p.gens)) || p.gens[index] != m.Gen() {
		return false
	}

	p.gens[index]++
	p.free.Enqueue(index)

	p.alive--

	return true
}

// report whether a member is alive.
//
// time: O(1)
func (p *Population[Tag]) Has(m Member[Tag]) bool {
	index := m.Index()
	return index < uint32(len(p.gens)) && p.gens[index] == m.Gen()
}

// return the number of alive members.
//
// time: O(1)
func (p *Population[Tag]) Alive() uint32 {
	return p.alive
}

// return the number of allocated slots.
//
// time: O(1)
func (p *Population[Tag]) Capacity() uint32 {
	return uint32(len(p.gens))
}
