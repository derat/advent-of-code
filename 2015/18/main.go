package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	lights := lib.InputLinesBytes("2015/18", '.', '#')
	lib.AssertEq(len(lights[0]), len(lights))
	dim := len(lights)

	const steps = 100 // given in challenge

	lights1 := lights
	for i := 0; i < steps; i++ {
		lights1 = update(lights1, false)
	}
	fmt.Println(count(lights1))

	// Part 2: Corners are stuck on.
	// This mutates the original lights, but whatever.
	lights2 := lights
	lights2[0][0] = '#'
	lights2[0][dim-1] = '#'
	lights2[dim-1][0] = '#'
	lights2[dim-1][dim-1] = '#'
	for i := 0; i < steps; i++ {
		lights2 = update(lights2, true)
	}
	fmt.Println(count(lights2))
}

// Performs a single update.
func update(lights [][]byte, cornersStuck bool) [][]byte {
	newLights := make([][]byte, len(lights))
	for r, row := range lights {
		newLights[r] = make([]byte, len(row))
		for c, ch := range row {
			if cornersStuck && (r == 0 || r == len(lights)-1) && (c == 0 || c == len(row)-1) {
				newLights[r][c] = '#'
				continue
			}

			cnt := neighbors(lights, r, c)
			on := (ch == '#' && (cnt == 2 || cnt == 3)) || (ch == '.' && cnt == 3)
			if on {
				newLights[r][c] = '#'
			} else {
				newLights[r][c] = '.'
			}
		}
	}
	return newLights
}

// Returns number of neighbors of the specified light that are turned on.
func neighbors(lights [][]byte, r, c int) int {
	on := func(r, c int) int {
		if r < 0 || r >= len(lights) || c < 0 || c >= len(lights[0]) {
			return 0
		}
		if lights[r][c] == '#' {
			return 1
		}
		return 0
	}
	return on(r-1, c-1) + on(r-1, c) + on(r-1, c+1) +
		on(r, c-1) + on(r, c+1) +
		on(r+1, c-1) + on(r+1, c) + on(r+1, c+1)
}

// Returns the number of lights that are turned on.
func count(lights [][]byte) int {
	var cnt int
	for _, r := range lights {
		for _, ch := range r {
			if ch == '#' {
				cnt++
			}
		}
	}
	return cnt
}
