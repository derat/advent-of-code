package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	grid := lib.InputByteGrid("2021/15")

	// Part 1: Find lowest cost from top left to bottom right.
	// Don't count cost of starting position unless you enter it.
	end := lib.PackInts(grid.MaxRow(), grid.MaxCol())
	cost := lib.AStar(
		[]interface{}{lib.PackInts(0, 0)},
		func(s interface{}) bool { return s.(uint64) == end },
		func(s interface{}, next map[interface{}]int) {
			r, c := lib.UnpackInt2(s.(uint64))
			for _, d := range [][2]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}} {
				nr, nc := r+d[0], c+d[1]
				if nr >= 0 && nr <= grid.MaxRow() && nc >= 0 && nc <= grid.MaxCol() {
					next[lib.PackInts(nr, nc)] = int(grid[nr][nc] - '0')
				}
			}
		},
		func(s interface{}) int {
			r, c := lib.UnpackInt2(s.(uint64))
			return (grid.MaxRow() - r) + (grid.MaxCol() - c)
		})
	fmt.Println(cost)

	// Part 2: "The entire cave is actually five times larger in both dimensions than you thought;
	// the area you originally scanned is just one tile in a 5x5 tile area that forms the full map."
	//
	// This was mostly just a copy-and-paste job. *shrug* I'm not sure if there's a smarter way to
	// handle this, but by initial hunch is that there isn't due to the way that the costs wrap
	// around when they go above 9.
	const repeats = 5
	maxr2, maxc2 := repeats*grid.Rows()-1, repeats*grid.Cols()-1
	end2 := lib.PackInts(maxr2, maxc2)
	cost2 := lib.AStar(
		[]interface{}{lib.PackInts(0, 0)},
		func(s interface{}) bool { return s.(uint64) == end2 },
		func(s interface{}, next map[interface{}]int) {
			r, c := lib.UnpackInt2(s.(uint64))
			for _, d := range [][2]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}} {
				nr, nc := r+d[0], c+d[1]
				if nr >= 0 && nr <= maxr2 && nc >= 0 && nc <= maxc2 {
					// "Your original map tile repeats to the right and downward; each time the tile
					// repeats to the right or downward, all of its risk levels are 1 higher than
					// the tile immediately up or left of it."
					cost := int(grid[nr%grid.Rows()][nc%grid.Cols()] - '0')
					cost += nr/grid.Rows() + nc/grid.Cols()
					// "However, risk levels above 9 wrap back around to 1."
					for cost > 9 {
						cost -= 9
					}
					next[lib.PackInts(nr, nc)] = cost
				}
			}
		},
		func(s interface{}) int {
			r, c := lib.UnpackInt2(s.(uint64))
			return (maxr2 - r) + (maxc2 - c)
		})
	fmt.Println(cost2)
}
