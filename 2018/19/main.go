package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	lines := lib.InputLines("2018/19")

	var ipr int
	lib.Extract(lines[0], `^#ip (\d+)$`, &ipr)

	var ins []instr
	for _, ln := range lines[1:] {
		var op string
		var in instr
		lib.Extract(ln, `^(\w+) (\d+) (\d+) (\d+)$`, &op, &in.a, &in.b, &in.c)
		var ok bool
		if in.op, ok = strOp[op]; !ok {
			lib.Panicf("Invalid op in %q", ln)
		}
		ins = append(ins, in)
	}

	// Part 1: Run the program and print the final value of register 0.
	var ip int
	var regs [6]int
	reg := func(i int) int { return regs[i] }
	for ip >= 0 && ip < len(ins) {
		regs[ipr] = ip
		in := &ins[ip]
		regs[in.c] = in.op.run(in.a, in.b, reg)
		ip = regs[ipr]
		ip++
	}
	fmt.Println(regs[0])
}

type instr struct {
	op      op
	a, b, c int
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

var strOp = map[string]op{
	"addr": addr,
	"addi": addi,
	"mulr": mulr,
	"muli": muli,
	"banr": banr,
	"bani": bani,
	"borr": borr,
	"bori": bori,
	"setr": setr,
	"seti": seti,
	"gtir": gtir,
	"gtri": gtri,
	"gtrr": gtrr,
	"eqir": eqir,
	"eqri": eqri,
	"eqrr": eqrr,
}
