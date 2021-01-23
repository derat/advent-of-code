package main

import (
	"fmt"
	"math"
	"sort"

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
	var maxCnt int // number of asteroids visible from sr,sc
	var sr, sc int // coordinates of monitoring station (needed for part 2)
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
			dr, dc := reduce(r0-r, c0-c)
			for r1, c1 := r0+dr, c0+dc; r1 >= 0 && c1 >= 0 && r1 < nrows && c1 < ncols; r1, c1 = r1+dr, c1+dc {
				delete(vis, lib.PackInts(r1, c1))
			}
		}
		if cnt := len(vis); cnt > maxCnt {
			sr, sc = r, c
			maxCnt = cnt
		}
	}
	fmt.Println(maxCnt)

	// Part 2: Vaporize asteroids in clockwise order starting from straight up.
	// Print 100*c+r of 200th to be vaporized.
	rem := make(map[uint64]struct{}, len(asts)-1)
	sp := lib.PackInts(sr, sc)
	for _, p := range asts {
		if p != sp {
			rem[p] = struct{}{}
		}
	}
	var nvap int // number of vaporized asteroids
	for len(rem) > 0 {
		sm := make(map[uint64]struct{}, len(rem))
		for p := range rem {
			r, c := lib.UnpackInt2(p)
			dr, dc := reduce(r-sr, c-sc)
			sm[lib.PackInts(dr, dc)] = struct{}{}
		}
		slopes := make([][2]int, 0, len(sm)) // dr,dc
		for p := range sm {
			dr, dc := lib.UnpackIntSigned2(p)
			slopes = append(slopes, [2]int{dr, dc})
		}
		sort.Slice(slopes, func(i, j int) bool {
			a, b := slopes[i], slopes[j]
			return sortKey(a[0], a[1]) < sortKey(b[0], b[1])
		})
	Loop:
		for _, sl := range slopes {
			dr, dc := sl[0], sl[1]
			for r, c := sr+dr, sc+dc; r >= 0 && c >= 0 && r < nrows && c < ncols; r, c = r+dr, c+dc {
				p := lib.PackInts(r, c)
				if _, ok := rem[p]; ok {
					delete(rem, p)
					nvap++
					if nvap == 200 {
						fmt.Println(100*c + r)
					}
					continue Loop
				}
			}
		}
	}
}

// reduce reduces slope dr,dc to its minimum integer factors.
func reduce(dr, dc int) (int, int) {
	if dr == 0 {
		return 0, dc / lib.Abs(dc)
	}
	if dc == 0 {
		return dr / lib.Abs(dr), 0
	}

	// I had a lot of trouble getting this right. Not only does a dr,dc of 4,2 need to be
	// converted to 2,1, but 6,4 needs to be converted to 3,2.
	for f := lib.Min(lib.Abs(dr), lib.Abs(dc)); f > 1; {
		if dr%f == 0 && dc%f == 0 {
			dr /= f
			dc /= f
			f = lib.Min(lib.Abs(dr), lib.Abs(dc))
		} else {
			f--
		}
	}
	return dr, dc
}

// sortKey returns a value to use for sorting dr,dc values clockwise starting from straight up.
func sortKey(dr, dc int) float64 {
	// Moving clockwise, atan2(dr,dc) gives us:
	//  -1, 0: -pi/2 (up)
	//   0, 1: 0     (right)
	//   1, 0: pi/2  (down)
	//   0,-1: pi    (left)
	tan := math.Atan2(float64(dr), float64(dc))

	// In the area from left to up (i.e. dr and dc both negative),
	// atan2 goes from -pi to -pi/2. Add 2*pi to those values so we
	// can just do direct comparisons to get the desired ordering.
	if tan < -math.Pi/2 {
		tan += 2 * math.Pi
	}

	return tan
}
