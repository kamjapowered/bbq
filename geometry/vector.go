package geometry

import (
	"math"

	c "github.com/vistormu/go-dsa/constraints"
)

// store a 3d vector where z can be left as zero for 2d usage
type Vector[T c.Number] struct {
	X, Y, Z T
}

// create a vector up until the specified dimension
func NewVector[T c.Number](v ...T) Vector[T] {
	vec := Vector[T]{}

	for i, value := range v {
		switch i {
		case 0:
			vec.X = value
		case 1:
			vec.Y = value
		case 2:
			vec.Z = value
		default:
			return vec
		}
	}

	return vec
}

// add v to the vector
func (a Vector[T]) Add(v Vector[T]) Vector[T] {
	return Vector[T]{X: a.X + v.X, Y: a.Y + v.Y, Z: a.Z + v.Z}
}

// subtract v from the vector
func (a Vector[T]) Sub(v Vector[T]) Vector[T] {
	return Vector[T]{X: a.X - v.X, Y: a.Y - v.Y, Z: a.Z - v.Z}
}

// scale the vector by s
func (a Vector[T]) Scale(s T) Vector[T] {
	return Vector[T]{X: a.X * s, Y: a.Y * s, Z: a.Z * s}
}

// compute the dot product with v
func (a Vector[T]) Dot(v Vector[T]) T {
	return a.X*v.X + a.Y*v.Y + a.Z*v.Z
}

// compute the cross product with v
func (a Vector[T]) Cross(v Vector[T]) Vector[T] {
	return Vector[T]{
		X: a.Y*v.Z - a.Z*v.Y,
		Y: a.Z*v.X - a.X*v.Z,
		Z: a.X*v.Y - a.Y*v.X,
	}
}

// compute squared length
func (a Vector[T]) LenSq() T {
	return a.Dot(a)
}

// compute length as float64
func (a Vector[T]) Len() float64 {
	return math.Sqrt(float64(a.LenSq()))
}

// compute a unit vector as float64 components
//
// return zero vector if length is zero
func (a Vector[T]) Norm() Vector[float64] {
	l := a.Len()
	if l == 0 {
		return Vector[float64]{}
	}
	inv := 1.0 / l
	return Vector[float64]{
		X: float64(a.X) * inv,
		Y: float64(a.Y) * inv,
		Z: float64(a.Z) * inv,
	}
}
