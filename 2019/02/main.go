package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInt64s("2019/2")

	// Part 1: Set pos 1 to 12 and pos 2 to 2, and print pos 0 after halt.
	vm := lib.NewIntcode(input)
	vm.Mem[1] = 12
	vm.Mem[2] = 2
	lib.Assert(vm.Run())
	fmt.Println(vm.Mem[0])

	// Part 2: Find noun (pos 1) and verb (pos 2) values in [0,99] that result
	// in pos 0 being 19690720.
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			vm := lib.NewIntcode(input)
			vm.Mem[1] = int64(noun)
			vm.Mem[2] = int64(verb)
			lib.Assert(vm.Run())
			if vm.Mem[0] == 19690720 {
				fmt.Println(100*noun + verb)
				break
			}
		}
	}
}
