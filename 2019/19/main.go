package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInt64s("2019/19")

	// Part 1: Count affected points in 50x50 grid.
	// I assumed that the program would read repeated coordinates,
	// but it seems like I need to run it once for each point...?
	grid := lib.NewByteGrid(50, 50, '.')
	for r := range grid {
		for c := range grid[r] {
			if check(input, r, c) {
				grid[r][c] = '#'
			}
		}
	}
	fmt.Println(grid.Count('#'))

	// Part 2: Fit a 100x100 square into the beam as close to 0,0 as possible.
	// Then print 10000*x + y for its top-left corner.
	const dim = 100

	// I'm sure there's a better way to do this. *shrug*
	// Walk down until the top-right corner of the square hits the beam.
	var r, c int
	for ; !check(input, r, c+dim-1); r++ {
	}
	// Follow the top of the beam down and to the right until the square fits.
	for {
		if check(input, r+dim-1, c) {
			break
		}
		if check(input, r, c+dim) {
			c++
		} else {
			r++
		}
	}
	// Now walk the square back up and to the left as far as we can.
	for {
		if check(input, r+dim-1, c-1) { // move left
			c--
		} else if check(input, r-1, c+dim-1) { // move up
			r--
		} else if check(input, r+dim-2, c-1) && check(input, r-1, c+dim-2) {
			c--
			r--
		} else {
			break
		}
	}
	fmt.Println(c*10_000 + r)
}

func check(prog []int64, r, c int) (hit bool) {
	vm := lib.NewIntcode(prog)
	vm.Start()
	vm.In <- int64(c)
	vm.In <- int64(r)
	v := <-vm.Out
	switch v {
	case 0:
		hit = false
	case 1:
		hit = true
	default:
		lib.Panicf("Invalid value %v", v)
	}
	lib.Assert(vm.Wait())
	return hit
}
