package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	words := lib.InputLines("2018/2")

	// Part 1: Multiple counts of words with two and with three duplicate letters.
	var cnt2, cnt3 int
	for _, word := range words {
		cnts := make(map[rune]int, 26)
		for _, b := range word {
			cnts[b]++
		}
		if lib.MapHasValue(cnts, 2) {
			cnt2++
		}
		if lib.MapHasValue(cnts, 3) {
			cnt3++
		}
	}
	fmt.Println(cnt2 * cnt3)

	// Part 2: Print common portions of two words that differ by exactly one letter.
	seen := make(map[string]struct{})
Loop:
	for _, word := range words {
		for i := range word {
			blank := word[:i] + "_" + word[i+1:]
			if _, ok := seen[blank]; ok {
				fmt.Println(word[:i] + word[i+1:])
				break Loop
			}
			seen[blank] = struct{}{}
		}
	}
}
