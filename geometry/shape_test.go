package geometry

import "testing"

func TestNewLine(t *testing.T) {
	p0, p1 := NewVector(1, 2, 3), NewVector(4, 6, 3)
	l := NewLine(p0, p1)

	if l.Point != p0 {
		t.Fatalf("Point = %v, want %v", l.Point, p0)
	}
	if want := NewVector(3, 4, 0); l.Direction != want {
		t.Fatalf("Direction = %v, want %v", l.Direction, want)
	}
}

func TestLineAt(t *testing.T) {
	// the line through (1,2,3) and (4,6,3), so direction (3,4,0)
	l := NewLine(NewVector(1, 2, 3), NewVector(4, 6, 3))

	tests := []struct {
		name string
		t    int
		want Vector[int]
	}{
		{"t=0 is the point", 0, NewVector(1, 2, 3)},
		{"t=1 is the second point", 1, NewVector(4, 6, 3)},
		{"t=2 runs past it", 2, NewVector(7, 10, 3)},
		// a line extends both ways
		{"t=-1 runs behind it", -1, NewVector(-2, -2, 3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := l.At(tt.t); got != tt.want {
				t.Fatalf("At(%d) = %v, want %v", tt.t, got, tt.want)
			}
		})
	}
}

func TestLineAtDegenerate(t *testing.T) {
	// two identical points give a zero direction, so every t maps to the point
	p := NewVector(5, 5, 5)
	l := NewLine(p, p)

	for _, param := range []int{-2, 0, 1, 7} {
		if got := l.At(param); got != p {
			t.Fatalf("At(%d) = %v, want %v", param, got, p)
		}
	}
}

func TestNewRay(t *testing.T) {
	origin, through := NewVector(2, 0, 0), NewVector(2, 3, 0)
	r := NewRay(origin, through)

	if r.Origin != origin {
		t.Fatalf("Origin = %v, want %v", r.Origin, origin)
	}
	if want := NewVector(0, 3, 0); r.Direction != want {
		t.Fatalf("Direction = %v, want %v", r.Direction, want)
	}
}

func TestRayAt(t *testing.T) {
	r := NewRay(NewVector(2.0, 0.0, 0.0), NewVector(2.0, 3.0, 0.0))

	tests := []struct {
		name string
		t    float64
		want Vector[float64]
	}{
		{"t=0 is the origin", 0, NewVector(2.0, 0.0, 0.0)},
		{"t=1 is the through point", 1, NewVector(2.0, 3.0, 0.0)},
		{"t=0.5 is halfway", 0.5, NewVector(2.0, 1.5, 0.0)},
		{"t=3 runs past it", 3, NewVector(2.0, 9.0, 0.0)},
		// At does no clamping, so a negative t lands behind the origin
		{"t=-1 is behind the origin", -1, NewVector(2.0, -3.0, 0.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkVec(t, "At", r.At(tt.t), tt.want)
		})
	}
}

func TestRayAndLineAgree(t *testing.T) {
	// both are a point plus t times a direction, so they trace the same points
	p0, p1 := NewVector(1, -2, 4), NewVector(-3, 5, 4)
	l := NewLine(p0, p1)
	r := NewRay(p0, p1)

	for _, param := range []int{-3, 0, 1, 2, 10} {
		if got, want := r.At(param), l.At(param); got != want {
			t.Fatalf("Ray.At(%d) = %v, Line.At = %v", param, got, want)
		}
	}
}

func TestSegment(t *testing.T) {
	start, end := NewVector(1, 1, 0), NewVector(4, 5, 0)
	s := NewSegment(start, end)

	if s.Start != start || s.End != end {
		t.Fatalf("NewSegment = %v, want start %v end %v", s, start, end)
	}

	if got, want := s.Direction(), NewVector(3, 4, 0); got != want {
		t.Fatalf("Direction = %v, want %v", got, want)
	}

	// the direction spans the segment, so start plus it is end
	if got := s.Start.Add(s.Direction()); got != s.End {
		t.Fatalf("Start + Direction = %v, want %v", got, s.End)
	}

	// its length is the distance between the endpoints
	if got := s.Direction().Len(); got != 5 {
		t.Fatalf("Direction length = %f, want 5", got)
	}

	// reversing the segment negates the direction
	rev := NewSegment(end, start)
	if got, want := rev.Direction(), s.Direction().Scale(-1); got != want {
		t.Fatalf("reversed Direction = %v, want %v", got, want)
	}

	// a degenerate segment has no direction
	if got := NewSegment(start, start).Direction(); got != (Vector[int]{}) {
		t.Fatalf("degenerate Direction = %v, want zero", got)
	}
}

func TestRect(t *testing.T) {
	r := NewRect(8, 6)

	if r.Width != 8 || r.Height != 6 {
		t.Fatalf("NewRect = %v, want 8x6", r)
	}
	if got := r.HalfWidth(); got != 4 {
		t.Fatalf("HalfWidth = %d, want 4", got)
	}
	if got := r.HalfHeight(); got != 3 {
		t.Fatalf("HalfHeight = %d, want 3", got)
	}
}

func TestNewSquare(t *testing.T) {
	s := NewSquare(5.0)

	if s.Width != 5 || s.Height != 5 {
		t.Fatalf("NewSquare = %v, want 5x5", s)
	}
	if got := s.HalfWidth(); got != 2.5 {
		t.Fatalf("HalfWidth = %f, want 2.5", got)
	}
	if s.HalfWidth() != s.HalfHeight() {
		t.Fatalf("a square has unequal halves: %f and %f", s.HalfWidth(), s.HalfHeight())
	}
}

func TestRectHalvesTruncateOnIntegers(t *testing.T) {
	// integer division truncates, so an odd side loses its half unit
	r := NewRect(5, 7)

	if got := r.HalfWidth(); got != 2 {
		t.Fatalf("HalfWidth = %d, want 2", got)
	}
	if got := r.HalfHeight(); got != 3 {
		t.Fatalf("HalfHeight = %d, want 3", got)
	}

	// the same sides as floats keep the half unit
	f := NewRect(5.0, 7.0)
	if got := f.HalfWidth(); got != 2.5 {
		t.Fatalf("float HalfWidth = %f, want 2.5", got)
	}
	if got := f.HalfHeight(); got != 3.5 {
		t.Fatalf("float HalfHeight = %f, want 3.5", got)
	}
}

func TestRectResize(t *testing.T) {
	r := NewRect(8, 6)
	r.Resize(10, 2)

	if r.Width != 10 || r.Height != 2 {
		t.Fatalf("after Resize = %v, want 10x2", r)
	}
	if got := r.HalfWidth(); got != 5 {
		t.Fatalf("HalfWidth = %d, want 5", got)
	}
	if got := r.HalfHeight(); got != 1 {
		t.Fatalf("HalfHeight = %d, want 1", got)
	}

	// resizing to zero is allowed
	r.Resize(0, 0)
	if r != (Rect[int]{}) {
		t.Fatalf("after Resize(0, 0) = %v, want zero", r)
	}
}
