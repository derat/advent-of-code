package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInt64s("2019/11")

	// Part 1: How many panels get painted at least once?
	panels := make(map[uint64]int64)
	run(input, panels)
	fmt.Println(len(panels))

	// Part 2: What's the 8-letter registration identifier that's
	// painted after starting on a single white panel?
	panels = map[uint64]int64{lib.PackInts(0, 0): 1}
	run(input, panels)

	rmin, rmax := math.MaxInt32, math.MinInt32
	cmin, cmax := math.MaxInt32, math.MinInt32
	for p := range panels {
		r, c := lib.UnpackIntSigned2(p)
		rmin = lib.Min(rmin, r)
		rmax = lib.Max(rmax, r)
		cmin = lib.Min(cmin, c)
		cmax = lib.Max(cmax, c)
	}
	nrows := rmax - rmin + 1
	ncols := cmax - cmin + 1
	grid := lib.NewBytes(nrows, ncols, ' ')
	for p, v := range panels {
		r, c := lib.UnpackIntSigned2(p)
		if v == 1 {
			grid[r][c] = '#'
		}
	}
	//fmt.Println(lib.DumpBytes(grid))
	fmt.Println(lib.OCR(grid, ' '))
}

func run(prog []int64, panels map[uint64]int64) {
	vm := lib.NewIntcode(prog)

	// Use unbuffered channels.
	vm.In = make(chan int64)
	vm.Out = make(chan int64)

	vm.Start()

	var r, c int
	dir := lib.Up
	color := panels[lib.PackInts(r, c)] // current color (0 is black, 1 is white)
	painting := true                    // waiting for paint command (as opposed to turn)

Loop:
	for {
		select {
		case v, ok := <-vm.Out:
			if !ok { // channel closed (program done)
				break Loop
			}
			if painting {
				panels[lib.PackInts(r, c)] = v
			} else { // turning
				switch v {
				case 0:
					dir = dir.Left()
				case 1:
					dir = dir.Right()
				default:
					lib.Panicf("Invalid turn direction %v", v)
				}
				r += dir.DR()
				c += dir.DC()
			}
			painting = !painting
			color = panels[lib.PackInts(r, c)]
		case vm.In <- color:
		}
	}
	lib.Assert(vm.Wait())
}
