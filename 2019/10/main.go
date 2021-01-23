package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputLinesBytes("2019/10", '#', '.')
	nrows, ncols := len(input), len(input[0])

	var asts []uint64
	for r, row := range input {
		for c, ch := range row {
			if ch == '#' {
				asts = append(asts, lib.PackInts(r, c))
			}
		}
	}

	// Part 1: Print maximum number of asteroids that can be seen from an asteroid.
	// For each asteroid, iterate over all other asteroids and cast a ray
	// through each to delete all asteroids behind it.
	var maxCnt int
	for _, p := range asts {
		r, c := lib.UnpackInt2(p)
		vis := make(map[uint64]struct{}, len(asts)-1)
		for _, a := range asts {
			if a != p {
				vis[a] = struct{}{}
			}
		}
		for p0 := range vis {
			r0, c0 := lib.UnpackInt2(p0)

			dr, dc := r0-r, c0-c

			// I had a lot of trouble getting this right. Not only does a dr,dc of 4,2 need to be
			// converted to 2,1, but 6,4 needs to be converted to 3,2.
			if dr == 0 {
				dc /= lib.Abs(dc)
			} else if dc == 0 {
				dr /= lib.Abs(dr)
			} else {
				for f := lib.Min(lib.Abs(dr), lib.Abs(dc)); f > 1; {
					if dr%f == 0 && dc%f == 0 {
						dr /= f
						dc /= f
						f = lib.Min(lib.Abs(dr), lib.Abs(dc))
					} else {
						f--
					}
				}
			}

			for r1, c1 := r0+dr, c0+dc; r1 >= 0 && c1 >= 0 && r1 < nrows && c1 < ncols; r1, c1 = r1+dr, c1+dc {
				delete(vis, lib.PackInts(r1, c1))
			}
		}
		maxCnt = lib.Max(maxCnt, len(vis))
	}
	fmt.Println(maxCnt)
}
