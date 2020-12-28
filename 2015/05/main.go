package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var nice int
	for _, ln := range lib.InputLines("2015/5") {
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
	}
	fmt.Println(nice)
}
