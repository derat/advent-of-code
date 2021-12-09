package main

import (
	"bytes"
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	init := lib.InputByteGrid("2018/18", '.', '#', '|')

	// Part 1: Print resource score after 10 minutes.
	grid := init.Copy()
	for i := 0; i < 10; i++ {
		grid = update(grid)
	}
	fmt.Println(score(grid))

	// Part 2: Print resource score after 1000000000 minutes.
	const need = 1_000_000_000
	join := func(b [][]byte) string { return string(bytes.Join(b, []byte{'\n'})) }
	grid = init.Copy()
	seen := map[string]int{join(grid): 0}
	for min := 1; true; min++ {
		grid = update(grid)
		state := join(grid)
		if m, ok := seen[state]; ok {
			dist := min - m
			rem := (need - min) % dist
			for i := 0; i < rem; i++ {
				grid = update(grid)
			}
			fmt.Println(score(grid))
			break
		}
		seen[state] = min
	}
}

func update(grid lib.ByteGrid) lib.ByteGrid {
	next := lib.NewByteGrid(len(grid), len(grid[0]), '.')
	for r, row := range grid {
		for c, ch := range row {
			nch := ch
			trees := grid.CountRect(r-1, c-1, r+1, c+1, '|')
			lumber := grid.CountRect(r-1, c-1, r+1, c+1, '#')

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

func score(grid lib.ByteGrid) int {
	return grid.Count('|') * grid.Count('#')
}
