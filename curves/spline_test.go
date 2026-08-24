package curves

import (
	stdmath "math"
	"testing"

	"kamjapowered.com/bbq/geometry"
	bmath "kamjapowered.com/bbq/math"
)

type vec = geometry.Vector[float64]

func vecEqual(a, b vec, e float64) bool {
	return bmath.AlmostEqual(a.X, b.X, e) &&
		bmath.AlmostEqual(a.Y, b.Y, e) &&
		bmath.AlmostEqual(a.Z, b.Z, e)
}

func TestLinearEndpoints(t *testing.T) {
	p0 := geometry.NewVector[float64](1, 2, 3)
	p1 := geometry.NewVector[float64](-4, 5, 0)
	l := NewLinear(p0, p1)

	if got := l.Compute(0); !vecEqual(got, p0, eps) {
		t.Errorf("Compute(0) = %v, want %v", got, p0)
	}
	if got := l.Compute(1); !vecEqual(got, p1, eps) {
		t.Errorf("Compute(1) = %v, want %v", got, p1)
	}
}

func TestLinearMatchesLerp(t *testing.T) {
	p0 := geometry.NewVector[float64](2, -1, 4)
	p1 := geometry.NewVector[float64](8, 3, -2)
	l := NewLinear(p0, p1)

	for _, u := range samples(50) {
		want := vec{
			X: p0.X + u*(p1.X-p0.X),
			Y: p0.Y + u*(p1.Y-p0.Y),
			Z: p0.Z + u*(p1.Z-p0.Z),
		}
		if got := l.Compute(u); !vecEqual(got, want, eps) {
			t.Fatalf("Compute(%v) = %v, want %v", u, got, want)
		}
	}
}

func TestBezierNoControlPoints(t *testing.T) {
	b := NewBezier[float64]()
	if got := b.Compute(0.5); got != (vec{}) {
		t.Errorf("Compute(0.5) = %v, want zero vector", got)
	}
}

func TestBezierSinglePointIsConstant(t *testing.T) {
	p := geometry.NewVector[float64](3, -7, 1)
	b := NewBezier(p)
	for _, u := range samples(20) {
		if got := b.Compute(u); !vecEqual(got, p, eps) {
			t.Fatalf("Compute(%v) = %v, want %v", u, got, p)
		}
	}
}

// a degree one bezier is exactly a linear segment
func TestBezierDegreeOneIsLinear(t *testing.T) {
	p0 := geometry.NewVector[float64](-2, 6, 3)
	p1 := geometry.NewVector[float64](5, -1, 0)
	b := NewBezier(p0, p1)
	l := NewLinear(p0, p1)

	for _, u := range samples(50) {
		if got, want := b.Compute(u), l.Compute(u); !vecEqual(got, want, eps) {
			t.Fatalf("Compute(%v) = %v, want %v", u, got, want)
		}
	}
}

// the curve must pass through the first and last control points
func TestBezierInterpolatesEnds(t *testing.T) {
	pts := []vec{
		geometry.NewVector[float64](0, 0),
		geometry.NewVector[float64](1, 4),
		geometry.NewVector[float64](3, -2),
		geometry.NewVector[float64](6, 1),
		geometry.NewVector[float64](7, 7),
		geometry.NewVector[float64](9, 0),
	}

	for n := 2; n <= len(pts); n++ {
		b := NewBezier(pts[:n]...)
		if got := b.Compute(0); !vecEqual(got, pts[0], eps) {
			t.Errorf("n = %d, Compute(0) = %v, want %v", n, got, pts[0])
		}
		if got := b.Compute(1); !vecEqual(got, pts[n-1], eps) {
			t.Errorf("n = %d, Compute(1) = %v, want %v", n, got, pts[n-1])
		}
	}
}

// the quadratic bernstein form, (1-u)^2 p0 + 2(1-u)u p1 + u^2 p2
func TestBezierQuadraticMatchesBernstein(t *testing.T) {
	p0 := geometry.NewVector[float64](0, 0)
	p1 := geometry.NewVector[float64](1, 2)
	p2 := geometry.NewVector[float64](2, 0)
	b := NewBezier(p0, p1, p2)

	for _, u := range samples(50) {
		v := 1 - u
		want := p0.Scale(v * v).Add(p1.Scale(2 * v * u)).Add(p2.Scale(u * u))
		if got := b.Compute(u); !vecEqual(got, want, eps) {
			t.Fatalf("Compute(%v) = %v, want %v", u, got, want)
		}
	}
}

// every point of the curve lies inside the bounding box of the control points
func TestBezierStaysInControlHull(t *testing.T) {
	pts := []vec{
		geometry.NewVector[float64](0, 0),
		geometry.NewVector[float64](1, 5),
		geometry.NewVector[float64](4, -3),
		geometry.NewVector[float64](6, 2),
	}
	b := NewBezier(pts...)

	minX, maxX := pts[0].X, pts[0].X
	minY, maxY := pts[0].Y, pts[0].Y
	for _, p := range pts[1:] {
		minX, maxX = stdmath.Min(minX, p.X), stdmath.Max(maxX, p.X)
		minY, maxY = stdmath.Min(minY, p.Y), stdmath.Max(maxY, p.Y)
	}

	for _, u := range samples(200) {
		p := b.Compute(u)
		if p.X < minX-eps || p.X > maxX+eps || p.Y < minY-eps || p.Y > maxY+eps {
			t.Fatalf("Compute(%v) = %v, outside hull x[%v, %v] y[%v, %v]", u, p, minX, maxX, minY, maxY)
		}
	}
}

// the constructor must copy its input so later edits cannot mutate the curve
func TestBezierCopiesInput(t *testing.T) {
	pts := []vec{
		geometry.NewVector[float64](0, 0),
		geometry.NewVector[float64](1, 1),
	}
	b := NewBezier(pts...)
	before := b.Compute(0.5)

	pts[1] = geometry.NewVector[float64](100, 100)
	if got := b.Compute(0.5); !vecEqual(got, before, eps) {
		t.Errorf("curve changed after mutating input, got %v, want %v", got, before)
	}
}

// repeated evaluation must not corrupt the stored control points
func TestBezierIsPure(t *testing.T) {
	b := NewBezier(
		geometry.NewVector[float64](0, 0),
		geometry.NewVector[float64](1, 3),
		geometry.NewVector[float64](2, 0),
	)

	first := b.Compute(0.3)
	for range 10 {
		b.Compute(0.7)
	}
	if got := b.Compute(0.3); !vecEqual(got, first, eps) {
		t.Errorf("Compute(0.3) = %v on re-evaluation, want %v", got, first)
	}
}

func TestCatmullRomTooFewPoints(t *testing.T) {
	pts := []vec{
		geometry.NewVector[float64](0, 0),
		geometry.NewVector[float64](1, 1),
		geometry.NewVector[float64](2, 0),
	}

	for n := 0; n < 4; n++ {
		s := NewCatmullRom(pts[:min(n, len(pts))]...)
		if got := s.Compute(0.5); got != (vec{}) {
			t.Errorf("n = %d, Compute(0.5) = %v, want zero vector", n, got)
		}
	}
}

// the spline must interpolate p1 through p_{n-2} at the segment knots
func TestCatmullRomInterpolatesKnots(t *testing.T) {
	pts := []vec{
		geometry.NewVector[float64](-1, 3),
		geometry.NewVector[float64](0, 0),
		geometry.NewVector[float64](1, 2),
		geometry.NewVector[float64](2, -1),
		geometry.NewVector[float64](3, 4),
		geometry.NewVector[float64](4, 0),
	}
	s := NewCatmullRom(pts...)

	segments := len(pts) - 3
	for k := 0; k <= segments; k++ {
		u := float64(k) / float64(segments)
		want := pts[k+1]
		if got := s.Compute(u); !vecEqual(got, want, eps) {
			t.Errorf("Compute(%v) = %v, want knot %d %v", u, got, k+1, want)
		}
	}
}

// evenly spaced collinear points give back the straight line
func TestCatmullRomCollinearIsLinear(t *testing.T) {
	pts := []vec{
		geometry.NewVector[float64](0, 0),
		geometry.NewVector[float64](1, 1),
		geometry.NewVector[float64](2, 2),
		geometry.NewVector[float64](3, 3),
	}
	s := NewCatmullRom(pts...)

	for _, u := range samples(50) {
		want := geometry.NewVector(1+u, 1+u)
		if got := s.Compute(u); !vecEqual(got, want, eps) {
			t.Fatalf("Compute(%v) = %v, want %v", u, got, want)
		}
	}
}

// the parameter must stay clamped for inputs at and beyond the ends
func TestCatmullRomClampsSegmentIndex(t *testing.T) {
	pts := []vec{
		geometry.NewVector[float64](0, 0),
		geometry.NewVector[float64](1, 1),
		geometry.NewVector[float64](2, 0),
		geometry.NewVector[float64](3, 1),
		geometry.NewVector[float64](4, 0),
	}
	s := NewCatmullRom(pts...)

	if got := s.Compute(1); !vecEqual(got, pts[len(pts)-2], eps) {
		t.Errorf("Compute(1) = %v, want %v", got, pts[len(pts)-2])
	}
	// just below the end must be close to the end
	if got := s.Compute(1 - 1e-9); !vecEqual(got, pts[len(pts)-2], 1e-6) {
		t.Errorf("Compute(1-) = %v, want ~%v", got, pts[len(pts)-2])
	}
}

// a parameter below the domain must clamp to the first segment instead of panicking
func TestCatmullRomNegativeParameter(t *testing.T) {
	pts := []vec{
		geometry.NewVector[float64](0, 0),
		geometry.NewVector[float64](1, 1),
		geometry.NewVector[float64](2, 2),
		geometry.NewVector[float64](3, 3),
	}
	s := NewCatmullRom(pts...)

	want := geometry.NewVector[float64](0.5, 0.5)
	if got := s.Compute(-0.5); !vecEqual(got, want, eps) {
		t.Errorf("Compute(-0.5) = %v, want %v", got, want)
	}
}

func TestCatmullRomCopiesInput(t *testing.T) {
	pts := []vec{
		geometry.NewVector[float64](0, 0),
		geometry.NewVector[float64](1, 1),
		geometry.NewVector[float64](2, 0),
		geometry.NewVector[float64](3, 1),
	}
	s := NewCatmullRom(pts...)
	before := s.Compute(0.5)

	pts[2] = geometry.NewVector[float64](50, 50)
	if got := s.Compute(0.5); !vecEqual(got, before, eps) {
		t.Errorf("spline changed after mutating input, got %v, want %v", got, before)
	}
}

func TestHermiteEndpoints(t *testing.T) {
	p0 := geometry.NewVector[float64](0, 1, 2)
	p1 := geometry.NewVector[float64](4, -3, 1)
	m0 := geometry.NewVector[float64](1, 1, 0)
	m1 := geometry.NewVector[float64](-2, 0, 3)
	h := NewHermite(p0, p1, m0, m1)

	if got := h.Compute(0); !vecEqual(got, p0, eps) {
		t.Errorf("Compute(0) = %v, want %v", got, p0)
	}
	if got := h.Compute(1); !vecEqual(got, p1, eps) {
		t.Errorf("Compute(1) = %v, want %v", got, p1)
	}
}

// the derivative at the ends must equal the given tangents
func TestHermiteTangents(t *testing.T) {
	p0 := geometry.NewVector[float64](0, 0)
	p1 := geometry.NewVector[float64](1, 2)
	m0 := geometry.NewVector[float64](3, -1)
	m1 := geometry.NewVector[float64](-1, 4)
	h := NewHermite(p0, p1, m0, m1)

	const d = 1e-4
	derivative := func(u float64) vec {
		return h.Compute(u + d).Sub(h.Compute(u - d)).Scale(1 / (2 * d))
	}

	if got := derivative(0); !vecEqual(got, m0, 1e-6) {
		t.Errorf("p'(0) = %v, want %v", got, m0)
	}
	if got := derivative(1); !vecEqual(got, m1, 1e-6) {
		t.Errorf("p'(1) = %v, want %v", got, m1)
	}
}

// tangents equal to p1 - p0 reduce the segment to a straight line
func TestHermiteWithChordTangentsIsLinear(t *testing.T) {
	p0 := geometry.NewVector[float64](-1, 2, 0)
	p1 := geometry.NewVector[float64](5, -4, 3)
	m := p1.Sub(p0)
	h := NewHermite(p0, p1, m, m)
	l := NewLinear(p0, p1)

	for _, u := range samples(50) {
		if got, want := h.Compute(u), l.Compute(u); !vecEqual(got, want, eps) {
			t.Fatalf("Compute(%v) = %v, want %v", u, got, want)
		}
	}
}

// zero tangents with equal endpoints give a constant, the basis is a partition of unity
func TestHermiteConstant(t *testing.T) {
	p := geometry.NewVector[float64](7, -2, 1)
	h := NewHermite(p, p, vec{}, vec{})

	for _, u := range samples(50) {
		if got := h.Compute(u); !vecEqual(got, p, eps) {
			t.Fatalf("Compute(%v) = %v, want %v", u, got, p)
		}
	}
}

// the vector curves are generic over float32 as well
func TestFloat32Splines(t *testing.T) {
	const e = float32(1e-5)
	p0 := geometry.NewVector[float32](0, 0)
	p1 := geometry.NewVector[float32](2, 4)

	l := NewLinear(p0, p1)
	if got := l.Compute(0.5); !bmath.AlmostEqual(got.X, 1, e) || !bmath.AlmostEqual(got.Y, 2, e) {
		t.Errorf("Linear.Compute(0.5) = %v, want (1, 2)", got)
	}

	b := NewBezier(p0, p1)
	if got := b.Compute(1); !bmath.AlmostEqual(got.X, p1.X, e) || !bmath.AlmostEqual(got.Y, p1.Y, e) {
		t.Errorf("Bezier.Compute(1) = %v, want %v", got, p1)
	}

	h := NewHermite(p0, p1, geometry.Vector[float32]{}, geometry.Vector[float32]{})
	if got := h.Compute(0); got != p0 {
		t.Errorf("Hermite.Compute(0) = %v, want %v", got, p0)
	}
}

func BenchmarkBezier(b *testing.B) {
	curve := NewBezier(
		geometry.NewVector[float64](0, 0),
		geometry.NewVector[float64](1, 3),
		geometry.NewVector[float64](2, -1),
		geometry.NewVector[float64](3, 0),
	)
	for b.Loop() {
		_ = curve.Compute(0.37)
	}
}
