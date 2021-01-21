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

	// Part 1:
	minDist := math.MaxInt32
	test := func(hs []hseg, vs []vseg) {
		for i, h := range hs {
			for j, v := range vs {
				if i == 0 && j == 0 {
					continue // skip intersection at 0,0
				}
				if x, y, ok := cross(h, v); ok {
					minDist = lib.Min(minDist, lib.Abs(x)+lib.Abs(y))
				}
			}
		}
	}
	test(hsegs1, vsegs2)
	test(hsegs2, vsegs1)
	lib.AssertLess(minDist, math.MaxInt32)
	fmt.Println(minDist)
}

func read(ln string) ([]hseg, []vseg) {
	var x, y int
	var hsegs []hseg
	var vsegs []vseg
	for _, s := range strings.Split(ln, ",") {
		var dir byte
		var sz int
		lib.Extract(s, `^([UDLR])(\d+)$`, &dir, &sz)
		switch dir {
		case 'U':
			vsegs = append(vsegs, vseg{x: x, ymin: y, ymax: y + sz})
			y += sz
		case 'D':
			vsegs = append(vsegs, vseg{x: x, ymin: y - sz, ymax: y})
			y -= sz
		case 'L':
			hsegs = append(hsegs, hseg{xmin: x - sz, xmax: x, y: y})
			x -= sz
		case 'R':
			hsegs = append(hsegs, hseg{xmin: x, xmax: x + sz, y: y})
			x += sz
		}
	}
	return hsegs, vsegs
}

type hseg struct{ xmin, xmax, y int }
type vseg struct{ x, ymin, ymax int }

func cross(h hseg, v vseg) (int, int, bool) {
	if v.x >= h.xmin && v.x <= h.xmax && h.y >= v.ymin && h.y <= v.ymax {
		return v.x, h.y, true
	}
	return 0, 0, false
}
