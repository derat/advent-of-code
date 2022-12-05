package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	pgs := lib.InputParagraphs("2022/5")
	lib.AssertEq(len(pgs), 2)

	// The first paragraph looks like this:
	//
	//      [D]
	//  [N] [C]
	//  [Z] [M] [P]
	//   1   2   3
	desc := pgs[0]
	lib.AssertLessEq(2, len(desc))
	stackNums := lib.ExtractInts(desc[len(desc)-1])
	stacks := make(map[int][]string)
	for i := len(desc) - 2; i >= 0; i-- {
		ln := desc[i]
		for j, n := range stackNums {
			if crate := string(ln[4*j+1]); crate != " " {
				stacks[n] = append(stacks[n], crate)
			}
		}
	}

	// Make a copy for part 2.
	stacks2 := make(map[int][]string, len(stacks))
	for k, v := range stacks {
		stacks2[k] = append([]string(nil), v...)
	}

	// The second paragraph looks like this:
	//  move 1 from 2 to 1
	//  move 3 from 1 to 3
	//  move 2 from 2 to 1
	//  move 1 from 1 to 2
	for _, ln := range pgs[1] {
		var cnt, src, dst int
		lib.Extract(ln, `^move (\d+) from (\d+) to (\d+)$`, &cnt, &src, &dst)

		// Part 1: crates are popped one at a time.
		for i := 0; i < cnt; i++ {
			idx := len(stacks[src]) - 1
			lib.Assert(idx >= 0)
			stacks[dst] = append(stacks[dst], stacks[src][idx])
			stacks[src] = stacks[src][:idx]
		}

		// Part 2: crates are moved in groups.
		ls := len(stacks2[src])
		lib.AssertLessEq(cnt, ls)
		stacks2[dst] = append(stacks2[dst], stacks2[src][ls-cnt:]...)
		stacks2[src] = stacks2[src][:ls-cnt]
	}

	printTop := func(st map[int][]string) {
		var top string
		for _, n := range stackNums {
			stack := st[n]
			lib.Assert(len(stack) >= 1)
			top += stack[len(stack)-1]
		}
		fmt.Println(top)
	}
	printTop(stacks)
	printTop(stacks2)
}
