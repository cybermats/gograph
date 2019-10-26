package main

import "sort"

func rankByCount(mapCount map[string]int) PairList {
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

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
