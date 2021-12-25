package main

import (
	"fmt"
	"runtime"

	"github.com/derat/advent-of-code/lib"
)

const (
	inp uint8 = iota
	add
	mul
	div
	mod
	eql
)

const (
	ninput = 14
	bsize  = 18
)

func main() {
	const (
		varArg    = ` ([w-z])`
		varNumArg = ` ([w-z]|-?\d+)`
	)
	var ins []lib.Instr
	for _, ln := range lib.InputLines("2021/24") {
		ins = append(ins, lib.NewInstr(ln, 'w', 'z', map[uint8]string{
			inp: "^inp" + varArg + `$`,
			add: "^add" + varArg + varNumArg + `$`,
			mul: "^mul" + varArg + varNumArg + `$`,
			div: "^div" + varArg + varNumArg + `$`,
			mod: "^mod" + varArg + varNumArg + `$`,
			eql: "^eql" + varArg + varNumArg + `$`,
		}))
	}

	// Here's my fast (~30 ms) rewritten solution based on manual analysis of the code.
	// It's not at all generalized, but I'm assuming that the constants in the code are
	// the only things that vary across different inputs.
	var bls []block
	for i := 0; i < ninput; i++ {
		start := i * bsize
		bl, err := newBlock(ins[start : start+bsize])
		if err != nil {
			lib.Panicf("Bad block %v: %v", i, err)
		}
		bls = append(bls, bl)
	}

	// Part 1: "To enable as many submarine features as possible, find the largest valid
	// fourteen-digit model number that contains no 0 digits. What is the largest model number
	// accepted by MONAD?"
	input, ok := findInput(0, bls, false)
	lib.Assert(ok)
	regs := make([]int64, 4)
	err := run(ins, regs, input) // check validity
	lib.Assertf(err == nil, fmt.Sprint(err))
	lib.AssertEq(regs[3], 0)
	fmt.Println(join(input))

	// Part 2: "What is the smallest model number accepted by MONAD?"
	input, ok = findInput(0, bls, true)
	lib.Assert(ok)
	regs = make([]int64, 4)
	err = run(ins, regs, input) // check validity
	lib.Assertf(err == nil, fmt.Sprint(err))
	lib.AssertEq(regs[3], 0)
	fmt.Println(join(input))

	// I can't say that I'm proud of this one. On my not-super-fast Chromebook, part 1 (with an
	// answer maybe 10% of the way into the input space) takes around half a minute and maybe 3 GB
	// of memory, and part 2 (near the middle of the input space) took closer to 10 minutes. Then I
	// did some refactoring, and now part 2 doesn't even finish before eating all my memory. At
	// least I got the stars, I guess.
	//
	// After I first read the problem and verified that testing all possible inputs in descending
	// order was too slow, my initial reaction was that I need to execute the instructions in
	// reverse to derive all 14 inputs, starting with the knowledge that z ends at 0.
	//
	// Looking at my input more closely, I saw that the same rough pattern is followed after each
	// 'inp' instruction. Digging further into that is probably the easiest way to solve this (maybe
	// even without writing any code?), but I decided to spend more time working on a generalized
	// solution first.
	//
	// It wasn't immediately obvious to me how I should track the values of registers while running
	// the program in reverse (e.g. if z ends at 0, then y is -z since my last instruction is "add z
	// y", and so on). I also still have nightmares about inverse modulo after Slam Shuffle
	// (https://adventofcode.com/2019/day/22).
	//
	// So, I ended up writing a recursive solution that branches over the 9 possible inputs every
	// time an 'inp' instruction is seen and caches states (denoted by the ip and register values)
	// that are doomed to fail. This is clearly dependent on there being a relatively small number
	// of possible register states at each instruction, though, which doesn't actually seem to be
	// the case.
	//
	// In any case, I'm checking in this crappy version of my code to make sure I don't lose
	// everything if the OOM killer kills Chrome -- I've had trouble with Crostini's data not being
	// synced to disk after system crashes before. :-/
	/*
		v := newVM(ins)
		res, input := v.run(state{}, false) // largest
		lib.Assert(res)
		lib.AssertEq(len(input), ninput)
		regs := make([]int64, 4)
		err := run(ins, regs, input) // check validity
		lib.Assertf(err == nil, fmt.Sprint(err))
		lib.AssertEq(regs[3], 0)
		fmt.Println(join(input))

		v.reset()
		res, input = v.run(state{}) // smallest
		lib.Assert(res)
		lib.AssertEq(len(input), ninput)
		regs = make([]int64, 4)
		err = run(ins, regs, input)
		lib.Assertf(err == nil, fmt.Sprint(err))
		lib.AssertEq(regs[3], 0)
		fmt.Println(join(input))
	*/

	// Here's the obvious approach of actually feeding all possible model numbers
	// into the program. As expected, the input space is way too big for this to
	// find an answer in a reasonable amount of time.
	/*
		var input []int64
		for _, ch := range []byte("99999999999999") {
			input = append(input, int64(ch-'0'))
		}
		lib.AssertEq(len(input), ninput)

		for {
			if regs, err := run(ins, input); err == nil && regs[3] == 0 {
				fmt.Println(input)
				break
			}
			lib.Assert(dec(input, len(input)-1))
		}
	*/
}

const (
	w = 0
	x = 1
	y = 2
	z = 3
)

// My input program consists of the following 18-instruction block repeated 14 times
// (once for each value read via the 'inp' instruction):
//
//   0  inp w
//   1  mul x 0
//   2  add x z
//   3  mod x 26
//   4  div z m   # m is 1 or 26
//   5  add x n   # n can be any value
//   6  eql x w
//   7  eql x 0
//   8  mul y 0
//   9  add y 25
//  10  mul y x
//  11  add y 1
//  12  mul z y
//  13  mul y 0
//  14  add y w
//  15  add y o   # o can be any value
//  16  mul y x
//  17  add z y
//
// If the input equals (z%26 + n), then:
//
//   6: x is 1
//   7: x is 0
//  11: y is 1
//  12: z is unchanged
//  16: y is 0
//  17: z is unchanged
//
// z = z/m
//
// Otherwise (input doesn't equal z%26 + n):
//
//   6: x is 0
//   7: x is 1
//  11: y is 26
//  12: z is z*26
//  16: y is input+o
//  17: z is z+input+o
//
// z = z/m * 26 + input + o
//
// Here are the m, n, and o constants from my input:
//
//  block   m   n   o
//      0   1  15   4
//      1   1  14  16
//      2   1  11  14
//      3  26 -13   3
//      4   1  14  11
//      5   1  15  13
//      6  26  -7  11
//      7   1  10   7
//      8  26 -12  12
//      9   1  15  15
//     10  26 -16  13
//     11  26  -9   1
//     12  26  -8  15
//     13  26  -8   4
//
// In the blocks where n > 9, there's no way for the input to match.
// m is also 1 in all of those, so we know we'll be be multiplying by 26
// without a divide to cancel it out. There are 7 of those blocks, so we
// need the other 7 with negative n values (which also all have m=26) to
// have matching input so we can get z back down toward 0.

type block struct{ m, n, o int64 }

// newBlock checks that the supplied instructions match the format described
// above and extracts the m, n, and o constants.
func newBlock(ins []lib.Instr) (block, error) {
	lib.AssertEq(len(ins), bsize)
	regs := make([]int64, 4)
	for i, want := range []struct {
		op     uint8
		ra, rb int
	}{
		{inp, w, -1},
		{mul, x, -1},
		{add, x, z},
		{mod, x, -1},
		{div, z, -1},
		{add, x, -1},
		{eql, x, w},
		{eql, x, -1},
		{mul, y, -1},
		{add, y, -1},
		{mul, y, x},
		{add, y, -1},
		{mul, z, y},
		{mul, y, -1},
		{add, y, w},
		{add, y, -1},
		{mul, y, x},
		{add, z, y},
	} {
		in := ins[i]
		lib.AssertEq(in.Op, want.op)
		lib.AssertEq(in.Ptr(0, regs), &regs[want.ra])
		if want.rb != -1 {
			lib.AssertEq(in.Ptr(1, regs), &regs[want.rb])
		}
	}

	// Check constants.
	lib.AssertEq(ins[1].Val(1, regs), 0)
	lib.AssertEq(ins[3].Val(1, regs), 26)
	lib.AssertEq(ins[7].Val(1, regs), 0)
	lib.AssertEq(ins[8].Val(1, regs), 0)
	lib.AssertEq(ins[9].Val(1, regs), 25)
	lib.AssertEq(ins[11].Val(1, regs), 1)
	lib.AssertEq(ins[13].Val(1, regs), 0)

	return block{
		m: ins[4].Val(1, regs),
		n: ins[5].Val(1, regs),
		o: ins[15].Val(1, regs),
	}, nil
}

// findInput finds the "best" (per the smallest arg) input that sets z to 0
// after executing the supplied blocks.
func findInput(z int64, bls []block, smallest bool) ([]int64, bool) {
	if len(bls) == 0 {
		return nil, z == 0
	}

	bl := bls[0]
	if bl.n > 9 {
		lib.AssertEq(bl.m, 1)
		var a, b, inc int64 = 9, 0, -1 // half-open so we can use != in loop
		if smallest {
			a, b, inc = 1, 10, 1
		}
		for i := a; i != b; i += inc {
			input, ok := findInput((z/bl.m)*26+i+bl.o, bls[1:], smallest)
			if ok {
				return append([]int64{i}, input...), true
			}
		}
		return nil, false
	}

	i := z%26 + bl.n
	if i < 1 || i > 9 {
		return nil, false
	}
	lib.AssertEq(bl.m, 26)
	input, ok := findInput(z/bl.m, bls[1:], smallest)
	if !ok {
		return nil, false
	}
	return append([]int64{i}, input...), true
}

type state struct {
	ip   int
	regs [4]int64
}

type vm struct {
	ins []lib.Instr
	bad map[state]struct{} // states that lead to failure for all inputs
}

func newVM(ins []lib.Instr) *vm { return &vm{ins, make(map[state]struct{})} }

func (v *vm) reset() {
	v.bad = make(map[state]struct{})
	runtime.GC()
}

func (v *vm) run(start state, smallest bool) (res bool, input []int64) {
	// Check if we already know that we're going to fail.
	if _, ok := v.bad[start]; ok {
		return false, nil
	}

	st := start
	regs := st.regs[:]
	for ; st.ip < len(v.ins); st.ip++ {
		in := v.ins[st.ip]
		switch in.Op {
		case inp:
			lib.AssertEq(in.Ptr(0, regs), &regs[0]) // storing in w
			regs[1] = 0
			regs[2] = 0

			var a, b, inc int64 = 9, 0, -1 // half-open so we can use != in loop
			if smallest {
				a, b, inc = 1, 10, 1
			}
			for i := a; i != b; i += inc {
				*in.Ptr(0, regs) = i
				if res, inp := v.run(state{st.ip + 1, st.regs}, smallest); res {
					return true, append([]int64{i}, inp...)
				}
			}
			// All inputs starting from this state failed.
			v.bad[start] = struct{}{}
			return false, nil
		case add:
			*in.Ptr(0, regs) = in.Val(0, regs) + in.Val(1, regs)
		case mul:
			*in.Ptr(0, regs) = in.Val(0, regs) * in.Val(1, regs)
		case div:
			a := in.Val(0, regs)
			b := in.Val(1, regs)
			if b == 0 {
				v.bad[start] = struct{}{}
				return false, nil
			}
			*in.Ptr(0, regs) = a / b
		case mod:
			a := in.Val(0, regs)
			b := in.Val(1, regs)
			if a < 0 || b <= 0 {
				v.bad[start] = struct{}{}
				return false, nil
			}
			*in.Ptr(0, regs) = a % b
		case eql:
			*in.Ptr(0, regs) = int64(lib.If(in.Val(0, regs) == in.Val(1, regs), 1, 0))
		}
	}

	// We're at the end of the program.
	if st.regs[3] == 0 {
		return true, nil // success!
	}
	v.bad[start] = struct{}{}
	return false, nil
}

func run(ins []lib.Instr, regs, input []int64) error {
	var ii int // index into input
	for i, in := range ins {
		switch in.Op {
		case inp:
			if ii >= len(input) {
				return fmt.Errorf("inst %d reads beyond %d-digit input", i, len(input))
			}
			*in.Ptr(0, regs) = input[ii]
			ii++
		case add:
			*in.Ptr(0, regs) = in.Val(0, regs) + in.Val(1, regs)
		case mul:
			*in.Ptr(0, regs) = in.Val(0, regs) * in.Val(1, regs)
		case div:
			a := in.Val(0, regs)
			b := in.Val(1, regs)
			if b == 0 {
				return fmt.Errorf("inst %d divides %v by 0", i, a)
			}
			*in.Ptr(0, regs) = a / b
		case mod:
			a := in.Val(0, regs)
			b := in.Val(1, regs)
			if a < 0 || b <= 0 {
				return fmt.Errorf("inst %d mods %v by %v", i, a, b)
			}
			*in.Ptr(0, regs) = a % b
		case eql:
			*in.Ptr(0, regs) = int64(lib.If(in.Val(0, regs) == in.Val(1, regs), 1, 0))
		}
	}
	return nil
}

// dec decrements n's i-th value (from the right), recursively borrowing from the left if needed.
// It returns false if num can't be decremented (i.e. the input was all 1s).
func dec(num []int64, i int) bool {
	if num[i]--; num[i] >= 1 {
		return true // successfully borrowed
	} else if i == 0 {
		return false // reached the end
	}
	num[i] = 9
	return dec(num, i-1) // borrow
}

func join(ns []int64) string {
	var s string
	for _, n := range ns {
		s += fmt.Sprint(n)
	}
	return s
}
