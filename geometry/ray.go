package geometry

import c "github.com/vistormu/go-dsa/constraints"

// store a ray as origin plus direction vector
type Ray[T c.Number] struct {
	Origin    Vector[T]
	Direction Vector[T]
}

// create a ray from two points
func NewRay[T c.Number](origin, through Vector[T]) Ray[T] {
	return Ray[T]{
		Origin:    origin,
		Direction: through.Sub(origin),
	}
}

// return a point at parameter t
func (r Ray[T]) At(t T) Vector[T] {
	return r.Origin.Add(r.Direction.Scale(t))
}
