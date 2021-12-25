package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	grid := lib.InputByteGrid("2021/25")

	// This was a pretty easy one. My input is 137 rows by 139 columns, so it seemed
	// clear that checking every grid and its west/north neighbors each step isn't the
	// way to go. The only points we need to track are empty ones with east-/south-moving
	// sea cucumbers to their west/north, and we can update the frontiers by checking
	// the neighbors after a sea cucumber moves.

	type point struct{ r, c int }
	get := func(p point) byte { return grid[p.r][p.c] }
	set := func(p point, ch byte) { grid[p.r][p.c] = ch }

	west := func(p point) point { return point{p.r, (p.c - 1 + grid.Cols()) % grid.Cols()} }
	east := func(p point) point { return point{p.r, (p.c + 1) % grid.Cols()} }
	north := func(p point) point { return point{(p.r - 1 + grid.Rows()) % grid.Rows(), p.c} }
	south := func(p point) point { return point{(p.r + 1) % grid.Rows(), p.c} }

	// add adds p to ef and sf if p is empty and its west/north neighbors have
	// east-/south-moving sea cucumbers, respectively.
	add := func(p point, ef, sf map[point]struct{}) {
		if get(p) != '.' {
			return
		}
		if wp := west(p); get(wp) == '>' {
			ef[p] = struct{}{}
		}
		if np := north(p); get(np) == 'v' {
			sf[p] = struct{}{}
		}
	}

	// Track the frontiers. I didn't think about it very hard, but I believe it might
	// be possible to just use a single frontier instead. Tracking the two directions
	// separately lets me put some asserts to the loop to make sure that I'm not keeping
	// any unnecessary points in the frontier, though.
	efront := make(map[point]struct{})
	sfront := make(map[point]struct{})
	for r, row := range grid {
		for c := range row {
			add(point{r, c}, efront, sfront)
		}
	}

	// Part 1: "What is the first step on which no sea cucumbers move?"
	for step := 1; true; step++ {
		var moves int

		nefront := make(map[point]struct{})
		for p := range efront {
			wp := west(p)
			lib.AssertEq(get(wp), '>')
			set(p, '>')
			set(wp, '.')
			moves++

			// This space is occupied now, but the space to the west has opened up.
			// We also need to check the space to the east, which we may be able to
			// move into next time.
			delete(sfront, p)
			add(wp, nefront, sfront)
			add(east(p), nefront, sfront)
		}
		efront = nefront

		nsfront := make(map[point]struct{})
		for p := range sfront {
			np := north(p)
			lib.AssertEq(get(np), 'v')
			set(p, 'v')
			set(np, '.')
			moves++

			// This space is occupied now, but the space above has opened up.
			// We also need to check the space to the south.
			delete(efront, p)
			add(np, efront, nsfront)
			add(south(p), efront, nsfront)
		}
		sfront = nsfront

		if moves == 0 {
			fmt.Println(step)
			break
		}
	}
}
