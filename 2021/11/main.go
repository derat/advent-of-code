package main

import (
	"fmt"
	"time"

	"github.com/derat/advent-of-code/lib"
)

const (
	animate = false
	delay   = 50 * time.Millisecond
)

func main() {
	grid := lib.InputByteGrid("2021/11")
	//fmt.Println("Before any steps:")
	//fmt.Println(grid.Dump())
	if animate {
		fmt.Print(lib.ClearScreen)
	}

	var nflashes int
	for step := 1; true; step++ {
		old := nflashes
		var pending [][2]int // pending flashes from this round

		inc := func(r, c int) {
			grid[r][c]++
			if grid[r][c] == '9'+1 {
				nflashes++
				pending = append(pending, [2]int{r, c})
			}
		}

		// "First, the energy level of each octopus increases by 1."
		grid.Iter(inc)
		if animate {
			dump(grid)
		}

		// "Then, any octopus with an energy level greater than 9 flashes. This increases the energy
		// level of all adjacent octopuses by 1, including octopuses that are diagonally adjacent.
		// If this causes an octopus to have an energy level greater than 9, it also flashes. This
		// process continues as long as new octopuses keep having their energy level increased
		// beyond 9. (An octopus can only flash at most once per step.)"
		for len(pending) > 0 {
			r, c := pending[0][0], pending[0][1]
			pending = pending[1:]

			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					if nr, nc := r+dr, c+dc; grid.InBounds(nr, nc) && !(nr == r && nc == c) {
						inc(nr, nc)
					}
				}
			}
		}
		if animate && nflashes > old {
			dump(grid)
		}

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
		if step == 100 && !animate {
			fmt.Println(nflashes)
		}

		// Part 2: "What is the first step during which all octopuses flash?"
		if nflashes == old+100 && !animate {
			fmt.Println(step)
			break
		}
	}
}

var (
	colors = []string{
		lib.Color(233),
		lib.Color(235),
		lib.Color(238),
		lib.Color(241),
		lib.Color(244),
		lib.Color(246),
		lib.Color(248),
		lib.Color(250),
		lib.Color(252),
		lib.Color(254),
	}
	flashColor = lib.Color(11)
)

func dump(grid lib.ByteGrid) {
	var s string
	s += lib.MoveHome
	for r, row := range grid {
		if r != 0 {
			s += "\n"
		}
		for _, ch := range row {
			if i := ch - '0'; i >= 0 && int(i) < len(colors) {
				s += colors[i] + string(ch) + " "
			} else {
				s += flashColor + "█ "
			}
		}
	}
	fmt.Print(s)
	time.Sleep(delay)
}
