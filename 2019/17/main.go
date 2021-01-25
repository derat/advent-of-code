package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInt64s("2019/17")

	vm := lib.NewIntcode(input)

	var inRow bool
	var grid [][]byte
	vm.OutFunc = func(v int64) {
		switch v {
		case 10: // newline
			inRow = false
		default:
			if !inRow {
				grid = append(grid, []byte{})
				inRow = true
			}
			r := len(grid) - 1
			grid[r] = append(grid[r], byte(v))
		}
	}
	lib.Assert(vm.Run())
	//fmt.Println(lib.DumpBytes(grid))

	nrows := len(grid)
	ncols := len(grid[0])
	for _, r := range grid {
		lib.AssertEq(len(r), ncols)
	}

	// Part 1: Print sum of r*c of intersections.
	var sum int
	for r := 1; r < nrows-1; r++ {
		for c := 1; c < ncols-1; c++ {
			if grid[r][c] == '#' && grid[r-1][c] == '#' && grid[r+1][c] == '#' &&
				grid[r][c-1] == '#' && grid[r][c+1] == '#' {
				sum += r * c
			}
		}
	}
	fmt.Println(sum)
}
