package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

const debug = false

func main() {
	input := lib.InputInt64s("2019/21")

	// Part 1: Write program to survey hull without falling into space,
	// then print hull damage (final non-ASCII output).
	vm := lib.NewIntcode(input)

	// If any of A, B, and C are false (there's a hole in front of us)
	// and D is true (we can safely jump), then jump.
	cmds := "NOT A T\n" +
		"NOT B J\n" +
		"OR T J\n" +
		"NOT C T\n" +
		"OR T J\n" +
		"AND D J\n" +
		"WALK\n"
	var nin int
	vm.InFunc = func() int64 {
		v := cmds[nin]
		nin++
		return int64(v)
	}
	vm.Start()
	for v := range vm.Out {
		if v <= 255 {
			if debug {
				fmt.Print(string(rune(v)))
			}
		} else {
			fmt.Println(v)
		}
	}
	lib.Assert(vm.Wait())
}
