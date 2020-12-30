package main

import (
	"fmt"
	"math/big"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var a, b big.Int // registers
	var ins []inst

	reg := func(r string) *big.Int {
		switch r {
		case "a":
			return &a
		case "b":
			return &b
		default:
			panic(fmt.Sprintf("Bad register %q", r))
		}
	}

	for _, ln := range lib.InputLines("2015/23") {
		var in inst
		var rest string
		lib.Extract(ln, `^(hlf|tpl|inc|jmp|jie|jio) (.+)$`, &in.op, &rest)
		switch in.op {
		case "hlf", "tpl", "inc":
			in.reg = reg(rest)
		case "jmp":
			lib.Extract(rest, `^([+-]\d+)$`, &in.off)
		case "jie", "jio":
			var r string
			lib.Extract(rest, `^([ab]), ([+-]\d+)$`, &r, &in.off)
			in.reg = reg(r)
		}
		ins = append(ins, in)
	}

	var ip int
	for ip >= 0 && ip < len(ins) {
		in := &ins[ip]
		switch in.op {
		case "hlf":
			in.reg.Div(in.reg, big.NewInt(2))
			ip++
		case "tpl":
			in.reg.Mul(in.reg, big.NewInt(3))
			ip++
		case "inc":
			in.reg.Add(in.reg, big.NewInt(1))
			ip++
		case "jmp":
			ip += in.off
		case "jie":
			if in.reg.Bit(0) == 0 {
				ip += in.off
			} else {
				ip++
			}
		case "jio": // sigh: "o" is for "one", not "odd"
			if in.reg.Cmp(big.NewInt(1)) == 0 {
				ip += in.off
			} else {
				ip++
			}
		default:
			panic(fmt.Sprintf("Bad op %q", in.op))
		}
	}
	fmt.Println(b.String())
}

type inst struct {
	op  string
	reg *big.Int
	off int
}
