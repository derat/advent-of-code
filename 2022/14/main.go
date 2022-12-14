package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputLines("2022/14")

	// Blegh, do one pass first to find the bounds.
	xmin, xmax, ymin, ymax := 500, 500, 0, 0
	for _, ln := range input {
		vals := lib.ExtractInts(ln)
		lib.Assert(len(vals)%2 == 0)
		for i := 0; i < len(vals); i += 2 {
			xmin = lib.Min(xmin, vals[i])
			xmax = lib.Max(xmax, vals[i])
			ymin = lib.Min(ymin, vals[i+1])
			ymax = lib.Max(ymax, vals[i+1])
		}
	}

	// For part 2, save the original abyss from part 1 and draw a floor
	// with added width to match the full height.
	abyss := ymax + 1
	ymax += 2
	height := ymax - ymin + 1
	xmin -= height
	xmax += height

	// Parse the input again to draw the rock paths.
	grid := lib.NewByteGrid(ymax-ymin+1, xmax-xmin+1, '.')
	for _, ln := range input {
		vals := lib.ExtractInts(ln)
		for i := 0; i < len(vals)-3; i += 2 {
			x0, y0, x1, y1 := vals[i], vals[i+1], vals[i+2], vals[i+3]
			grid.SetRect(y0-ymin, x0-xmin, y1-ymin, x1-xmin, '#')
		}
	}

	// Draw the floor for part 2.
	grid.SetRect(ymax-ymin, 0, ymax-ymin, xmax-xmin, '#')

	get := func(x, y int) byte { return grid.Get(y-ymin, x-xmin, '.') }
	set := func(x, y int, ch byte) { grid[y-ymin][x-xmin] = ch }

	var sand int
	var reachedAbyss bool // hit the original abyss for part 1
	for {
		x, y := 500, 0
		if get(x, y) == 'o' { // source is blocked (part 2)
			fmt.Println(sand)
			break
		}

	FallLoop:
		for {
			if y == abyss && !reachedAbyss { // fell into the abyss (part 1)
				reachedAbyss = true
				fmt.Println(sand)
			}

			switch {
			case get(x, y+1) == '.': // can fall straight down
				y++
			case get(x-1, y+1) == '.': // can fall down and left
				x--
				y++
			case get(x+1, y+1) == '.': // can fall down and right
				x++
				y++
			default: // stops in current location
				set(x, y, 'o')
				sand++
				break FallLoop
			}
		}
	}
}
