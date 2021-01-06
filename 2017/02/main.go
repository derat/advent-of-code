package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var sum, sum2 int
	for _, ln := range lib.InputLines("2017/2") {
		nums := lib.ExtractInts(ln)
		sum += lib.Max(nums...) - lib.Min(nums...)

		// Part 2
	Loop:
		for i, a := range nums {
			for j, b := range nums {
				if i != j && a%b == 0 {
					sum2 += a / b
					break Loop
				}
			}
		}
	}
	fmt.Println(sum)
	fmt.Println(sum2)
}
