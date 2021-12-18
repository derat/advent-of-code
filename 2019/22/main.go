package main

import (
	"fmt"
	"math/big"

	"github.com/derat/advent-of-code/lib"
)

const (
	dump = false
	test = false
)

func main() {
	// This reminds me of 2016/21, except that one had operations based
	// on values (e.g. "swap letter a with letter b"), while this one
	// only has operations based on position.
	//
	// The obvious optimization would be to just track the position of
	// card 2019 through the operations, but I suspect that part 2 is going
	// to be something like "Print the position of card 2019 after a million
	// operations", so I'm going to get the final positions of all cards.
	//
	// (Well, that assumption was wrong. See below.)
	var ins []instr
	for _, ln := range lib.InputLines("2019/22") {
		var in instr
		switch {
		case ln == "deal into new stack":
			in.op = dealnew
		case lib.TryExtract(ln, `^cut (-?\d+)$`, &in.v):
			in.op = cut
		case lib.TryExtract(ln, `^deal with increment (\d+)$`, &in.v):
			in.op = dealinc
		default:
			lib.Panicf("Invalid line %q", ln)
		}
		ins = append(ins, in)
	}

	shuffle := func(pos, ncards int64) int64 {
		for _, in := range ins {
			pos = in.do(pos, ncards)
		}
		return pos
	}
	unshuffle := func(pos, ncards int64) int64 {
		for i := len(ins) - 1; i >= 0; i-- {
			pos = ins[i].undo(pos, ncards)
		}
		return pos
	}

	// Part 1: With a deck of 10007 cards, print card 2019's position after shuffling.
	fmt.Println(shuffle(2019, 10_007))

	// Part 2: Print the value of the card at position 2020 after shuffling
	// 119315717514047 cards 101741582076661 times.
	const (
		epos      = 2020
		ncards    = 119315717514047
		nshuffles = 101741582076661
	)

	x := big.NewInt(1) // scalar
	y := big.NewInt(0) // constant
	for _, in := range ins {
		switch in.op {
		case dealnew:
			x.Mul(x, big.NewInt(-1))
			y.Mul(y, big.NewInt(-1))
			y.Add(y, big.NewInt(ncards-1))
		case cut:
			y.Sub(y, big.NewInt(in.v))
		case dealinc:
			x.Mul(x, big.NewInt(in.v))
			y.Mul(y, big.NewInt(in.v))
		}
	}
	x.Mod(x, big.NewInt(ncards))
	y.Mod(y, big.NewInt(ncards))
	if dump {
		fmt.Println("x", x.Int64())
		fmt.Println("y", y.Int64())
	}

	x0 := big.NewInt(1) // scalar
	y0 := big.NewInt(0) // constant
	for i := len(ins) - 1; i >= 0; i-- {
		in := ins[i]
		switch in.op {
		case dealnew: // same as forward
			x0.Mul(x0, big.NewInt(-1))
			y0.Sub(y0, big.NewInt(ncards-1))
			y0.Mul(y0, big.NewInt(-1))
		case cut:
			y0.Add(y0, big.NewInt(in.v))
		case dealinc:
			inv := big.NewInt(in.v)
			inv.ModInverse(inv, big.NewInt(ncards))
			x0.Mul(x0, inv)
			y0.Mul(y0, inv)
		}
	}
	x0.Mod(x0, big.NewInt(ncards))
	y0.Mod(y0, big.NewInt(ncards))
	if dump {
		fmt.Println("x0", x0.Int64())
		fmt.Println("y0", y0.Int64())
	}

	// Drop out constant term per https://math.stackexchange.com/a/2388143:
	//
	//  f(n) = x*f(n-1) + y
	//	f(n-1) = x*f(n-2) + y                 prev in sequence
	//	f(n) - f(n-1) = x*f(n-1) - x*f(n-2)   subtract second equation from first
	//	f(n) = (x+1)*f(n-1) - x*f(n-2)        rearrange
	//
	// For my input:
	//
	//  f(n+2) = 7277816997830721537*f(n+1) - 7277816997830721536*f(n)
	//    f(0) = 2020
	//    f(1) = 14699430743457602759107 = 61829916141776 (mod ncards)
	//
	// This is a homogeneous linear recurrence relation of order 2 with constant coefficients.
	// Per https://study.com/academy/lesson/how-to-solve-linear-recurrence-relations.html:
	//
	// First, find characteristic equation:
	//
	//	s^i = (x+1)*s^(i-1) - x*s^(i-2)       transform to variable s raised to a power
	//	s^2 = (x+1)*s - x                     divide by s^(i-2)
	//  s^2 - (x+1)*s + x = 0                 make equal to 0
	//
	// I started trying to solve this in code with the quadratic formula, but that was very
	// painful. Plugging it into Wolfram Alpha yields a closed-form solution for f(n):
	//
	//	f(n+2) = 40286879916730*f(n+1) - 40286879916729*f(n)
	//	f(0) = 2020
	//	f(1) = 81416758296639728
	//	f(n) = (20354189574159427 * 40286879916729^n - 9315216211787) / 10071719979182
	//
	// Reverse shuffle:
	//
	//	f(n+2) = 70994688272735*f(n+1) - 70994688272734*f(n)
	//	f(0) = 2020
	//	f(1) = 14976082037420
	//  f(n) = (20 * (374402050885 * 2^(n+1) * 35497344136367^n + 7169714711444263)) / 70994688272733

	// shufflen returns the position of card 2020 after nshuffles forward shuffles.
	shufflen := func(ncards, nshuffles int64) int64 {
		nc := big.NewInt(ncards)

		v := big.NewInt(40286879916729)
		v.Exp(v, big.NewInt(nshuffles), nc)
		v.Mul(v, big.NewInt(20354189574159427))
		v.Sub(v, big.NewInt(9315216211787))

		inv := big.NewInt(10071719979182)
		inv.ModInverse(inv, nc)
		v.Mul(v, inv)

		v.Mod(v, nc)
		return v.Int64()
	}

	// unshufflen returns the position of card 2020 after nshuffles reverse shuffles.
	unshufflen := func(ncards, nshuffles int64) int64 {
		nc := big.NewInt(ncards)

		a := big.NewInt(2)
		a.Exp(a, big.NewInt(nshuffles+1), nc)

		b := big.NewInt(35497344136367)
		b.Exp(b, big.NewInt(nshuffles), nc)

		v := big.NewInt(374402050885)
		v.Mul(v, a)
		v.Mul(v, b)
		v.Add(v, big.NewInt(7169714711444263))
		v.Mul(v, big.NewInt(20))

		inv := big.NewInt(70994688272733)
		inv.ModInverse(inv, nc)
		v.Mul(v, inv)

		v.Mod(v, big.NewInt(ncards))
		return v.Int64()
	}

	if test {
		pos := int64(epos)
		fmt.Println("shuffle from", pos)
		for i := int64(1); i < 10; i++ {
			pos = shuffle(pos, ncards)
			fmt.Println(i, pos)
			fmt.Println(i, shufflen(ncards, i))
		}
		fmt.Println()

		pos = epos
		fmt.Println("unshuffle from", pos)
		for i := int64(1); i < 10; i++ {
			pos = unshuffle(pos, ncards)
			fmt.Println(i, pos)
			fmt.Println(i, unshufflen(ncards, i))
		}
	}

	// Shuffling ncards-1 times gets us back to the original 2020 position.
	lib.AssertEq(shufflen(ncards, ncards-1), int64(epos))

	// We want to figure out where the card at the 2020 position was nshuffles ago.
	fmt.Println(unshufflen(ncards, nshuffles))
}

type op int

const (
	dealnew op = iota // "deal into new stack" (reverse)
	cut               // "cut n" (rotate left)
	dealinc           // "deal with increment n"
)

type instr struct {
	op op
	v  int64 // only for cut and dealinc
}

func (in *instr) do(pos, ncards int64) int64 {
	switch in.op {
	case dealnew:
		return ncards - pos - 1
	case cut:
		return (pos - in.v + ncards) % ncards
	case dealinc:
		// Use big.Int to avoid overflow.
		res := big.NewInt(pos)
		res.Mul(res, big.NewInt(in.v))
		res.Mod(res, big.NewInt(ncards))
		return res.Int64()
	}
	lib.Panicf("Invalid op %v", in.op)
	return -1
}

func (in *instr) undo(pos, ncards int64) int64 {
	switch in.op {
	case dealnew:
		return ncards - pos - 1
	case cut:
		return (pos + in.v + ncards) % ncards
	case dealinc:
		res := big.NewInt(0)
		nc := big.NewInt(ncards)
		res.ModInverse(big.NewInt(in.v), nc)
		res.Mul(res, big.NewInt(pos))
		res.Mod(res, nc)
		return res.Int64()
	}
	lib.Panicf("Invalid op %v", in.op)
	return -1
}
