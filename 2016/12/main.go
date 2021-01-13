package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var ins []lib.Instr
	for _, ln := range lib.InputLines("2016/12") {
		ins = append(ins, lib.NewInstr(ln, 'a', 'd', map[uint8]string{
			cpy: `^cpy ([a-d]|-?\d+) ([a-d])$`,
			inc: `^inc ([a-d])$`,
			dec: `^dec ([a-d])$`,
			jnz: `^jnz ([a-d]|-?\d+) ([a-d]|-?\d+)$`,
		}))
	}

	// Part 1
	regs := make([]int64, 4)
	run(ins, regs)
	fmt.Println(regs[0])

	// Part 2: Initialize register c to 1 instead of to 0.
	regs = []int64{0, 0, 1, 0}
	run(ins, regs)
	fmt.Println(regs[0])
}

func run(ins []lib.Instr, regs []int64) {
	var ip int64
	for ip >= 0 && ip < int64(len(ins)) {
		jumped := false
		in := &ins[ip]
		switch in.Op {
		case cpy:
			*in.Ptr(1, regs) = in.Val(0, regs)
		case inc:
			*in.Ptr(0, regs)++
		case dec:
			*in.Ptr(0, regs)--
		case jnz:
			if in.Val(0, regs) != 0 {
				ip += in.Val(1, regs)
				jumped = true
			}
		default:
			panic(fmt.Sprintf("Invalid op %v", in.Op))
		}
		if !jumped {
			ip++
		}
	}
}

const (
	cpy = iota
	inc
	dec
	jnz
)
