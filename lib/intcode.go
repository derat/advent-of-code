// Copyright 2022 Daniel Erat.
// All rights reserved.

package lib

// Intcode runs Intcode instructions.
type Intcode struct {
	Mem     map[int64]int64
	In, Out chan int64
	InFunc  func() int64 // used instead of In if non-nil
	OutFunc func(int64)  // used instead of Out if non-nil
	done    chan bool
	halt    bool
}

// NewIntcode returns a new Intcode VM with a copy of the supplied initial memory.
func NewIntcode(init []int64) *Intcode {
	vm := &Intcode{
		Mem: make(map[int64]int64, len(init)),
		In:  make(chan int64),
		Out: make(chan int64),
	}
	for addr, val := range init {
		vm.Mem[int64(addr)] = val
	}
	return vm
}

// Start starts the VM in a goroutine.
// Wait can be used to wait for the VM to halt.
func (vm *Intcode) Start() {
	Assertf(vm.done == nil, "Already running")
	vm.done = make(chan bool, 1)
	go func() {
		vm.done <- vm.Run()
		close(vm.done)
	}()
}

// Wait waits for the previously-started VM to halt.
// Its return value is the same as that of Run.
func (vm *Intcode) Wait() bool {
	Assertf(vm.done != nil, "Not started")
	return <-vm.done
}

// Halt makes the VM exit with success before running the next instruction.
func (vm *Intcode) Halt() {
	vm.halt = true
}

// Run synchronously runs the VM to completion.
// It returns true if hlt was executed and false if the VM crashed due
// to an invalid opcode or bad memory access.
func (vm *Intcode) Run() (halted bool) {
	defer func() {
		if r := recover(); r == nil {
			halted = true
		}
		close(vm.Out)
	}()

	var modeDiv = []int64{100, 1000, 10000}

	var ip int64 // instruction start index
	var rb int64 // relative base
	var op int64 // opcode (including mode)
	var sz int   // instruction size (including opcode)

	// Gets the (mode-appropriate) 1-indexed argument.
	get := func(arg int) int64 {
		Assert(arg > 0)
		sz = Max(arg+1, sz)
		addr := ip + int64(arg)
		Assertf(addr >= 0, "Bad arg addr %v (ip %v)", addr, ip)
		v := vm.Mem[addr]

		mode := (op / modeDiv[arg-1]) % 10
		switch mode {
		case 0: // position mode: address
			Assertf(v >= 0, "Bad pos read addr %v", v)
			return vm.Mem[v]
		case 1: // immediate mode: literal value
			return v
		case 2: // relative mode: offset from relative base
			addr := rb + v
			Assertf(addr >= 0, "Bad rel read addr %v (base %v, offset %v)", addr, rb, v)
			return vm.Mem[addr]
		default:
			Panicf("Invalid mode %d", mode)
		}
		return 0 // unreached
	}

	// Sets the 1-indexed argument to the supplied value.
	set := func(arg int, val int64) {
		Assert(arg > 0)
		sz = Max(arg+1, sz)
		aaddr := ip + int64(arg)
		Assertf(aaddr >= 0, "Bad arg addr %v (ip %v)", aaddr, ip)
		saddr := vm.Mem[aaddr]

		mode := (op / modeDiv[arg-1]) % 10
		switch mode {
		case 0: // position mode: address
		case 2: // relative mode: offset from relative base
			saddr += rb
		default:
			Panicf("Invalid mode %d", mode)
		}

		Assertf(saddr >= 0, "Bad set addr %v (ip %v)", saddr, ip)
		vm.Mem[saddr] = val
	}

	for {
		// Handle requests from outside to stop running.
		if vm.halt {
			return
		}

		Assertf(ip >= 0, "Bad ip %v", ip)
		op = vm.Mem[ip]
		sz = 1

		switch op % modeDiv[0] {
		case 1: // add first two args and save to third arg
			set(3, get(1)+get(2))
		case 2: // multiply first two args and save to third arg
			set(3, get(1)*get(2))
		case 3: // read input and save to first arg
			var val int64
			if vm.InFunc != nil {
				val = vm.InFunc()
			} else {
				val = <-vm.In
			}
			set(1, val)
		case 4: // write first arg to output
			val := get(1)
			if vm.OutFunc != nil {
				vm.OutFunc(val)
			} else {
				vm.Out <- val
			}
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
			set(3, int64(If(get(1) < get(2), 1, 0)))
		case 8: // store 1 in third arg if first is equal to second
			set(3, int64(If(get(1) == get(2), 1, 0)))
		case 9: // adjust relative base by first arg
			rb += get(1)
		case 99: // halt the program
			return
		default:
			Panicf("Invalid op %d", op%modeDiv[0])
		}

		ip += int64(sz)
	}
}
