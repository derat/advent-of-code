package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

const debug = false

func main() {
	// This code started from 2016/23/main.go.
	// I'm removing the "tgl" instruction since it doesn't appear in my input,
	// and updated the instr struct to not hold strings and run() to handle
	// registers more efficiently.
	var ins []instr
	for _, ln := range lib.InputLines("2016/25") {
		var op, r1, r2 string
		var v1, v2 int64
		lib.Extract(ln, `^(cpy|inc|dec|jnz|out) (?:([a-d])|(-?\d+))(?: ([a-d])| (-?\d+))?$`,
			&op, &r1, &v1, &r2, &v2)
		ins = append(ins, newInstr(op, r1, r2, v1, v2))
	}

	// I spent way too long looking at the input and trying to figure out the answer
	// by hand. I let my code iterate over different inputs for a few minutes without
	// any luck. Much later, I realized that I was failing to increment the ip in the
	// handler for the new 'out' instruction. :-(
	for init := int64(1); true; init++ {
		// The first 8 instructions of my input have the effect of leaving a
		// at its initial value and setting d to that value plus 2548. If I'd written
		// working optimization code for 2016/23, I could use it here instead of
		// hardcoding this.
		if success := run(ins[8:], []int64{init, 0, 0, int64(init + 2548)}); success {
			fmt.Println(init)
			break
		}
	}
}

type op uint8

const (
	cpy op = iota
	inc
	dec
	jnz
	out
)

type instr struct {
	op     op
	r1, r2 int
	v1, v2 int64
}

func newInstr(op, r1, r2 string, v1, v2 int64) instr {
	in := instr{v1: v1, v2: v2}

	switch op {
	case "cpy":
		in.op = cpy
	case "inc":
		in.op = inc
	case "dec":
		in.op = dec
	case "jnz":
		in.op = jnz
	case "out":
		in.op = out
	}

	reg := func(s string) int {
		if s == "" {
			return -1
		}
		if s[0] >= 'a' && s[0] <= 'd' {
			return int(s[0] - 'a')
		}
		lib.Panicf("Invalid register %q", s)
		return -1
	}

	in.r1 = reg(r1)
	in.r2 = reg(r2)

	return in
}

func run(ins []instr, regs []int64) bool {
	seen := make(map[uint64]struct{}) // packed ints: ip, a, b, c, d

	var ip int64
	var outs int64
	for ip >= 0 && ip < int64(len(ins)) {
		state := lib.PackInts(int(ip), int(regs[0]), int(regs[1]), int(regs[2]), int(regs[3]))
		if _, ok := seen[state]; ok {
			return outs >= 2 && outs%2 == 0
		}
		seen[state] = struct{}{}

		in := &ins[ip]

		if debug {
			fmt.Printf("%2d %v %+v\n", ip, regs, in)
		}

		switch in.op {
		case cpy:
			if in.r2 >= 0 {
				if in.r1 >= 0 {
					regs[in.r2] = regs[in.r1]
				} else {
					regs[in.r2] = in.v1
				}
			}
			ip++
		case inc:
			if in.r1 >= 0 {
				regs[in.r1]++
			}
			ip++
		case dec:
			if in.r1 >= 0 {
				regs[in.r1]--
			}
			ip++
		case jnz:
			v := in.v1
			if in.r1 >= 0 {
				v = regs[in.r1]
			}
			o := in.v2
			if in.r2 >= 0 {
				o = regs[in.r2]
			}
			if v != 0 {
				ip += o
			} else {
				ip++
			}
		case out:
			v := in.v1
			if in.r1 >= 0 {
				v = regs[in.r1]
			}
			ip++

			// When we haven't seen any outputs, we need the first one to be 0.
			// After we've seen the first output (i.e. 0), we need to output 1.
			if v != outs%2 {
				return false
			}
			outs++
		default:
			lib.Panicf("Invalid op %q", in.op)
		}
	}

	return true
}
