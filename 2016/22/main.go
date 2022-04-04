package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var nodes [][]node                                 // indexes are x, y
	for _, ln := range lib.InputLines("2016/22")[2:] { // skip dumb header
		var x, y int
		var n node
		lib.Extract(ln, `^/dev/grid/node-x(\d+)-y(\d+)\s+(\d+)T\s+(\d+)T\s+(\d+)T\s+(\d+)%$`,
			&x, &y, &n.size, &n.used, &n.avail, &n.pct)
		for x >= len(nodes) {
			nodes = append(nodes, nil)
		}
		for y >= len(nodes[x]) {
			nodes[x] = append(nodes[x], node{})
		}
		nodes[x][y] = n
	}

	var cnt int
	for x, ns := range nodes {
		for y, a := range ns {
			for x2, ns2 := range nodes {
				for y2, b := range ns2 {
					if x == x2 && y == y2 {
						continue // don't compare nodes against themselves
					}
					if a.used == 0 {
						continue // skip empty source nodes
					}
					if a.used <= b.avail {
						cnt++
					}
				}
			}
		}
	}
	fmt.Println(cnt)

	// Part 2: Find minimum number of moves to get data from y=0 and max x to (0, 0).
	//
	// My input contains:
	// - 32 nodes with 501-510T size and 490-499T used
	// - 888 nodes with 85-94T size and 64-73T used
	// - 1 node with 93T size and 0T used
	//
	// We can only move entire chunks of data, so the big nodes seem useless --
	// they don't have enough available space to hold any other nodes' data,
	// and no nodes have enough space to hold their data. I think we can disregard
	// them entirely.
	//
	// Among the rest of the nodes, the smallest nonzero usage (64T) is more than
	// half the size of the largest capacity (94T), so we'll never be able to pack
	// two or more nodes' data into a single node.
	//
	// So I think the numbers are essentially irrelevant. We should ignore the
	// larger nodes, and among the 889 smaller nodes, we just have a single empty
	// space that we're pushing around to try to get the target data to (0, 0).
	//
	// This seems to make things easy! Our state can be represented by the location
	// of the single empty node and the location of the data from the target node.
	// This means that we'll eliminate a bunch of essentially-equivalent states.

	// Find the nonempty node with smallest amount of data.
	minUsed := math.MaxInt32
	for _, ns := range nodes {
		for _, n := range ns {
			if n.used > 0 {
				minUsed = lib.Min(minUsed, n.used)
			}
		}
	}
	lib.AssertLess(minUsed, math.MaxInt32)

	// Find nodes that are capable of holding that amount of data.
	var maxAvail int
	var init []state
	for x, ns := range nodes {
		for y, n := range ns {
			if n.avail >= minUsed && n.avail > maxAvail {
				maxAvail = n.avail
				init = append(init, state{x, y, len(nodes) - 1, 0}) // data in orig position
			}
		}
	}
	lib.AssertEq(len(init), 1) // only one space

	nx, ny := len(nodes), len(nodes[0]) // nodes in each direction

	steps := lib.AStar(init,
		func(s state) bool { return s.dx == 0 && s.dy == 0 },
		func(s state, m map[state]int) {
			for _, pos := range [][2]int{
				{s.sx - 1, s.sy},
				{s.sx + 1, s.sy},
				{s.sx, s.sy - 1},
				{s.sx, s.sy + 1},
			} {
				sx0, sy0 := pos[0], pos[1]
				// Ignore nodes whose data is perpetually stuck because it's larger
				// than the maximum amount that was initially available.
				if sx0 >= 0 && sx0 < nx && sy0 >= 0 && sy0 < ny && nodes[sx0][sy0].used <= maxAvail {
					// If the space moved to the position where the target data was located,
					// the target data is now in the space's old position.
					swapped := sx0 == s.dx && sy0 == s.dy
					dx0, dy0 := lib.If(swapped, s.sx, s.dx), lib.If(swapped, s.sy, s.dy)
					m[state{sx0, sy0, dx0, dy0}] = 1
				}
			}
		},
		func(s state) int {
			// Use the max of the data's Manhattan distance from the space and the space's distance
			// from (0, 0) as a lower bound of the required moves.
			return lib.Max(lib.Abs(s.dx-s.sx)+lib.Abs(s.dy-s.sy), s.sx+s.sy)
		})
	fmt.Println(steps)
}

type node struct {
	size, used, avail, pct int
}

type state struct {
	sx, sy int // location of empty space
	dx, dy int // location of data
}
