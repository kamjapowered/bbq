package geometry

import c "github.com/vistormu/go-dsa/constraints"

// store a polygon as points in local space
//
// interpret as 2d on the xy plane
type Polygon[T c.Number] struct {
	Points []Vector[T]
}

// create an empty polygon
func NewPolygon[T c.Number]() Polygon[T] {
	return Polygon[T]{Points: make([]Vector[T], 0)}
}

// add a point
func (p *Polygon[T]) Add(v Vector[T]) {
	p.Points = append(p.Points, v)
}

// clear all points
func (p *Polygon[T]) Clear() {
	p.Points = p.Points[:0]
}
