package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var possible int
	for _, ln := range lib.InputLines("2016/3") {
		sides := lib.ExtractInts(ln)
		lib.AssertEq(len(sides), 3)
		sum := lib.Sum(sides...)
		valid := true
		for _, s := range sides {
			if sum-s <= s { // or sum <= 2*s
				valid = false
				break
			}
		}
		if valid {
			possible++
		}
	}
	fmt.Println(possible)
}
