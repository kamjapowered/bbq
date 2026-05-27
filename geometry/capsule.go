package geometry

import c "github.com/vistormu/go-dsa/constraints"

// store a capsule as a segment plus radius
type Capsule[T c.Number] struct {
	Segment Segment[T]
	Radius  T
}

// create a capsule
func NewCapsule[T c.Number](seg Segment[T], radius T) Capsule[T] {
	return Capsule[T]{Segment: seg, Radius: radius}
}
