package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	grid := lib.InputByteGrid("2022/12")

	sr, sc := grid.Find('S')
	lib.Assert(grid.InBounds(sr, sc))
	grid[sr][sc] = 'a'

	er, ec := grid.Find('E')
	lib.Assert(grid.InBounds(er, ec))
	grid[er][ec] = 'z'

	type pos struct{ r, c int }
	search := func(starts []pos) int {
		return lib.AStar(
			starts,
			func(st pos) bool { return st.r == er && st.c == ec },
			func(st pos, next map[pos]int) {
				sh := grid[st.r][st.c]
				for _, n := range []pos{{st.r - 1, st.c}, {st.r + 1, st.c}, {st.r, st.c - 1}, {st.r, st.c + 1}} {
					if nh := grid.Get(n.r, n.c, 0); nh != 0 && nh <= sh+1 {
						next[n] = 1
					}
				}
			},
			func(st pos) int { return lib.Abs(er-st.r) + lib.Abs(ec-st.c) })
	}

	// Part 1: Minimum steps from starting position.
	fmt.Println(search([]pos{{sr, sc}}))

	// Part 2: Minimum steps from any position with height 'a'.
	var starts []pos
	grid.Iter(func(r, c int) {
		if grid[r][c] == 'a' {
			starts = append(starts, pos{r, c})
		}
	})
	fmt.Println(search(starts))
}
