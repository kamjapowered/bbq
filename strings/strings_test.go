package strings

import (
	"testing"
)

func TestDistances(t *testing.T) {
	tests := []struct {
		name string
		fn   DistanceFn
	}{
		{"levenshtein", LevenshteinDistance},
		{"Damerau", DamerauLevenshteinDistance},
		{"Jaro", JaroDistance},
		{"JaroWinkler", JaroWinklerDistance},
		{"TrigramJaccard", TrigramJaccardDistance},
	}

	haystack := []string{"apple", "banana", "carrot"}
	query := "carote"

	for _, tt := range tests {
		matches := FuzzyFind(query, haystack, tt.fn, 3)
		if matches[0].String != "carrot" {
			t.Fatalf("%s failed\n", tt.name)
		}
	}
}

func BenchmarkDistances(b *testing.B) {
	tests := []struct {
		name string
		fn   DistanceFn
	}{
		{"levenshtein", LevenshteinDistance},
		{"hamming", HammingDistance},
		{"Damerau", DamerauLevenshteinDistance},
		{"Jaro", JaroDistance},
		{"JaroWinkler", JaroWinklerDistance},
		{"TrigramJaccard", TrigramJaccardDistance},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(sb *testing.B) {
			for b.Loop() {
				_ = tt.fn("encyclopaedia", "encyclopedia")
			}
		})
	}
}
