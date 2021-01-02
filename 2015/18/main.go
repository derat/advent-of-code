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
	fmt.Println(lib.CountBytesFull(lights1, '#'))

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
	fmt.Println(lib.CountBytesFull(lights2, '#'))
}

// Performs a single update.
func update(lights [][]byte, cornersStuck bool) [][]byte {
	newLights := lib.NewBytes(len(lights), len(lights[0]), '.')
	for r, row := range lights {
		for c, ch := range row {
			if cornersStuck && (r == 0 || r == len(lights)-1) && (c == 0 || c == len(row)-1) {
				newLights[r][c] = '#'
				continue
			}

			cnt := lib.CountBytes(lights, r-1, c-1, r+1, c+1, '#')
			on := (ch == '#' && (cnt == 3 || cnt == 4)) || (ch == '.' && cnt == 3)
			if on {
				newLights[r][c] = '#'
			} else {
				newLights[r][c] = '.'
			}
		}
	}
	return newLights
}
