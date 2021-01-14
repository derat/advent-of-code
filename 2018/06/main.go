package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var points [][2]int // row, column
	rmin, rmax := math.MaxInt32, -1
	cmin, cmax := math.MaxInt32, -1
	for _, ln := range lib.InputLines("2018/6") {
		var r, c int
		lib.Extract(ln, `^(\d+), (\d+)$`, &c, &r) // col first to match example
		points = append(points, [2]int{r, c})
		rmin, rmax = lib.Min(rmin, r), lib.Max(rmax, r)
		cmin, cmax = lib.Min(cmin, c), lib.Max(cmax, c)
	}

	// Create a grid with a one-point boundary around the edges
	// and adjust points to fit in it.
	nrows := rmax - rmin + 3
	ncols := cmax - cmin + 3
	for i := range points {
		points[i][0] -= rmin - 1
		points[i][1] -= cmin - 1
	}

	// For each position, find the closest point(s).
	coords := make([][]coord, nrows)
	for r := range coords {
		coords[r] = make([]coord, ncols)
		for c := range coords[r] {
			co := &coords[r][c]
			*co = coord{math.MaxInt32, -1, 0}
			for id, p := range points {
				dist := lib.Abs(p[0]-r) + lib.Abs(p[1]-c)
				if dist < co.min {
					co.min = dist
					co.point = id
				} else if dist == co.min {
					co.point = -1
				}
				co.dists += dist
			}
		}
	}

	// Part 1: Print size of largest non-infinite region.
	sizes := make([]int, len(points))
	inf := make([]bool, len(points))
	for r, row := range coords {
		for c, co := range row {
			if co.point < 0 {
				continue
			}
			sizes[co.point]++
			if r == 0 || r == nrows-1 || c == 0 || c == ncols-1 {
				inf[co.point] = true
			}
		}
	}
	var max int
	for i, sz := range sizes {
		if sz > max && !inf[i] {
			max = sz
		}
	}
	fmt.Println(max)

	// Part 2: Print size of region of all locations with summed dists < 10000.
	var cnt int
	for _, row := range coords {
		for _, co := range row {
			if co.dists < 10_000 {
				cnt++
			}
		}
	}
	fmt.Println(cnt)
}

type coord struct {
	min   int // minimum Manhattan distance to a point
	point int // index of point with min distance (-1 if multiple)
	dists int // sum of dists to all points
}
