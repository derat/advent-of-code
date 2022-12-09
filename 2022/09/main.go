package main

import (
	"fmt"
	"time"

	"github.com/derat/advent-of-code/lib"
)

const (
	// This animation looks like garbage.
	anim  = false
	delay = 100 * time.Millisecond
)

func main() {
	knots := make([][2]int, 10)
	seen := make(map[[2]int]struct{})
	seen2 := make(map[[2]int]struct{})

	if anim {
		fmt.Println(dump(knots))
	}

	// Record knots' initial positions.
	seen[knots[1]] = struct{}{}
	seen2[knots[9]] = struct{}{}

	for _, ln := range lib.InputLines("2022/9") {
		var dir string
		var steps, sx, sy int
		lib.Extract(ln, `^([LRUD])\s+(\d+)$`, &dir, &steps)
		switch dir {
		case "L":
			sx = -steps
		case "R":
			sx = steps
		case "U":
			sy = -steps
		case "D":
			sy = steps
		default:
			lib.Panicf("Invalid direction %q", dir)
		}

		// Move the head one step at a time and then move the knots if needed.
		for sx != 0 || sy != 0 {
			ux := unit(sx)
			knots[0][0] += ux
			sx -= ux

			uy := unit(sy)
			knots[0][1] += uy
			sy -= uy

			for i := 1; i < len(knots); i++ {
				knots[i][0], knots[i][1] = move(knots[i-1][0], knots[i-1][1], knots[i][0], knots[i][1])
			}

			seen[knots[1]] = struct{}{}
			seen2[knots[9]] = struct{}{}

			if anim {
				time.Sleep(delay)
				fmt.Println(lib.ClearScreen)
				fmt.Println(dump(knots))
			}
		}
	}

	fmt.Println(len(seen))
	fmt.Println(len(seen2))
}

func unit(v int) int {
	if v == 0 {
		return 0
	}
	return v / lib.Abs(v)
}

func move(hx, hy, tx, ty int) (int, int) {
	dx := hx - tx
	dy := hy - ty
	if lib.Abs(dx) > 1 || lib.Abs(dy) > 1 {
		tx += unit(dx)
		ty += unit(dy)
	}
	return tx, ty
}

func dump(knots [][2]int) string {
	var xmin, xmax, ymin, ymax int
	for _, p := range knots {
		xmin = lib.Min(xmin, p[0])
		xmax = lib.Max(xmax, p[0])
		ymin = lib.Min(ymin, p[1])
		ymax = lib.Max(ymax, p[1])
	}
	grid := lib.NewByteGrid(ymax-ymin+1, xmax-xmin+1, '.')
	for i, p := range knots {
		grid[p[1]-ymin][p[0]-xmin] = '1' + byte(i)
	}
	return grid.Dump()
}
