package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// First, read in all the clay regions so we can determine the region we care about.
	var clay []region
	bnd := region{math.MaxInt32, -1, math.MaxInt32, -1}
	for _, ln := range lib.InputLines("2018/17") {
		var reg region
		switch {
		case lib.ExtractMaybe(ln, `^y=(\d+), x=(\d+)\.\.(\d+)$`, &reg.r0, &reg.c0, &reg.c1):
			reg.r1 = reg.r0
		case lib.ExtractMaybe(ln, `^x=(\d+), y=(\d+)\.\.(\d+)$`, &reg.c0, &reg.r0, &reg.r1):
			reg.c1 = reg.c0
		default:
			lib.Panicf("Bad line %q", ln)
		}
		bnd.r0 = lib.Min(bnd.r0, reg.r0)
		bnd.r1 = lib.Max(bnd.r1, reg.r1)
		bnd.c0 = lib.Min(bnd.c0, reg.c0)
		bnd.c1 = lib.Max(bnd.c1, reg.c1)
		clay = append(clay, reg)
	}

	// Leave a blank space on the top and both sides.
	nrows := bnd.r1 - bnd.r0 + 2
	ncols := bnd.c1 - bnd.c0 + 3

	// Fill in the clay.
	grid := lib.NewByteGrid(nrows, ncols, '.')
	for _, reg := range clay {
		grid.SetRect(reg.r0-bnd.r0+1, reg.c0-bnd.c0+1, reg.r1-bnd.r0+1, reg.c1-bnd.c0+1, '#')
	}

	scol := 500 - bnd.c0 + 1 // 500 is given in puzzle
	grid[0][scol] = '+'
	drop(grid, 1, scol)
	//fmt.Println(grid.Dump())

	// Part 1: Number of tiles reachable by water.
	fmt.Println(grid.CountRect(1, 0, nrows-1, ncols-1, '~', '|'))

	// Part 2: Number of tiles containing stationary water.
	fmt.Println(grid.CountRect(1, 0, nrows-1, ncols-1, '~'))
}

// drop drops water in grid starting at r, c.
func drop(grid [][]byte, r, c int) {
	for r < len(grid)-1 {
		switch grid[r+1][c] {
		case '.', '|':
			// If there's sand or moving water below us, add moving water and move down.
			grid[r][c] = '|'
			r++
		case '#', '~':
			// When we hit clay or stationary water, search to both sides.
			lclay, lcol := look(grid, r, c, -1)
			rclay, rcol := look(grid, r, c, 1)

			if lclay && rclay {
				// If there was clay on both sides, fill the row with stationary
				// water and recurse from one row up to figure out what happens
				// on that row.
				for c0 := lcol + 1; c0 < rcol; c0++ {
					grid[r][c0] = '~'
				}
				drop(grid, r-1, c)
			} else {
				// If either side went over an edge, fill the row with running
				// water and recurse on those side(s).
				for c0 := lcol + 1; c0 < rcol; c0++ {
					grid[r][c0] = '|'
				}
				if !lclay {
					drop(grid, r, lcol)
				}
				if !rclay {
					drop(grid, r, rcol)
				}
			}
			return
		}
	}

	// Make the water flow down the (uncounted) bottom row.
	grid[r][c] = '|'
}

// look searches sideways in grid from r, c in the direction specified by dc.
// If it hits clay, it returns true and the clay's column.
// If it reaches a column unsupported by clay or stationary water,
// it returns false and the unsupported column.
func look(grid [][]byte, r, c, dc int) (clay bool, col int) {
	row := grid[r]
	for c += dc; true; c += dc {
		lib.Assertf(c >= 0 && c < len(row), "Bad column %d", c)
		if row[c] == '#' {
			return true, c
		} else if down := grid[r+1][c]; down != '#' && down != '~' {
			return false, c
		}
	}
	panic("Unreached")
}

type region struct{ r0, r1, c0, c1 int }
