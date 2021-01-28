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
	const (
		epos      = 2020
		ncards    = 119315717514047
		nshuffles = 101741582076661
	)

	// Hmm. So I was right that there would be a bunch of shuffles, but there
	// are also a bunch of cards now. I think I should run a reverse shuffle until
	// I find a loop.
	//
	// Okay, that doesn't seem work: even after a few million reverse shuffles,
	// I haven't seen us end up back at position 2020. I don't know what to do next.
	// The operations in my input are complicated enough that it doesn't seem possible
	// to analyze them.
	/*
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
	*/

	// When using prime numbers of cards around 10000 and 20000, I can see that we always seem to
	// return to the original sequence after n-1 shuffles, and also (in some cases) after numbers of
	// shuffles that evenly divide n-1. I looked at the factors of the n-1 values but didn't see any
	// patterns. I think that the number of shuffles gien in the problem will need to be close to
	// one of these loop points, but I'm not sure how to figure out what the loop point is without
	// knowing the number of loops.
	//
	// Furthermore, the only factors of ncards-1 are 2 and 59657858757023, so I'm actually not
	// convinced that there's a loop point near nshuffles. For the much smaller card counts that I
	// looked at, there were always either 1 or 2 loops when the only factors were 2 and a prime.
	/*
		for _, ncards := range []int64{9833, 9839, 9851, 9857, 9859, 9871, 9883, 9887, 9901, 9907, 9923,
			9929, 9931, 9941, 9949, 9967, 9973, 10007, 10009, 10037, 10039, 10061, 10067, 10069, 10079,
			10091, 10093, 10099, 10103, 10111, 10133, 10139, 10141, 10151, 10159, 10163, 10169, 10177,
			19949, 19961, 19963, 19973, 19979, 19991, 19993, 19997} {
			pos := int64(epos)
			var nshufs int64
			var loops []int64
			for nshufs <= ncards && len(loops) < 10 {
				nshufs++
				pos = unshuffle(pos, ncards)
				if pos == epos {
					loops = append(loops, nshufs)
				}
			}
			fmt.Printf("%5d: %0.3f %v\n", ncards, float64(ncards)/float64(loops[0]), loops)
		}
	*/
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
