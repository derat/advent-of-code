package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var counts [][256]int
	for _, ln := range lib.InputLines("2016/6") {
		if counts == nil {
			counts = make([][256]int, len(ln))
		}
		for i, c := range ln {
			counts[i][byte(c)]++
		}
	}
	for i := range counts {
		var max int
		var best byte
		for ch, cnt := range counts[i] {
			if cnt > max {
				max = cnt
				best = byte(ch)
			}
		}
		fmt.Printf("%c", best)
	}
	fmt.Println()
}
