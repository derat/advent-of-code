package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	grid := lib.InputByteGrid("2021/11")
	//fmt.Println("Before any steps:")
	//fmt.Println(grid.Dump())

	var nflashes int
	var inc func(r, c int)
	inc = func(r, c int) {
		// "First, the energy level of each octopus increases by 1."
		grid[r][c]++

		// "Then, any octopus with an energy level greater than 9 flashes. This increases the
		// energy level of all adjacent octopuses by 1, including octopuses that are diagonally
		// adjacent. If this causes an octopus to have an energy level greater than 9, it also
		// flashes. This process continues as long as new octopuses keep having their energy
		// level increased beyond 9. (An octopus can only flash at most once per step.)"
		if grid[r][c] == '9'+1 {
			nflashes++
			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					if nr, nc := r+dr, c+dc; grid.InBounds(nr, nc) && !(nr == r && nc == c) {
						inc(nr, nc)
					}
				}
			}
		}
	}

	for step := 1; true; step++ {
		old := nflashes
		grid.Iter(inc)

		// "Finally, any octopus that flashed during this step has its energy level set to 0,
		// as it used all of its energy to flash."
		grid.Iter(func(r, c int) {
			if grid[r][c] > '9' {
				grid[r][c] = '0'
			}
		})

		//fmt.Printf("\nAfter step %d:\n", step)
		//fmt.Println(grid.Dump())

		// Part 1: "How many total flashes are there after 100 steps?"
		if step == 100 {
			fmt.Println(nflashes)
		}

		// Part 2: "What is the first step during which all octopuses flash?"
		if nflashes == old+100 {
			fmt.Println(step)
			break
		}
	}
}
