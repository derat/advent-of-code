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

	// Part 1: Print last snd value at time of first rcv.
	var write []int
	vm := newVM(ins, &write, nil, -1)
	for !vm.blocked {
		vm.tick()
	}
	fmt.Println((write)[len(write)-1])

	// Part 2: Run two programs in parallel until both stop
	// and then print number of snds in program 1.
	var q0, q1 []int // named after writer
	vm0 := newVM(ins, &q0, &q1, 0)
	vm1 := newVM(ins, &q1, &q0, 1)
	for vm0.runnable() || vm1.runnable() {
		if vm0.runnable() {
			vm0.tick()
		}
		if vm1.runnable() {
			vm1.tick()
		}
	}
	fmt.Println(vm1.nsnd)
}

type vm struct {
	regs [26]int
	ins  []instr
	ip   int

	write, read *[]int // snd and rcv queues

	id      int  // program ID; -1 for part 1
	nsnd    int  // number of snd calls
	oob     bool // ip went out of bounds
	blocked bool // waiting on rcv
}

func newVM(ins []instr, write, read *[]int, id int) *vm {
	vm := &vm{
		ins:   ins,
		write: write,
		read:  read,
		id:    id,
	}
	if id >= 0 {
		vm.regs['p'-'a'] = id
	}
	return vm
}

func (vm *vm) runnable() bool {
	return !vm.oob && (!vm.blocked || len(*vm.read) > 0)
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
	case snd:
		*vm.write = append(*vm.write, vm.get(in.r1, in.v1))
		vm.nsnd++
	case set:
		vm.regs[in.r1] = vm.get(in.r2, in.v2)
	case add:
		vm.regs[in.r1] += vm.get(in.r2, in.v2)
	case mul:
		vm.regs[in.r1] *= vm.get(in.r2, in.v2)
	case mod:
		vm.regs[in.r1] %= vm.get(in.r2, in.v2)
	case rcv:
		if vm.id < 0 {
			// Part 1: Just block on the first rcv with a nonzero value.
			if vm.get(in.r1, in.v1) != 0 {
				vm.blocked = true
				return
			}
		} else {
			// Part 2: Read from the queue.
			if len(*vm.read) == 0 {
				vm.blocked = true
				return
			}
			vm.regs[in.r1] = (*vm.read)[0]
			*vm.read = (*vm.read)[1:]
		}
	case jgz:
		if vm.get(in.r1, in.v1) > 0 {
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

// I'm very tired of writing code like this and should try to add a
// library function the next time that this comes up.
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
