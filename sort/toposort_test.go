package sort

import "testing"

func TestTopoSort(t *testing.T) {
	nodes := []string{"parse", "compile", "link", "test", "package"}
	edges := []Edge[string]{
		{From: "parse", To: "compile"},
		{From: "compile", To: "link"},
		{From: "compile", To: "test"},
		{From: "link", To: "package"},
		{From: "test", To: "package"},
	}

	order, ok := TopoSort(nodes, edges)
	if !ok {
		t.Fatalf("expected acyclic graph")
	}

	assertTopoOrder(t, order, edges)
}

func TestTopoSortPreservesInputOrderForReadyNodes(t *testing.T) {
	nodes := []string{"a", "b", "c", "d"}
	edges := []Edge[string]{
		{From: "b", To: "d"},
	}

	order, ok := TopoSort(nodes, edges)
	if !ok {
		t.Fatalf("expected acyclic graph")
	}

	expected := []string{"a", "b", "c", "d"}
	assertEqual(t, order, expected)
}

func TestTopoSortIgnoresDuplicateAndUnknownEdges(t *testing.T) {
	nodes := []string{"a", "b", "c"}
	edges := []Edge[string]{
		{From: "a", To: "b"},
		{From: "a", To: "b"},
		{From: "a", To: "missing"},
		{From: "missing", To: "b"},
		{From: "c", To: "c"},
	}

	order, ok := TopoSort(nodes, edges)
	if !ok {
		t.Fatalf("expected acyclic graph")
	}

	expected := []string{"a", "c", "b"}
	assertEqual(t, order, expected)
}

func TestTopoSortDetectsCycle(t *testing.T) {
	nodes := []string{"a", "b", "c"}
	edges := []Edge[string]{
		{From: "a", To: "b"},
		{From: "b", To: "c"},
		{From: "c", To: "a"},
	}

	order, ok := TopoSort(nodes, edges)
	if ok {
		t.Fatalf("expected cycle")
	}
	if order != nil {
		t.Fatalf("expected nil order on cycle, got %v", order)
	}
}

func assertTopoOrder[T comparable](t *testing.T, order []T, edges []Edge[T]) {
	t.Helper()

	pos := make(map[T]int, len(order))
	for i, v := range order {
		pos[v] = i
	}

	for _, edge := range edges {
		if pos[edge.From] >= pos[edge.To] {
			t.Fatalf("expected %v before %v in %v", edge.From, edge.To, order)
		}
	}
}

func assertEqual[T comparable](t *testing.T, got, expected []T) {
	t.Helper()

	if len(got) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}

	for i, v := range got {
		if v != expected[i] {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	}
}
