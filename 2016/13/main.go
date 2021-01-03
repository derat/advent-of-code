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

	// Part 2: Count locations reachable in at most 50 steps.
	// Sigh, just redo the search using BFS, I guess.
	// What was the point of optimizing the first part?
	seen := map[uint64]struct{}{key(sx, sy): struct{}{}}
	todo := map[uint64]struct{}{key(sx, sy): struct{}{}}
	for i := 0; i <= 50; i++ {
		newTodo := make(map[uint64]struct{})
		for k := range todo {
			seen[k] = struct{}{}
			x, y := unkey(k)
			for _, n := range [][2]int{{x - 1, y}, {x + 1, y}, {x, y - 1}, {x, y + 1}} {
				nx, ny := n[0], n[1]
				if nx < 0 || ny < 0 || wall(nx, ny) {
					continue
				}
				nk := key(nx, ny)
				if _, ok := seen[nk]; !ok {
					newTodo[nk] = struct{}{}
				}
			}
		}
		todo = newTodo
	}
	fmt.Println(len(seen))
}

type node struct {
	x, y, priority int
}

func key(x, y int) uint64 {
	return lib.PackInts([]int{x, y}, 32)
}

func unkey(k uint64) (x, y int) {
	v := lib.UnpackInts(k, 32, 2)
	return v[0], v[1]
}
