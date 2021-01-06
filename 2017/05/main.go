package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	jumps := lib.InputInts("2017/5")
	fmt.Println(run(jumps, false))
	fmt.Println(run(jumps, true))
}

func run(jumps []int, part2 bool) int {
	jumps = append([]int(nil), jumps...) // preserve orig
	var cnt int
	for i := 0; i >= 0 && i < len(jumps); {
		off := jumps[i]
		if part2 && off >= 3 {
			jumps[i]--
		} else {
			jumps[i]++
		}
		i += off
		cnt++
	}
	return cnt
}
