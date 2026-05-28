package math

import c "github.com/kamjapowered/bbq/constraints"

// map v from [inMin, inMax] to [outMin, outMax]
//
// return outMin if inMin equals inMax
//
// time: O(1)
//
//microwave:export
func Remap[T c.Float](v, inMin, inMax, outMin, outMax T) T {
	if inMin == inMax {
		return outMin
	}
	inRange := inMax - inMin
	outRange := outMax - outMin
	return (v-inMin)*outRange/inRange + outMin
}
