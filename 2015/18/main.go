package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	lights := lib.InputLinesBytes("2015/18", '.', '#')
	lib.AssertEq(len(lights[0]), len(lights))

	const steps = 100 // given in challenge
	for i := 0; i < steps; i++ {
		newLights := make([][]byte, len(lights))
		for r, row := range lights {
			newLights[r] = make([]byte, len(row))
			for c, ch := range row {
				cnt := count(lights, r, c)
				on := (ch == '#' && (cnt == 2 || cnt == 3)) || (ch == '.' && cnt == 3)
				if on {
					newLights[r][c] = '#'
				} else {
					newLights[r][c] = '.'
				}
			}
		}
		lights = newLights
	}

	var cnt int
	for _, r := range lights {
		for _, ch := range r {
			if ch == '#' {
				cnt++
			}
		}
	}
	fmt.Println(cnt)
}

// Returns number of neighbors of the specified light that are turned on.
func count(lights [][]byte, r, c int) int {
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
