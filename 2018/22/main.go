package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	lines := lib.InputLines("2018/22")
	var depth, tr, tc int
	lib.Extract(lines[0], `^depth: (\d+)$`, &depth)
	lib.Extract(lines[1], `^target: (\d+),(\d+)$`, &tc, &tr)

	ero := func(geo int) int { return (geo + depth) % 20183 }
	risk := func(ero int) int { return ero % 3 }

	geos := make([][]int, tr+1)
	var total int // total risk level
	for r := range geos {
		geos[r] = make([]int, tc+1)
		for c := range geos[r] {
			var geo int
			switch {
			case r == 0 && c == 0, r == tr && c == tc:
				geo = 0
			case r == 0:
				geo = c * 16807
			case c == 0:
				geo = r * 48271
			default:
				geo = ero(geos[r][c-1]) * ero(geos[r-1][c])
			}
			geos[r][c] = geo
			total += risk(ero(geo))
		}
	}

	// Part 1: Print total risk level for 0,0 to target coords.
	fmt.Println(total)
}
