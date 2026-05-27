package strings

func JaroDistance(source, target string) int {
	if source == target {
		return 0
	}
	sLen, tLen := len(source), len(target)
	if sLen == 0 || tLen == 0 {
		return 1000
	}
	matchWindow := max(sLen, tLen)/2 - 1
	sMatches := make([]bool, sLen)
	tMatches := make([]bool, tLen)

	matches := 0
	for i := range sLen {
		start := max(0, i-matchWindow)
		end := min(i+matchWindow+1, tLen)
		for j := start; j < end; j++ {
			if tMatches[j] || source[i] != target[j] {
				continue
			}
			sMatches[i], tMatches[j] = true, true
			matches++
			break
		}
	}
	if matches == 0 {
		return 1000
	}
	// count transpositions
	transpositions := 0
	j := 0
	for i := range sLen {
		if !sMatches[i] {
			continue
		}
		for !tMatches[j] {
			j++
		}
		if source[i] != target[j] {
			transpositions++
		}
		j++
	}
	jaro := (float64(matches)/float64(sLen) +
		float64(matches)/float64(tLen) +
		(float64(matches)-float64(transpositions)/2.0)/float64(matches)) / 3.0
	return int((1.0-jaro)*1000 + 0.5)
}
