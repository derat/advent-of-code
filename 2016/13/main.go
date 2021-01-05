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

	ts := lib.PackInt2(tx, ty) // target state

	min := lib.AStar(
		[]uint64{lib.PackInt2(sx, sy)},
		func(s uint64) bool { return s == ts },
		func(s uint64) (ns []uint64) {
			x, y := lib.UnpackInt2(s)
			for _, n := range [][2]int{{x + 1, y}, {x, y + 1}, {x - 1, y}, {x, y - 1}} {
				if n[0] >= 0 && n[1] >= 0 && !wall(n[0], n[1]) {
					ns = append(ns, lib.PackInt2(n[0], n[1]))
				}
			}
			return ns
		},
		func(s uint64) int {
			x, y := lib.UnpackInt2(s)
			return lib.Abs(tx-x) + lib.Abs(tx-y)
		})
	fmt.Println(min)

	// Part 2: Count locations reachable in at most 50 steps.
	// Sigh, just redo the search using BFS, I guess.
	// What was the point of optimizing the first part?
	seen := map[uint64]struct{}{lib.PackInt2(sx, sy): struct{}{}}
	todo := map[uint64]struct{}{lib.PackInt2(sx, sy): struct{}{}}
	for i := 0; i <= 50; i++ {
		newTodo := make(map[uint64]struct{})
		for k := range todo {
			seen[k] = struct{}{}
			x, y := lib.UnpackInt2(k)
			for _, n := range [][2]int{{x - 1, y}, {x + 1, y}, {x, y - 1}, {x, y + 1}} {
				nx, ny := n[0], n[1]
				if nx < 0 || ny < 0 || wall(nx, ny) {
					continue
				}
				nk := lib.PackInt2(nx, ny)
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
