package geometry

import c "github.com/vistormu/go-dsa/constraints"

// store an ellipse in local space using radii
//
// interpret as 2d on the xy plane
type Ellipse[T c.Number] struct {
	RadiusX T
	RadiusY T
}

// create an ellipse from radii
func NewEllipse[T c.Number](rx, ry T) Ellipse[T] {
	return Ellipse[T]{RadiusX: rx, RadiusY: ry}
}

// create an ellipse from height and width
func NewEllipseHW[T c.Number](h, w T) Ellipse[T] {
	return Ellipse[T]{RadiusX: w / 2, RadiusY: h / 2}
}

// create a circle from radius
func NewCircle[T c.Number](r T) Ellipse[T] {
	return Ellipse[T]{RadiusX: r, RadiusY: r}
}

// create a circle from diameter
func NewCircleD[T c.Number](d T) Ellipse[T] {
	return Ellipse[T]{RadiusX: d / 2, RadiusY: d / 2}
}
