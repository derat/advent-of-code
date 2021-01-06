package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	jumps := lib.InputInts("2017/5")
	fmt.Println(run(jumps))
}

func run(jumps []int) int {
	var cnt int
	for i := 0; i >= 0 && i < len(jumps); {
		off := jumps[i]
		jumps[i]++
		i += off
		cnt++
	}
	return cnt
}
