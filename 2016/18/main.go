package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const (
		nrows  = 40
		nrows2 = 400000

		trap = '^'
		safe = '.'

		left   = 0x1
		center = 0x2
		right  = 0x4
	)

	first := lib.InputLinesBytes("2016/18", '.', '^')[0]

	cnt := lib.CountBytesFull([][]byte{first}, '.')
	prev := first
	for r := 1; r < nrows2; r++ {
		// Part 1
		if r == nrows {
			fmt.Println(cnt)
		}

		// This is just the naive approach of carrying forward the previous row,
		// but I wonder if we could optimize things further by computing the *next*
		// row as we determine the current row's tiles, which would then feed into
		// the row after that, and so on. That's maybe more interesting if we're
		// computing just a single tile from a given row in the far-distant future.
		// I also didn't spend any time looking for other patterns.
		row := make([]byte, len(first))
		for i := range row {
			st := lib.If(i > 0 && prev[i-1] == trap, left, 0) |
				lib.If(prev[i] == trap, center, 0) |
				lib.If(i < len(prev)-1 && prev[i+1] == trap, right, 0)
			if st == left|center || st == center|right || st == left || st == right {
				row[i] = trap
			} else {
				row[i] = safe
				cnt++
			}
		}
		prev = row
	}

	// Part 2
	fmt.Println(cnt)
}
