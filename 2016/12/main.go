package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var ins []instr
	for _, ln := range lib.InputLines("2016/12") {
		var in instr
		lib.Extract(ln, `^(cpy|inc|dec|jnz) (?:([a-d])|(-?\d+))(?: ([a-d])| (-?\d+))?$`,
			&in.op, &in.r1, &in.v1, &in.r2, &in.v2)
		ins = append(ins, in)
	}

	// Part 1
	a, _, _, _ := run(ins, 0, 0, 0, 0)
	fmt.Println(a)

	// Part 2: Initialize register c to 1 instead of to 0.
	a, _, _, _ = run(ins, 0, 0, 1, 0)
	fmt.Println(a)
}

type instr struct {
	op     string
	r1, r2 string
	v1, v2 int64
}

func run(ins []instr, a, b, c, d int64) (int64, int64, int64, int64) {
	reg := func(n string) *int64 {
		switch n {
		case "a":
			return &a
		case "b":
			return &b
		case "c":
			return &c
		case "d":
			return &d
		default:
			panic(fmt.Sprintf("Invalid register %q", n))
		}
	}

	var ip int64
	for ip >= 0 && ip < int64(len(ins)) {
		in := &ins[ip]
		switch in.op {
		case "cpy":
			if in.r1 != "" {
				*reg(in.r2) = *reg(in.r1)
			} else {
				*reg(in.r2) = in.v1
			}
			ip++
		case "inc":
			*reg(in.r1)++
			ip++
		case "dec":
			*reg(in.r1)--
			ip++
		case "jnz":
			v := in.v1
			if in.r1 != "" {
				v = *reg(in.r1)
			}
			if v != 0 {
				ip += in.v2
			} else {
				ip++
			}
		default:
			panic(fmt.Sprintf("Invalid op %q", in.op))
		}
	}

	return a, b, c, d
}
