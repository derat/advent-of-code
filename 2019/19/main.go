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
	grid := lib.NewBytes(50, 50, '.')
	for r := range grid {
		for c := range grid[r] {
			vm := lib.NewIntcode(input)
			vm.Start()
			vm.In <- int64(c)
			vm.In <- int64(r)
			v := <-vm.Out
			switch v {
			case 0:
				grid[r][c] = '.'
			case 1:
				grid[r][c] = '#'
			default:
				lib.Panicf("Invalid value %v", v)
			}
			lib.Assert(vm.Wait())
		}
	}
	//fmt.Println(lib.DumpBytes(grid))
	fmt.Println(lib.CountBytesFull(grid, '#'))
}
