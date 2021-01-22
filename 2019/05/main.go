package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInts("2019/5")

	// Part 1: Run diagnostic program with input of 1 and print diagnostic code
	// (the final output after a stream of 0s).
	vm := newVM(input)
	vm.in <- 1
	ch := make(chan bool)
	go func() {
		ch <- vm.run()
	}()
	var last int
	for v := range vm.out {
		last = v
	}
	lib.Assert(<-ch)
	fmt.Println(last)
}

type vm struct {
	mem     map[int]int
	in, out chan int
}

func newVM(init []int) *vm {
	vm := &vm{
		mem: make(map[int]int, len(init)),
		in:  make(chan int, 1),
		out: make(chan int, 1),
	}
	for addr, val := range init {
		vm.mem[addr] = val
	}
	return vm
}

func (vm *vm) get(addr, mode int) int {
	switch mode {
	case imm:
		return addr
	case pos:
		val, ok := vm.mem[addr]
		if !ok {
			panic("Invalid read")
		}
		return val
	default:
		panic("Invalid mode")
	}
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
		close(vm.out)
	}()

	var ip int
	var sz int // size of current instruction (including op)

	// These functions read the specified number of parameters following ip.
	params1 := func() int {
		sz = 2
		return vm.get(ip+1, pos)
	}
	params3 := func() (int, int, int) {
		sz = 4
		return vm.get(ip+1, pos), vm.get(ip+2, pos), vm.get(ip+3, pos)
	}

	for {
		sz = 1 // number of consumed ints (including op)
		in := vm.get(ip, pos)
		op := in % 100
		modes := []int{
			(in / 100) % 10,
			(in / 1000) % 10,
			(in / 10000) % 10,
		}

		switch op {
		case add:
			a0, a1, a2 := params3()
			vm.set(a2, vm.get(a0, modes[0])+vm.get(a1, modes[1]))
		case mul:
			a0, a1, a2 := params3()
			vm.set(a2, vm.get(a0, modes[0])*vm.get(a1, modes[1]))
		case inp:
			vm.set(params1(), <-vm.in)
		case out:
			vm.out <- vm.get(params1(), modes[0])
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
	inp = 3
	out = 4
	hlt = 99
)

const (
	pos = 0
	imm = 1
)
