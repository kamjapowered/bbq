package strings

import (
	"sort"
)

type DistanceFn func(query, candidate string) int

type Match struct {
	String string
	Score  int
	Index  int
}

func FuzzyFind(query string, haystack []string, distance DistanceFn, maxResults int) []Match {
	results := make([]Match, 0, len(haystack))
	for i, cand := range haystack {
		cost := distance(query, cand)
		results = append(results, Match{cand, cost, i})
	}

	sort.Slice(results, func(i, j int) bool { return results[i].Score < results[j].Score })

	if maxResults > 0 && len(results) > maxResults {
		results = results[:maxResults]
	}

	return results
}
