package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInts("2019/2")

	// Part 1: Set pos 1 to 12 and pos 2 to 2, and print pos 0 after halt.
	prog1 := append([]int(nil), input...)
	prog1[1] = 12
	prog1[2] = 2
	run(prog1)
	fmt.Println(prog1[0])
}

func run(prog []int) {
	var ip int
	for {
		op, in1, in2, out := prog[ip], prog[ip+1], prog[ip+2], prog[ip+3]
		switch op {
		case add:
			prog[out] = prog[in1] + prog[in2]
		case mul:
			prog[out] = prog[in1] * prog[in2]
		case hlt:
			return
		}
		ip += 4
	}
}

const (
	add = 1
	mul = 2
	hlt = 99
)
