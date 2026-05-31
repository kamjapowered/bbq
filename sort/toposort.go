package sort

// describe a directed dependency from one node to another
//
// From must appear before To in a topological order
//
//microwave:export
type Edge[T comparable] struct {
	From T
	To   T
}

// return a topological order for nodes
//
// nodes must not contain duplicates
// edges whose endpoints are not present in nodes are ignored
// duplicate edges are collapsed
// self edges are ignored
// return false if the graph contains a cycle
//
// time: O(n + e)
//
//microwave:export
func TopoSort[T comparable](nodes []T, edges []Edge[T]) ([]T, bool) {
	n := len(nodes)

	idx := make(map[T]int, n)
	for i, node := range nodes {
		idx[node] = i
	}

	adj := make([][]int, n)
	indeg := make([]int, n)
	seen := make(map[[2]int]struct{}, len(edges))

	for _, edge := range edges {
		from, ok := idx[edge.From]
		if !ok {
			continue
		}

		to, ok := idx[edge.To]
		if !ok {
			continue
		}

		if from == to {
			continue
		}

		key := [2]int{from, to}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		adj[from] = append(adj[from], to)
		indeg[to]++
	}

	ready := make([]int, 0, n)
	for i := range n {
		if indeg[i] == 0 {
			ready = append(ready, i)
		}
	}

	order := make([]T, 0, n)
	for head := 0; head < len(ready); head++ {
		v := ready[head]
		order = append(order, nodes[v])

		for _, w := range adj[v] {
			indeg[w]--
			if indeg[w] == 0 {
				ready = append(ready, w)
			}
		}
	}

	if len(order) != n {
		return nil, false
	}

	return order, true
}
