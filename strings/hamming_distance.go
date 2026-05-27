package strings

func HammingDistance(source, target string) int {
	if len(source) != len(target) {
		panic("HammingDistance: inputs must be equal length")
	}

	distance := 0
	for i := range source {
		if source[i] != target[i] {
			distance++
		}
	}
	return distance
}
