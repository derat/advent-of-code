package lib

// Intcode runs Intcode instructions.
type Intcode struct {
	Mem     map[int]int
	In, Out chan int
	done    chan bool
}

// NewIntcode returns a new Intcode VM with a copy of the supplied initial memory.
func NewIntcode(init []int) *Intcode {
	vm := &Intcode{
		Mem: make(map[int]int, len(init)),
		In:  make(chan int, 1),
		Out: make(chan int, 1),
	}
	for addr, val := range init {
		vm.Mem[addr] = val
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

	var modeDiv = []int{100, 1000, 10000}

	var ip int // instruction start index
	var op int // opcode (including mode)
	var sz int // instruction size (including opcode)

	// Gets the (mode-appropriate) 1-indexed argument.
	get := func(arg int) int {
		Assert(arg > 0)
		sz = Max(arg+1, sz)
		v, ok := vm.Mem[ip+arg]
		Assertf(ok, "Bad read %v", ip+arg)

		mode := (op / modeDiv[arg-1]) % 10
		switch mode {
		case 0: // position mode: address
			vp, ok := vm.Mem[v]
			Assertf(ok, "Bad read %v", v)
			return vp
		case 1: // immediate mode: literal value
			return v
		default:
			Panicf("Invalid mode %d", mode)
		}
		return 0 // unreached
	}

	// Sets the 1-indexed argument to the supplied value.
	set := func(arg, val int) {
		Assert(arg > 0)
		sz = Max(arg+1, sz)
		addr, ok := vm.Mem[ip+arg] // always treated as an address
		Assertf(ok, "Bad read %v", ip+arg)
		vm.Mem[addr] = val
	}

	for {
		var ok bool
		op, ok = vm.Mem[ip]
		Assertf(ok, "Bad ip %v", ip)
		sz = 1

		switch op % modeDiv[0] {
		case 1: // add first two args and save to third arg
			set(3, get(1)+get(2))
		case 2: // multiply first two args and save to third arg
			set(3, get(1)*get(2))
		case 3: // read input and save to first arg
			set(1, <-vm.In)
		case 4: // write first arg to output
			vm.Out <- get(1)
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
			set(3, If(get(1) < get(2), 1, 0))
		case 8: // store 1 in third arg if first is equal to second
			set(3, If(get(1) == get(2), 1, 0))
		case 99: // halt the program
			return
		default:
			Panicf("Invalid op %d", op%modeDiv[0])
		}

		ip += sz
	}
}
