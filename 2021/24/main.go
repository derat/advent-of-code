package main

import (
	"fmt"

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

	v := newVM(ins)
	res, input := v.run(state{}, false /* smallest */)
	lib.Assert(res)
	lib.AssertEq(len(input), 14)
	regs, err := run(ins, input) // check validity
	lib.Assertf(err == nil, fmt.Sprint(err))
	lib.AssertEq(regs[3], 0)
	fmt.Println(join(input))

	return // TODO: Get part 2 working again.

	v.reset()
	res, input = v.run(state{}, true /* smallest */)
	lib.Assert(res)
	lib.AssertEq(len(input), 14)
	regs, err = run(ins, input)
	lib.Assertf(err == nil, fmt.Sprint(err))
	lib.AssertEq(regs[3], 0)
	fmt.Println(join(input))

	/*
		// Here's the obvious approach of actually feeding all possible model numbers
		// into the program. As expected, the input space is way too big for this to
		// find an answer in a reasonable amount of time.
		var input []int64
		for _, ch := range []byte("99999999999999") {
			input = append(input, int64(ch-'0'))
		}
		lib.AssertEq(len(input), 14)

		for {
			if regs, err := run(ins, input); err == nil && regs[3] == 0 {
				fmt.Println(input)
				break
			}
			lib.Assert(dec(input, len(input)-1))
		}
	*/
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

func (v *vm) reset() { v.bad = make(map[state]struct{}) }

func (v *vm) run(start state, smallest bool) (res bool, input []int64) {
	st := start
	regs := st.regs[:]
	for ; st.ip < len(v.ins); st.ip++ {
		// Check if we already know that we're going to fail.
		if _, ok := v.bad[start]; ok {
			return false, nil
		}

		in := v.ins[st.ip]
		switch in.Op {
		case inp:
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

func run(ins []lib.Instr, input []int64) ([]int64, error) {
	regs := make([]int64, 4) //  w, x, y, z
	var ii int               // index into input
	for i, in := range ins {
		switch in.Op {
		case inp:
			if ii >= len(input) {
				return nil, fmt.Errorf("inst %d reads beyond %d-digit input", i, len(input))
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
				return nil, fmt.Errorf("inst %d divides %v by 0", i, a)
			}
			*in.Ptr(0, regs) = a / b
		case mod:
			a := in.Val(0, regs)
			b := in.Val(1, regs)
			if a < 0 || b <= 0 {
				return nil, fmt.Errorf("inst %d mods %v by %v", a, b)
			}
			*in.Ptr(0, regs) = a % b
		case eql:
			*in.Ptr(0, regs) = int64(lib.If(in.Val(0, regs) == in.Val(1, regs), 1, 0))
		}
	}
	return regs, nil
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
