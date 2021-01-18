package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	lines := lib.InputLines("2018/19")

	var ipr int
	lib.Extract(lines[0], `^#ip (\d+)$`, &ipr)

	var ins []instr
	for _, ln := range lines[1:] {
		var op string
		var in instr
		lib.Extract(ln, `^(\w+) (\d+) (\d+) (\d+)$`, &op, &in.a, &in.b, &in.c)
		var ok bool
		if in.op, ok = strOp[op]; !ok {
			lib.Panicf("Invalid op in %q", ln)
		}
		ins = append(ins, in)
	}

	// Part 1: Run the program and print the final value of register 0.
	fmt.Println(run(ins, ipr, 0)[0])

	// Part 2: This time, initialize register 0 to 1.
	//
	// My (annotated) input:
	//
	//     #ip 4
	//   0 addi 4 16 4  # jump to 17
	//   1 seti 1  7 2  # r2 = 1
	//   2 seti 1  1 5  # r5 = 1
	//   3 mulr 2  5 3  # r3 = r2*r5
	//   4 eqrr 3  1 3  # r3 = 1 (r3 == r1) or 0 (r3 != r1)
	//   5 addr 3  4 4  # jump to 6 + r3
	//   6 addi 4  1 4  # jump to 8
	//   7 addr 2  0 0  # r0 += r2
	//   8 addi 5  1 5  # r5 += 1
	//   9 gtrr 5  1 3  # r3 = 1 (r5 > r1) or 0 (r5 <= r1)
	//  10 addr 4  3 4  # jump to 11 + r3
	//  11 seti 2  7 4  # jump to 3
	//  12 addi 2  1 2  # r2 += 1
	//  13 gtrr 2  1 3  # r3 = 1 (r2 > r1) or 0 (r2 <= r1)
	//  14 addr 3  4 4  # jump to 15 + r3
	//  15 seti 1  3 4  # jump to 2
	//  16 mulr 4  4 4  # jump to 257 (exit)
	//  17 addi 1  2 1  # r1 += 2
	//  18 mulr 1  1 1  # r1 *= r1
	//  19 mulr 4  1 1  # r1 *= 19
	//  20 muli 1 11 1  # r1 *= 11
	//  21 addi 3  3 3  # r3 += 3
	//  22 mulr 3  4 3  # r3 *= 22
	//  23 addi 3  9 3  # r3 += 9
	//  24 addr 1  3 1  # r1 += r3
	//  25 addr 4  0 4  # jump to 26 + r0
	//  26 seti 0  1 4  # jump to 1
	//  27 setr 4  9 3  # r3 = 27
	//  28 mulr 3  4 3  # r3 *= 28
	//  29 addr 4  3 3  # r3 += 29
	//  30 mulr 4  3 3  # r3 *= 30
	//  31 muli 3 14 3  # r3 *= 3
	//  32 mulr 3  4 3  # r3 *= 32
	//  33 addr 1  3 1  # r1 += r3
	//  34 seti 0  6 0  # r0 = 0
	//  35 seti 0  7 4  # jump to 1
	//
	//  r1 = 911            // 17-24
	//  if r0 == 1 {        // 25-26
	//    r1 += 2260800     // 27-33
	//    r0 = 0            // 34
	//  }
	//  r2 = 1              // 1
	//  for r2 <= r1 {      // 13-15
	//    r5 = 1            // 2
	//    for r5 <= r1 {    // 9-11
	//      r3 = r2 * r5    // 3
	//      if r3 == r1 {   // 4-6
	//        r0 += r2      // 7
	//      }
	//      r5 += 1         // 8
	//    }
	//    r2 += 1           // 12
	//  }
	fmt.Println(run(ins, ipr, 1)[0])
}

func run(ins []instr, ipr, r0 int) [6]int {
	var ip int
	var regs [6]int
	regs[0] = r0
	reg := func(i int) int { return regs[i] }
	for ip >= 0 && ip < len(ins) {
		// Optimized version of 1-16 for part 2.
		if ip == 1 {
			for n := 1; n <= regs[1]; n++ {
				if regs[1]%n == 0 {
					regs[0] += n
				}
			}
			break
		}

		regs[ipr] = ip
		in := &ins[ip]
		regs[in.c] = in.op.run(in.a, in.b, reg)
		ip = regs[ipr]
		ip++
	}
	return regs
}

type instr struct {
	op      op
	a, b, c int
}

type op int

const (
	addr op = iota
	addi
	mulr
	muli
	banr
	bani
	borr
	bori
	setr
	seti
	gtir
	gtri
	gtrr
	eqir
	eqri
	eqrr

	nops = int(eqrr) + 1
)

// run returns the value to store in c when o is called
// with the supplied a and b values and register state.
func (o op) run(va, vb int, reg func(int) int) int {
	switch o {
	case addr:
		return reg(va) + reg(vb)
	case addi:
		return reg(va) + vb
	case mulr:
		return reg(va) * reg(vb)
	case muli:
		return reg(va) * vb
	case banr:
		return reg(va) & reg(vb)
	case bani:
		return reg(va) & vb
	case borr:
		return reg(va) | reg(vb)
	case bori:
		return reg(va) | vb
	case setr:
		return reg(va)
	case seti:
		return va
	case gtir:
		return lib.If(va > reg(vb), 1, 0)
	case gtri:
		return lib.If(reg(va) > vb, 1, 0)
	case gtrr:
		return lib.If(reg(va) > reg(vb), 1, 0)
	case eqir:
		return lib.If(va == reg(vb), 1, 0)
	case eqri:
		return lib.If(reg(va) == vb, 1, 0)
	case eqrr:
		return lib.If(reg(va) == reg(vb), 1, 0)
	default:
		panic(fmt.Sprintf("Invalid opcode %d", o))
	}
}

var strOp = map[string]op{
	"addr": addr,
	"addi": addi,
	"mulr": mulr,
	"muli": muli,
	"banr": banr,
	"bani": bani,
	"borr": borr,
	"bori": bori,
	"setr": setr,
	"seti": seti,
	"gtir": gtir,
	"gtri": gtri,
	"gtrr": gtrr,
	"eqir": eqir,
	"eqri": eqri,
	"eqrr": eqrr,
}
