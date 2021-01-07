package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var banks []byte
	for _, v := range lib.InputInts("2017/6") {
		lib.AssertLess(v, 256)
		banks = append(banks, byte(v))
	}

	// For part 1, we need to find the number of redistribution cycles before a state is repeated.
	// There are 16 banks holding a total of 111 blocks in my input. Two of the banks start with 15
	// blocks; all the others have fewer. Some banks' block counts will go above 16, so we need more
	// than 4 bits for each (and thus can't fit the full state in 64 bits). Just use strings, I guess.
	seen := make(map[string]int) // values are step counts when state was seen

	var idx int
	var max byte
	for i, v := range banks {
		if v > max {
			idx, max = i, v
		}
	}

	var steps, steps2 int
	for {
		// Bail out if we've already seen this state.
		st := string(banks)
		if _, ok := seen[st]; ok {
			steps2 = steps - seen[st]
			break
		}
		seen[st] = steps

		// Redistribute the blocks.
		amt := banks[idx]
		banks[idx] = 0

		max = 0
		start := idx
		for i := 0; i < len(banks); i++ {
			j := (start - i + len(banks)) % len(banks)
			n := amt / byte(len(banks)-i)
			banks[j] += n
			amt -= n
			// Ties go to the lowest-indexed bank.
			if banks[j] > max || (banks[j] == max && j < idx) {
				idx, max = j, banks[j]
			}
		}

		steps++
	}
	fmt.Println(steps)
	fmt.Println(steps2)
}
