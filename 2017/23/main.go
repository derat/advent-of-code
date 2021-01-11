package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var ins []instr
	for _, ln := range lib.InputLines("2017/23") {
		ins = append(ins, newInstr(ln))
	}

	// Part 1:
	vm := newVM(ins)
	for !vm.oob {
		vm.tick()
	}
	fmt.Println(vm.nmul)
}

// This is just a hacked-up copy of 2017/18.
type vm struct {
	regs [26]int
	ins  []instr
	ip   int

	nmul int  // number of mul calls
	oob  bool // ip went out of bounds
}

func newVM(ins []instr) *vm {
	vm := &vm{ins: ins}
	return vm
}

func (vm *vm) get(b int, v int) int {
	if b >= 0 {
		return vm.regs[b]
	}
	return v
}

func (vm *vm) tick() {
	if vm.ip < 0 || vm.ip >= len(vm.ins) {
		vm.oob = true
		return
	}

	var jumped bool

	in := &vm.ins[vm.ip]
	switch in.op {
	case set:
		vm.regs[in.r1] = vm.get(in.r2, in.v2)
	case sub:
		vm.regs[in.r1] -= vm.get(in.r2, in.v2)
	case mul:
		vm.regs[in.r1] *= vm.get(in.r2, in.v2)
		vm.nmul++
	case jnz:
		if vm.get(in.r1, in.v1) != 0 {
			vm.ip += vm.get(in.r2, in.v2)
			jumped = true
		}
	default:
		lib.Panicf("Invalid op %d", in.op)
	}

	if !jumped {
		vm.ip++
	}
}

type op int

const (
	set op = iota
	sub
	mul
	jnz
)

type instr struct {
	op     op
	r1, r2 int
	v1, v2 int
}

func newInstr(ln string) instr {
	const re = `(?:([a-z])|(-?\d+))` // matches register or constant

	var op, r1, r2 string
	in := instr{r1: -1, r2: -1}
	switch {
	case lib.ExtractMaybe(ln, `^(set|sub|mul) ([a-z]) `+re+`$`, &op, &r1, &r2, &in.v2):
	case lib.ExtractMaybe(ln, `^(jnz) `+re+` `+re+`$`, &op, &r1, &in.v1, &r2, &in.v2):
	default:
		lib.Panicf("Bad instruction %q", ln)
	}

	switch op {
	case "set":
		in.op = set
	case "sub":
		in.op = sub
	case "mul":
		in.op = mul
	case "jnz":
		in.op = jnz
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
