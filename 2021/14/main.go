package main

import (
	"fmt"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	pgs := lib.InputParagraphs("2021/14")
	lib.AssertEq(len(pgs), 2)
	lib.AssertEq(len(pgs[0]), 1)

	initial := pgs[0][0]
	rules := make(map[string]string)
	for _, ln := range pgs[1] {
		var from, to string
		lib.Extract(ln, `^([A-Z][A-Z]) -> ([A-Z])$`, &from, &to)
		rules[from] = to
	}

	type key struct {
		pair   string // two letters
		rounds int    // number of rounds of expansion
	}
	cache := make(map[key]map[rune]int) // values are counts of inserted characters

	// Returns the counts of inserted characters after expanding the supplied two-letter
	// string n times.
	var expand func(string, int) map[rune]int
	expand = func(from string, n int) map[rune]int {
		lib.AssertEq(len(from), 2)
		lib.AssertLess(0, n)

		// I initially did something like this but without caching (instead just passing
		// a single map down when recursing), but it was too slow for part 2.
		k := key{from, n}
		if c, ok := cache[k]; ok {
			return c
		}

		to, ok := rules[from]
		lib.Assert(ok)
		counts := map[rune]int{rune(to[0]): 1}
		if n > 1 {
			for _, pair := range []string{from[:1] + to, to + from[1:]} {
				for r, n := range expand(pair, n-1) {
					counts[r] += n
				}
			}
		}
		cache[k] = counts
		return counts
	}

	// Returns the diff between the counts of the most-common letter and least-common
	// letter after the specified number of rounds.
	diff := func(rounds int) int {
		counts := make(map[rune]int)
		for _, ch := range initial {
			counts[rune(ch)]++
		}
		for i := range initial[:len(initial)-1] {
			for r, n := range expand(initial[i:i+2], rounds) {
				counts[r] += n
			}
		}
		vals := lib.MapVals(counts)
		sort.Ints(vals)
		return vals[len(vals)-1] - vals[0]
	}

	// Part 1 is 10 rounds, part 2 is 40 rounds.
	fmt.Println(diff(10))
	fmt.Println(diff(40))
}
