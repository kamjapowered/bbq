package strings

func DamerauLevenshteinDistance(source, target string) int {
	ns, nt := len(source), len(target)
	const inf = int(^uint(0) >> 1) // max int

	// allocate +2 to simplify sentinel row / column
	d := make([][]int, ns+2)
	for i := range d {
		d[i] = make([]int, nt+2)
	}
	maxDist := ns + nt
	d[0][0] = maxDist
	for i := 0; i <= ns; i++ {
		d[i+1][0] = maxDist
		d[i+1][1] = i
	}
	for j := 0; j <= nt; j++ {
		d[0][j+1] = maxDist
		d[1][j+1] = j
	}

	lastRow := make(map[byte]int)
	for i := 1; i <= ns; i++ {
		db := 0
		for j := 1; j <= nt; j++ {
			k := lastRow[target[j-1]]
			l := db
			cost := 1
			if source[i-1] == target[j-1] {
				cost = 0
				db = j
			}
			d[i+1][j+1] = min(
				d[i][j]+cost,              // substitution
				d[i+1][j]+1,               // insertion
				d[i][j+1]+1,               // deletion
				d[k][l]+(i-k-1)+1+(j-l-1), // transposition
			)
		}
		lastRow[source[i-1]] = i
	}
	return d[ns+1][nt+1]
}
