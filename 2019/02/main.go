package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInts("2019/2")

	// Part 1: Set pos 1 to 12 and pos 2 to 2, and print pos 0 after halt.
	vm1 := newVM(input)
	vm1.set(1, 12)
	vm1.set(2, 2)
	lib.Assert(vm1.run())
	fmt.Println(vm1.get(0))

	// Part 2: Find noun (pos 1) and verb (pos 2) values in [0,99] that result
	// in pos 0 being 19690720.
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			vm2 := newVM(input)
			vm2.set(1, noun)
			vm2.set(2, verb)
			if !vm2.run() {
				continue
			}
			if vm2.get(0) == 19690720 {
				fmt.Println(100*noun + verb)
				break
			}
		}
	}
}

type vm struct {
	// I'm guessing that we'll need to support writing outside the initial program's
	// address space, and that we may need a sparse representation of memory. Who
	// knows if I'll be right, though...
	mem map[int]int
}

func newVM(init []int) *vm {
	vm := &vm{
		mem: make(map[int]int, len(init)),
	}
	for addr, val := range init {
		vm.mem[addr] = val
	}
	return vm
}

func (vm *vm) get(addr int) int {
	val, ok := vm.mem[addr]
	if !ok {
		panic("Invalid read")
	}
	return val
}

func (vm *vm) set(addr, val int) {
	vm.mem[addr] = val
}

// run runs the VM to completion. It returns true if hlt was executed
// and false if the VM crashed due to an invalid opcode or bad memory
// access.
func (vm *vm) run() (halted bool) {
	defer func() {
		if r := recover(); r == nil {
			halted = true
		}
	}()

	var ip int
	var sz int // size of current instruction (including op)

	// params3 reads three params following ip.
	params3 := func() (int, int, int) {
		sz = 4
		return vm.get(ip + 1), vm.get(ip + 2), vm.get(ip + 3)
	}

	for {
		sz = 1
		switch vm.mem[ip] {
		case add:
			in1, in2, out := params3()
			vm.set(out, vm.get(in1)+vm.get(in2))
		case mul:
			in1, in2, out := params3()
			vm.set(out, vm.get(in1)*vm.get(in2))
		case hlt:
			return
		default:
			panic("Invalid op")
		}
		ip += sz
	}
}

const (
	add = 1
	mul = 2
	hlt = 99
)
