package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInts("2019/7")

	// Part 1: Find maximum output signal.
	fmt.Println(run(input, []int{0, 1, 2, 3, 4}, false))

	// Part 2: Feed e's output into a's input.
	fmt.Println(run(input, []int{5, 6, 7, 8, 9}, true))
}

func run(input, vals []int, feedback bool) int {
	pch := make(chan []int)
	go lib.Perms(vals, pch)

	// The problem just calls for running each amplifier in parallel
	// and manually copying values, but this is Go, so why not use
	// goroutines and channels?
	var max int
	for phases := range pch {
		amps := make([]*vm, len(phases))
		for i := range amps {
			amps[i] = newVM(input)
		}

		// Wire up the amplifiers.
		for i, a := range amps {
			if i == 0 {
				if feedback {
					a.in = amps[len(amps)-1].out
				}
			} else {
				a.in = amps[i-1].out
			}
		}

		// Start the amplifiers and feed them their phase signals.
		for i, a := range amps {
			a.start()
			a.in <- phases[i]
		}

		// Send the input signal to the first and read the output from the last.
		amps[0].in <- 0
		lib.Assert(amps[0].wait())
		if out := <-amps[len(amps)-1].out; out > max {
			max = out
		}
	}
	return max
}

type vm struct {
	mem     map[int]int
	in, out chan int
	done    chan bool
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

// start starts the VM in a goroutine.
// wait can be used to wait for the VM to halt.
func (vm *vm) start() {
	lib.Assertf(vm.done == nil, "Already running")
	vm.done = make(chan bool, 1)
	go func() {
		vm.done <- vm.run()
		close(vm.done)
	}()
}

// wait waits for the previously-started VM to halt.
// Its return value is the same as that of run.
func (vm *vm) wait() bool {
	lib.Assertf(vm.done != nil, "Not started")
	return <-vm.done
}

// run synchronously runs the VM to completion.
// It returns true if hlt was executed and false if the VM crashed due
// to an invalid opcode or bad memory access.
func (vm *vm) run() (halted bool) {
	defer func() {
		if r := recover(); r == nil {
			halted = true
		}
		close(vm.out)
	}()

	var ip int // index of current instruction
	var sz int // size of current instruction (including op)

	// These functions read the specified number of parameters following ip.
	params1 := func() int {
		sz = 2
		return vm.get(ip+1, pos)
	}
	params2 := func() (int, int) {
		sz = 3
		return vm.get(ip+1, pos), vm.get(ip+2, pos)
	}
	params3 := func() (int, int, int) {
		sz = 4
		return vm.get(ip+1, pos), vm.get(ip+2, pos), vm.get(ip+3, pos)
	}

	for {
		sz = 1 // number of consumed ints (including op)
		in := vm.get(ip, pos)
		op := in % 100
		m := []int{ // mode for first, second, third arg
			(in / 100) % 10,
			(in / 1000) % 10,
			(in / 10000) % 10,
		}

		switch op {
		case add:
			a0, a1, a2 := params3()
			vm.set(a2, vm.get(a0, m[0])+vm.get(a1, m[1]))
		case mul:
			a0, a1, a2 := params3()
			vm.set(a2, vm.get(a0, m[0])*vm.get(a1, m[1]))
		case inp:
			vm.set(params1(), <-vm.in)
		case out:
			vm.out <- vm.get(params1(), m[0])
		case jit:
			a0, a1 := params2()
			if v := vm.get(a0, m[0]); v != 0 {
				ip = vm.get(a1, m[1])
				sz = 0
			}
		case jif:
			a0, a1 := params2()
			if v := vm.get(a0, m[0]); v == 0 {
				ip = vm.get(a1, m[1])
				sz = 0
			}
		case slt:
			a0, a1, a2 := params3()
			val := lib.If(vm.get(a0, m[0]) < vm.get(a1, m[1]), 1, 0)
			vm.set(a2, val)
		case seq:
			a0, a1, a2 := params3()
			val := lib.If(vm.get(a0, m[0]) == vm.get(a1, m[1]), 1, 0)
			vm.set(a2, val)
		case hlt:
			return
		default:
			panic("Invalid op")
		}
		ip += sz
	}
}

const (
	add = 1  // add two args and save to third arg
	mul = 2  // multiply two args and save to third
	inp = 3  // read input and save to arg
	out = 4  // write arg to output
	jit = 5  // if first arg is non-zero, set ip to second arg
	jif = 6  // if first arg is zero, set ip to second arg
	slt = 7  // if first arg is less than second, store 1 in third; otherwise 0
	seq = 8  // if first arg is equal to second, store 1 in third; otherwise 0
	hlt = 99 // stop the program
)

const (
	pos = 0 // position mode
	imm = 1 // immediate mode
)
