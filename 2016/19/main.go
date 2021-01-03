package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	min := 1
	max := lib.InputInts("2016/19")[0]
	for step := 2; min < max; step *= 2 {
		// If we ended on the last elf, then it'll take the first elf's presents next round.
		// Otherwise, the second-to-last elf takes the last elf's presents.
		if (max-min)%step == 0 {
			min += step
		} else {
			max -= step / 2
		}
	}
	fmt.Println(min)
}
