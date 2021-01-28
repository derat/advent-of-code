package main

import (
	"fmt"
	"math/big"

	"github.com/derat/advent-of-code/lib"
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
		case lib.ExtractMaybe(ln, `^cut (-?\d+)$`, &in.v):
			in.op = cut
		case lib.ExtractMaybe(ln, `^deal with increment (\d+)$`, &in.v):
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
	//
	// Hmm. So I was right that there would be a bunch of shuffles, but there
	// are also a bunch of cards now. I think I should run a reverse shuffle until
	// I find a loop.
	//
	// Okay, that doesn't seem work: even after a few million reverse shuffles,
	// I haven't seen us end up back at position 2020. I don't know what to do next.
	// The operations in my input are complicated enough that it doesn't seem possible
	// to analyze them.
	const (
		epos      = 2020
		ncards    = 119315717514047
		nshuffles = 101741582076661
	)

	pos := int64(epos)
	var nshufs int
	for {
		nshufs++
		pos = unshuffle(pos, ncards)
		if pos == epos {
			break
		}
	}
	// This doesn't get reached. :-/
	fmt.Println("looped after", nshufs)
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
