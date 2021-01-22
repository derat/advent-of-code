package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInt64s("2019/9")

	// Part 1: Write 1 and print BOOST keycode (only output).
	vm := lib.NewIntcode(input)
	vm.Start()
	vm.In <- 1
	var vals []int64
	for v := range vm.Out {
		vals = append(vals, v)
	}
	lib.AssertEq(len(vals), 1)
	fmt.Println(vals[0])

	// Part 2: Write 2 and print coordinates of distress signal (only output).
	vm = lib.NewIntcode(input)
	vm.Start()
	vm.In <- 2
	fmt.Println(<-vm.Out)
}
