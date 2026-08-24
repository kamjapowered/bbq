package geometry

import (
	"slices"
	"testing"
)

func TestNewPath(t *testing.T) {
	p := NewPath[int]()

	if len(p.Points) != 0 {
		t.Fatalf("Points = %v, want empty", p.Points)
	}
	if p.Points == nil {
		t.Fatal("Points is nil, want an allocated slice")
	}
	if p.Closed {
		t.Fatal("Closed = true, want false")
	}
}

func TestPathAppend(t *testing.T) {
	p := NewPath[int]()
	p.Append(NewVector(1, 2, 3))
	p.AppendXY(4, 5)
	p.Append(Vector[int]{})

	want := []Vector[int]{
		{X: 1, Y: 2, Z: 3},
		// AppendXY leaves z at zero
		{X: 4, Y: 5, Z: 0},
		{},
	}

	if !slices.Equal(p.Points, want) {
		t.Fatalf("Points = %v, want %v", p.Points, want)
	}
}

func TestPathAppendOnZeroValue(t *testing.T) {
	// the zero path has a nil slice, and append must still grow it
	var p Path[float64]
	p.AppendXY(1, 2)

	if len(p.Points) != 1 || p.Points[0] != (Vector[float64]{X: 1, Y: 2}) {
		t.Fatalf("Points = %v, want one point at (1, 2)", p.Points)
	}

	// and clear must survive the nil slice too
	var q Path[float64]
	q.Clear()
	if len(q.Points) != 0 {
		t.Fatalf("Points = %v, want empty", q.Points)
	}
}

func TestPathClear(t *testing.T) {
	p := NewPath[int]()
	for i := range 4 {
		p.AppendXY(i, i)
	}
	p.Closed = true

	before := cap(p.Points)
	p.Clear()

	if len(p.Points) != 0 {
		t.Fatalf("Points = %v, want empty", p.Points)
	}
	// clear empties without giving back the allocation
	if got := cap(p.Points); got != before {
		t.Fatalf("capacity = %d, want %d", got, before)
	}
	// clear speaks for the points, not for the closed flag
	if !p.Closed {
		t.Fatal("Closed = false, want it left alone")
	}

	// the path is reusable afterwards
	p.AppendXY(9, 9)
	if len(p.Points) != 1 || p.Points[0] != (Vector[int]{X: 9, Y: 9}) {
		t.Fatalf("Points = %v, want one point at (9, 9)", p.Points)
	}
}

func TestPathIter(t *testing.T) {
	p := NewPath[int]()
	want := []Vector[int]{{X: 1}, {X: 2}, {X: 3}, {X: 4}}
	for _, v := range want {
		p.Append(v)
	}

	got := slices.Collect(p.Iter())
	if !slices.Equal(got, want) {
		t.Fatalf("Iter = %v, want %v", got, want)
	}

	// iterating does not consume
	if again := slices.Collect(p.Iter()); !slices.Equal(again, want) {
		t.Fatalf("second Iter = %v, want %v", again, want)
	}
}

func TestPathIterBreak(t *testing.T) {
	p := NewPath[int]()
	for i := range 5 {
		p.AppendXY(i, 0)
	}

	var got []Vector[int]
	for v := range p.Iter() {
		if len(got) == 2 {
			break
		}
		got = append(got, v)
	}

	want := []Vector[int]{{X: 0}, {X: 1}}
	if !slices.Equal(got, want) {
		t.Fatalf("Iter stopped with %v, want %v", got, want)
	}

	// breaking out leaves the path intact
	if len(p.Points) != 5 {
		t.Fatalf("Points = %v, want five points", p.Points)
	}
}

func TestPathIterEmpty(t *testing.T) {
	// an empty path yields nothing, whether it was never filled, was cleared,
	// or is the zero value
	empty := NewPath[int]()

	cleared := NewPath[int]()
	cleared.AppendXY(1, 1)
	cleared.Clear()

	var zero Path[int]

	for name, p := range map[string]*Path[int]{
		"new":     &empty,
		"cleared": &cleared,
		"zero":    &zero,
	} {
		t.Run(name, func(t *testing.T) {
			count := 0
			for range p.Iter() {
				count++
			}
			if count != 0 {
				t.Fatalf("Iter yielded %d points, want 0", count)
			}
		})
	}
}

func TestNewPolygon(t *testing.T) {
	p := NewPolygon[int]()

	if len(p.Points) != 0 {
		t.Fatalf("Points = %v, want empty", p.Points)
	}
	if p.Points == nil {
		t.Fatal("Points is nil, want an allocated slice")
	}
}

func TestPolygonAppendAndClear(t *testing.T) {
	p := NewPolygon[float64]()
	want := []Vector[float64]{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}}
	for _, v := range want {
		p.Append(v)
	}

	if !slices.Equal(p.Points, want) {
		t.Fatalf("Points = %v, want %v", p.Points, want)
	}

	before := cap(p.Points)
	p.Clear()

	if len(p.Points) != 0 {
		t.Fatalf("Points = %v, want empty", p.Points)
	}
	if got := cap(p.Points); got != before {
		t.Fatalf("capacity = %d, want %d", got, before)
	}

	// and it refills
	p.Append(Vector[float64]{X: 2, Y: 2})
	if len(p.Points) != 1 || p.Points[0] != (Vector[float64]{X: 2, Y: 2}) {
		t.Fatalf("Points = %v, want one point at (2, 2)", p.Points)
	}
}

func TestPolygonAppendOnZeroValue(t *testing.T) {
	var p Polygon[int]
	p.Append(NewVector(1, 2, 0))

	if len(p.Points) != 1 {
		t.Fatalf("Points = %v, want one point", p.Points)
	}

	p.Clear()
	if len(p.Points) != 0 {
		t.Fatalf("Points = %v, want empty", p.Points)
	}
}
