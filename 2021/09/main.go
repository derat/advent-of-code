package main

import (
	"fmt"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	grid := lib.InputByteGrid("2021/9")

	// Part 1: Sum of of the risk levels of all low points.
	var sum int
	var lows [][2]int // [r, c]
	grid.Iter(func(r, c int) {
		if ch := grid[r][c]; ch < grid.Get(r-1, c, '9') && ch < grid.Get(r+1, c, '9') &&
			ch < grid.Get(r, c-1, '9') && ch < grid.Get(r, c+1, '9') {
			sum += 1 + int(ch-'0')
			lows = append(lows, [2]int{r, c})
		}
	})
	fmt.Println(sum)

	// Part 2: Product of the areas of the three largest basins.
	used := lib.NewByteGrid(grid.Rows(), grid.Cols(), 0)
	var fill func(r, c int) int
	fill = func(r, c int) int {
		area := 1
		used[r][c] = 1
		for _, p := range [][2]int{{r - 1, c}, {r + 1, c}, {r, c - 1}, {r, c + 1}} {
			pr, pc := p[0], p[1]
			if h := grid.Get(pr, pc, '9'); h != '9' && used.Get(pr, pc, 0) == 0 {
				// Recursion seems to work fine for my input, but another option would
				// be to keep a 'frontier' list of points and loop until it's empty.
				area += fill(pr, pc)
			}
		}
		return area
	}

	var areas []int
	for _, p := range lows {
		// I think that this 'used' check to make sure that the low point hasn't already been used
		// in a basin isn't necessary. The problem text says:
		//
		// "A basin is all locations that eventually flow downward to a single low point. Therefore,
		// every low point has a basin, although some basins are very small. Locations of height 9
		// do not count as being in any basin, and all other locations will always be part of
		// exactly one basin."
		//
		// In other words, I don't think that the input can contain a basin with multiple low
		// points:
		//
		//  99999
		//  92329
		//  99999
		//
		// Also note that input like this apparently doesnt't contain any low points (and therefore
		// also doesn't contain any basins):
		//
		//  9999
		//  9229
		//  9999
		if r, c := p[0], p[1]; used.Get(r, c, 0) == 0 {
			areas = append(areas, fill(r, c))
		}
	}
	lib.AssertLess(2, len(areas))
	sort.Sort(sort.Reverse(sort.IntSlice(areas)))
	fmt.Println(areas[0] * areas[1] * areas[2])
}
