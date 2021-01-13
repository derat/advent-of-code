package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var ins []lib.Instr
	for _, ln := range lib.InputLines("2017/18") {
		ins = append(ins, lib.NewInstr(ln, 'a', 'z', map[uint8]string{
			snd: `^snd ([a-z]|-?\d+)$`,
			set: `^set ([a-z]) ([a-z]|-?\d+)$`,
			add: `^add ([a-z]) ([a-z]|-?\d+)$`,
			mul: `^mul ([a-z]) ([a-z]|-?\d+)$`,
			mod: `^mod ([a-z]) ([a-z]|-?\d+)$`,
			rcv: `^rcv ([a-z]|-?\d+)$`,
			jgz: `^jgz ([a-z]|-?\d+) ([a-z]|-?\d+)$`,
		}))
	}

	// Part 1: Print last snd value at time of first rcv.
	var write []int64
	vm := newVM(ins, &write, nil, -1)
	for !vm.blocked {
		vm.tick()
	}
	fmt.Println((write)[len(write)-1])

	// Part 2: Run two programs in parallel until both stop
	// and then print number of snds in program 1.
	var q0, q1 []int64 // named after writer
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
	regs []int64
	ins  []lib.Instr
	ip   int

	write, read *[]int64 // snd and rcv queues

	id      int64 // program ID; -1 for part 1
	nsnd    int   // number of snd calls
	oob     bool  // ip went out of bounds
	blocked bool  // waiting on rcv
}

func newVM(ins []lib.Instr, write, read *[]int64, id int64) *vm {
	vm := &vm{
		regs:  make([]int64, 26),
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

func (vm *vm) tick() {
	if vm.ip < 0 || vm.ip >= len(vm.ins) {
		vm.oob = true
		return
	}

	var jumped bool

	in := &vm.ins[vm.ip]
	switch in.Op {
	case snd:
		*vm.write = append(*vm.write, in.Val(0, vm.regs))
		vm.nsnd++
	case set:
		*in.Ptr(0, vm.regs) = in.Val(1, vm.regs)
	case add:
		*in.Ptr(0, vm.regs) += in.Val(1, vm.regs)
	case mul:
		*in.Ptr(0, vm.regs) *= in.Val(1, vm.regs)
	case mod:
		*in.Ptr(0, vm.regs) %= in.Val(1, vm.regs)
	case rcv:
		if vm.id < 0 {
			// Part 1: Just block on the first rcv with a nonzero value.
			if in.Val(0, vm.regs) != 0 {
				vm.blocked = true
				return
			}
		} else {
			// Part 2: Read from the queue.
			if len(*vm.read) == 0 {
				vm.blocked = true
				return
			}
			*in.Ptr(0, vm.regs) = (*vm.read)[0]
			*vm.read = (*vm.read)[1:]
		}
	case jgz:
		if in.Val(0, vm.regs) > 0 {
			vm.ip += int(in.Val(1, vm.regs))
			jumped = true
		}
	default:
		lib.Panicf("Invalid op %d", in.Op)
	}

	if !jumped {
		vm.ip++
	}
}

const (
	snd = iota
	set
	add
	mul
	mod
	rcv
	jgz
)
