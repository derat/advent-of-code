package main

import (
	"fmt"
	"math"

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

	// Part 2:
	//
	// Here's my input (with instruction numbers added):
	//
	//   0 set b 65
	//   1 set c b
	//   2 jnz a 2
	//   3 jnz 1 5
	//   4 mul b 100
	//   5 sub b -100000
	//   6 set c b
	//   7 sub c -17000
	//   8 set f 1
	//   9 set d 2
	//  10 set e 2
	//  11 set g d
	//  12 mul g e
	//  13 sub g b
	//  14 jnz g 2
	//  15 set f 0
	//  16 sub e -1
	//  17 set g e
	//  18 sub g b
	//  19 jnz g -8
	//  20 sub d -1
	//  21 set g d
	//  22 sub g b
	//  23 jnz g -13
	//  24 jnz f 2
	//  25 sub h -1
	//  26 set g b
	//  27 sub g c
	//  28 jnz g 2
	//  29 jnz 1 3
	//  30 sub b -17
	//  31 jnz 1 -23
	//
	// Here's the whole thing manually converted to pseudocode:
	//
	//  b = 65                // 0
	//  c = 65                // 1
	//  if a != 0 {           // 2
	//    b = 106500          // 4-5
	//    c = 123500          // 6-7
	//  }
	//
	//  for {                 // 8-31
	//    for {               // 8-23
	//      f = 1             // 8
	//      d = 2             // 9
	//      for {             // 10-23
	//        e = 2           // 10
	//        for {           // 11-19
	//          if d*e == b { // 11-14
	//            f = 0       // 15
	//          }
	//          e++           // 16
	//          if e == b {   // 17-19
	//            break
	//          }
	//        }
	//        d++             // 20
	//        if d == b {     // 21-23
	//          break
	//        }
	//      }
	//    }
	//
	//    if f == 0 {         // 24
	//      h++               // 25
	//    }
	//    if b == c {         // 26-28
	//      break             // 29
	//    }
	//    b += 17             // 30
	//  }
	//
	// Note that the 'g' register is only used for comparing values.
	// More tersely in Go:
	//
	// Loop:
	//	for b := 106500; b <= 123500; b += 17 {
	//		f := false
	//		for d := 2; d < b; d++ {
	//			for e := 2; e < b; e++ {
	//				if d*e == b {
	//					h++
	//					continue Loop
	//				}
	//			}
	//		}
	//	}
	//
	// So we're counting the number of primes between 106500 and 123500,
	// using a step of 17. Below, I just manually skip the 8-23 loop when we
	// hit instruction 8 and manually set register 'h' to 0 if 'b' is prime
	// and to 1 otherwise.

	prime := func(n int) bool {
		root := int(math.Sqrt(float64(n)))
		for i := 2; i <= root; i++ {
			if n%i == 0 {
				return true
			}
		}
		return false
	}

	vm = newVM(ins)
	vm.regs[0] = 1
	for !vm.oob {
		if vm.ip == 8 {
			vm.regs[5] = lib.If(prime(vm.regs[1]), 0, 1)
			vm.ip = 24
			continue
		}
		vm.tick()
	}
	fmt.Println(vm.regs[7])
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
