package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var enclosed, overlapping int
	for _, ln := range lib.InputLines("2022/4") {
		var a1, a2, b1, b2 int
		lib.Extract(ln, `^(\d+)-(\d+),(\d+)-(\d+)$`, &a1, &a2, &b1, &b2)
		if (a1 <= b1 && a2 >= b2) || (b1 <= a1 && b2 >= a2) {
			enclosed++
		}
		if (a1 <= b1 && a2 >= b1) || (a1 <= b2 && a2 >= b2) || (a1 >= b1 && a2 <= b2) {
			overlapping++
		}
	}
	fmt.Println(enclosed)
	fmt.Println(overlapping)
}
