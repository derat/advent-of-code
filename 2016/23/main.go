package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

const debug = false

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

	// Part 2: Initialize reg a to 12 instead of 7.
	a, _, _, _ = run(ins, 12, 0, 0, 0)
	fmt.Println(a)
}

type instr struct {
	op     string
	r1, r2 string
	v1, v2 int64

	adds map[string]int64
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
	//optimize(ins)

	var ip int64
	for ip >= 0 && ip < int64(len(ins)) {
		in := &ins[ip]

		if debug {
			fmt.Printf("a=%d b=%d c=%d d=%d\n", a, b, c, d)
			for i, in := range ins {
				if int64(i) == ip {
					fmt.Print("* ")
				} else {
					fmt.Print("  ")
				}
				fmt.Printf("%2d %+v\n", i, in)
			}
			fmt.Println()
		}

		// Disgusting: I just manually identified a few loops in my input and
		// replaced them with multiplications here. Luckily, this code doesn't
		// seem to be modified (much) by the tgl instruction.
		if ip == 4 && ops(ins[4:10], "cpy", "inc", "dec", "jnz", "dec", "jnz") {
			// 4 cpy b c
			// 5 inc a
			// 6 dec c
			// 7 jnz c -2
			// 8 dec d
			// 9 jnz d -5
			*reg("a") += *reg("d") * *reg("b")
			*reg("c") = 0
			*reg("d") = 0
			ip = 10
			continue
		} else if ip == 13 && ops(ins[13:16], "dec", "inc", "jnz") {
			// 13 dec d
			// 14 inc c
			// 15 jnz d -2
			*reg("c") += *reg("d")
			*reg("d") = 0
			ip = 16
			continue
		} else if ip == 20 && ops(ins[20:26], "cpy", "inc", "dec", "jnz", "dec", "jnz") {
			// 20 cpy 97 d (initially jnz)
			// 21 inc a
			// 22 dec d (initially inc)
			// 23 jnz d -2
			// 24 dec c (initially inc)
			// 25 jnz c -5
			*reg("a") += 97 * *reg("c")
			*reg("d") = 0
			*reg("c") = 0
			ip = 26
			continue
		}

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
				// Here's where we'd skip over loops if optimize() worked correctly.
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
			//optimize(ins)
			ip++
		default:
			panic(fmt.Sprintf("Invalid op %q", in.op))
		}
	}

	return a, b, c, d
}

// ops returns true if the supplied instructions have the supplied ops.
func ops(ins []instr, ops ...string) bool {
	lib.AssertEq(len(ins), len(ops))
	for i := range ins {
		if ins[i].op != ops[i] {
			return false
		}
	}
	return true
}

// optimize attempts to find loops in ins.
// Before getting it to recognize nested loops, I gave up and just manually
// identified the loops in my input, though.
func optimize(ins []instr) {
	// Clear existing optimizations first.
	for i := range ins {
		ins[i].adds = nil
	}

	for {
		found := false

	Loop:
		for i := range ins {
			// Skip already-optimized instructions.
			if ins[i].adds != nil {
				continue
			}

			// We're only interested in jnz ops that check a register and jump back by a constant.
			if ins[i].op != "jnz" || ins[i].r1 == "" || ins[i].r2 != "" || ins[i].v2 >= 0 {
				continue
			}

			// Count the increments and decrements that happen each time through.
			adds := make(map[string]int64, 4)
			for j := i + int(ins[i].v2); j < i; j++ {
				switch {
				case ins[j].op == "inc":
					adds[ins[j].r1] += 1
				case ins[j].op == "dec":
					adds[ins[j].r1] -= 1
				// This is the spot where we'd identify nested loops.
				//case ins[j].op == "jnz" && ins[j].adds != nil:
				default:
					// If we found a non-inc/dec, give up.
					continue Loop
				}
			}
			ins[i].adds = adds
			found = true
		}

		if !found {
			break
		}
	}
}
