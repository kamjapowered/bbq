package set

import "testing"

type tag struct{}

func collectMembers(p *Population[tag]) []Member[tag] {
	var out []Member[tag]
	for m := range p.Iter() {
		out = append(out, m)
	}
	return out
}

func contains(ms []Member[tag], target Member[tag]) bool {
	for _, m := range ms {
		if m == target {
			return true
		}
	}
	return false
}

func TestPopulation_AddStartsAtGenOneAndAlive(t *testing.T) {
	p := NewPopulation[tag]()
	if p.Alive() != 0 || p.Capacity() != 0 {
		t.Fatal("expected empty population")
	}

	m := p.Add()
	if m.Zero() {
		t.Fatal("member 0 must never be produced")
	}
	if m.Index() != 0 {
		t.Fatalf("expected index 0, got %d", m.Index())
	}
	if m.Gen() != 1 {
		t.Fatalf("expected gen 1, got %d", m.Gen())
	}
	if !p.Has(m) {
		t.Fatal("expected member to be alive")
	}
	if p.Alive() != 1 || p.Capacity() != 1 {
		t.Fatalf("wrong counts: alive=%d capacity=%d", p.Alive(), p.Capacity())
	}
}

func TestPopulation_RemoveInvalidatesAndReuseIncrementsGen(t *testing.T) {
	p := NewPopulation[tag]()
	m1 := p.Add()

	if !p.Remove(m1) {
		t.Fatal("expected remove to succeed")
	}
	if p.Has(m1) {
		t.Fatal("expected stale member to be dead")
	}
	if p.Remove(m1) {
		t.Fatal("expected second remove to fail")
	}
	if p.Alive() != 0 {
		t.Fatalf("expected 0 alive, got %d", p.Alive())
	}

	m2 := p.Add()
	if m2.Index() != m1.Index() {
		t.Fatalf("expected index reuse: got %d want %d", m2.Index(), m1.Index())
	}
	if m2.Gen() != m1.Gen()+1 {
		t.Fatalf("expected gen increment: got %d want %d", m2.Gen(), m1.Gen()+1)
	}
	if p.Has(m1) {
		t.Fatal("old member must stay dead after slot reuse")
	}
	if !p.Has(m2) {
		t.Fatal("new member must be alive")
	}
}

func TestPopulation_IterYieldsOnlyAlive(t *testing.T) {
	p := NewPopulation[tag]()
	m0 := p.Add()
	m1 := p.Add()
	m2 := p.Add()

	if got := collectMembers(p); len(got) != 3 {
		t.Fatalf("expected 3 live members, got %d", len(got))
	}

	// remove the middle one
	p.Remove(m1)
	got := collectMembers(p)
	if len(got) != 2 {
		t.Fatalf("expected 2 live members, got %d", len(got))
	}
	if contains(got, m1) {
		t.Fatal("removed member must not be yielded")
	}
	if !contains(got, m0) || !contains(got, m2) {
		t.Fatal("expected m0 and m2 to be yielded")
	}

	// reuse the freed slot
	m1b := p.Add()
	if m1b.Index() != m1.Index() {
		t.Fatal("expected freed slot to be reused")
	}
	got = collectMembers(p)
	if len(got) != 3 {
		t.Fatalf("expected 3 live members after reuse, got %d", len(got))
	}
	if !contains(got, m1b) {
		t.Fatal("expected reused member to be yielded")
	}
	if contains(got, m1) {
		t.Fatal("stale member must not be yielded after reuse")
	}
}

func TestPopulation_IterEmpty(t *testing.T) {
	p := NewPopulation[tag]()
	if got := collectMembers(p); len(got) != 0 {
		t.Fatalf("expected no members, got %d", len(got))
	}
}

func TestPopulation_IterEarlyStop(t *testing.T) {
	p := NewPopulation[tag]()
	p.Add()
	p.Add()
	p.Add()

	count := 0
	for range p.Iter() {
		count++
		break
	}
	if count != 1 {
		t.Fatalf("expected early stop after 1, got %d", count)
	}
}
