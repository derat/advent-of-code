package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

const debug = false

func main() {
	lines := lib.InputLines("2018/21")

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

	// Here's the annotated and then Go version of my input:
	//
	//     #ip 3
	//   0 seti 123 0 1       # r1 = 123
	//   1 bani 1 456 1       # r1 &= 456
	//   2 eqri 1 72 1        # r1 = 1 (r1 == 72) or 0 (r1 != 72)
	//   3 addr 1 3 3         # jmp 4 + r1
	//   4 seti 0 0 3         # jmp 1
	//   5 seti 0 7 1         # r1 = 0
	//   6 bori 1 65536 4     # r4 = r1 | 65536
	//   7 seti 3798839 3 1   # r1 = 3798839
	//   8 bani 4 255 5       # r5 = r4 & 255
	//   9 addr 1 5 1         # r1 += r5
	//  10 bani 1 16777215 1  # r1 &= 16777215
	//  11 muli 1 65899 1     # r1 *= 65899
	//  12 bani 1 16777215 1  # r1 &= 16777215
	//  13 gtir 256 4 5       # r5 = 1 (r4 < 256) or 0 (r4 >= 256)
	//  14 addr 5 3 3         # jmp 15 + r5
	//  15 addi 3 1 3         # jmp 17
	//  16 seti 27 6 3        # jmp 28
	//  17 seti 0 2 5         # r5 = 0
	//  18 addi 5 1 2         # r2 = r5 + 1
	//  19 muli 2 256 2       # r2 *= 256
	//  20 gtrr 2 4 2         # r2 = 1 (r2 > r4) or 0 (r2 <= r4)
	//  21 addr 2 3 3         # jmp 22 + r2
	//  22 addi 3 1 3         # jmp 24
	//  23 seti 25 3 3        # jmp 26
	//  24 addi 5 1 5         # r5 += 1
	//  25 seti 17 1 3        # jmp 18
	//  26 setr 5 6 4         # r4 = r5
	//  27 seti 7 8 3         # jmp 8
	//  28 eqrr 1 0 5         # r5 = 1 (r0 == r1) or 0 (r0 != r1)
	//  29 addr 5 3 3         # jmp exit (r5 == 1) or 30 (r5 == 0)
	//  30 seti 5 6 3         # jmp 6
	//
	//  for 123&456 != 72 {}  // 0-4
	//  r1 = 0                // 5
	// Loop:
	//  for {                 // 6-27
	//    r4 = r1 | 65536     // 6
	//    r1 = 3798839        // 7
	//    for {               // 8-27
	//      r5 = r4 & 255     // 8
	//      r1 += r5          // 9
	//      r1 &= 16777215    // 10
	//      r1 *= 65899       // 11
	//      r1 &= 16777215    // 12
	//      if r4 < 256 {     // 13-16
	//        if r0 == r1 {   // 28
	//          exit          // 29
	//        } else {
	//          continue Loop // 30
	//        }
	//      }
	//      // Sets r4 to r4 / 256.
	//      r5 = 0            // 17
	//      for ; ; r5 += 1 { // 24
	//        r2 = r5 + 1     // 18
	//        r2 *= 256       // 19
	//        if r2 > r4 {    // 20-23, 25
	//          break
	//        }
	//      }
	//      r4 = r5           // 26
	//    }
	//  }

	// For part 1 (smallest r0 leading to earliest termination), I ran the actual instructions but
	// with the 17-26 loop optimized to a single divide. I iterated over increasing r0 values with a
	// limit on the total number of instructions and submitted a few wrong answers until I got to
	// the right one.
	//
	// For part 2, I wised up and just ran a Go version of the code until it encounters a loop,
	// paying attention to r1's value at the point where it's compared against r0 to exit the
	// program. I probably could've instead just used the optimized code from part 1.
	var r1, r4 int
	var fastest, slowest int
	seen := make(map[uint64]struct{})
	seen0 := make(map[int]struct{})
Loop:
	for {
		st := lib.PackInts(r1, r4)
		if _, ok := seen[st]; ok {
			break
		}
		seen[st] = struct{}{}

		r4 = r1 | 65536
		r1 = 3798839
		for {
			r1 += r4 & 255
			r1 &= 16777215
			r1 *= 65899
			r1 &= 16777215
			if r4 < 256 {
				if fastest == 0 {
					fastest = r1
				}
				if _, ok := seen0[r1]; !ok {
					slowest = r1
					seen0[r1] = struct{}{}
				}
				continue Loop
			}
			r4 /= 256
		}
	}

	fmt.Println(fastest) // part 1
	fmt.Println(slowest) // part 2
}

func run(ins []instr, ipr, r0, ops int) ([6]int, int) {
	var ip int
	var regs [6]int
	regs[0] = r0
	reg := func(i int) int { return regs[i] }
	var nops int
	for ip >= 0 && ip < len(ins) {
		regs[ipr] = ip
		in := &ins[ip]
		regs[in.c] = in.op.exec(in.a, in.b, reg)

		if debug {
			fmt.Printf("%2d [%1d %12d %1d %2d %8d %3d]\n",
				ip, regs[0], regs[1], regs[2], regs[3], regs[4], regs[5])
		}

		ip = regs[ipr]
		ip++

		nops++
		if ops > 0 && nops == ops {
			return regs, -1
		}
	}
	return regs, nops
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

// exec returns the value to store in c when o is called
// with the supplied a and b values and register state.
func (o op) exec(va, vb int, reg func(int) int) int {
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
