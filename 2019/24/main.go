package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputLinesBytes("2019/24", '#', '.')

	// Part 1: Print biodiversity score after loop detected.
	grid := lib.CopyBytes(input)
	seen := map[uint64]struct{}{bio(grid): struct{}{}}
	for {
		next := lib.NewBytes(len(grid), len(grid[0]), '.')
		for r, row := range grid {
			for c, ch := range row {
				cnt := lib.CountBytes(grid, r-1, c, r-1, c, '#') +
					lib.CountBytes(grid, r+1, c, r+1, c, '#') +
					lib.CountBytes(grid, r, c-1, r, c-1, '#') +
					lib.CountBytes(grid, r, c+1, r, c+1, '#')
				switch ch {
				case '#':
					if cnt == 1 { // bug dies unless there's exactly one adjacent bug
						next[r][c] = '#'
					}
				case '.': // empty space infested if exactly 1 or 2 adjacent bugs
					if cnt == 1 || cnt == 2 {
						next[r][c] = '#'
					}
				default:
					lib.Panicf("Invalid byte %q", ch)
				}
			}
		}
		grid = next
		score := bio(grid)
		if _, ok := seen[score]; ok {
			fmt.Println(score)
			break
		} else {
			seen[score] = struct{}{}
		}
	}
}

func bio(grid [][]byte) uint64 {
	var score uint64
	pts := uint64(1)
	for _, row := range grid {
		for _, ch := range row {
			if ch == '#' {
				score += pts
			}
			pts <<= 1
		}
	}
	return score
}
