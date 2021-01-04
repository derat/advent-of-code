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
	frontier := lib.NewHeap(func(a, b interface{}) bool { return a.(item).pri < b.(item).pri })
	costs := make(map[uint64]int) // number of steps to get to state
	for x, ns := range nodes {
		for y, n := range ns {
			if n.avail >= minUsed && n.avail > maxAvail {
				maxAvail = n.avail
				st := pack(x, y, len(nodes)-1, 0) // data in orig position
				frontier.Insert(item{st, 0})
				costs[st] = 0
			}
		}
	}
	lib.AssertEq(frontier.Len(), 1) // only one space

	nx, ny := len(nodes), len(nodes[0]) // nodes in each direction

	for frontier.Len() != 0 {
		cur := frontier.Pop().(item)
		cost := costs[cur.state]
		sx, sy, dx, dy := unpack(cur.state)

		// Check if we're done, i.e. the data has reached (0, 0).
		if dx == 0 && dy == 0 {
			fmt.Println(cost)
			break
		}

		// Try moving the space.
		for _, pos := range [][2]int{
			{sx - 1, sy},
			{sx + 1, sy},
			{sx, sy - 1},
			{sx, sy + 1},
		} {
			sx0, sy0 := pos[0], pos[1]
			if sx0 < 0 || sx0 >= nx || sy0 < 0 || sy0 >= ny {
				continue
			}

			// Ignore nodes whose data is perpetually stuck because it's larger
			// than the maximum amount that was initially available.
			if nodes[sx0][sy0].used > maxAvail {
				continue
			}

			// If the space moved to the position where the target data was located,
			// the target data is now in the space's old position.
			swapped := sx0 == dx && sy0 == dy
			dx0, dy0 := lib.If(swapped, sx, dx), lib.If(swapped, sy, dy)

			p := pack(sx0, sy0, dx0, dy0)
			newCost := cost + 1
			if old, ok := costs[p]; ok && old <= newCost {
				continue // already in state with equal or lower cost
			}

			// Use the data's Manhattan distance from the space and then the space's
			// distance from (0, 0) as a lower bound of the required moves.
			est := lib.Abs(dx0-sx0) + lib.Abs(dy0-sy0) + sx0 + sy0
			frontier.Insert(item{p, newCost + est})
			costs[p] = newCost
		}
	}
}

// node holds information provided about a node.
type node struct {
	size, used, avail, pct int
}

// item describes an item in the priority queue.
type item struct {
	state uint64
	pri   int
}

// pack packs the location of the empty space and the data.
func pack(sx, sy, dx, dy int) uint64 {
	return lib.PackInts([]int{sx, sy, dx, dy}, 16)
}

// unpack unpacks the location of the empty space and the data.
func unpack(p uint64) (sx, sy, dx, dy int) {
	vals := lib.UnpackInts(p, 16, 4)
	return vals[0], vals[1], vals[2], vals[3]
}
