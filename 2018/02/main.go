package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var cnt2, cnt3 int
	for _, ln := range lib.InputLines("2018/2") {
		cnts := make(map[rune]int, 26)
		for _, b := range ln {
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
}
