package geometry

import c "github.com/kamjapowered/bbq/constraints"

// store an ellipse in local space using radii
//
// interpret as 2d on the xy plane
//
//microwave:export
type Ellipse[T c.Number] struct {
	RadiusX T
	RadiusY T
}

// create an ellipse from radii
//
//microwave:export
func NewEllipse[T c.Number](rx, ry T) Ellipse[T] {
	return Ellipse[T]{RadiusX: rx, RadiusY: ry}
}

// create an ellipse from height and width
//
//microwave:export
func NewEllipseHW[T c.Number](h, w T) Ellipse[T] {
	return Ellipse[T]{RadiusX: w / 2, RadiusY: h / 2}
}

// create a circle from radius
//
//microwave:export
func NewCircle[T c.Number](r T) Ellipse[T] {
	return Ellipse[T]{RadiusX: r, RadiusY: r}
}

// create a circle from diameter
//
//microwave:export
func NewCircleD[T c.Number](d T) Ellipse[T] {
	return Ellipse[T]{RadiusX: d / 2, RadiusY: d / 2}
}
