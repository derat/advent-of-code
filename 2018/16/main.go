package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var samples []sample
	var prog []instr
	pgs := lib.InputParagraphs("2018/16")
	for i, pg := range pgs {
		if i < len(pgs)-1 {
			lib.AssertEq(len(pg), 3)
			var sa sample
			lib.Extract(pg[0], `^Before:\s*\[(\d+),\s*(\d+),\s*(\d+),\s*(\d+)\]$`,
				&sa.before[0], &sa.before[1], &sa.before[2], &sa.before[3])
			sa.in = newInstr(pg[1])
			lib.Extract(pg[2], `^After:\s*\[(\d+),\s*(\d+),\s*(\d+),\s*(\d+)\]$`,
				&sa.after[0], &sa.after[1], &sa.after[2], &sa.after[3])
			samples = append(samples, sa)
		} else {
			for _, ln := range pg {
				prog = append(prog, newInstr(ln))
			}
		}
	}

	// Part 1: Print number of samples where instruction behaves like three or more opcodes.
	var cnt int
	for _, sa := range samples {
		if len(sa.ops()) >= 3 {
			cnt++
		}
	}
	fmt.Println(cnt)
}

type instr struct {
	op, a, b, c int
}

func newInstr(s string) instr {
	var in instr
	lib.Extract(s, `^(\d+)\s+(\d+)\s+(\d+)\s+(\d+)$`, &in.op, &in.a, &in.b, &in.c)
	return in
}

type op int

const (
	addr op = iota
	addi
	mulr
	muli
	banr
	bani
	borr
	bori
	setr
	seti
	gtir
	gtri
	gtrr
	eqir
	eqri
	eqrr
)

type sample struct {
	before, after [4]int
	in            instr
}

func (s *sample) ops() []op {
	va, vb := s.in.a, s.in.b
	reg := func(idx int) int {
		if idx < 0 || idx >= len(s.before) {
			panic(0)
		}
		return s.before[idx]
	}

	// I couldn't think of a more concise way to do this.
	// Each computation happens within a function so we
	// can recover from panics in the case of an invalid register reference.
	var ops []op
	for op, fc := range map[op]func() int{
		addr: func() int { return reg(va) + reg(vb) },
		addi: func() int { return reg(va) + vb },
		mulr: func() int { return reg(va) * reg(vb) },
		muli: func() int { return reg(va) * vb },
		banr: func() int { return reg(va) & reg(vb) },
		bani: func() int { return reg(va) & vb },
		borr: func() int { return reg(va) | reg(vb) },
		bori: func() int { return reg(va) | vb },
		setr: func() int { return reg(va) },
		seti: func() int { return va },
		gtir: func() int { return lib.If(va > reg(vb), 1, 0) },
		gtri: func() int { return lib.If(reg(va) > vb, 1, 0) },
		gtrr: func() int { return lib.If(reg(va) > reg(vb), 1, 0) },
		eqir: func() int { return lib.If(va == reg(vb), 1, 0) },
		eqri: func() int { return lib.If(reg(va) == vb, 1, 0) },
		eqrr: func() int { return lib.If(reg(va) == reg(vb), 1, 0) },
	} {
		func() {
			defer recover()
			c := fc() // panics on invalid register access

			for i, v := range s.after {
				if (i == s.in.c && v != c) || (i != s.in.c && v != s.before[i]) {
					return
				}

			}
			ops = append(ops, op)
		}()
	}
	return ops
}
