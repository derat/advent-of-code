package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const (
		nrows = 40

		trap = '^'
		safe = '.'

		left   = 0x1
		center = 0x2
		right  = 0x4
	)

	first := lib.InputLinesBytes("2016/18", '.', '^')[0]

	rows := make([][]byte, nrows)
	rows[0] = first
	for r := 1; r < nrows; r++ {
		rows[r] = make([]byte, len(rows[0]))
		prev := rows[r-1]
		for i := range rows[r] {
			st := lib.If(i > 0 && prev[i-1] == trap, left, 0) |
				lib.If(prev[i] == trap, center, 0) |
				lib.If(i < len(prev)-1 && prev[i+1] == trap, right, 0)
			if st == left|center || st == center|right || st == left || st == right {
				rows[r][i] = trap
			} else {
				rows[r][i] = safe
			}
		}
	}
	fmt.Println(lib.CountBytesFull(rows, '.'))
}
