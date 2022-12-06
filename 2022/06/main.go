package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.Input("2022/6")

	// Part 1: 1-based index of start of first sequence of 4 unique chars.
	counts := make(map[byte]int)
	start := -1
	for i := range input {
		counts[input[i]]++
		if i >= 4 {
			sub(counts, input[i-4])
		}
		if len(counts) == 4 {
			start = i
			break
		}
	}
	lib.Assert(start >= 0)
	fmt.Println(start + 1)

	// Part 2: 14 unique chars.
	counts2 := make(map[byte]int)
	start2 := -1
	for i := range input {
		counts2[input[i]]++
		if i >= 14 {
			sub(counts2, input[i-14])
		}
		if len(counts2) == 14 {
			start2 = i
			break
		}
	}
	lib.Assert(start2 >= 0)
	fmt.Println(start2 + 1)
}

func sub(m map[byte]int, ch byte) {
	if m[ch] > 1 {
		m[ch]--
	} else {
		delete(m, ch)
	}
}
