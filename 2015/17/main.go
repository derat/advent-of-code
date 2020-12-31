package main

import (
	"fmt"
	"math/bits"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const total = 150 // given in puzzle
	conts := lib.InputInts("2015/17")
	combos := lib.FindCombos(conts, 1<<len(conts)-1, total)
	fmt.Println(len(combos))

	// Part 2
	minConts := len(conts) + 1 // min number of containers used
	numCombos := 0             // combos using minConts
	for _, combo := range combos {
		if cnt := bits.OnesCount64(combo); cnt < minConts {
			minConts = cnt
			numCombos = 1
		} else if cnt == minConts {
			numCombos++
		}
	}
	fmt.Println(numCombos)
}
