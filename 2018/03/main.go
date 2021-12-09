package main

import (
	"fmt"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	claims := make(map[int]claim)
	for _, ln := range lib.InputLines("2018/3") {
		var id, w, h int
		var cl claim
		lib.Extract(ln, `^#(\d+) @ (\d+),(\d+): (\d+)x(\d+)$`, &id, &cl.l, &cl.t, &w, &h)
		cl.b, cl.r = cl.t+h-1, cl.l+w-1
		claims[id] = cl
	}

	// Part 1: Print number of points with two or more overlapping claims.
	const dim = 1000
	grid := lib.NewByteGrid(dim, dim, 0)
	for _, cl := range claims {
		for r := cl.t; r <= cl.b; r++ {
			for c := cl.l; c <= cl.r; c++ {
				if cnt := grid[r][c]; cnt < 255 {
					grid[r][c]++
				}
			}
		}
	}
	var cnt int
	for r := 0; r < dim; r++ {
		for c := 0; c < dim; c++ {
			if grid[r][c] >= 2 {
				cnt++
			}
		}
	}
	fmt.Println(cnt)

	// Part 2: Print ID of the only non-overlapped claim.

	// This is the boring approach of just checking each claim against the grid
	// that was computed for part 1.
	/*
		Loop:
			for id, cl := range claims {
				for r := cl.t; r <= cl.b; r++ {
					for c := cl.l; c <= cl.r; c++ {
						if cnt := grid[r][c]; cnt != 1 {
							continue Loop
						}
					}
				}
				fmt.Println(id)
			}
	*/

	// First, sort the claims by increasing top edges.
	// Then, walk through them top to bottom, checking for possible intersections.
	cols := make(map[int]int, len(claims))
	for id := range claims {
		cols[id] = 0
	}
	sorted := lib.MapIntKeys(claims)
	sort.Slice(sorted, func(i, j int) bool { return claims[sorted[i]].t < claims[sorted[j]].t })
	for i, id := range sorted {
		cl := claims[id]
		for _, oid := range sorted[i+1:] {
			ocl := claims[oid]
			if ocl.t > cl.b {
				break
			}
			if ocl.r >= cl.l && ocl.l <= cl.r {
				cols[id]++
				cols[oid]++
			}
		}
	}
	for id, cnt := range cols {
		if cnt == 0 {
			fmt.Println(id)
			break
		}
	}

	// The other obvious approach is to iterate through all the rectangles, comparing each
	// against the earlier ones and splitting them into smaller rectangles as intersections
	// are found. It's not obvious to me that that'd be much faster, though and the code
	// seems way more complicated.
	//
	// I feel like there must be a better way to do this, perhaps by sorting the claims in
	// multiple ways and then walking through them in a single pass. I wasn't able to come
	// up with anything that works, though.
}

type claim struct {
	t, b, l, r int
}
