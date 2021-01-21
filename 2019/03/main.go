package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	lines := lib.InputLines("2019/3")
	lib.AssertEq(len(lines), 2)
	hsegs1, vsegs1 := read(lines[0])
	hsegs2, vsegs2 := read(lines[1])

	// Part 1: Minimum Manhattan distance from origin of an intersection.
	minDist := math.MaxInt32
	test := func(hs []hseg, vs []vseg) {
		for i, h := range hs {
			for j, v := range vs {
				if i == 0 && j == 0 {
					continue // skip intersection at 0,0
				}
				if x, y, _, ok := cross(h, v); ok {
					minDist = lib.Min(minDist, lib.Abs(x)+lib.Abs(y))
				}
			}
		}
	}
	test(hsegs1, vsegs2)
	test(hsegs2, vsegs1)
	lib.AssertLess(minDist, math.MaxInt32)
	fmt.Println(minDist)

	// Part 2: Minimum summed distanced traveled at an intersection.
	minSteps := math.MaxInt32
	test2 := func(hs []hseg, vs []vseg) {
		for i, h := range hs {
			for j, v := range vs {
				if i == 0 && j == 0 {
					continue // skip intersection at 0,0
				}
				if _, _, steps, ok := cross(h, v); ok {
					minSteps = lib.Min(minSteps, steps)
				}
			}
		}
	}
	test2(hsegs1, vsegs2)
	test2(hsegs2, vsegs1)
	lib.AssertLess(minSteps, math.MaxInt32)
	fmt.Println(minSteps)
}

func read(ln string) ([]hseg, []vseg) {
	var x, y, dist int
	var hsegs []hseg
	var vsegs []vseg
	for _, s := range strings.Split(ln, ",") {
		var dir byte
		var sz int
		lib.Extract(s, `^([UDLR])(\d+)$`, &dir, &sz)
		switch dir {
		case 'U':
			vsegs = append(vsegs, vseg{
				x:    x,
				ymin: y,
				ymax: y + sz,
				dmin: dist,
				dmax: dist + sz,
			})
			y += sz
			dist += sz
		case 'D':
			vsegs = append(vsegs, vseg{
				x:    x,
				ymin: y - sz,
				ymax: y,
				dmin: dist + sz,
				dmax: dist,
			})
			y -= sz
			dist += sz
		case 'L':
			hsegs = append(hsegs, hseg{
				xmin: x - sz,
				xmax: x,
				y:    y,
				dmin: dist + sz,
				dmax: dist,
			})
			x -= sz
			dist += sz
		case 'R':
			hsegs = append(hsegs, hseg{
				xmin: x,
				xmax: x + sz,
				y:    y,
				dmin: dist,
				dmax: dist + sz,
			})
			x += sz
			dist += sz
		}
	}
	return hsegs, vsegs
}

type hseg struct{ xmin, xmax, y, dmin, dmax int }
type vseg struct{ x, ymin, ymax, dmin, dmax int }

func cross(h hseg, v vseg) (x, y, dist int, ok bool) {
	if v.x >= h.xmin && v.x <= h.xmax && h.y >= v.ymin && h.y <= v.ymax {
		dist := h.dmin + (h.dmax-h.dmin)*(v.x-h.xmin)/(h.xmax-h.xmin) +
			v.dmin + (v.dmax-v.dmin)*(h.y-v.ymin)/(v.ymax-v.ymin)
		return v.x, h.y, dist, true
	}
	return 0, 0, 0, false
}
