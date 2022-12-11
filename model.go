package prom_query_tool

import (
	"sort"
	"time"
)

type Point[V any] struct {
	Timestamp time.Time
	Value     V
}

type TupleList[L comparable, V any] struct {
	sortFunc func(l1, l2 L) bool
	List     []Tuple[L, V]
}

func (t *TupleList[L, V]) Len() int {
	return len(t.List)
}

func (t *TupleList[L, V]) Less(i, j int) bool {
	return t.sortFunc(t.List[i].Label, t.List[j].Label)
}

func (t *TupleList[L, V]) Swap(i, j int) {
	s := t.List[i]
	t.List[i] = t.List[j]
	t.List[j] = s
}

func (t *TupleList[L, V]) ListData(sortFuncs ...func(l1, l2 L) bool) []Tuple[L, V] {
	if len(sortFuncs) != 0 {
		t.sortFunc = sortFuncs[0]
		sort.Sort(t)
	}
	return t.List
}

type Tuple[L comparable, V any] struct {
	Label L
	Value V
}
