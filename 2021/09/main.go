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
		h := grid[r][c]
		for _, p := range [][2]int{{r - 1, c}, {r + 1, c}, {r, c - 1}, {r, c + 1}} {
			pr, pc := p[0], p[1]
			if ph := grid.Get(pr, pc, '9'); ph >= h && ph != '9' && used.Get(pr, pc, 0) == 0 {
				used[pr][pc] = 1
				// Recursion seems to work fine for my input, but another option would
				// be to keep a 'frontier' list of points and loop until it's empty.
				area += fill(pr, pc)
			}
		}
		return area
	}

	var areas []int
	for _, p := range lows {
		r, c := p[0], p[1]
		areas = append(areas, fill(r, c))
	}
	lib.AssertLess(2, len(areas))
	sort.Sort(sort.Reverse(sort.IntSlice(areas)))
	fmt.Println(areas[0] * areas[1] * areas[2])
}
