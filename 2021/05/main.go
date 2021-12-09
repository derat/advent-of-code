package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const dim = 1000
	grid := lib.NewByteGrid(dim, dim, '.')
	grid2 := lib.NewByteGrid(dim, dim, '.')

	drawVent := func(b [][]byte, x1, y1, x2, y2 int) {
		grid.IterLine(y1, x1, y2, x2, func(r, c int) {
			if ch := b[r][c]; ch == '.' {
				b[r][c] = '1'
			} else {
				lib.Assert(ch >= '1' && ch <= '8')
				b[r][c] = ch + 1
			}
		})
	}

	for _, ln := range lib.InputLines("2021/5") {
		// Input is "x1,y1 -> x2,y2", i.e. col,row.
		var x1, y1, x2, y2 int
		lib.Extract(ln, `^(\d+),(\d+) -> (\d+),(\d+)$`, &x1, &y1, &x2, &y2)
		lib.AssertLess(x1, dim)
		lib.AssertLess(y1, dim)
		lib.AssertLess(x2, dim)
		lib.AssertLess(y2, dim)

		// Part 1: Only consider horizontal or vertical lines.
		if x1 == x2 || y1 == y2 {
			drawVent(grid, x1, y1, x2, y2)
		}
		// Part 2: Also consider diagonal lines.
		drawVent(grid2, x1, y1, x2, y2)
	}

	countOverlaps := func(b lib.ByteGrid) int {
		var cnt int
		b.IterRect(0, 0, dim-1, dim-1, func(r, c int) {
			if ch := b[r][c]; ch >= '2' && ch <= '9' {
				cnt++
			}
		})
		return cnt
	}
	fmt.Println(countOverlaps(grid))
	fmt.Println(countOverlaps(grid2))
}
