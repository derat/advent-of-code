package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.ExtractDigits(lib.Input("2019/16"))

	base := []int{0, 1, 0, -1}        // base pattern given in puzzle
	pats := make([][]int, len(input)) // pattern for each digit
	for i := range pats {
		pats[i] = make([]int, len(input))
		rep := i + 1 // number of times each value in base is repeated
		for j := range pats[i] {
			pats[i][j] = base[((j+1)/rep)%len(base)]
		}
	}

	// Phase 1: Print first eight digits after 100 phases.
	vals := append([]int(nil), input...)
	for ph := 0; ph < 100; ph++ {
		next := make([]int, len(vals))
		for i := range vals {
			pat := pats[i]
			var sum int
			for j := range vals {
				sum += vals[j] * pat[j]
			}
			next[i] = lib.Abs(sum) % 10
		}
		vals = next
	}
	for _, v := range vals[:8] {
		fmt.Print(v)
	}
	fmt.Println()
}
