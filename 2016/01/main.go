package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var x, y int
	dx, dy := 0, 1

	var twice bool
	var tx, ty int
	visited := map[uint64]struct{}{
		lib.PackInts(0, 0): struct{}{},
	}

	for _, s := range strings.Split(lib.InputLines("2016/1")[0], ", ") {
		var dir string
		var dist int
		lib.Extract(s, `^([LR])(\d+)$`, &dir, &dist)
		mult := lib.If(dir == "L", 1, -1)
		if dx != 0 {
			dy = mult * dx
			dx = 0
		} else {
			dx = mult * -dy
			dy = 0
		}

		for i := 0; i < dist; i++ {
			x, y = x+dx, y+dy

			// Part 2: Find first location visited twice.
			// I was initially thrown because I thought that only ending points count as visits.
			// The problem actually wants the first point where the path intersects itself.
			if !twice {
				key := lib.PackInts(x, y)
				if _, ok := visited[key]; ok {
					tx, ty = x, y
					twice = true
				} else {
					visited[key] = struct{}{}
				}
			}
		}
	}

	fmt.Println(lib.Abs(x) + lib.Abs(y))
	lib.Assert(twice)
	fmt.Println(lib.Abs(tx) + lib.Abs(ty))
}
