package strings

func TrigramJaccardDistance(source, target string) int {
	toSet := func(s string) map[string]struct{} {
		m := make(map[string]struct{}, max(1, len(s)-2))
		if len(s) < 3 {
			m[s] = struct{}{}
			return m
		}
		for i := 0; i < len(s)-2; i++ {
			m[s[i:i+3]] = struct{}{}
		}
		return m
	}
	sSet, tSet := toSet(source), toSet(target)

	intersection, union := 0, len(sSet)+len(tSet)
	for trig := range sSet {
		if _, ok := tSet[trig]; ok {
			intersection++
			union-- // was counted twice
		}
	}
	if union == 0 { // identical and shorter than 3 chars
		return 0
	}
	jaccard := float64(intersection) / float64(union)
	return int((1.0-jaccard)*1000 + 0.5)
}
