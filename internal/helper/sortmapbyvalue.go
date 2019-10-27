package helper

import "sort"

// SortMapByValue is a helper function which takes a
// map and outputs a sorted list of pairs, ordered by
// the value in the maps.
func SortMapByValue(mapCount map[string]int) PairList {
	pl := make(PairList, 0, len(mapCount))
	for k, v := range mapCount {
		pl = append(pl, Pair{k, v})
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

// Pair represents a single, sorted, item from a map
type Pair struct {
	Key   string
	Value int
}

// PairList represents a sorted map
type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool {
	if p[i].Value == p[j].Value {
		return p[i].Key < p[j].Key
	}
	return p[i].Value < p[j].Value
}
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
