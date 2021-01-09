package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// Given in puzzle:
	const (
		mod    = 2147483647
		rounds = 40_000_000
	)
	factors := [2]int64{16807, 48271}

	init := lib.InputInts("2017/15")
	lib.AssertEq(len(init), 2)
	a, b := int64(init[0]), int64(init[1])

	// Part 1: Number of pairs with matching bottom 16 bits.
	// Just using a naive approach here.
	var match int
	for i := 0; i < rounds; i++ {
		a = (a * factors[0]) % mod
		b = (b * factors[1]) % mod
		if a&0xffff == b&0xffff {
			match++
		}
	}
	fmt.Println(match)
}
