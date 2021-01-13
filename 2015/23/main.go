package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var ins []lib.Instr
	for _, ln := range lib.InputLines("2015/23") {
		ins = append(ins, lib.NewInstr(ln, 'a', 'b', map[uint8]string{
			hlf: `^hlf ([ab])$`,
			tpl: `^tpl ([ab])$`,
			inc: `^inc ([ab])$`,
			jmp: `^jmp ([+-]\d+)$`,
			jie: `^jie ([ab]), ([+-]\d+)$`,
			jio: `^jio ([ab]), ([+-]\d+)$`,
		}))
	}

	regs := make([]int64, 2)
	run := func() {
		var ip int
		for ip >= 0 && ip < len(ins) {
			in := &ins[ip]
			switch in.Op {
			case hlf:
				*in.Ptr(0, regs) /= 2
				ip++
			case tpl:
				*in.Ptr(0, regs) *= 3
				ip++
			case inc:
				*in.Ptr(0, regs) += 1
				ip++
			case jmp:
				ip += int(in.Val(0, regs))
			case jie:
				if in.Val(0, regs)&1 == 0 {
					ip += int(in.Val(1, regs))
				} else {
					ip++
				}
			case jio: // sigh: "o" is for "one", not "odd"
				if in.Val(0, regs) == 1 {
					ip += int(in.Val(1, regs))
				} else {
					ip++
				}
			default:
				panic(fmt.Sprintf("Bad op %v", in.Op))
			}
		}
	}

	run()
	fmt.Println(regs[1])

	// Part 2: Register 'a' starts 1 instead of as 0
	regs[0] = 1
	regs[1] = 0
	run()
	fmt.Println(regs[1])
}

const (
	hlf = iota
	tpl
	inc
	jmp
	jie
	jio
)
