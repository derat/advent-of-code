package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

// Outer map is keyed by starting location ID.
// Inner map is from bitfield of visited location IDs to minimum cost.
type costMap map[uint8]map[uint8]int

func main() {
	ids := make(map[string]uint8) // from location to ID
	edges := make(map[uint8]int)  // from two OR-ed location IDs to cost
	var all uint8

	for _, ln := range lib.InputLines("2015/9") {
		var a, b string
		var cost int
		lib.Extract(ln, `^(\w+) to (\w+) = (\d+)$`, &a, &b, &cost)

		var aid, bid uint8
		if aid = ids[a]; aid == 0 {
			aid = 1 << len(ids)
			ids[a] = aid
			all |= aid
		}
		if bid = ids[b]; bid == 0 {
			bid = 1 << len(ids)
			ids[b] = bid
			all |= bid
		}
		edges[aid|bid] = cost
	}
	lib.AssertLessEq(len(ids), 8)

	costs := make(costMap)

	// Returns the minimum cost to travel from start to the locations in todo.
	var minCost func(start, todo uint8) int
	minCost = func(start, todo uint8) int {
		// Check for already-computed cost.
		if m, ok := costs[start]; ok {
			if c, ok := m[todo]; ok {
				return c
			}
		}

		// Consider all of start's edges.
		best := -1
		for i := 0; i < 8; i++ {
			var next uint8 = 1 << i
			if next&todo == 0 { // don't need to visit next
				continue
			}
			cost, ok := edges[start|next]
			if !ok { // no edge between start and next
				continue
			}
			newTodo := todo & ^next
			if newTodo != 0 {
				cost += minCost(next, newTodo)
			}
			if best == -1 || cost < best {
				best = cost
			}
		}

		// Cache the answer.
		m := costs[start]
		if m == nil {
			m = make(map[uint8]int)
			costs[start] = m
		}
		m[todo] = best

		return best
	}

	// Try starting from each location to find the minimum cost.
	min := -1
	for _, id := range ids {
		if c := minCost(id, all & ^id); min == -1 || c < min {
			min = c
		}
	}
	fmt.Println(min)
}
