package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var nice, nice2 int
	for _, ln := range lib.InputLines("2015/5") {
		// Part 1:
		var vowels, reps, bad int
		for i, ch := range ln {
			if ch == 'a' || ch == 'e' || ch == 'i' || ch == 'o' || ch == 'u' {
				vowels++
			}
			if i > 0 {
				if ln[i-1] == byte(ch) {
					reps++
				}
				if ls := ln[i-1 : i+1]; ls == "ab" || ls == "cd" || ls == "pq" || ls == "xy" {
					bad++
				}
			}
		}
		if vowels >= 3 && reps > 0 && bad == 0 {
			nice++
		}

		// Part 2:
		var reps2 int                 // repeats with one char between
		pairs := make(map[string]int) // non-overlapping two-char pairs
		var addedPair bool            // added to pairs last time
		for i, ch := range ln {
			if i > 0 {
				pair := ln[i-1 : i+1]
				if pair[0] == pair[1] && i > 1 && ln[i-2] == byte(ch) && addedPair {
					// Avoid adding an overlapping pair for e.g. "aaa".
					addedPair = false
				} else {
					pairs[pair]++
					addedPair = true
				}
			}
			if i > 1 && ln[i-2] == byte(ch) {
				reps2++
			}

		}
		if reps2 > 0 && lib.Max(lib.MapVals(pairs)...) > 1 {
			nice2++
		}
	}
	fmt.Println(nice)
	fmt.Println(nice2)
}
