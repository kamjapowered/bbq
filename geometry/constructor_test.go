package geometry

import "testing"

func TestNewArrow(t *testing.T) {
	start, end := NewVector(0, 0, 0), NewVector(3, 4, 0)
	a := NewArrow(start, end, 2, 1)

	want := Arrow[int]{Start: start, End: end, HeadLength: 2, HeadWidth: 1}
	if a != want {
		t.Fatalf("NewArrow = %v, want %v", a, want)
	}
}

func TestNewCapsule(t *testing.T) {
	s := NewSegment(NewVector(0, 0, 0), NewVector(0, 10, 0))
	c := NewCapsule(s, 3)

	if c.Segment != s {
		t.Fatalf("Segment = %v, want %v", c.Segment, s)
	}
	if c.Radius != 3 {
		t.Fatalf("Radius = %d, want 3", c.Radius)
	}
}

func TestEllipse(t *testing.T) {
	tests := []struct {
		name string
		got  Ellipse[float64]
		want Ellipse[float64]
	}{
		{"radii", NewEllipse(3.0, 4.0), Ellipse[float64]{RadiusX: 3, RadiusY: 4}},
		// NewEllipseHW takes height first, and halves both sides
		{"height and width", NewEllipseHW(8.0, 6.0), Ellipse[float64]{RadiusX: 3, RadiusY: 4}},
		{"circle", NewCircle(5.0), Ellipse[float64]{RadiusX: 5, RadiusY: 5}},
		{"circle from diameter", NewCircleD(5.0), Ellipse[float64]{RadiusX: 2.5, RadiusY: 2.5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}
