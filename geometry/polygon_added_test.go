package geometry

import "testing"

func TestPolygonAppendXY(t *testing.T) {
	var p Polygon[float64]
	p.AppendXY(1, 2)
	p.AppendXY(3, 4)

	want := []Vector[float64]{{X: 1, Y: 2}, {X: 3, Y: 4}}
	if len(p.Points) != len(want) {
		t.Fatalf("Points = %v, want %v", p.Points, want)
	}
	for i, v := range want {
		if p.Points[i] != v {
			t.Fatalf("Points[%d] = %v, want %v", i, p.Points[i], v)
		}
	}
}

func TestPolygonIter(t *testing.T) {
	var p Polygon[float64]

	// a zero value polygon yields nothing
	for range p.Iter() {
		t.Fatal("Iter yielded from an empty polygon")
	}

	p.AppendXY(1, 2)
	p.AppendXY(3, 4)
	p.AppendXY(5, 6)

	var got []Vector[float64]
	for v := range p.Iter() {
		got = append(got, v)
	}
	if len(got) != 3 || got[0].X != 1 || got[2].Y != 6 {
		t.Fatalf("Iter yielded %v", got)
	}

	// an early break stops iteration and does not mutate
	count := 0
	for range p.Iter() {
		count++
		break
	}
	if count != 1 || len(p.Points) != 3 {
		t.Fatalf("early break: count = %d, points = %d", count, len(p.Points))
	}
}
