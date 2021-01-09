package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var ins []instr
	for _, ln := range lib.InputLines("2017/18") {
		ins = append(ins, newInstr(ln))
	}
	fmt.Println(run(ins))
}

func run(ins []instr) int {
	regs := make([]int, 26)
	get := func(b int, v int) int {
		if b >= 0 {
			return regs[b]
		}
		return v
	}

	var ip, lastSnd int
	for ip >= 0 && ip < len(ins) {
		jumped := false
		in := &ins[ip]
		switch in.op {
		case snd:
			lastSnd = get(in.r1, in.v1)
		case set:
			regs[in.r1] = get(in.r2, in.v2)
		case add:
			regs[in.r1] += get(in.r2, in.v2)
		case mul:
			regs[in.r1] *= get(in.r2, in.v2)
		case mod:
			regs[in.r1] %= get(in.r2, in.v2)
		case rcv:
			if get(in.r1, in.v1) != 0 {
				return lastSnd
			}
		case jgz:
			if get(in.r1, in.v1) > 0 {
				ip += get(in.r2, in.v2)
				jumped = true
			}
		default:
			lib.Panicf("Invalid op %d", in.op)
		}
		if !jumped {
			ip++
		}
	}
	panic("Didn't execute rcv")
}

type op int

const (
	snd op = iota
	set
	add
	mul
	mod
	rcv
	jgz
)

type instr struct {
	op     op
	r1, r2 int
	v1, v2 int
}

// I'm very tired of writing code like this and should add a library function
// the next time that this comes up.
func newInstr(ln string) instr {
	const re = `(?:([a-z])|(-?\d+))` // matches register or constant

	var op, r1, r2 string
	in := instr{r1: -1, r2: -1}
	switch {
	case lib.ExtractMaybe(ln, `^(snd|rcv) `+re+`$`, &op, &r1, &in.v1):
	case lib.ExtractMaybe(ln, `^(set|add|mul|mod) ([a-z]) `+re+`$`, &op, &r1, &r2, &in.v2):
	case lib.ExtractMaybe(ln, `^(jgz) `+re+` `+re+`$`, &op, &r1, &in.v1, &r2, &in.v2):
	default:
		lib.Panicf("Bad instruction %q", ln)
	}

	switch op {
	case "snd":
		in.op = snd
	case "set":
		in.op = set
	case "add":
		in.op = add
	case "mul":
		in.op = mul
	case "mod":
		in.op = mod
	case "rcv":
		in.op = rcv
	case "jgz":
		in.op = jgz
	default:
		lib.Panicf("Invalid op %q", op)
	}

	if r1 != "" {
		in.r1 = int(r1[0] - 'a')
	}
	if r2 != "" {
		in.r2 = int(r2[0] - 'a')
	}

	return in
}
