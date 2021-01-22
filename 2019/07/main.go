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

// vm runs Intcode instructions.
type vm struct {
	mem     map[int]int
	in, out chan int
	done    chan bool
}

// newVM returns a new vm with a copy of the supplied initial memory.
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

	var modeDiv = []int{100, 1000, 10000}

	var ip int // instruction start index
	var op int // opcode (including mode)
	var sz int // instruction size (including opcode)

	// Gets the (mode-appropriate) 1-indexed argument.
	get := func(arg int) int {
		lib.Assert(arg > 0)
		sz = lib.Max(arg+1, sz)
		v, ok := vm.mem[ip+arg]
		lib.Assertf(ok, "Bad read %v", ip+arg)

		mode := (op / modeDiv[arg-1]) % 10
		switch mode {
		case 0: // position mode: address
			vp, ok := vm.mem[v]
			lib.Assertf(ok, "Bad read %v", v)
			return vp
		case 1: // immediate mode: literal value
			return v
		default:
			lib.Panicf("Invalid mode %d", mode)
		}
		return 0 // unreached
	}

	// Sets the 1-indexed argument to the supplied value.
	set := func(arg, val int) {
		lib.Assert(arg > 0)
		sz = lib.Max(arg+1, sz)
		addr, ok := vm.mem[ip+arg] // always treated as an address
		lib.Assertf(ok, "Bad read %v", ip+arg)
		vm.mem[addr] = val
	}

	for {
		var ok bool
		op, ok = vm.mem[ip]
		lib.Assertf(ok, "Bad ip %v", ip)
		sz = 1

		switch op % modeDiv[0] {
		case 1: // add first two args and save to third arg
			set(3, get(1)+get(2))
		case 2: // multiply first two args and save to third arg
			set(3, get(1)*get(2))
		case 3: // read input and save to first arg
			set(1, <-vm.in)
		case 4: // write first arg to output
			vm.out <- get(1)
		case 5: // jump to second arg if first arg is nonzero
			if addr := get(2); get(1) != 0 {
				ip = addr
				sz = 0 // don't advance ip
			}
		case 6: // jump to second arg if first arg is zero
			if addr := get(2); get(1) == 0 {
				ip = addr
				sz = 0 // don't advance ip
			}
		case 7: // store 1 in third arg if first is less than second
			set(3, lib.If(get(1) < get(2), 1, 0))
		case 8: // store 1 in third arg if first is equal to second
			set(3, lib.If(get(1) == get(2), 1, 0))
		case 99: // halt the program
			return
		default:
			lib.Panicf("Invalid op %d", op%modeDiv[0])
		}

		ip += sz
	}
}
