package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var lines [][]int
	for _, ln := range lib.InputLines("2016/3") {
		lines = append(lines, lib.ExtractInts(ln))
	}

	// Part 1: Each line specifies triangle's sides
	var poss int
	for _, sides := range lines {
		if valid(sides...) {
			poss++
		}
	}
	fmt.Println(poss)

	// Part 2: Each group of 3 going down each column specifies triangle's sides
	var poss2 int
	lib.AssertEq(len(lines)%3, 0)
	for i := 0; i < len(lines); i += 3 {
		for j := 0; j < len(lines[i]); j++ {
			if valid(lines[i][j], lines[i+1][j], lines[i+2][j]) {
				poss2++
			}
		}
	}
	fmt.Println(poss2)
}

func valid(sides ...int) bool {
	lib.AssertEq(len(sides), 3)
	sum := lib.Sum(sides...)
	for _, s := range sides {
		if sum-s <= s { // or sum <= 2*s
			return false
		}
	}
	return true
}
