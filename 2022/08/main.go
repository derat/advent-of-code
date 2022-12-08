package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputByteGrid("2022/8")

	// For a given point (r, c) at height h, count trees in direction d
	// up to but not including the first tree of height >= h.
	type view struct {
		r, c int
		h    byte
		d    lib.Dir
	}
	dists := make(map[view]int)

	var getDist func(r, c int, h byte, d lib.Dir) int
	getDist = func(r, c int, h byte, d lib.Dir) int {
		v := view{r, c, h, d}
		if n, ok := dists[v]; ok {
			return n
		}
		r0, c0 := r+d.DR(), c+d.DC()
		if !input.InBounds(r0, c0) {
			return 0
		}
		if h0 := input[r0][c0]; h0 >= h {
			return 0
		}
		n := getDist(r0, c0, h, d) + 1
		dists[v] = n
		return n
	}

	var vis, score int
	for r := 0; r < input.Rows(); r++ {
		for c := 0; c < input.Cols(); c++ {
			h := input[r][c]

			left := getDist(r, c, h, lib.Left)
			right := getDist(r, c, h, lib.Right)
			up := getDist(r, c, h, lib.Up)
			down := getDist(r, c, h, lib.Down)

			outLeft := left == c
			outRight := right == input.MaxCol()-c
			outUp := up == r
			outDown := down == input.MaxRow()-r

			// Part 1: count trees that are visible from outside.
			if outLeft || outRight || outUp || outDown {
				vis++
			}

			// Part 2: find best scenic score (multiplied visible trees in each direction).
			// Add 1 to the count if it doesn't reach outside to include the blocking tree.
			add := func(cnt int, out bool) int { return cnt + lib.If(!out, 1, 0) }
			sc := add(left, outLeft) * add(right, outRight) * add(up, outUp) * add(down, outDown)
			if sc > score {
				score = sc
			}
		}
	}
	fmt.Println(vis)
	fmt.Println(score)
}
