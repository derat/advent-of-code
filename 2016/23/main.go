package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// This code initially comes from 2016/12/main.go.
	var ins []instr
	for _, ln := range lib.InputLines("2016/23") {
		var in instr
		lib.Extract(ln, `^(cpy|inc|dec|jnz|tgl) (?:([a-d])|(-?\d+))(?: ([a-d])| (-?\d+))?$`,
			&in.op, &in.r1, &in.v1, &in.r2, &in.v2)
		ins = append(ins, in)
	}

	a, _, _, _ := run(ins, 7, 0, 0, 0)
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
			lib.Panicf("Invalid register %q", n)
			return nil
		}
	}

	ins = append([]instr(nil), ins...)

	var ip int64
	for ip >= 0 && ip < int64(len(ins)) {
		in := &ins[ip]

		switch in.op {
		case "cpy":
			if in.r2 != "" {
				if in.r1 != "" {
					*reg(in.r2) = *reg(in.r1)
				} else {
					*reg(in.r2) = in.v1
				}
			}
			ip++
		case "inc":
			if in.r1 != "" {
				*reg(in.r1)++
			}
			ip++
		case "dec":
			if in.r1 != "" {
				*reg(in.r1)--
			}
			ip++
		case "jnz":
			v := in.v1
			if in.r1 != "" {
				v = *reg(in.r1)
			}
			o := in.v2
			if in.r2 != "" {
				o = *reg(in.r2)
			}
			if v != 0 {
				ip += o
			} else {
				ip++
			}
		case "tgl":
			i := ip + *reg(in.r1)
			if i >= 0 && i < int64(len(ins)) {
				ti := &ins[i]
				switch ti.op {
				case "inc":
					ti.op = "dec"
				case "dec", "tgl":
					ti.op = "inc"
				case "jnz":
					ti.op = "cpy"
				case "cpy":
					ti.op = "jnz"
				default:
					lib.Panicf("Can't toggle %q", ti.op)
				}
			}
			ip++
		default:
			panic(fmt.Sprintf("Invalid op %q", in.op))
		}
	}

	return a, b, c, d
}
