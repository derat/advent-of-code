package main

import (
	"fmt"
	"math/bits"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const (
		sx = 1
		sy = 1
		tx = 31
		ty = 39
	)

	magic := lib.InputInts("2016/13")[0]
	wall := func(x, y int) bool {
		n := x*x + 3*x + 2*x*y + y + y*y + magic
		return bits.OnesCount64(uint64(n))%2 == 1
	}

	// https://www.redblobgames.com/pathfinding/a-star/introduction.html
	frontier := lib.NewHeap(func(a, b interface{}) bool { return a.(node).priority < b.(node).priority })
	frontier.Insert(node{sx, sy, 0})
	costs := map[uint64]int{key(sx, sy): 0}

	for frontier.Len() != 0 {
		cur := frontier.Pop().(node)
		cost := costs[key(cur.x, cur.y)]

		// Check if we're done.
		if cur.x == tx && cur.y == ty {
			fmt.Println(cost)
			break
		}

		for _, next := range []node{
			node{cur.x + 1, cur.y, 0},
			node{cur.x, cur.y + 1, 0},
			node{cur.x - 1, cur.y, 0},
			node{cur.x, cur.y - 1, 0},
		} {
			if next.x < 0 || next.y < 0 || wall(next.x, next.y) {
				continue
			}

			newCost := cost + 1
			k := key(next.x, next.y)
			if old, ok := costs[k]; ok && old <= newCost {
				continue // already visited with equal or lower cost
			}

			// Use the Manhattan distance to the target as a lower bound of the remaining cost.
			est := lib.Abs(tx-next.x) + lib.Abs(tx-next.y)
			next.priority = newCost + est
			frontier.Insert(next)
			costs[k] = newCost
		}
	}
}

type node struct {
	x, y, priority int
}

func key(x, y int) uint64 {
	return lib.PackInts([]int{x, y}, 32)
}
