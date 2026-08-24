package set

import "testing"

func TestMemberPacksIndexAndGen(t *testing.T) {
	tests := []struct {
		name       string
		index, gen uint32
	}{
		{"zero index and gen", 0, 0},
		{"index only", 7, 0},
		{"gen only", 0, 3},
		{"both", 42, 9},
		{"max index", ^uint32(0), 1},
		{"max gen", 1, ^uint32(0)},
		{"both max", ^uint32(0), ^uint32(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMember[tag](tt.index, tt.gen)

			if got := m.Index(); got != tt.index {
				t.Fatalf("Index = %d, want %d", got, tt.index)
			}
			if got := m.Gen(); got != tt.gen {
				t.Fatalf("Gen = %d, want %d", got, tt.gen)
			}
		})
	}
}

func TestMemberZero(t *testing.T) {
	var m Member[tag]
	if !m.Zero() {
		t.Fatal("the zero value reports Zero = false")
	}
	if got := newMember[tag](0, 0); !got.Zero() {
		t.Fatal("member(0, 0) reports Zero = false")
	}

	// anything in either half makes it non zero
	for _, m := range []Member[tag]{newMember[tag](1, 0), newMember[tag](0, 1), newMember[tag](1, 1)} {
		if m.Zero() {
			t.Fatalf("member index=%d gen=%d reports Zero = true", m.Index(), m.Gen())
		}
	}
}

func TestMemberIdentityIsPerIndexAndGen(t *testing.T) {
	a := newMember[tag](5, 2)
	b := newMember[tag](5, 2)
	c := newMember[tag](5, 3)
	d := newMember[tag](6, 2)

	if a != b {
		t.Fatal("two members with the same index and gen differ")
	}
	if a == c {
		t.Fatal("a later generation of the same index compares equal")
	}
	if a == d {
		t.Fatal("the same generation of a different index compares equal")
	}
}
