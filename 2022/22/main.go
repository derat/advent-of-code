package main

import (
	"fmt"
	"regexp"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	pgs := lib.InputParagraphs("2022/22")
	lib.AssertEq(len(pgs), 2)

	// The input lines annoyingly aren't all the same width, so manually create a grid filled with
	// spaces and copy the input into it.
	var mc int
	for _, ln := range pgs[0] {
		mc = lib.Max(mc, len(ln))
	}
	grid := lib.NewByteGrid(len(pgs[0]), mc, ' ')
	for r, ln := range pgs[0] {
		copy(grid[r], []byte(ln))
	}

	var cmin, cmax []int // min and max columns for each row
	for r := range grid {
		min, max := 0, grid.MaxCol()
		for ; min < grid.Cols() && grid[r][min] == ' '; min++ {
		}
		for ; max >= 0 && grid[r][max] == ' '; max-- {
		}
		lib.AssertInRange(min, 0, grid.MaxCol())
		lib.AssertInRange(max, 0, grid.MaxCol())
		cmin = append(cmin, min)
		cmax = append(cmax, max)
	}

	var rmin, rmax []int // min and max rows for each column
	for c := 0; c < grid.Cols(); c++ {
		min, max := 0, grid.MaxRow()
		for ; min < grid.Rows() && grid[min][c] == ' '; min++ {
		}
		for ; max >= 0 && grid[max][c] == ' '; max-- {
		}
		lib.AssertInRange(min, 0, grid.MaxRow())
		lib.AssertInRange(max, 0, grid.MaxRow())
		rmin = append(rmin, min)
		rmax = append(rmax, max)
	}

	move := func(r, c int, dir lib.Dir, steps int) (int, int) {
		for ; steps > 0; steps-- {
			nr, nc := r+dir.DR(), c+dir.DC()
			if dir.DR() != 0 { // vertical
				if nr < rmin[nc] {
					nr = rmax[nc] // wrapped top
				} else if nr > rmax[nc] {
					nr = rmin[nc] // wrapped bottom
				}
			} else { // horizontal
				if nc < cmin[nr] {
					nc = cmax[nr] // wrapped left
				} else if nc > cmax[nr] {
					nc = cmin[nr] // wrapped bottom
				}
			}
			ch := grid[nr][nc]
			switch ch {
			case '#':
				break // stopped by wall
			case '.':
				r, c = nr, nc // empty space
			default:
				lib.Panicf("Moved to %q at row %d, col %d", ch, nr, nc)
			}
		}
		return r, c
	}

	// "You begin the path in the leftmost open tile of the top row of tiles.
	// Initially, you are facing to the right (from the perspective of how the map is drawn)."
	var r, c int
	dir := lib.Right
	for ; c < grid.Rows() && grid[r][c] != '.'; c++ {
	}
	lib.AssertLess(c, grid.Rows())

	lib.AssertEq(len(pgs[1]), 1)
	for _, cmd := range lib.Tokenize(pgs[1][0], regexp.MustCompile(`^\d+|[LR]`)) {
		switch cmd {
		case "L":
			dir = dir.Left()
		case "R":
			dir = dir.Right()
		default:
			steps := lib.ExtractInt(cmd)
			r, c = move(r, c, dir, steps)
		}
	}

	// Part 1: Generate a password from the final position.
	// "Rows start from 1 at the top and count downward; columns start from 1 at the left and count
	// rightward. ... The final password is the sum of 1000 times the row, 4 times the column, and
	// the facing."
	fmt.Println(1000*(r+1) + 4*(c+1) + dirValue[dir])
}

// "Facing is 0 for right (>), 1 for down (v), 2 for left (<), and 3 for up (^)."
var dirValue = map[lib.Dir]int{
	lib.Right: 0,
	lib.Down:  1,
	lib.Left:  2,
	lib.Up:    3,
}
