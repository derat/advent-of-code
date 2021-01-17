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

	// Part 2: Determine number of each op and print register 0 after executing program.
	unk := make(map[op]map[int]struct{})  // possible opcodes for each op
	runk := make(map[int]map[op]struct{}) // possible ops for each opcode
	for _, sa := range samples {
		code := sa.in.op
		for _, o := range sa.ops() {
			unk[o] = lib.AddSet(unk[o], code).(map[int]struct{})
			runk[code] = lib.AddSet(runk[code], o).(map[op]struct{})
		}
	}

	// Identify ops that start out with a single possible code.
	var next []op
	for o, codes := range unk {
		if len(codes) == 1 {
			next = append(next, o)
		}
	}

	// Repeatedly match ops with codes until there are none left.
	ops := make(map[int]op) // codes to ops
	for len(next) > 0 {
		o := next[0]
		next = next[1:]

		if _, ok := unk[o]; !ok {
			continue // already handled
		}

		lib.AssertEq(len(unk[o]), 1)
		code := lib.MapSomeKey(unk[o]).(int)
		ops[code] = o

		delete(unk, o)
		for bo := range runk[code] {
			delete(unk[bo], code)
			if len(unk[bo]) == 1 {
				next = append(next, bo)
			}
		}
		delete(runk, code)
	}
	lib.AssertEq(len(ops), nops)

	// Run the program.
	var regs [4]int
	reg := func(i int) int { return regs[i] }
	for _, in := range prog {
		op, ok := ops[in.op]
		lib.Assertf(ok, "Invalid opcode %d", in.op)
		regs[in.c] = op.run(in.a, in.b, reg)
	}
	fmt.Println(regs[0])
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

	nops = int(eqrr) + 1
)

// run returns the value to store in c when o is called
// with the supplied a and b values and register state.
func (o op) run(va, vb int, reg func(int) int) int {
	switch o {
	case addr:
		return reg(va) + reg(vb)
	case addi:
		return reg(va) + vb
	case mulr:
		return reg(va) * reg(vb)
	case muli:
		return reg(va) * vb
	case banr:
		return reg(va) & reg(vb)
	case bani:
		return reg(va) & vb
	case borr:
		return reg(va) | reg(vb)
	case bori:
		return reg(va) | vb
	case setr:
		return reg(va)
	case seti:
		return va
	case gtir:
		return lib.If(va > reg(vb), 1, 0)
	case gtri:
		return lib.If(reg(va) > vb, 1, 0)
	case gtrr:
		return lib.If(reg(va) > reg(vb), 1, 0)
	case eqir:
		return lib.If(va == reg(vb), 1, 0)
	case eqri:
		return lib.If(reg(va) == vb, 1, 0)
	case eqrr:
		return lib.If(reg(va) == reg(vb), 1, 0)
	default:
		panic(fmt.Sprintf("Invalid opcode %d", o))
	}
}

type sample struct {
	before, after [4]int
	in            instr
}

func (s *sample) ops() []op {
	reg := func(idx int) int {
		if idx < 0 || idx >= len(s.before) {
			panic(fmt.Sprintf("Invalid register access %d", idx))
		}
		return s.before[idx]
	}

	var ops []op
	for i := 0; i < int(nops); i++ {
		op := op(i)
		func() {
			defer recover()
			c := op.run(s.in.a, s.in.b, reg) // panics on invalid register access

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
