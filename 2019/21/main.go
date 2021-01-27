package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

const debug = false

func main() {
	input := lib.InputInt64s("2019/21")

	// Part 1: Write program to survey hull without falling into space,
	// then print hull damage (final non-ASCII output).

	// If any of A, B, and C are false (there's a hole in front of us)
	// and D is true (we can safely jump), then jump.
	run(input, []string{
		"NOT A T",
		"NOT B J",
		"OR T J",
		"NOT C T",
		"OR T J",
		"AND D J",
		"WALK",
	})

	// Part 2: Use RUN instead of WALK. Registers E-I indicate ground 5-9 tiles away.
	//
	// Three things need to be true to jump:
	//	- Any of A-C are false (otherwise why bother?).
	//	- D is true (it's safe to jump).
	//	- E and H aren't both false (or we'll be stuck after jumping).
	run(input, []string{
		"NOT T T",
		"AND A T",
		"AND B T",
		"AND C T",
		"NOT T T", // T = !A || !B || !C = !(A && B && C)
		"OR E J",
		"OR H J",  // J = E || H
		"AND T J", // both of the previous conditions need to hold
		"AND D J", // set J to false if we aren't able to jump
		"RUN",
	})
}

func run(prog []int64, cmds []string) {
	vm := lib.NewIntcode(prog)

	in := strings.Join(cmds, "\n") + "\n"
	var nin int
	vm.InFunc = func() int64 {
		v := in[nin]
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
