package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	pgs := lib.InputParagraphs("2021/13")
	lib.AssertEq(len(pgs), 2)

	type dot struct{ x, y int }
	dots := make(map[dot]struct{}, len(pgs[0]))
	for _, ln := range pgs[0] {
		nums := lib.ExtractInts(ln)
		lib.AssertEq(len(nums), 2)
		dots[dot{nums[0], nums[1]}] = struct{}{}
	}

	for i, ln := range pgs[1] {
		var axis string
		var fold int
		lib.Extract(ln, `^fold along (x|y)=(\d+)$`, &axis, &fold)
		switch axis {
		case "x":
			for d := range dots {
				lib.Assert(d.x != fold)
				if d.x > fold {
					delete(dots, d)
					nx := d.x - 2*(d.x-fold)
					lib.AssertLessEq(0, nx)
					dots[dot{nx, d.y}] = struct{}{}
				}
			}
		case "y":
			for d := range dots {
				lib.Assert(d.y != fold)
				if d.y > fold {
					delete(dots, d)
					ny := d.y - 2*(d.y-fold)
					lib.AssertLessEq(0, ny)
					dots[dot{d.x, ny}] = struct{}{}
				}
			}
		}

		// Part 1: Find number of dots after first fold.
		if i == 0 {
			fmt.Println(len(dots))
		}
	}

	// Part 2: Read the code (eight capital letters) after completing all folds.
	var xmax, ymax int
	for d := range dots {
		xmax = lib.Max(xmax, d.x)
		ymax = lib.Max(ymax, d.y)
	}
	b := lib.NewByteGrid(ymax+1, xmax+1, ' ')
	for d := range dots {
		b[d.y][d.x] = 'X'
	}
	//fmt.Println(b.Dump())
	fmt.Println(lib.OCR([][]byte(b), ' '))
}
