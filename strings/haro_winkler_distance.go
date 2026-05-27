package strings

func JaroWinklerDistance(source, target string) int {
	const prefixScale = 0.1 // Winkler’s constant
	jaroScaled := 1.0 - float64(JaroDistance(source, target))/1000.0
	prefix := 0
	for i := 0; i < min(4, min(len(source), len(target))); i++ {
		if source[i] == target[i] {
			prefix++
		} else {
			break
		}
	}
	jw := jaroScaled + float64(prefix)*prefixScale*(1.0-jaroScaled)
	return int((1.0-jw)*1000 + 0.5)
}
