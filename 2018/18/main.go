package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	grid := lib.InputLinesBytes("2018/18", '.', '#', '|')
	for i := 0; i < 10; i++ {
		grid = update(grid)
	}
	fmt.Println(lib.CountBytesFull(grid, '|') * lib.CountBytesFull(grid, '#'))
}

func update(grid [][]byte) [][]byte {
	next := lib.NewBytes(len(grid), len(grid[0]), '.')
	for r, row := range grid {
		for c, ch := range row {
			nch := ch
			trees := lib.CountBytes(grid, r-1, c-1, r+1, c+1, '|')
			lumber := lib.CountBytes(grid, r-1, c-1, r+1, c+1, '#')

			switch ch {
			case '.':
				if trees >= 3 {
					nch = '|'
				}
			case '|':
				if lumber >= 3 {
					nch = '#'
				}
			case '#':
				if lumber >= 2 && trees >= 1 { // we're counting ourselves
					nch = '#'
				} else {
					nch = '.'
				}
			default:
				lib.Panicf("Bad acre %q", ch)
			}
			next[r][c] = nch
		}
	}
	return next
}
