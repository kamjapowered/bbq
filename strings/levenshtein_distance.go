package strings

func LevenshteinDistance(source, target string) int {
	ns, nt := len(source), len(target)

	previous := make([]int, nt+1)
	for j := 0; j <= nt; j++ {
		previous[j] = j
	}

	current := make([]int, nt+1)
	for i := 1; i <= ns; i++ {
		current[0] = i
		for j := 1; j <= nt; j++ {
			substitutionCost := 0
			if source[i-1] != target[j-1] {
				substitutionCost = 1
			}
			current[j] = min(
				previous[j]+1,                  // deletion
				current[j-1]+1,                 // insertion
				previous[j-1]+substitutionCost, // substitution
			)
		}
		previous, current = current, previous
	}
	return previous[nt]
}
