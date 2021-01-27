package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const ncards = 10_007

	// Part 1: Print the position of card 2019 after shuffling.
	//
	// This reminds me of 2016/21, except that one had operations based
	// on values (e.g. "swap letter a with letter b"), while this one
	// only has operations based on position.
	//
	// The obvious optimization would be to just track the position of
	// card 2019 through the operations, but I suspect that part 2 is going
	// to be something like "Print the position of card 2019 after a million
	// operations", so I'm going to get the final positions of all cards.
	cards := make([]int, ncards)
	lib.SetAscInt(cards, 0)

	for _, ln := range lib.InputLines("2019/22") {
		var v int
		switch {
		case ln == "deal into new stack":
			lib.Reverse(cards)
		case lib.ExtractMaybe(ln, `^cut (-?\d+)$`, &v):
			lib.RotateSlice(cards, -v)
		case lib.ExtractMaybe(ln, `^deal with increment (\d+)$`, &v):
			// Not sure if there's a more-efficient way to do this using swaps...
			dealt := make([]int, ncards)
			for i := range cards {
				dealt[i*v%len(cards)] = cards[i]
			}
			cards = dealt
		default:
			lib.Panicf("Invalid line %q", ln)
		}
	}
	for i, v := range cards {
		if v == 2019 {
			fmt.Println(i)
			break
		}
	}
}
