package main

import (
	"fmt"
	"sync"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInt64s("2019/23")

	const ncomps = 50

	// Part 1: Print Y address of first packet sent to address 255.
	var vms []*lib.Intcode
	done := make(chan struct{}, 1)
	bus := newBus(ncomps)
	bus.snoop = func(addr int, x, y int64) {
		if addr == 255 {
			fmt.Println(y)
			bus.snoop = nil
			close(done)
		}
	}
	for addr := 0; addr < ncomps; addr++ {
		vm := makeComp(input, addr, bus)
		vm.Start()
		vms = append(vms, vm)
	}
	<-done
	for _, vm := range vms {
		vm.Halt()
	}

	// Part 2: Print first Y value delivered by NAT to address 0 twice in a row.
	// This code is complete garbage.
	vms = nil
	done = make(chan struct{}, 1)
	bus = newBus(ncomps)

	var nat [2]int64 // last packet sent to address 255
	var haveNat bool
	bus.snoop = func(addr int, x, y int64) {
		if addr == 255 {
			nat = [2]int64{x, y}
			haveNat = true
		}
	}
	var lastY int64
	bcast := make(chan [2]int64)
	bus.idle = func() {
		if haveNat {
			bcast <- nat
			if nat[1] == lastY {
				fmt.Println(nat[1])
				close(bcast)
				close(done)
			}
			lastY = nat[1]
			haveNat = false
		}
	}
	go func() {
		for pkg := range bcast {
			bus.send(0, pkg[0], pkg[1], 0)
		}
	}()

	for addr := 0; addr < ncomps; addr++ {
		vm := makeComp(input, addr, bus)
		vm.Start()
		vms = append(vms, vm)
	}
	<-done
	for _, vm := range vms {
		vm.Halt()
	}
}

func makeComp(prog []int64, addr int, bus *bus) *lib.Intcode {
	vm := lib.NewIntcode(prog)

	inState := boot
	var inY int64 // Y value from packet being read
	vm.InFunc = func() int64 {
		switch inState {
		case boot:
			inState = idle
			return int64(addr)
		case idle:
			x, y, ok := bus.recv(addr)
			if !ok { // no queued packets
				return -1
			}
			inState = reading
			inY = y
			return x
		case reading:
			inState = idle
			return inY
		default:
			lib.Panicf("Invalid input state %v", inState)
		}
		return -1
	}

	outState := idle
	var outAddr int
	var outX int64
	vm.OutFunc = func(v int64) {
		switch outState {
		case idle:
			outAddr = int(v)
			outState = writeX
		case writeX:
			outX = v
			outState = writeY
		case writeY:
			bus.send(outAddr, outX, v, addr)
			outState = idle
		default:
			lib.Panicf("Invalid output state %v", inState)
		}
	}

	return vm
}

type state int

const (
	boot state = iota
	idle
	reading // reading packet: sent X, next is Y
	writeX  // writing packet: waiting for X
	writeY  // writing packet: waiting for Y
)

type bus struct {
	packets [][][2]int64
	active  uint64
	mu      sync.Mutex
	snoop   func(int, int64, int64)
	idle    func()
}

func newBus(ncomps int) *bus {
	return &bus{
		packets: make([][][2]int64, ncomps),
		active:  (1 << ncomps) - 1,
	}
}

func (b *bus) send(addr int, x, y int64, self int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.active |= 1 << self
	if b.snoop != nil {
		b.snoop(addr, x, y)
	}
	if addr < len(b.packets) {
		b.active |= 1 << addr
		b.packets[addr] = append(b.packets[addr], [2]int64{x, y})
	}
}

func (b *bus) recv(addr int) (x, y int64, ok bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.packets[addr]) == 0 {
		b.active &= ^(1 << addr)
		if b.active == 0 {
			b.idle()
		}
		return 0, 0, false
	}

	b.active |= 1 << addr
	x, y = b.packets[addr][0][0], b.packets[addr][0][1]
	b.packets[addr] = b.packets[addr][1:]
	return x, y, true
}
