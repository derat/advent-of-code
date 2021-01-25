package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.ExtractDigits(lib.Input("2019/16"))

	base := []int{0, 1, 0, -1}        // base pattern given in puzzle
	pats := make([][]int, len(input)) // pattern for each digit
	for i := range pats {
		pats[i] = make([]int, len(input))
		rep := i + 1 // number of times each value in base is repeated
		for j := range pats[i] {
			pats[i][j] = base[((j+1)/rep)%len(base)]
		}
	}

	// Part 1: Print first eight digits after 100 phases.
	vals := append([]int(nil), input...)
	for ph := 0; ph < 100; ph++ {
		next := make([]int, len(vals))
		for i := range vals {
			pat := pats[i]
			var sum int
			for j := range vals {
				sum += vals[j] * pat[j]
			}
			next[i] = lib.Abs(sum) % 10
		}
		vals = next
	}
	for _, v := range vals[:8] {
		fmt.Print(v)
	}
	fmt.Println()

	// Part 2: Repeat input signal 10_000 times and run 100 phases, then print 8 digits at
	// offset specified by first 7 digits of original input.

	// My input is 650 digits, so it gets expanded to 6_500_000 digits.
	// The first 7 digits of my input, representing the message offset, are 5_976_521.
	// The numbers starting at that position will be generated from patterns starting with
	// 5_976_521 zeros followed by 5_976_522 ones. As a result, we can ignore everything
	// before the offset -- it won't have any effect on the digits that we care about.
	// We also don't need to worry about the later 0s or -1s in the pattern, since they
	// would appear beyond the end of the number for these digits.
	var offset int
	for _, v := range input[:7] {
		offset = 10*offset + int(v)
	}
	vals = make([]int, 10_000*len(input)-offset)
	for i := range vals {
		vals[i] = input[(offset+i)%len(input)]
	}

	for ph := 0; ph < 100; ph++ {
		next := make([]int, len(vals))

		// The first number just uses the sum of all the others.
		var sum int
		for _, v := range vals {
			sum += v
		}
		next[0] = lib.Abs(sum) % 10

		// For the following numbers, take the previous sum (before preserving only the last
		// digit) and subtract the value added to it from the previous column. The naive
		// approach of recomputing the sum for each digit is too slow.
		for i := 1; i < len(vals); i++ {
			sum -= vals[i-1]
			next[i] = lib.Abs(sum) % 10
		}
		vals = next
	}
	for _, v := range vals[:8] {
		fmt.Print(v)
	}
	fmt.Println()
}
