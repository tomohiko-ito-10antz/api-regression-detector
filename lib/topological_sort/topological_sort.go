package topological_sort

type Set[T comparable] map[T]*any

func (s Set[T]) Push(v T) {
	s[v] = nil
}

func (s Set[T]) Pop() (T, bool) {
	var v T
	if len(s) == 0 {
		return v, false
	}

	for si := range s {
		v = si
		break
	}

	delete(s, v)

	return v, true
}

type Graph[T comparable] map[T]Set[T]

func NewGraph[T comparable]() Graph[T] {
	return Graph[T]{}
}

func (o Graph[T]) Arrow(from T, to T) {
	if o[from] == nil {
		o[from] = Set[T]{}
	}

	o[from].Push(to)
}

func Perform[T comparable](graph Graph[T]) (map[T]int, bool) {
	// count in-degrees
	inDegrees := map[T]int{}
	for u := range graph {
		inDegrees[u] = 0
	}
	for _, vs := range graph {
		for v := range vs {
			inDegrees[v]++
		}
	}

	// init queue
	q := Set[T]{}
	for u := range graph {
		if inDegrees[u] == 0 {
			q.Push(u)
		}
	}

	// topological sort
	order := map[T]int{}
	ord := 0
	for len(q) > 0 {
		u, _ := q.Pop()
		order[u] = ord
		ord++
		for v := range graph[u] {
			inDegrees[v]--
			if inDegrees[v] == 0 {
				q.Push(v)
			}
		}
	}

	if ord < len(graph) {
		return nil, false
	}

	return order, true
}
