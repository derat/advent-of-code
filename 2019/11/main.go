package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInt64s("2019/11")

	// Part 1: How many panels get painted at least once?
	vm := lib.NewIntcode(input)

	// Use unbuffered channels.
	vm.In = make(chan int64)
	vm.Out = make(chan int64)

	vm.Start()

	var r, c int
	dir := lib.Up
	painted := make(map[uint64]int64)
	var color int64  // current color (0 is black, 1 is white)
	painting := true // waiting for paint command (as opposed to turn)

Loop:
	for {
		select {
		case v, ok := <-vm.Out:
			if !ok { // channel closed (program done)
				break Loop
			}
			if painting {
				painted[lib.PackInts(r, c)] = v
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
			color = painted[lib.PackInts(r, c)]
		case vm.In <- color:
		}
	}
	lib.Assert(vm.Wait())
	fmt.Println(len(painted))
}
