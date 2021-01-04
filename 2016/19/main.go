package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	elves := lib.InputInts("2016/19")[0]

	// Part 1
	min := 1
	max := elves
	for step := 2; min < max; step *= 2 {
		// If we ended on the last elf, then it'll take the first elf's presents next round.
		// Otherwise, the second-to-last elf takes the last elf's presents.
		if (max-min)%step == 0 {
			min += step
		} else {
			max -= step / 2
		}
	}
	fmt.Println(min)

	// Part 2: Each elf takes presents from the elf across the circle (rounding down)
	// instead of from the elf to their left.
	//
	// This approach (using a linked list of spans) is awful. It starts out fast, but it
	// slows to a crawl before too long, presumably due to fragmentation as more and more
	// spans are created. It took about 40 minutes to reduce my input of ~3M elves in half
	// to ~1.5M elves, and probably close to 2 hours for the whole thing.
	//
	// Presumably there's a solution here that's closer to what I did in part 1, but it
	// didn't leap out at me after staring at some example sequences.
	spn := &span{min: 1, max: elves} // current span
	spn.prev, spn.next = spn, spn
	elf := 1     // current elf
	cnt := elves // elves still in the game
	for {
		off := cnt / 2 // offset to elf whose presents are gonna get ganked
		ts, te := find(spn, elf, off)

		if ts.min == te {
			ts.min++
		} else if ts.max == te {
			ts.max--
		} else {
			ns := &span{min: te + 1, max: ts.max, prev: ts, next: ts.next}
			ts.max = te - 1
			ts.next.prev = ns
			ts.next = ns
			if spn == ts && elf > te {
				spn = ns
			}
		}
		if ts.min > ts.max {
			ts.prev.next = ts.next
			ts.next.prev = ts.prev
		}

		cnt--
		if cnt == 1 {
			break
		}
		spn, elf = find(spn, elf, 1) // move to next elf
	}
	lib.AssertEq(spn, spn.next)
	lib.AssertEq(spn.min, spn.max)
	fmt.Println(spn.min)
}

type span struct {
	min, max   int
	prev, next *span
}

func (s *span) num() int {
	return s.max - s.min + 1
}

func (s *span) String() string {
	return fmt.Sprintf("[%d, %d]", s.min, s.max)
}

// find performs a linear scan to advance by the specified offset
// from startElf in sp.
func find(sp *span, startElf, offset int) (s *span, elf int) {
	lib.Assertf(offset >= 0, "Negative offset %d", offset)
	rem := offset

	// Check the starting span first.
	lib.Assertf(sp.min <= startElf && sp.max >= startElf,
		"Elf %d not in starting span %s", startElf, sp)
	if elf := startElf + rem; elf <= sp.max {
		return sp, elf
	}
	rem -= sp.max - startElf + 1

	for sp = sp.next; true; sp = sp.next {
		if sp.min+rem <= sp.max {
			return sp, sp.min + rem
		}
		rem -= sp.num() // go on to next span
	}
	panic("Not reached")
}
